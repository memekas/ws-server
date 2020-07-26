package models

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// RabbitMQ connection
type RabbitMQ struct {
	con *amqp.Connection
}

// Init - open connection
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

// Get connection
func (r *RabbitMQ) Get() *amqp.Connection {
	return r.con
}

// Close connection
func (r *RabbitMQ) Close() error {
	return r.con.Close()
}

// CreateQueue -
func (r *RabbitMQ) CreateQueue() (amqp.Queue, error) {
	ch, err := r.con.Channel()
	if err != nil {
		return amqp.Queue{}, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	return q, nil
}

// TryRabbit -
func TryRabbit(log *logrus.Logger) {
	conn, err := amqp.Dial("amqp://admin:mysecretpassword@localhost:5672/notification")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}
