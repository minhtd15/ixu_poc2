package controller

import (
	"CAS/dto"
	"CAS/service"
	"encoding/json"
	"log"
	"net/http"
)

type cashController struct {
	CashService *service.CashServiceDB
}

func NewCashController(cashService *service.CashServiceDB) *cashController {
	return &cashController{
		CashService: cashService,
	}
}

func (cc *cashController) CashController(w http.ResponseWriter, r *http.Request) {
	// convert json request to transaction struct
	transaction := dto.TransactionRequestDto{}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Fatalf("Error converting json to transaction object")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checkAccountBalance, err := cc.CashService.CheckBalance(transaction.SenderId, transaction.AmountMoney)
	if err != nil {
		response := "Error in checking the customer balance"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		log.Fatalf("Error checking customer's %v account balance", transaction.SenderId)
	}

	if !checkAccountBalance {
		response := "Not enough balance"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	// transfer the money to recipient and commission members
	transfer := cc.CashService.TransferMoney(transaction.Sender, transaction.Recipient, transaction.AmountMoney)

}
