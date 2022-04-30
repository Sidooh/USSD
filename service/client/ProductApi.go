package client

import (
	"bytes"
	"encoding/json"
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

type UtilityAccount struct {
	Id            int
	Provider      string
	AccountNumber string `json:"account_number"`
}

func InitProductClient() *ProductApiClient {
	client := ProductApiClient{}
	client.ApiClient.init(os.Getenv("PRODUCTS_URL"))
	return &client
}

func (p *ProductApiClient) BuyAirtime(request *AirtimePurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/airtime", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductApiClient) PayUtility(request UtilityPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/utility", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductApiClient) GetAirtimeAccounts(id string, response interface{}) error {
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/airtime-accounts", nil).send(apiResponse)
	if err != nil {
		return err
	}

	// TODO: Can we get rid of this round trip?
	dbByte, err := json.Marshal(apiResponse.Data)
	err = json.Unmarshal(dbByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductApiClient) GetUtilityAccounts(id string, response interface{}) error {
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/utility-accounts", nil).send(apiResponse)
	if err != nil {
		return err
	}

	// TODO: Can we get rid of this round trip?
	dbByte, err := json.Marshal(apiResponse.Data)
	err = json.Unmarshal(dbByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductApiClient) PurchaseVoucher(request *VoucherPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/voucher/top-up", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}
