package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
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

type Subscription struct {
	Id        int
	Amount    string
	Status    string
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

const timeFormat = `2006-01-02 15:04:05`

//type Time struct {
//	time.Time
//}
//
//func (t Time) UnmarshalJSON(b []byte) error {
//	ret, err := time.Parse(timeFormat, string(b))
//	if err != nil {
//		return err
//	}
//	t = Time{ret}
//	return nil
//}

func InitProductClient() *ProductApiClient {
	client := ProductApiClient{}
	client.ApiClient.init(os.Getenv("PRODUCTS_URL"))
	client.client.Timeout = 40 * time.Second
	return &client
}

func (p *ProductApiClient) BuyAirtime(request *PurchaseRequest) error {
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

func (p *ProductApiClient) GetSubscription(id string, response interface{}) error {
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/current-subscription", nil).send(apiResponse)
	if err != nil {
		return err
	}

	// TODO: Can we get rid of this round trip?
	dbByte, err := json.Marshal(apiResponse.Data)
	err = json.Unmarshal(dbByte, response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductApiClient) PurchaseSubscription(request *SubscriptionPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/subscription", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}
