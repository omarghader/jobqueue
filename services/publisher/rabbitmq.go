package publisher

import (
	"fmt"

	// "github.com/omarghader/jobqueue/protocol"
	"github.com/gogo/protobuf/proto"
	"github.com/streadway/amqp"
	instabotProtocol "omarghader.com/instagram/instabot/protocol"
)

// -----------------------------------------------------------------------------
func (s *publisherService) connect() error {
	var err error

	// Connection to RabbitMQ
	s.connection, err = amqp.Dial(s.amqpURI)
	if err != nil {
		return err
	}

	// Create a channel
	s.channel, err = s.connection.Channel()
	if err != nil {
		return err
	}

	// Publish confirm
	err = s.channel.Confirm(false)
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
func (s *publisherService) PublishRabbitmq(excahngeName, routingKey string, job *instabotProtocol.Event) error {

	// Routing key
	// var routingKey = "persist"
	var contentType = "job"

	// // Serialize object
	out, err := proto.Marshal(job)
	if err != nil {
		return err
	}

	// routingKey = strings.ToLower(protocol.EventType_name[int32(job.Type)])

	err = s.channel.Publish(
		excahngeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     fmt.Sprintf("application/x-protobuf;%s", contentType),
			ContentEncoding: "",
			Body:            out,
			DeliveryMode:    amqp.Transient,
			Priority:        0, // 0-9
		})

	if err != nil {
		return err
	}

	return nil
}
