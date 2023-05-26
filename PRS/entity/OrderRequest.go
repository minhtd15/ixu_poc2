package entity

import (
	_ "PRS/service"
	_ "database/sql"
	_ "log"
)

type orderRequest struct {
	UserID      int    `json:"userID"`
	ProductID   string `json:"productID"`
	AmountOrder int    `json:"amount"`
}

func NewOrderRequest(UserID int, ProductID string, AmountOrder int) *orderRequest {
	return &orderRequest{
		UserID:      UserID,
		ProductID:   ProductID,
		AmountOrder: AmountOrder,
	}
}

//func (r *orderRequest) GetQuantityInStock(db *sql.DB) (int, error) {
//	inStock, err := service.GetQuantityInStock(r.ProductID, db)
//	if err != nil {
//		log.Fatalf("Error getting quantity in stock")
//	}
//	return inStock, nil
//}
//
//func (r *orderRequest) GetPriceEach(db *sql.DB) (float64, error) {
//	priceEach, err := service.GetPriceEach(r.ProductID, db)
//	if err != nil {
//		log.Fatalf("Error get each price")
//	}
//	return priceEach, err
//}
//
//func (r *orderRequest) UpdateQuantityInStock(db *sql.DB) error {
//	err := service.UpdateQuantityInStock(r.ProductID, r.AmountOrder, db)
//	if err != nil {
//		log.Fatalf("Error update quantity in stock")
//		return err
//	}
//	return nil
//}
