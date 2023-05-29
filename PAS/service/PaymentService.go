package service

import (
	"database/sql"
	"fmt"
	_ "log"
)

type PaymentService struct {
	db *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}
func (c *PaymentService) GetBalance(userID int) (float64, error) {
	fmt.Println("Connected to Oracle")
	// logical solve
	var balance float64
	err := c.db.QueryRow("select BALANCE from SYSTEM.PAYMENTDB where USERID = ?", userID).Scan(&balance)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if err != nil {
		return 0.0, fmt.Errorf("Failed to commit transaction: %v", err)
	}

	return balance, err
}

func (c *PaymentService) UpdateBalance(balance float64, userID int) (float64, error) {
	tx, err := c.db.Begin()
	_, err = c.db.Exec("UPDATE SYSTEM.PAYMENTDB SET BALANCE = ? WHERE USERID = ?", balance, userID)
	if err != nil {
		tx.Rollback()
		return 0.0, fmt.Errorf("failed to deduct balance: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0.0, fmt.Errorf("Failed to update balance: %v", err)
	}
	return balance, nil
}
