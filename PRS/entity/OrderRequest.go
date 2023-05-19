package entity

import (
	"PRS/service"
	"log"
)

type OrderRequest struct {
	UserID      int    `json:"userID"`
	ProductID   string `json:"productID"`
	AmountOrder int    `json:"amount"`
}

func (r OrderRequest) GetQuantityInStock() (int, error) {
	inStock, err := service.GetQuantityInStock(r.ProductID)
	if err != nil {
		log.Fatalf("Error getting quantity in stock")
	}
	return inStock, nil
}

func (r OrderRequest) GetPriceEach() (float64, error) {
	priceEach, err := service.GetPriceEach(r.ProductID)
	if err != nil {
		log.Fatalf("Error get each price")
	}
	return priceEach, err
}

func (r OrderRequest) UpdateQuantityInStock() error {
	err := service.UpdateQuantityInStock(r.ProductID, r.AmountOrder)
	if err != nil {
		log.Fatalf("Error update quantity in stock")
		return err
	}
	return nil
}
