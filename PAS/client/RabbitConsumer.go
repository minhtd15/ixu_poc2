package client

import (
	"github.com/streadway/amqp"
	"log"
)

func RabbitConsumer() <-chan amqp.Delivery {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ server: %v", err)
	}
	defer conn.Close()

	// create channel to send and receive message
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open the channel: %v", err)
	}
	defer conn.Close()

	// declare queue to receive message from service order
	q, err := ch.QueueDeclare(
		"payment_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// register consumer for queue and process message when received
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}
	return msgs
}
