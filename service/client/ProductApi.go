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

func (p *ProductApiClient) GetAirtimeAccounts(id string, response interface{}) error {
	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/airtime-accounts", nil).send(response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductApiClient) GetUtilityAccounts(id string, response interface{}) error {
	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/utility-accounts", nil).send(response)
	if err != nil {
		return err
	}

	return nil
}
