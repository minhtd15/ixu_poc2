package controller

import (
	"PAS/entity"
	"PAS/service"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex

type paymentController struct {
	PaymentService *service.PaymentService
	DB             *sql.DB
}

func NewPaymentController(paymentService *service.PaymentService, db *sql.DB) *paymentController {
	return &paymentController{
		PaymentService: paymentService,
		DB:             db,
	}
}

func (pc *paymentController) BalanceController(w http.ResponseWriter, r *http.Request) {
	order := entity.DeductRequest{}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Fatalf("Error converting json to object: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := pc.PaymentService.CheckBalance(order.UserID, order.TotalMoneyOrder)
	if err != nil {
		log.Printf("Error checking balance which involve get balance and compare with the amount of money")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if status {
		response := entity.PaymentResponse{
			Status: "Enough balance",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		response := entity.PaymentResponse{
			Status: "Not enough balance",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (pc *paymentController) PaymentController(w http.ResponseWriter, r *http.Request) {
	order := entity.DeductRequest{}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Fatalf("Error converting json to object: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Deduct
	err := pc.PaymentService.UpdateBalance(order.UserID, order.TotalMoneyOrder)
	if err != nil {
		log.Printf("Error updating balance of customer %v", order.UserID)
		http.Error(w, "Cannot deduct the balance in customer account", http.StatusBadRequest)
		return
	}

	// return successful response for the service order
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payment successful"))
}
