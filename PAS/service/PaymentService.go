package service

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB = ConnectToDB()

func ConnectToDB() *sql.DB {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
		return nil
	}
	defer db.Close()

	return db
}

func GetBalance(userID int) (float64, error) {
	fmt.Println("Connected to Oracle")
	// logical solve
	var balance float64
	err := db.QueryRow("select BALANCE from SYSTEM.PAYMENTDB where USERID = ?", userID).Scan(&balance)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if err != nil {
		return 0.0, fmt.Errorf("Failed to commit transaction: %v", err)
	}

	return balance, err
}

func UpdateBalance(balance float64, userID int) (float64, error) {
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
