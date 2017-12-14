package services

import (
	"time"

	"github.com/omarghader/jobqueue/dao"
	"github.com/omarghader/jobqueue/protocol"
	instabotProtocol "omarghader.com/instagram/instabot/protocol"
)

// BackoffFunc is a callback that returns a backoff. It is configurable
// via the SetBackoff option in the manager. The BackoffFunc is used to
// vary the timespan between retries of failed jobs.
type BackoffFunc func(attempts int) time.Duration

type Processor func(...interface{}) error

type SchedulerService interface {
	// SetLogger specifies the logger to use when e.g. reporting errors.
	SetLogger(logger Logger)

	// SetStore specifies the backing Store implementation for the manager.
	SetStore(store dao.Store)

	// SetBackoffFunc specifies the backoff tion that returns the time span
	// between retries of failed jobs. Exponential backoff is used by default.
	SetBackoffFunc(fn BackoffFunc)

	// SetConcurrency sets the maximum number of workers that will be run at
	// the same time, for a given rank. Concurrency must be greater or equal
	// to 1 and is 5 by default.
	SetConcurrency(rank, n int)

	// Register registers a topic and the associated processor for jobs with
	// that topic.
	Register(topic string, p Processor) error

	// -- Start and Stop --

	// Start runs the manager. Use Stop, Close, or CloseWithTimeout to stop it.
	Start() error

	// Stop stops the manager. It waits for working jobs to finish.
	Stop() error

	// Close is an alias to Stop. It stops the manager and waits for working
	// jobs to finish.
	Close() error

	// CloseWithTimeout stops the manager. It waits for the specified timeout,
	// then closes down, even if there are still jobs working. If the timeout
	// is negative, the manager waits forever for all working jobs to end.
	CloseWithTimeout(timeout time.Duration) error

	// Add gives the manager a new job to execute. If Add returns nil, the caller
	// can be sure the job is stored in the backing store. It will be picked up
	// by the scheduler at a later time.
	Add(job *protocol.Job) error

	// -- Stats, Lookup and List --

	// Stats returns current statistics about the job queue.
	Stats(request *protocol.StatsRequest) (*protocol.Stats, error)

	// Lookup returns the job with the specified identifer.
	// If no such job exists, ErrNotFound is returned.
	Lookup(id string) (*protocol.Job, error)

	// LookupByCorrelationID returns the details of jobs by their correlation identifier.
	// If no such job could be found, an empty array is returned.
	LookupByCorrelationID(correlationID string) ([]*protocol.Job, error)

	// List returns all jobs matching the parameters in the request.
	List(request *protocol.ListRequest) (*protocol.ListResponse, error)

	// -- Scheduler --

	// schedule periodically picks up waiting jobs and passes them to idle workers.
	Schedule()
}

type PublisherService interface {
	PublishRabbitmq(excahngeName, routingKey string, event *instabotProtocol.Event) error
	PublishWebhook(job *protocol.Event) error
}

type Logger interface {
	Printf(format string, v ...interface{})
}
