package controller

import (
	"PAS/entity"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex

type deductRequestTmp struct {
	UserID     int
	TotalOrder float64
}

var reqTmp = deductRequestTmp{}
var req = entity.NewDeductRequest(reqTmp.UserID, reqTmp.TotalOrder)

func HandleDeduct(w http.ResponseWriter, r *http.Request) {

	if err := json.NewDecoder(r.Body).Decode(&reqTmp); err != nil {
		log.Fatalf("Error converting json to object: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Deduct
	if err := deductBalance(); err != nil {
		log.Fatalf("Error deducting balance: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return result
	response := entity.PaymentResponse{
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

func deductBalance() error {
	// sử dụng lock để đảm bảo chỉ có một goroutine có thể thực hiện lệnh trừ tiền tại một thời điểm.
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("da tru")
	balance, err := req.GetBalance()

	fmt.Printf("Khách hàng muốn mua hàng có giá trị %v và số dư tài khoản của khách hàng là %v \n", req.TotalOrder, balance)
	if err != nil {
		log.Fatalf("Cannot find user who has ID: %v", req.UserID)
		return UserNotFound
	}

	if balance < req.TotalOrder {
		log.Fatalf("Do no have enough money for the purchase")
		return InsufficientBalance
	}
	balance -= req.TotalOrder
	fmt.Println(balance)

	balance, err = req.UpdateBalance(balance)
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
