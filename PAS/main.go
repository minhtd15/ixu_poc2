package main

import (
	"PAS/client"
	"PAS/entity"
	"PAS/service"
	"database/sql"
	"encoding/json"
	_ "github.com/gorilla/mux"
	"log"
	_ "net/http"
)

var db *sql.DB

func main() {

	msgs := client.RabbitConsumer()

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var req entity.DeductRequest
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Deduct Balance from user's account
			err = service.HandleDeduct(req.UserID, req.TotalOrder)

			// send a success message to RabbitMQ queue
			successResp := entity.PaymentResponse{
				Success: true,
				Message: "Payment processed successfully",
			}

			successRespBytes, err := json.Marshal(successResp)
			if err != nil {
				log.Printf("Failed to marshal response: %v", err)
				continue
			}
			err = client.SendMessage(successRespBytes)
			if err != nil {
				log.Printf("Failed to send success message: %v", err)
				continue
			}
		}
	}()

	log.Printf("Waiting for messages")
	<-forever
}
