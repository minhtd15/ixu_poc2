package service

import (
	_ "database/sql"
	"fmt"
)

func GetPriceEach(productID string) (float64, error) {
	db, err := ConnectToDB()
	fmt.Println("Connected to Oracle")

	// logical solving
	var priceEach float64
	err = db.QueryRow("SELECT priceEach from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)

	if err != nil {
		return 0, fmt.Errorf("Failed to check the price of each product %v", err)
	}

	if err != nil {
		fmt.Printf("foul 22 %v", err)
		return 0, fmt.Errorf("Failed to check quantity in stock: %v", err)
	}

	fmt.Println("we went here")
	return priceEach, err
}
