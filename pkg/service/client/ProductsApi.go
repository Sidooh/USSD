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

type SubscriptionTypeApiResponse struct {
	ApiResponse
	Data *SubscriptionType `json:"data"`
}

func InitProductClient() *ProductsApiClient {
	client := ProductsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_PRODUCTS_API_URL"))
	client.client.Timeout = 40 * time.Second
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

func (p *ProductsApiClient) GetAirtimeAccounts(id string, response interface{}) error {
	apiResponse := new(ApiResponse)

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

func (p *ProductsApiClient) GetUtilityAccounts(id string) error {
	apiResponse := new(ApiResponse)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/utility-accounts", nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
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

func (p *ProductsApiClient) GetSubscription(id string, response interface{}) error {
	apiResponse := new(ApiResponse)

	return cached("subscription_"+id, response, func() (interface{}, error) {
		err := p.newRequest(http.MethodGet, "/accounts/"+id+"/current-subscription", nil).send(apiResponse)
		if err != nil {
			return nil, err
		}

		cache.Set("subscription_"+id, response, 24*time.Hour)

		return response, nil
	})
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
	apiResponse := new(SubscriptionTypeApiResponse)
	var subscriptionType *SubscriptionType

	err := cache.Get("default_subscription", subscriptionType)
	if err == nil {
		return subscriptionType, nil
	}

	err = p.newRequest(http.MethodGet, "/subscription-types/default", nil).send(apiResponse)
	if err != nil {
		return &SubscriptionType{}, err
	}

	cache.Set("default_subscription", subscriptionType, 28*24*time.Hour)

	return apiResponse.Data, nil
}

func (p *ProductsApiClient) FetchAccountEarnings(id string) error {
	apiResponse := new(ApiResponse)

	err := p.newRequest(http.MethodGet, "/accounts/"+id+"/earnings", nil).send(&apiResponse)
	if err != nil {
		return err
	}

	return nil
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

func (p *ProductsApiClient) FetchEarningRates() error {
	apiResponse := new(ApiResponse)

	err := p.newRequest(http.MethodGet, "/earnings/rates", nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}
