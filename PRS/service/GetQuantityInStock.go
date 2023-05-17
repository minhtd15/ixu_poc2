package service

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
)

func GetQuantityInStock(productID string) (int, error) {
	db, err := ConnectToDB()

	fmt.Println("Connected to Oracle")
	tx, err := db.Begin()

	// logical solving
	var inStock int
	err = db.QueryRow("SELECT QuantityInStock from SYSTEM.STOCK where productID = ?", productID).Scan(&inStock)

	if err != nil {
		//tx.Rollback()
		return 0, fmt.Errorf("Failed to check the Quantity in Stock: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("Failed to check quantity in stock: %v", err)
	}

	return inStock, err
}
func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		fmt.Println(err)
		return db, nil
	}
	defer db.Close()

	return db, nil
}
