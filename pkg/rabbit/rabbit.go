package rabbit

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func (r *RabbitMQ) Init() error {
	rURI := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		os.Getenv("RABBITMQ_DEFAULT_USER"),
		os.Getenv("RABBITMQ_DEFAULT_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		os.Getenv("RABBITMQ_DEFAULT_VHOST"),
	)
	con, err := amqp.Dial(rURI)
	if err != nil {
		return err
	}
	r.con = con
	return nil
}

func (r *RabbitMQ) Get() *amqp.Connection {
	return r.con
}

func (r *RabbitMQ) Close() error {
	return r.con.Close()
}
