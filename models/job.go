package models

const (
	// models.Waiting for executing.
	Waiting string = "waiting"
	// Working is the state for currently executing jobs.
	Working string = "working"
	// Succeeded without errors.
	Succeeded string = "succeeded"
	// Failed even after retries.
	Failed string = "failed"
)

type Job struct {
	ID               string `bson:"_id"`
	Topic            string
	State            string
	Args             *string
	Rank             int
	Priority         int64
	Retry            int
	MaxRetry         int    `bson:"max_retry"`
	CorrelationGroup string `bson:"correlation_group"`
	CorrelationID    string `bson:"correlation_id"`
	Created          int64
	Started          int64
	Completed        int64
	LastMod          int64        `bson:"last_mod"`
	FirstStart       int64        `bson:"first_start"`
	Delay            int64        `bson:"delay"`
	Repeat           int64        `bson:"repeat"`
	NbSuccess        int64        `bson:"nb_sucess"`
	NbError          int64        `bson:"nb_error"`
	Notification     Notification `bson:"notification"`
}

// ListRequest specifies a filter for listing jobs.
type ListRequest struct {
	Topic            string // filter by topic
	CorrelationGroup string // filter by correlation group
	CorrelationID    string // filter by correlation identifier
	State            string // filter by job state
	Limit            int    // maximum number of jobs to return
	Offset           int    // number of jobs to skip (for pagination)
}

// ListResponse is the outcome of invoking List on the Store.
type ListResponse struct {
	Total int   // total number of jobs found, excluding pagination
	Jobs  []Job // list of jobs
}
