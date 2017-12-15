package scheduler

import (
	"fmt"
	"math"
	"time"

	"github.com/omarghader/jobqueue/mapper"
	"github.com/omarghader/jobqueue/protocol"

	"github.com/omarghader/jobqueue/models"
)

// worker is a single instance processing jobs.
type worker struct {
	m    *manager
	jobc <-chan *protocol.Job
}

// newWorker creates a new worker. It spins up a new goroutine that waits
// on jobc for new jobs to process.
func newWorker(m *manager, jobc <-chan *protocol.Job) *worker {
	w := &worker{m: m, jobc: jobc}
	go w.run()
	return w
}

// run is the main goroutine in the worker. It listens for new jobs, then
// calls process.
func (w *worker) run() {
	defer w.m.workersWg.Done()
	for {
		select {
		case job, more := <-w.jobc:
			if !more {
				// jobc has been closed
				return
			}
			err := w.process(job)
			if err != nil {
				w.m.logger.Printf("jobqueue: job %v failed: %v", job.ID, err)
			}
		}
	}
}

// process runs a single job.
func (w *worker) process(job *protocol.Job) error {
	defer func() {
		w.m.mu.Lock()
		w.m.working[job.Rank]--
		w.m.mu.Unlock()
	}()

	w.processOneTopic(job)

	return nil
}
func (w *worker) processOneTopic(job *protocol.Job) error {
	// Find the topic
	w.m.mu.Lock()
	p, found := w.m.tm[job.Topic]
	w.m.mu.Unlock()
	if !found {
		return fmt.Errorf("no processor found for topic %s", job.Topic)
	}

	w.m.testJobStarted() // testing hook

	args := []interface{}{}
	args = append(args, job.Notification)
	args = append(args, job.Args)
	// Execute the job
	err := p(args...)
	fmt.Println("jobqueue run process")
	if err != nil {
		w.m.logger.Printf("jobqueue: Job %v failed with: %v", job.ID, err)
		job.NbError++

		if job.Repeat > 0 {
			// Successfully restart the job
			w.jobRepeat(job, false) //not infinite Loop
			return nil
		}
		if job.Repeat < 0 {
			// Successfully restart the job infinitely
			w.jobRepeat(job, true)
			return nil
		}

		if job.Retry >= job.MaxRetry {
			// Failed
			w.jobFailed(job)
			return nil
		}

		// Retry
		w.jobRetry(job)
		return nil
	}

	job.NbSuccess++
	if job.Repeat > 0 {
		// Successfully restart the job
		w.jobRepeat(job, false) //not infinite Loop
		return nil
	}
	if job.Repeat < 0 {
		// Successfully restart the job infinitely
		w.jobRepeat(job, true)
		return nil
	}

	w.jobSucceeded(job)
	return nil
}

func (w *worker) jobSucceeded(job *protocol.Job) error {
	// Successfully executed the job
	job.State = models.Succeeded
	job.Completed = time.Now().UnixNano()
	j, err := mapper.NewJob(job)
	if err != nil {
		return err
	}
	err = w.m.st.Update(j)
	if err != nil {
		return err
	}
	w.m.testJobSucceeded()
	return nil
}

func (w *worker) jobFailed(job *protocol.Job) error {
	// Failed
	w.m.testJobFailed() // testing hook
	job.State = models.Failed
	job.Completed = time.Now().UnixNano()
	j, err := mapper.NewJob(job)
	if err != nil {
		return err
	}
	err = w.m.st.Update(j)
	if err != nil {
		return err
	}
	return nil
}

func (w *worker) jobRetry(job *protocol.Job) error {
	w.m.testJobRetry() // testing hook
	job.Priority = -time.Now().Add(w.m.backoff(job.Retry)).UnixNano()
	job.State = models.Waiting
	job.Retry++
	j, err := mapper.NewJob(job)
	if err != nil {
		return err
	}
	err = w.m.st.Update(j)
	if err != nil {
		return err
	}
	return nil
}

func (w *worker) jobRepeat(job *protocol.Job, infiniteLoop bool) error {
	// Successfully restart the job
	job.State = models.Waiting
	// newPriority := -time.Now().Add(time.Duration(job.Delay) * time.Second).UnixNano()
	// newPriority := -time.Unix(0, -job.Priority).Add(time.Duration(job.Delay) * time.Second).UnixNano()
	// fmt.Printf("%d, %d, diff: %d\n", job.Priority, newPriority, job.Priority-newPriority)

	// IF A Pending job is still not completed, then continue
	newPriority := nextRun(job)
	now := time.Now()
	fmt.Println(newPriority, time.Unix(0, newPriority).Sub(now).Seconds())

	if newPriority == 0 {
		w.jobSucceeded(job)
		return nil
	}
	job.Priority = newPriority
	if !infiniteLoop {
		job.Repeat--
	}
	j, err := mapper.NewJob(job)
	if err != nil {
		return err
	}
	err = w.m.st.Update(j)
	if err != nil {
		return err
	}
	w.m.testJobSucceeded()

	return nil
}

func nextRun(job *protocol.Job) int64 {
	//Calculate Next time to run the job
	newPriority := -time.Unix(0, -job.Priority).Add(time.Duration(job.Delay) * time.Second).UnixNano()
	now := time.Now()

	// if nextRun is already Passed
	if now.UnixNano() > -newPriority {
		fmt.Println("newPritory is passed")
		// if job must be repeated
		if job.Repeat > 0 {
			// calculate maximal expectedPrority
			expectedPrority := -time.Unix(0, -job.Priority).Add(time.Duration(job.Delay*job.Repeat) * time.Second).UnixNano()

			// if all iterations time have been passed, then recaulculate
			if now.UnixNano() > -expectedPrority {
				diff := (now.UnixNano() - expectedPrority) / (10 ^ 6*job.Delay)
				fmt.Println("bigger", diff)
				return 0
			} else {
				fmt.Println("JOB REPEAT")

				diff := (time.Unix(0, -expectedPrority).Sub(now).Seconds()) / float64(job.Delay)
				nbSkipped := math.Ceil(diff)
				fmt.Println("less", diff, nbSkipped, (time.Unix(0, -job.Priority).Sub(now).Seconds())/float64(job.Delay))
				effectiveNextPriority := -time.Unix(0, -now.UnixNano()).Add(time.Duration(job.Delay) * time.Second).UnixNano()

				return effectiveNextPriority
			}

		} else {
			// We run the process with a delay next time
			expectedPrority := -now.Add(time.Duration(job.Delay) * time.Second).UnixNano()
			return expectedPrority
		}

	} else {
		fmt.Println("New Priority is not passed")
		return newPriority
	}
}
