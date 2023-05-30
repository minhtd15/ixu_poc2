package main

import (
	"PAS/controller"
	"PAS/service"
	"database/sql"
	_ "github.com/godror/godror"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	/* msgs := client.RabbitConsumer()

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
			err = controller.HandleDeduct(req.UserID, req.TotalOrder)

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
	<-forever */
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	paymentService := service.NewPaymentService(db)

	paymentController := controller.NewPaymentController(paymentService, db)

	r := mux.NewRouter()
	r.HandleFunc("/payment/deduct", paymentController.PaymentController).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", r))
	log.Printf("Payment completed")

}
