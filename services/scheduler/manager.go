// Copyright 2016-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package scheduler

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/omarghader/jobqueue/mapper"

	"github.com/omarghader/jobqueue/dao"
	"github.com/omarghader/jobqueue/models"
	"github.com/omarghader/jobqueue/protocol"
	"github.com/omarghader/jobqueue/services"

	"github.com/satori/go.uuid"
)

const (
	defaultConcurrency = 5
)

func nop() {}

// Manager schedules job executing. Create a new manager via New.
type manager struct {
	logger  services.Logger
	st      dao.Store // persistent storage
	backoff services.BackoffFunc

	mu          sync.Mutex                    // guards the following block
	tm          map[string]services.Processor // maps topic to processor
	concurrency map[int]int                   // number of parallel workers
	working     map[int]int                   // number of busy workers
	started     bool
	workers     map[int][]*worker
	stopSched   chan struct{} // stop signal for scheduler
	workersWg   sync.WaitGroup
	jobc        map[int]chan *protocol.Job

	testManagerStarted   func() // testing hook
	testManagerStopped   func() // testing hook
	testSchedulerStarted func() // testing hook
	testSchedulerStopped func() // testing hook
	testJobAdded         func() // testing hook
	testJobScheduled     func() // testing hook
	testJobStarted       func() // testing hook
	testJobRetry         func() // testing hook
	testJobFailed        func() // testing hook
	testJobSucceeded     func() // testing hook
}

// New creates a new manager. Pass options to Manager to configure it.
func NewService() services.SchedulerService {
	m := &manager{
		logger: stdLogger{},
		// st:                   mongodb.NewStore(),
		backoff:              exponentialBackoff,
		tm:                   make(map[string]services.Processor),
		concurrency:          map[int]int{0: defaultConcurrency},
		working:              map[int]int{0: 0},
		testManagerStarted:   nop,
		testManagerStopped:   nop,
		testSchedulerStarted: nop,
		testSchedulerStopped: nop,
		testJobAdded:         nop,
		testJobScheduled:     nop,
		testJobStarted:       nop,
		testJobRetry:         nop,
		testJobFailed:        nop,
		testJobSucceeded:     nop,
	}
	// for _, opt := range options {
	// 	opt(m)
	// }
	return m
}

// -- Configuration --

// ManagerOption is the signature of an options provider.
// SetLogger specifies the logger to use when e.g. reporting errors.
func (m *manager) SetLogger(logger services.Logger) {
	m.logger = logger

}

// SetStore specifies the backing Store implementation for the manager.
func (m *manager) SetStore(store dao.Store) {
	m.st = store
}

// SetBackoffFunc specifies the backoff function that returns the time span
// between retries of failed jobs. Exponential backoff is used by default.
func (m *manager) SetBackoffFunc(fn services.BackoffFunc) {
	if fn != nil {
		m.backoff = fn
	} else {
		m.backoff = exponentialBackoff
	}
}

// SetConcurrency sets the maximum number of workers that will be run at
// the same time, for a given rank. Concurrency must be greater or equal
// to 1 and is 5 by default.
func (m *manager) SetConcurrency(rank, n int) {
	if n <= 1 {
		m.concurrency[rank] = 1
	} else {
		m.concurrency[rank] = n
	}

}

// Register registers a topic and the associated processor for jobs with
// that topic.
func (m *manager) Register(topic string, p services.Processor) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.tm[topic]; found {
		return fmt.Errorf("jobqueue: topic %s already registered", topic)
	}
	m.tm[topic] = p
	return nil
}

// -- Start and Stop --

// Start runs the manager. Use Stop, Close, or CloseWithTimeout to stop it.
func (m *manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.started {
		return errors.New("jobqueue: manager already started")
	}

	// Initialize Store
	err := m.st.Start()
	if err != nil {
		return err
	}

	m.jobc = make(map[int]chan *protocol.Job)
	m.workers = make(map[int][]*worker)
	for rank, concurrency := range m.concurrency {
		m.jobc[rank] = make(chan *protocol.Job, concurrency)
		m.workers[rank] = make([]*worker, concurrency)
		for i := 0; i < m.concurrency[rank]; i++ {
			m.workersWg.Add(1)
			m.workers[rank][i] = newWorker(m, m.jobc[rank])
		}
	}

	m.stopSched = make(chan struct{})
	go m.Schedule()

	m.started = true

	m.testManagerStarted() // testing hook

	return nil
}

// Stop stops the manager. It waits for working jobs to finish.
func (m *manager) Stop() error {
	return m.Close()
}

// Close is an alias to Stop. It stops the manager and waits for working
// jobs to finish.
func (m *manager) Close() error {
	return m.CloseWithTimeout(-1 * time.Second)
}

