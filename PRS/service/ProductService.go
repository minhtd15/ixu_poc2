package service

import (
	_ "PRS/entity"
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

func GetPriceEach(productID string) (float64, error) {
	log.Printf("Start to get price of the product: %v", productID)

	// logical solving
	var priceEach float64
	err := db.QueryRow("SELECT priceEach from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)

	if err != nil {
		log.Fatalf("Failed to check the price of each product %v", err)
		return 0, err
	}

	log.Printf("Get price success \nproduct: %v \n price: %v \n", productID, priceEach)
	return priceEach, err
}

func GetQuantityInStock(productID string) (int, error) {
	log.Printf("Start to get quantity in stock of product: %v", productID)
	// logical solving
	var inStock int
	err := db.QueryRow("SELECT QuantityInStock from SYSTEM.STOCK where productID = ?", productID).Scan(&inStock)

	if err != nil {
		log.Fatalf("failed to check the quantity in stock: %v", err)
		return 0, err
	}

	log.Printf("Get quantity success \nproduct: %v \n quantity", productID, inStock)
	return inStock, err
}

func UpdateQuantityInStock(productID string, amountOrder int) error {
	tx, err := db.Begin()

	inStock, err := GetQuantityInStock(productID)
	if err != nil {
		fmt.Errorf("failed to get quantity in stock to update the new quantity")
	}
	rs := inStock - amountOrder
	_, err = db.Exec("UPDATE SYSTEM.STOCK SET QuantityInStock = ? WHERE productID = ?", rs, productID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error running the SQL: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error update the quantity in stock: %v", err)
		return err
	}
	return nil
}
