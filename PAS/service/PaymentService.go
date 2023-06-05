package service

import (
	"database/sql"
	"fmt"
	"log"
	_ "log"
	"sync"
)

type PaymentService struct {
	db   *sql.DB
	lock sync.Mutex
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}
func (c *PaymentService) CheckBalance(userID int, totalMoneyOrder float64) (bool, error) {
	fmt.Println("Connected to Oracle")
	// logical solve
	var balance float64
	err := c.db.QueryRow("select BALANCE from SYSTEM.PAYMENTDB where USERID = ?", userID).Scan(&balance)
	if err != nil {
		log.Printf("Error query row to get balance of customer's account")
		return false, err
	}

	if balance < totalMoneyOrder {
		log.Printf("Not enough money in the bank acount of the customer's: %v", userID)
		return false, nil
	}
	return true, nil
}

func (c *PaymentService) UpdateBalance(userID int, totalMoneyOrder float64) error {
	log.Printf("Start the transaction")
	tx, err := c.db.Begin()
	if err != nil {
		log.Fatalf("Cannot start a transaction: %v", err)
		return err
	}

	// Lock row
	_, err = tx.Exec("SELECT * FROM SYSTEM.PAYMENTDB WHERE USERID = ? FOR UPDATE", userID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Cannot lock the row: %v", err)
		return err
	}

	// get the quantity in stock
	var balance float64
	err = c.db.QueryRow("SELECT BALANCE from SYSTEM.PAYMENTDB where USERID = ? FOR UPDATE", userID).Scan(&balance)

	_, err = tx.Exec("UPDATE SYSTEM.PAYMENTDB SET BALANCE = ? WHERE USERID = ?", balance-totalMoneyOrder, userID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error running the SQL: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error committing the transaction: %v", err)
		return err
	}

	return nil
}