// CloseWithTimeout stops the manager. It waits for the specified timeout,
// then closes down, even if there are still jobs working. If the timeout
// is negative, the manager waits forever for all working jobs to end.
func (m *manager) CloseWithTimeout(timeout time.Duration) error {
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	// Stop accepting new jobs
	m.stopSched <- struct{}{}
	<-m.stopSched
	close(m.stopSched)
	m.mu.Lock()
	for rank := range m.jobc {
		close(m.jobc[rank])
	}
	m.mu.Unlock()

	// Wait for all workers to complete?
	if timeout.Nanoseconds() < 0 {
		// Yes: Wait forever
		m.workersWg.Wait()
		m.testManagerStopped() // testing hook
		return nil
	}

	// Wait with timeout
	complete := make(chan struct{}, 1)
	go func() {
		// Stop workers
		m.workersWg.Wait()
		close(complete)
	}()
	var err error
	select {
	case <-complete: // Completed in time
	case <-time.After(timeout):
		err = errors.New("jobqueue: close timed out")
	}

	m.mu.Lock()
	m.started = false
	m.mu.Unlock()
	m.testManagerStopped() // testing hook
	return err
}

// -- Add --

// Add gives the manager a new job to execute. If Add returns nil, the caller
// can be sure the job is stored in the backing store. It will be picked up
// by the scheduler at a later time.
func (m *manager) Add(job *protocol.Job) error {
	if job.Topic == "" {
		return errors.New("jobqueue: no topic specified")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	job.State = models.Waiting
	job.Retry = 0
	job.Priority = -time.Now().UnixNano()
	if job.FirstStart != 0 {
		if job.FirstStart < time.Now().UnixNano() {
			return errors.New("FirstStart cannot be a past time! ")
		}
		job.Priority = -job.FirstStart
	}
	job.FirstStart = job.Priority
	job.Created = time.Now().UnixNano()

	topics := strings.Split(job.Topic, ",")
	for _, topic := range topics {
		_, found := m.tm[topic]
		if !found {
			return fmt.Errorf("jobqueue: topic %s not registered", topic)
		}
		job.Topic = topic

		job.ID = uuid.NewV4().String()
		// job.State = models.Waiting
		// job.Retry = 0
		// job.Priority = -time.Now().UnixNano()
		// if job.FirstStart != 0 {
		// 	if job.FirstStart < time.Now().UnixNano() {
		// 		return errors.New("FirstStart cannot be a past time! ")
		// 	}
		// 	job.Priority = -job.FirstStart
		// }
		// job.FirstStart = job.Priority
		// job.Created = time.Now().UnixNano()
		j, err := mapper.NewJob(job)
		if err != nil {
			return err
		}
		err = m.st.Create(j)
		if err != nil {
			return err
		}
		m.testJobAdded() // testing hook
	}
	return nil
}

// -- Stats, Lookup and List --

// Stats returns current statistics about the job queue.
func (m *manager) Stats(request *protocol.StatsRequest) (*protocol.Stats, error) {
	req, err := mapper.NewStatusRequest(request)
	if err != nil {
		return nil, err
	}
	jobs, err := m.st.Stats(req)
	if err != nil {
		return nil, err
	}
	return mapper.ToStats(jobs)
}

// Lookup returns the job with the specified identifer.
// If no such job exists, ErrNotFound is returned.
func (m *manager) Lookup(id string) (*protocol.Job, error) {
	job, err := m.st.Lookup(id)
	if err != nil {
		return nil, err
	}
	return mapper.ToJob(job)
}

// LookupByCorrelationID returns the details of jobs by their correlation identifier.
// If no such job could be found, an empty array is returned.
func (m *manager) LookupByCorrelationID(correlationID string) ([]*protocol.Job, error) {

	jobs, err := m.st.LookupByCorrelationID(correlationID)
	if err != nil {
		return nil, err
	}
	return mapper.ToJobs(jobs)
}

// List returns all jobs matching the parameters in the request.
func (m *manager) List(request *protocol.ListRequest) (*protocol.ListResponse, error) {
	req, err := mapper.NewListRequest(request)
	if err != nil {
		return nil, err
	}
	jobs, err := m.st.List(req)
	if err != nil {
		return nil, err
	}
	return mapper.ToListResponse(jobs)
}

// -- Scheduler --

// schedule periodically picks up waiting jobs and passes them to idle workers.
func (m *manager) Schedule() {
	m.testSchedulerStarted()       // testing hook
	defer m.testSchedulerStopped() // testing hook

	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			// Fill up available worker slots with jobs
			for {
				job, err := m.st.Next()
				if err == dao.ErrNotFound {
					break
				}
				if err != nil {
					m.logger.Printf("jobqueue: error picking next job to schedule: %v", err)
					break
				}
				if job == nil {
					break
				}
				m.mu.Lock()
				concurrency := m.concurrency[job.Rank]
				working := m.working[job.Rank]
				m.mu.Unlock()
				if working >= concurrency {
					// All workers busy
					break
				}
				m.mu.Lock()
				rank := job.Rank
				m.working[rank]++
				job.State = models.Working
				job.Started = time.Now().UnixNano()
				err = m.st.Update(job)
				if err != nil {
					m.mu.Unlock()
					m.logger.Printf("jobqueue: error updating job: %v", err)
					break
				}
				m.mu.Unlock()
				m.testJobScheduled()
				j, err := mapper.ToJob(job)
				if err != nil {
					m.logger.Printf("Error on mapping : %+v", err)
					continue
				}
				m.jobc[rank] <- j
			}
		case <-m.stopSched:
			m.stopSched <- struct{}{}
			return
		}
	}
}
