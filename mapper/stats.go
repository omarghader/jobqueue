package mapper

import (
	"github.com/omarghader/jobqueue/models"
	"github.com/omarghader/jobqueue/protocol"
)

// ToNotification will tranform Model to Protocol
func ToNotification(n *models.Notification) (*protocol.Notification, error) {

	return &protocol.Notification{
		AmqpExchangeName:       n.AmqpExchangeName,
		AmqpExchangeRoutingKey: n.AmqpExchangeRoutingKey,
		WebhookURLs:            n.WebhookURLs,
	}, nil
}

// NewListResponse will tranform Protocol to Model
func NewNotification(n *protocol.Notification) (*models.Notification, error) {
	return &models.Notification{
		AmqpExchangeName:       n.AmqpExchangeName,
		AmqpExchangeRoutingKey: n.AmqpExchangeRoutingKey,
		WebhookURLs:            n.WebhookURLs,
	}, nil
}
