package entity

import (
	"PRS/service"
	"database/sql"
	"log"
)

type OrderRequest struct {
	UserID      int    `json:"userID"`
	ProductID   string `json:"productID"`
	AmountOrder int    `json:"amount"`
}

type OrderService interface {
	GetQuantityInStock() (int, error)
	GetPriceEach() (float64, error)
	UpdateQuantityInStock() error
}

func NewOrderRequest(UserID int, ProductID string, AmountOrder int) *OrderRequest {
	return &OrderRequest{
		UserID:      UserID,
		ProductID:   ProductID,
		AmountOrder: AmountOrder,
	}
}

func (r *OrderRequest) GetQuantityInStock(db *sql.DB) (int, error) {
	inStock, err := service.GetQuantityInStock(r.ProductID, db)
	if err != nil {
		log.Fatalf("Error getting quantity in stock")
	}
	return inStock, nil
}

func (r *OrderRequest) GetPriceEach(db *sql.DB) (float64, error) {
	priceEach, err := service.GetPriceEach(r.ProductID, db)
	if err != nil {
		log.Fatalf("Error get each price")
	}
	return priceEach, err
}

func (r *OrderRequest) UpdateQuantityInStock(db *sql.DB) error {
	err := service.UpdateQuantityInStock(r.ProductID, r.AmountOrder, db)
	if err != nil {
		log.Fatalf("Error update quantity in stock")
		return err
	}
	return nil
}
