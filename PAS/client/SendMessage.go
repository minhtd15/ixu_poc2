package client

import (
	"github.com/streadway/amqp"
	"log"
)

func SendMessage(responseBytes []byte) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"payment_response",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        responseBytes,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Sent message to RabbitMQ: ", string(responseBytes))
	return nil
}
