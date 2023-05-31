package service

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

type ProductService struct {
	db   *sql.DB
	lock sync.Mutex
}
type UpdateResult struct {
	QuantityInStock int
	Err             error
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (c *ProductService) GetPriceEach(productID string) (float64, error) {
	log.Printf("Start to get price of the product: %v", productID)

	// logical solving
	var priceEach float64
	err := c.db.QueryRow("SELECT priceEach from SYSTEM.STOCK where productID = ?", productID).Scan(&priceEach)

	if err != nil {
		log.Fatalf("Failed to check the price of each product %v", err)
		return 0, err
	}

	log.Printf("Get price success \nproduct: %v \n price: %v \n", productID, priceEach)
	return priceEach, err
}

func (c *ProductService) UpdateQuantityInStock(productID string, amountOrder int) (*UpdateResult, error) {
	tx, err := c.db.Begin()
	if err != nil {
		log.Fatalf("Cannot start a transaction: %v", err)
		return nil, err
	}

	// Lock row
	_, err = tx.Exec("SELECT * FROM SYSTEM.STOCK WHERE productID = ? FOR UPDATE", productID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Cannot lock the row: %v", err)
		return nil, err
	}

	// get the quantity in stock
	var inStock int
	err = c.db.QueryRow("SELECT QuantityInStock from SYSTEM.STOCK where productID = ? FOR UPDATE", productID).Scan(&inStock)

	if amountOrder > inStock {
		log.Fatalf("The products in stock are not enough")
		return nil, fmt.Errorf("The products in stock are not enough")
	}

	_, err = tx.Exec("UPDATE SYSTEM.STOCK SET QuantityInStock = ? WHERE productID = ?", inStock-amountOrder, productID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error running the SQL: %v", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error committing the transaction: %v", err)
		return nil, err
	}

	return &UpdateResult{QuantityInStock: inStock - amountOrder, Err: nil}, nil
}
