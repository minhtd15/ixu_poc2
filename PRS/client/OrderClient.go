package client

import (
	"PRS/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OrderClient struct {
	BaseURL string
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{
		BaseURL: baseURL,
	}
}

func (oc *OrderClient) DoDeduct(request entity.BillRequest, w http.ResponseWriter) (*http.Response, error) {
	requestBody, _ := json.Marshal(request)
	URL := NewOrderClient(oc.BaseURL)

	url := fmt.Sprintf("%s/payment/deduct-balance", URL.BaseURL)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Cannot connect to client to deduct balance")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Cannot enact PUT request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}

func (oc *OrderClient) CheckBalance(request entity.BillRequest, w http.ResponseWriter) (bool, error) {
	requestBody, _ := json.Marshal(request)
	URL := NewOrderClient(oc.BaseURL)

	url := fmt.Sprintf("%s/payment/check-balance", URL.BaseURL)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Cannot connect to client to deduct balance")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Cannot enact GET request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false, err
	}
	defer resp.Body.Close()

	var response entity.CheckBalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	if response.Status == "Enough balance" {
		return true, nil
	} else {
		return false, nil
	}
}
