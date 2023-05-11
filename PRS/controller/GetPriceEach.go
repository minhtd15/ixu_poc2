package controller

import (
	_ "database/sql"
	"fmt"
)

func GetPriceEach(productID string) (float64, error) {
	db, err := ConnectToDB()

	fmt.Println("COnnected to Oracle")
	tx, err := db.Begin()

	// logical solving
	var priceEach float64
	err = db.QueryRow("SELECT PRICEEACH from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("Failed to check the Quantity in Stock: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("Failed to check quantity in stock: %v", err)
	}

	return priceEach, err
}
