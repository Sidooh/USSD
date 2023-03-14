package client

import (
	"USSD.sidooh/pkg/cache"
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

type ProductsApiClient struct {
	ApiClient
}

type SubscriptionApiResponse struct {
	ApiResponse
	Data *Subscription `json:"data"`
}

type SubscriptionTypeApiResponse struct {
	ApiResponse
	Data *SubscriptionType `json:"data"`
}

type AccountEarningsApiResponse struct {
	ApiResponse
	Data []EarningAccount `json:"data"`
}

type EarningRatesApiResponse struct {
	ApiResponse
	Data *map[string]EarningRate `json:"data"`
}

type UtilityAccountApiResponse struct {
	ApiResponse
	Data []UtilityAccount `json:"data"`
}

func InitProductClient() *ProductsApiClient {
	client := ProductsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_PRODUCTS_API_URL"))
	return &client
}

func (p *ProductsApiClient) BuyAirtime(request *PurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = p.newRequest(http.MethodPost, "/products/airtime", dataBytes).send(nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) PayUtility(request UtilityPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = p.newRequest(http.MethodPost, "/products/utility", dataBytes).send(nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) PayMerchant(request MerchantPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = p.newRequest(http.MethodPost, "/products/merchant", dataBytes).send(nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) GetAirtimeAccounts(id string) ([]UtilityAccount, error) {
	res := new(UtilityAccountApiResponse)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/airtime-accounts", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (p *ProductsApiClient) GetUtilityAccounts(id string) ([]UtilityAccount, error) {
	res := new(UtilityAccountApiResponse)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/utility-accounts", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (p *ProductsApiClient) PurchaseVoucher(request *VoucherPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	// TODO: Check/profile if this is better than new(ApiResponse)
	var response = ApiResponse{}
	err = p.newRequest(http.MethodPost, "/products/voucher", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) GetSubscription(id string) (*Subscription, error) {
	res := new(SubscriptionApiResponse)

	subscription, err := cache.Get[Subscription]("subscription_" + id)
	if err == nil && subscription != nil {
		return subscription, nil
	}

	err = p.newRequest(http.MethodGet, "/accounts/"+id+"/current-subscription", nil).send(res)
	if err != nil {
		return nil, err
	}

	cache.Set("subscription_"+id, res.Data, 24*time.Hour)

	return res.Data, nil
}

func (p *ProductsApiClient) PurchaseSubscription(request *SubscriptionPurchaseRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	cache.Remove("subscription_" + strconv.Itoa(request.AccountId))

	var response = ApiResponse{}
	err = p.newRequest(http.MethodPost, "/products/subscription", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) GetSubscriptionType() (*SubscriptionType, error) {
	res := new(SubscriptionTypeApiResponse)

	subscriptionType, err := cache.Get[SubscriptionType]("default_subscription")
	if err == nil && subscriptionType != nil {
		return subscriptionType, nil
	}

	err = p.newRequest(http.MethodGet, "/subscription-types/default", nil).send(res)
	if err != nil {
		return nil, err
	}

	cache.Set("default_subscription", res.Data, 28*24*time.Hour)

	return res.Data, nil
}

func (p *ProductsApiClient) FetchAccountEarnings(id string) ([]EarningAccount, error) {
	res := new(AccountEarningsApiResponse)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/earnings", nil).send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (p *ProductsApiClient) WithdrawEarnings(request *EarningsWithdrawalRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = p.newRequest(http.MethodPost, "/products/withdraw", dataBytes).send(nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsApiClient) FetchEarningRates() (*map[string]EarningRate, error) {
	res := new(EarningRatesApiResponse)

	err := p.newRequest(http.MethodGet, "/earnings/rates", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
