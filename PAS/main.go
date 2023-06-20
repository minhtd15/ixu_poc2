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

func main() {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	paymentService := service.NewPaymentService(db)
	paymentController := controller.NewPaymentController(paymentService, db)

	r := mux.NewRouter()
	r.HandleFunc("/payment/deduct-balance", paymentController.PaymentController).Methods("PUT")
	r.HandleFunc("/payment/check-balance", paymentController.BalanceController).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", r))
	log.Printf("Payment completed")

}
