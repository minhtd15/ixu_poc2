package client

import (
	"PRS/entity"
	"bytes"
	"encoding/json"
	"net/http"
)

type OrderClient struct {
	BaseURL string
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{
		BaseURL: "http://localhost:8081",
	}
}

func (oc *OrderClient) DoOrder(request entity.BillRequest, w http.ResponseWriter) (*http.Response, error) {
	requestBody, _ := json.Marshal(request)
	resp, err := http.Post(oc.BaseURL+"/payment/deduct", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return resp, err
}
