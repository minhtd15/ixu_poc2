package entity

type OrderRequest struct {
	UserID      int    `json:"userID"`
	ProductID   string `json:"productID"`
	AmountOrder int    `json:"amount"`
}
