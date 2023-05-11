package controller

import (
	"fmt"
)

func UpdateBalance(balance float64, userID int) (float64, error) {
	db, _ := ConnectToDB()
	tx, err := db.Begin()
	_, err = db.Exec("UPDATE SYSTEM.PAYMENTDB SET BALANCE = ? WHERE USERID = ?", balance, userID)
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
