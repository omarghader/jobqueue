package mapper

import (
	"encoding/json"

	"github.com/omarghader/jobqueue/models"
	"github.com/omarghader/jobqueue/protocol"
)

// ToJob will tranform Model to Protocol
func ToJob(j *models.Job) (*protocol.Job, error) {
	var args []interface{}
	if j.Args != nil && *j.Args != "" {
		if err := json.Unmarshal([]byte(*j.Args), &args); err != nil {
			return nil, err
		}
	}

	notif, err := ToNotification(&j.Notification)
	if err != nil {
		return nil, err
	}
	job := &protocol.Job{
		ID:               j.ID,
		Topic:            j.Topic,
		State:            j.State,
		Args:             args,
		Rank:             j.Rank,
		Priority:         j.Priority,
		Retry:            j.Retry,
		MaxRetry:         j.MaxRetry,
		CorrelationGroup: j.CorrelationGroup,
		CorrelationID:    j.CorrelationID,
		Created:          j.Created,
		Started:          j.Started,
		Completed:        j.Completed,
		LastMod:          j.LastMod,
		FirstStart:       j.FirstStart,
		Delay:            j.Delay,
		Repeat:           j.Repeat,
		NbSuccess:        j.NbSuccess,
		NbError:          j.NbError,
		Notification:     *notif,
	}
	return job, nil
}

// NewJob will tranform Protocol to Model
func NewJob(job *protocol.Job) (*models.Job, error) {
	var args *string
	if job.Args != nil {
		v, err := json.Marshal(job.Args)
		if err != nil {
			return nil, err
		}
		s := string(v)
		args = &s
	}

	notif, err := NewNotification(&job.Notification)
	if err != nil {
		return nil, err
	}
	return &models.Job{
		ID:               job.ID,
		Topic:            job.Topic,
		State:            job.State,
		Args:             args,
		Rank:             job.Rank,
		Priority:         job.Priority,
		Retry:            job.Retry,
		MaxRetry:         job.MaxRetry,
		CorrelationGroup: job.CorrelationGroup,
		CorrelationID:    job.CorrelationID,
		Created:          job.Created,
		Started:          job.Started,
		Completed:        job.Completed,
		LastMod:          job.LastMod,
		FirstStart:       job.FirstStart,
		Delay:            job.Delay,
		Repeat:           job.Repeat,
		NbSuccess:        job.NbSuccess,
		NbError:          job.NbError,
		Notification:     *notif,
	}, nil
}

// ToJobs will tranform Model to Protocol
func ToJobs(j []*models.Job) ([]*protocol.Job, error) {
	jobs := []*protocol.Job{}
	for _, v := range j {
		job, err := ToJob(v)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// NewJobs will tranform Model to Protocol
func NewJobs(j []*protocol.Job) ([]*models.Job, error) {
	jobs := []*models.Job{}
	for _, v := range j {
		job, err := NewJob(v)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

//---------------------------------------------------------------------

// ToListResponse will tranform Model to Protocol
func ToListResponse(res *models.ListResponse) (*protocol.ListResponse, error) {

	array := []protocol.Job{}
	for _, job := range res.Jobs {
		j, err := ToJob(&job)
		if err != nil {
			return nil, err
		}
		array = append(array, *j)
	}
	return &protocol.ListResponse{
		Jobs:  array,
		Total: res.Total,
	}, nil
}

// NewListResponse will tranform Protocol to Model
func NewListResponse(res *protocol.ListResponse) (*models.ListResponse, error) {
	array := []models.Job{}
	for _, job := range res.Jobs {
		j, err := NewJob(&job)
		if err != nil {
			return nil, err
		}
		array = append(array, *j)
	}
	return &models.ListResponse{
		Jobs:  array,
		Total: res.Total,
	}, nil
}

///---------------------------------------------------------------

// ToListRequest will tranform Model to Protocol
func ToListRequest(req *models.ListRequest) (*protocol.ListRequest, error) {

	return &protocol.ListRequest{
		CorrelationGroup: req.CorrelationGroup,
		CorrelationID:    req.CorrelationID,
		Limit:            req.Limit,
		Offset:           req.Offset,
		State:            req.State,
		Topic:            req.Topic,
	}, nil
}

// NewListRequest will tranform Protocol to Model
func NewListRequest(req *protocol.ListRequest) (*models.ListRequest, error) {
	return &models.ListRequest{
		CorrelationGroup: req.CorrelationGroup,
		CorrelationID:    req.CorrelationID,
		Limit:            req.Limit,
		Offset:           req.Offset,
		State:            req.State,
		Topic:            req.Topic,
	}, nil
}
