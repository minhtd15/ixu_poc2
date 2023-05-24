package client

import (
	"PRS/entity"
	"bytes"
	"encoding/json"
	"net/http"
)

func OrderClient(order entity.BillRequest, w http.ResponseWriter) (*http.Response, error) {
	// Chuyển đổi dữ liệu order thành JSON
	requestBody, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// Tạo request POST đến dịch vụ thanh toán
	resp, err := http.Post("http://localhost:8081/payment/deduct", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return resp, err
}
