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

// CheckCustomerBalance check the balance of the customer's account and return the value of total money ordered by the customer
func (c *ProductService) CheckCustomerBalance(userID int, productID string, totalOrder int) (*CustomerAccountMsg, error) {
	log.Printf("start to check the balance of the customer %v, and quantity order is: %v", userID, totalOrder)

	var balance float64
	err := c.db.QueryRow("select BALANCE from SYSTEM.PAYMENTDB where USERID = ?", userID).Scan(&balance)

	if err != nil {
		// log fail to check customer's balance
		log.Fatalf("Fail to check the customer %v balance: %v", userID, err)
		return nil, err
	}

	priceEach, err := c.getPriceEach(productID)
	totalMoneyCustomerOrder := priceEach * float64(totalOrder) // tong so tien ma khach hang order
	if err != nil {
		log.Fatalf("error getting price of each product while checking customer balance")
		return nil, err
	}

	if balance < totalMoneyCustomerOrder {
		log.Fatalf("Customer %v does not have enough money in the account", userID)
		return nil, fmt.Errorf("customer %v does not have enough money in the account", userID)
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

func (c *ProductService) getPriceEach(productID string) (float64, error) {
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
