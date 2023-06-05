package service

import (
	"database/sql"
	"fmt"
	"log"
)

type ProductService struct {
	db *sql.DB
}
type UpdateQuantityResult struct {
	QuantityInStock int
	Err             error
}
type CustomerAccountMsg struct {
	TotalMoneyOrdered float64
	Err               error
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

// CheckTotalPurchase CheckCustomerBalance check the balance of the customer's account and return the value of total money ordered by the customer
func (c *ProductService) CheckTotalPurchase(productID string, totalOrder int) (*CustomerAccountMsg, error) {
	log.Printf("Start to get price of the product: %v", productID)

	// logical solving
	var priceEach float64
	err := c.db.QueryRow("SELECT priceEach from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)

	if err != nil {
		log.Fatalf("Failed to check the price of each product %v", err)
		return nil, err
	}

	totalMoneyCustomerOrder := priceEach * float64(totalOrder) // tong so tien ma khach hang order
	if err != nil {
		log.Fatalf("error getting price of each product while checking customer balance")
		return nil, err
	}

	return &CustomerAccountMsg{
		TotalMoneyOrdered: totalMoneyCustomerOrder,
		Err:               nil,
	}, nil
}

func (c *ProductService) UpdateQuantityInStock(productID string, amountOrder int) error {
	log.Printf("Start the transaction")
	tx, err := c.db.Begin()
	if err != nil {
		log.Fatalf("Cannot start a transaction: %v", err)
		return err
	}

	// Lock row
	_, err = tx.Exec("SELECT * FROM SYSTEM.STOCK WHERE productID = ? FOR UPDATE", productID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Cannot lock the row: %v", err)
		return err
	}

	// get the quantity in stock
	var inStock int
	err = c.db.QueryRow("SELECT QuantityInStock from SYSTEM.STOCK where productID = ? FOR UPDATE", productID).Scan(&inStock)

	if amountOrder > inStock {
		log.Fatalf("The products in stock are not enough")
		return fmt.Errorf("The products in stock are not enough")
	}

	_, err = tx.Exec("UPDATE SYSTEM.STOCK SET QuantityInStock = ? WHERE productID = ?", inStock-amountOrder, productID)
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
