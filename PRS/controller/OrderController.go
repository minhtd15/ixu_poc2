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
	UserID      int
	ProductID   string
	AmountOrder int
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

	orderTmp := entity.NewOrderRequest(order.UserID, order.ProductID, order.AmountOrder)

	// DB check whether the quantity in stock
	inStock, err := oc.ProductService.GetQuantityInStock(orderTmp.ProductID)
	if err != nil {
		log.Fatalf("failed to get quantity in stock: %v", err)
	}

	// check the quantity in stock
	if orderTmp.AmountOrder > inStock {
		result := "The products in stock are not enough"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	priceEach, err := oc.ProductService.GetPriceEach(orderTmp.ProductID)
	if err != nil {
		log.Fatalf("failed to get correspond price of the product: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if there is enough quantity in stock, multiply the total amount of money that the customer want to buy
	var totalOrder = entity.BillRequest{
		UserID:     order.UserID,
		TotalOrder: float64(inStock) * priceEach,
	}

	// connect to client to connect to payment service to subtract the balance in the user's account
	resp, err := oc.OrderClient.DoOrder(totalOrder, w)
	if err != nil {
		log.Fatalf("Error connecting to payment service")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update the quantity in stock
	if err := oc.ProductService.UpdateQuantityInStock(orderTmp.ProductID, inStock-orderTmp.AmountOrder); err != nil {
		log.Fatalf("Error update the quantity in stock")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if resp.StatusCode == http.StatusPaymentRequired {
		// Rollback transaction
		if err := oc.ProductService.UpdateQuantityInStock(orderTmp.ProductID, inStock+orderTmp.AmountOrder); err != nil {
			log.Fatalf("Error update the quantity in stock after failed to deduct the balance in customer's account")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Not enough balance
		errorMessage := "Insufficient balance"
		http.Error(w, errorMessage, http.StatusPaymentRequired)

		return
	}
	/*// Send message to RabbitMQ
	orderBytes, err := json.Marshal(totalOrder)
	if err != nil {
		log.Fatalf()
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	} */
}
