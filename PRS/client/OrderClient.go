package client

import (
	"PRS/entity"
	"net/http"
	"strconv"
)

func OrderClient(order entity.BillRequest, w http.ResponseWriter) (*http.Response, error) {
	resp, err := http.Get("http://localhost:8081/payment/deduct/" + strconv.Itoa(order.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return resp, err
}
