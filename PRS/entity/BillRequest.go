package entity

type BillRequest struct {
	UserID          int     `json:"userID"`
	TotalMoneyOrder float64 `json:"totalMoneyOrder"`
}
