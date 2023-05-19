package controller

import (
	"PAS/entity"
	"PAS/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex

//func HandleDeduct(userID int, amount float64) error {
//	// Deduct
//	if err := deductBalance(userID, amount); err != nil {
//		return fmt.Errorf("Failed to deduct and update balance: %v", err)
//	}
//	// return result
//	log.Printf("Deducted %v from account %v", amount, userID)
//	return nil
//}

func HandleDeduct(w http.ResponseWriter, r *http.Request) {
	var req entity.DeductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatalf("Error converting json to object: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Deduct
	if err := deductBalance(req.UserID, req.TotalOrder); err != nil {
		log.Fatalf("Error deducting balance: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return result
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order successful"))
}

func deductBalance(userID int, amount float64) error {
	// sử dụng lock để đảm bảo chỉ có một goroutine có thể thực hiện lệnh trừ tiền tại một thời điểm.
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("da tru")
	balance, err := service.GetBalance(userID)

	fmt.Printf("Khách hàng muốn mua hàng có giá trị %v và số dư tài khoản của khách hàng là %v \n", amount, balance)
	if err != nil {
		return UserNotFound
	}

	if balance < amount {
		return InsufficientBalance
	}
	balance -= amount
	fmt.Println(balance)

	balance, err = service.UpdateBalance(balance, userID)
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
