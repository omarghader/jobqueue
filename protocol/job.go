package protocol

// Job is a task that needs to be executed.
type Job struct {
	ID               string        `json:"id"`        // internal identifier
	Topic            string        `json:"topic"`     // topic to find the correct processor
	State            string        `json:"state"`     // current state
	Args             []interface{} `json:"args"`      // arguments to pass to processor
	Rank             int           `json:"rank"`      // jobs with higher ranks get executed earlier
	Priority         int64         `json:"prio"`      // priority (highest gets executed first)
	Retry            int           `json:"retry"`     // current number of retries
	MaxRetry         int           `json:"maxretry"`  // maximum number of retries
	CorrelationGroup string        `json:"cgroup"`    // external group
	CorrelationID    string        `json:"cid"`       // external identifier
	Created          int64         `json:"created"`   // time when Add was called (in UnixNano)
	Started          int64         `json:"started"`   // time when the job was started (in UnixNano)
	Completed        int64         `json:"completed"` // time when job reached either state Succeeded or Failed (in UnixNano)
	LastMod          int64         `json:"last_mod"`
	FirstStart       int64         `json:"first_start"`
	Delay            int64         `json:"delay"`
	Repeat           int64         `json:"repeat"`
	NbSuccess        int64         `json:"nb_sucess"`
	NbError          int64         `json:"nb_error"`
	Notification     Notification  `json:"notification"`
}

// ListRequest specifies a filter for listing jobs.
type ListRequest struct {
	Topic            string `json:"topic"`             // filter by topic
	CorrelationGroup string `json:"correlation_group"` // filter by correlation group
	CorrelationID    string `json:"correlation_id"`    // filter by correlation identifier
	State            string `json:"state"`             // filter by job state
	Limit            int    `json:"limit"`             // maximum number of jobs to return
	Offset           int    `json:"offset"`            // number of jobs to skip (for pagination)
}

// ListResponse is the outcome of invoking List on the Store.
type ListResponse struct {
	Total int   // total number of jobs found, excluding pagination
	Jobs  []Job // list of jobs
}
