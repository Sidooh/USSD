package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type ProductsApiClient struct {
	ApiClient
}

type UtilityAccount struct {
	Id            int
	Provider      string
	AccountNumber string `json:"account_number"`
}

type SubscriptionType struct {
	Id       int
	Title    string
	Price    int
	Duration int
	Active   bool
}

type Subscription struct {
	Id        int
	Amount    string
	Status    string
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type EarningRate struct {
	Type  string
	Value float64
}

//const timeFormat = `2006-01-02 15:04:05`

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

func InitProductClient() *ProductsApiClient {
	client := ProductsApiClient{}
	client.ApiClient.init(os.Getenv("PRODUCTS_URL"))
	client.client.Timeout = 40 * time.Second
	return &client
}

func (p *ProductsApiClient) BuyAirtime(request *PurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/airtime", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) PayUtility(request UtilityPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/utility", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) GetAirtimeAccounts(id string, response interface{}) error {
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

func (p *ProductsApiClient) GetUtilityAccounts(id string, response interface{}) error {
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

func (p *ProductsApiClient) PurchaseVoucher(request *VoucherPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/vouchers/top-up", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) GetSubscription(id string, response interface{}) error {
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

func (p *ProductsApiClient) PurchaseSubscription(request *SubscriptionPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/subscriptions", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) GetSubscriptionType(response interface{}) error {
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, "/products/subscription-types/default", nil).send(apiResponse)
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

func (p *ProductsApiClient) FetchAccountEarnings(id string, response interface{}) error {
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/earnings", nil).send(&apiResponse)
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

func (p *ProductsApiClient) WithdrawEarnings(request *EarningsWithdrawalRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	var response = Response{}
	err = p.newRequest(http.MethodPost, "/products/withdraw", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) FetchEarningRates(response interface{}) error {
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, "/products/earnings/rates", nil).send(apiResponse)
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
