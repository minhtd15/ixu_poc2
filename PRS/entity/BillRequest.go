package entity

type BillRequest struct {
	UserID     int     `json:"userID"`
	TotalOrder float64 `json:"totalOrder"`
}
