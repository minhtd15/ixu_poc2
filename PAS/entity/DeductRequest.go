package entity

import (
	"PAS/service"
	"fmt"
	"log"
)

type DeductRequest struct {
	UserID     int     `json:"userID"`
	TotalOrder float64 `json:"totalOrder"`
}

func (r *DeductRequest) GetBalance() (float64, error) {
	balance, err := service.GetBalance(r.UserID)
	if err != nil {
		log.Fatalf("Error getting balance for the customer who has ID: %v", r.UserID)
		fmt.Errorf("cannot get the balance")
	}
	return balance, nil
}

func (r *DeductRequest) UpdateBalance(balance float64) (float64, error) {
	rs, err := service.UpdateBalance(balance, r.UserID)
	if err != nil {
		log.Fatalf("Cannot update the balance for the customer that have ID: %v", r.UserID)
		fmt.Errorf("acnnot update the balance")
	}
	return rs, nil
}
