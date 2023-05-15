package controller

import (
	_ "database/sql"
	"fmt"
)

func GetPriceEach(productID string) (float64, error) {
	fmt.Println("yess")
	db, err := ConnectToDB()

	fmt.Println("yess yess")
	fmt.Println("Connected to Oracle")
	//tx, err := db.Begin()
	//
	//// logical solving
	//var priceEach float64
	//err = db.QueryRow("SELECT PRICEEACH from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)
	//if err != nil {
	//	fmt.Printf("%v loi day nay", err)
	//	return 0, fmt.Errorf("Failed to check the Quantity in Stock: %v", err)
	//}
	//
	//fmt.Println("yess yess yesss")
	//err = tx.Commit()
	//if err != nil {
	//	return 0, fmt.Errorf("Failed to check quantity in stock: %v", err)
	//}

	tx, err := db.Begin()
	fmt.Println("hello this is MEEE")

	// logical solving
	var priceEach float64
	err = db.QueryRow("SELECT priceEach from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)

	fmt.Println("stocky baby")
	if err != nil {
		return 0, fmt.Errorf("Failed to check the price of each product %v", err)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("foul 22 %v", err)
		return 0, fmt.Errorf("Failed to check quantity in stock: %v", err)
	}

	fmt.Println("we went here")
	return priceEach, err
}
