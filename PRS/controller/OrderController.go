package controller

import (
	"PRS/client"
	"PRS/entity"
	"PRS/service"
	"encoding/json"
	"log"
	"net/http"
)

type orderRequestTmp struct {
	UserID           int
	ProductID        string
	TotalAmountOrder int // so luong san pham khach hang order
}

type orderController struct {
	ProductService *service.ProductService
	OrderClient    *client.OrderClient
}

func NewOrderController(productService *service.ProductService, orderClient *client.OrderClient) *orderController {
	return &orderController{
		ProductService: productService,
		OrderClient:    orderClient,
	}
}

func (oc *orderController) OrderController(w http.ResponseWriter, r *http.Request) {
	order := orderRequestTmp{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// calculate the total money that customer purchase by calculate priceEach * numberOrder
	TotalMoneyOrdered, err := oc.ProductService.CheckTotalPurchase(order.ProductID, order.TotalAmountOrder)
	if err != nil {
		http.Error(w, "error calculating the total money that the customer want to purchase", http.StatusBadRequest)
		return
	}

	// get the total value of money that customer ordered to send to payment service
	totalMoneyOrder := entity.BillRequest{order.UserID, TotalMoneyOrdered}

	// call the service payment to check the balance of the customer's account
	hasEnoughBalance, err := oc.OrderClient.CheckBalance(totalMoneyOrder, w)
	if err != nil {
		log.Fatalf("Error checking balance in payment service")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !hasEnoughBalance {
		// Not enough money
		response := entity.CheckBalanceResponse{
			Status: "Not enough balance",
		}
		log.Printf("customer's account does not have enough money")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return

	}

	// if the error occurs, the quantity in stock will be rehabilitated no matter what
	defer func() {
		if err != nil {
			err = oc.ProductService.UpdateQuantityInStock(order.ProductID, -order.TotalAmountOrder)
			if err != nil {
				log.Fatalf("Error updating the quantity in stock: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}()

	// update the quantity in stock
	err = oc.ProductService.UpdateQuantityInStock(order.ProductID, order.TotalAmountOrder)
	if err != nil {
		http.Error(w, "Cannot deduct the amount of products in stock", http.StatusBadRequest)
		return
	}

	// connect to client to connect to payment service to subtract the balance in the user's account
	resp, err := oc.OrderClient.DoDeduct(totalMoneyOrder, w)
	if err != nil {
		// if the connection to client failed, the quantity in stock would be rehabilitated as the previous amount
		err = oc.ProductService.UpdateQuantityInStock(order.ProductID, -order.TotalAmountOrder)
		log.Fatalf("Error connecting to payment service")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		// Rollback transaction after the failure in subtract in customer's account
		err = oc.ProductService.UpdateQuantityInStock(order.ProductID, -order.TotalAmountOrder)
		if err != nil {
			log.Fatalf("Error update the quantity in stock after failed to deduct the balance in customer's account")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Not enough balance
		errorMessage := "Insufficient balance"
		http.Error(w, errorMessage, http.StatusPaymentRequired)

		return
	}
}
