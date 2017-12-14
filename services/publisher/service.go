package publisher

import (
	"strings"

	"github.com/omarghader/jobqueue/services"

	"github.com/streadway/amqp"
)

type publisherService struct {
	webhookURLs []string
	amqpURI     string
	connection  *amqp.Connection
	channel     *amqp.Channel
}

// NewService returns an publisher publisher service instance
func NewService(amqpURI string, webhookURLs []string) (services.PublisherService, error) {
	srv := &publisherService{
		webhookURLs: webhookURLs,
		amqpURI:     amqpURI,
	}

	if len(strings.TrimSpace(amqpURI)) > 0 {
		// Connect to RabbitMQ
		err := srv.connect()
		if err != nil {
			return nil, err
		}
	}

	return srv, nil
}
