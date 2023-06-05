package entity

type DeductRequest struct {
	UserID          int     `json:"userID"`
	TotalMoneyOrder float64 `json:"totalOrder"`
}
