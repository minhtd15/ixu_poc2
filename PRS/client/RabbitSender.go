package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func RabbitSender(orderBytes []byte, w http.ResponseWriter) error {
	fmt.Println("RabbitMQ Connector")
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ server: %v", err)
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ Instance")

	// create channel to send and receive message
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open the channel: %v", err)
		panic(err)
	}
	defer conn.Close()

	// declare queue to receive message from service order
	q, err := ch.QueueDeclare(
		"order_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		panic(err)
	}

	fmt.Println(q)

	// Publish the message to the queue
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        orderBytes,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}

	log.Printf("Sent message: %s", orderBytes)
	return nil
}
