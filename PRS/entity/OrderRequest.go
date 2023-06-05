package entity

import (
	_ "database/sql"
	_ "log"
)

type orderRequest struct {
	UserID     int    `json:"userID"`
	ProductID  string `json:"productID"`
	TotalOrder int    `json:"amount"`
}
