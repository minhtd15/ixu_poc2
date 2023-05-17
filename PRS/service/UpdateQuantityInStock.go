package service

import "fmt"

func UpdateQuantityInStock(productID string, amountOrder int) error {
	db, _ := ConnectToDB()
	tx, err := db.Begin()

	inStock, err := GetQuantityInStock(productID)
	if err != nil {
		fmt.Errorf("failed to get quantity in stock to update the new quantity")
	}
	rs := inStock - amountOrder
	_, err = db.Exec("UPDATE SYSTEM.STOCK SET QuantityInStock = ? WHERE productID = ?", rs, productID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct quantity in stock: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to update quantity in stock: %v", err)
	}
	return nil
}
