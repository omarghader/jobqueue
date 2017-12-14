package protocol

type Notification struct {
	AmqpExchangeName       string `json:"amqp_exchange_name"`        // internal identifier
	AmqpExchangeRoutingKey string `json:"amqp_exchange_routing_key"` // internal identifier
	WebhookURLs            string `json:"webhook_urls"`              // internal identifier
}
