package entity

type DeductRequest struct {
	UserID     int     `json:"userID"`
	TotalOrder float64 `json:"totalOrder"`
}
