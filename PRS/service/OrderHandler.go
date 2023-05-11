package service

import (
	"PRS/client"
	"PRS/controller"
	"PRS/entity"
	"encoding/json"
	"net/http"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	order := entity.OrderRequest{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// DB check whether the quantity in stock
	inStock, err := controller.GetQuantityInStock(order.ProductID)
	priceEach, err := controller.GetPriceEach(order.ProductID)

	if order.AmountOrder > inStock {
		w.Write([]byte("The products in stock are not enough"))
		return
	}

	// if there is enough quantity in stock, mutiply the total amount of money that the customer want to buy
	var totalOrder = entity.BillRequest{
		UserID:     order.UserID,
		TotalOrder: float64(inStock) * priceEach,
	}

	// Send message to RabbitMQ
	orderBytes, err := json.Marshal(totalOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.RabbitSender(orderBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return successful
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order successful"))
}
