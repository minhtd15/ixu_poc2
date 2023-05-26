package service

import (
	_ "PRS/entity"
	"database/sql"
	"log"
)

type ProductService struct {
	db *sql.DB
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

func (c *ProductService) GetQuantityInStock(productID string) (int, error) {
	log.Printf("Start to get quantity in stock of product: %v", productID)
	// logical solving
	var inStock int
	err := c.db.QueryRow("SELECT QuantityInStock from SYSTEM.STOCK where productID = ?", productID).Scan(&inStock)

	if err != nil {
		log.Fatalf("failed to check the quantity in stock: %v", err)
		return 0, err
	}

	log.Printf("Get quantity success \nproduct: %v \n quantity", productID, inStock)
	return inStock, err
}

func (c *ProductService) UpdateQuantityInStock(productID string, amountOrder int) error {
	tx, err := c.db.Begin()

	inStock, err := c.GetQuantityInStock(productID)
	if err != nil {
		log.Fatalf("Cannot get the quantity in stock")
		return err
	}
	rs := inStock - amountOrder
	_, err = c.db.Exec("UPDATE SYSTEM.STOCK SET QuantityInStock = ? WHERE productID = ?", rs, productID)
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
