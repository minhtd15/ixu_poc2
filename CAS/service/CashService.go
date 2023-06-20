package service

import (
	"database/sql"
	"log"
)

type CashServiceDB struct {
	db *sql.DB
}

func NewCashServiceDB(db *sql.DB) *CashServiceDB {
	return &CashServiceDB{
		db: db,
	}
}

func (c *CashServiceDB) CheckBalance(senderId int, amountMoney float64) (bool, error) {
	// logical solve
	var balance float64
	err := c.db.QueryRow("select BALANCE from SYSTEM.PAYMENTDB where USERID = ?", senderId).Scan(&balance)
	if err != nil {
		log.Printf("Error query row to get balance of customer's account")
		return false, err
	}

	if balance < amountMoney {
		log.Printf("Not enough money in the bank acount to enact transaction: %v", senderId)
		return false, nil
	}
	return true, nil
}

func (c *CashServiceDB) UpdateQuantityInStock(sender int, recipient int, amountMoney float64) error {

}
