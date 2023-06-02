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

	// check the balance of the customer's account and return the value of total money ordered by the customer
	tmp, err := oc.ProductService.CheckCustomerBalance(order.UserID, order.ProductID, order.TotalAmountOrder)
	if err != nil {
		http.Error(w, "Cannot check the balance of the customer to check if it's enough", http.StatusBadRequest)
		return
	}

	// update so luong hang cua khach hang
	err = oc.ProductService.UpdateQuantityInStock(order.ProductID, order.TotalAmountOrder)
	if err != nil {
		err = oc.ProductService.UpdateQuantityInStock(order.ProductID, -order.TotalAmountOrder)
		http.Error(w, "Cannot deduct the amount of products in stock", http.StatusBadRequest)
	}

	totalMoneyOrder := entity.BillRequest{
		order.UserID,
		tmp.TotalMoneyOrdered,
	}

	// connect to client to connect to payment service to subtract the balance in the user's account
	resp, err := oc.OrderClient.DoOrder(totalMoneyOrder, w)
	if err != nil {
		log.Fatalf("Error connecting to payment service")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// if the connection to client failed, the quantity in stock would be rehabilitated as the previous amount
		err = oc.ProductService.UpdateQuantityInStock(order.ProductID, -order.TotalAmountOrder)
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
