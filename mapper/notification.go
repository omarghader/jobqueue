package mapper

import (
	"github.com/omarghader/jobqueue/models"
	"github.com/omarghader/jobqueue/protocol"
)

// ToStats will tranform Model to Protocol
func ToStats(s *models.Stats) (*protocol.Stats, error) {

	return &protocol.Stats{
		Failed:    s.Failed,
		Succeeded: s.Succeeded,
		Waiting:   s.Waiting,
		Working:   s.Working,
	}, nil
}

// NewListResponse will tranform Protocol to Model
func NewStats(s *protocol.Stats) (*models.Stats, error) {
	return &models.Stats{
		Failed:    s.Failed,
		Succeeded: s.Succeeded,
		Waiting:   s.Waiting,
		Working:   s.Working,
	}, nil
}

//-----------------------------------------------

// ToStatsRequest will tranform Model to Protocol
func ToStatsRequest(s *models.StatsRequest) (*protocol.StatsRequest, error) {

	return &protocol.StatsRequest{
		CorrelationGroup: s.CorrelationGroup,
		Topic:            s.Topic,
	}, nil
}

// NewStatusRequest will tranform Protocol to Model
func NewStatusRequest(s *protocol.StatsRequest) (*models.StatsRequest, error) {
	return &models.StatsRequest{
		CorrelationGroup: s.CorrelationGroup,
		Topic:            s.Topic,
	}, nil
}
