package controller

import (
	"PAS/service"
	"fmt"
	"log"
)

func HandleDeduct(userID int, amount float64) error {
	// Deduct
	if err := deductBalance(userID, amount); err != nil {
		return fmt.Errorf("Failed to deduct and update balance: %v", err)
	}
	// return result
	log.Printf("Deducted %v from account %v", amount, userID)
	return nil
}

func deductBalance(userID int, amount float64) error {
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
