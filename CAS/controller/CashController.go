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
		log.Fatalf("Error checking customer's %v account balance", transaction.SenderId)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// transfer the money to recipient and commission members
	transfer := cc.CashService.TransferMoney(transaction.Sender, transaction.Recipient, transaction.AmountMoney)

}



