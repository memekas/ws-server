package rabbit

import "github.com/streadway/amqp"

type RabbitMQ struct {
	con *amqp.Connection
}
