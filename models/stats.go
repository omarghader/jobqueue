package models

// Stats returns statistics about the job queue.
type Stats struct {
	Waiting   int `json:"waiting"`   // number of jobs waiting to be executed
	Working   int `json:"working"`   // number of jobs currently being executed
	Succeeded int `json:"succeeded"` // number of successfully completed jobs
	Failed    int `json:"failed"`    // number of failed jobs (even after retries)
}

// StatsRequest returns information about the number of managed jobs.
type StatsRequest struct {
	Topic            string // filter by topic
	CorrelationGroup string // filter by correlation group
}
