package client

import (
	"bytes"
	"net/http"
	"os"
)

type ProductApiClient struct {
	ApiClient
}

const (
	CONSUMER = "CONSUMER"
	MPESA    = "MPESA"
	VOUCHER  = "VOUCHER"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func InitProductClient() *ProductApiClient {
	client := ProductApiClient{}
	client.ApiClient.init(os.Getenv("PRODUCTS_URL"))
	return &client
}

func (p *ProductApiClient) BuyAirtime(request AirtimePurchaseRequest) error {
	jsonData, err := request.Marshal()
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/airtime", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}
