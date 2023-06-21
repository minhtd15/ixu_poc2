package service

import (
	"database/sql"
)

type PaymentHistoryService struct {
	db *sql.DB
}

func NewPaymentHistoryService(db *sql.DB) *PaymentHistoryService {
	return &PaymentHistoryService{
		db: db,
	}
}
