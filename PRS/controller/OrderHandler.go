package controller

import (
	"PRS/client"
	"PRS/entity"
	"PRS/service"
	"encoding/json"
	"fmt"
	"log"
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
	inStock, err := service.GetQuantityInStock(order.ProductID)
	if err != nil {
		fmt.Errorf("failed to get quantity in stock")
	}

	priceEach, err := service.GetPriceEach(order.ProductID)
	if err != nil {
		fmt.Errorf("failed to get each of product's price")
	}

	// check the quantity in stock
	if order.AmountOrder > inStock {
		result := "The products in stock are not enough"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

	}

	// if there is enough quantity in stock, multiply the total amount of money that the customer want to buy
	var totalOrder = entity.BillRequest{
		UserID:     order.UserID,
		TotalOrder: float64(inStock) * priceEach,
	}

	err = service.UpdateQuantityInStock(order.ProductID, order.AmountOrder)

	// Send message to RabbitMQ
	orderBytes, err := json.Marshal(totalOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.RabbitSender(orderBytes, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msgs := client.ResponseConsumer()

	for msg := range msgs {
		response := entity.PaymentResponse{}
		err := json.Unmarshal(msg.Body, &response)
		if err != nil {
			log.Printf("Failed to unmarshal response: %s", err.Error())
			continue
		}

		// Check the response status
		if response.Success {
			// Payment successful
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Order successful"))
		} else {
			// Payment failed
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Order failed"))
		}

		// Exit the loop since we have received the response
		break
	}
}
