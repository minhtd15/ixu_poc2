package controller

import (
	"PAS/entity"
	"PAS/service"
	"database/sql"
	"encoding/json"
	"fmt"
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

func (pc *paymentController) PaymentController(w http.ResponseWriter, r *http.Request) {
	order := entity.DeductRequest{}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Fatalf("Error converting json to object: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return failure
	response := entity.PaymentResponse{
		Message: "Order failed",
	}

	// Deduct
	if err := deductBalance(pc, order); err != nil {
		log.Fatalf("Error deducting balance: %v", err)
		http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		return
	}

	// return result
	response = entity.PaymentResponse{
		Message: "Order successful",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error converting response to json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func deductBalance(pc *paymentController, order entity.DeductRequest) error {
	// sử dụng lock để đảm bảo chỉ có một goroutine có thể thực hiện lệnh trừ tiền tại một thời điểm.
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("da tru")
	balance, err := pc.PaymentService.GetBalance(order.UserID)

	fmt.Printf("Khách hàng muốn mua hàng có giá trị %v và số dư tài khoản của khách hàng là %v \n", order.TotalOrder, balance)
	if err != nil {
		log.Fatalf("Cannot find user who has ID: %v", order.UserID)
		return UserNotFound
	}

	if balance < order.TotalOrder {
		log.Fatalf("Do not have enough money for the purchase")
		return InsufficientBalance
	}
	balance -= order.TotalOrder
	fmt.Println(balance)

	balance, err = pc.PaymentService.UpdateBalance(balance, order.UserID)
	if err != nil {
		return err
	}

	return nil
}

var (
	UserNotFound        = &Error{"User not found"}
	InsufficientBalance = &Error{"Insufficient balance"}
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}
