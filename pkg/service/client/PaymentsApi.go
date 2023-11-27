package client

import (
	"USSD.sidooh/pkg/cache"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type PaymentsApiClient struct {
	ApiClient
}

type VoucherBalancesApiResponse struct {
	ApiResponse
	Data *[]Balance `json:"data"`
}

type FLoatBalanceApiResponse struct {
	ApiResponse
	Data *Balance `json:"data"`
}

type ChargesApiResponse struct {
	ApiResponse
	Data *[]AmountCharge `json:"data"`
}

type MerchantSearchApiResponse struct {
	ApiResponse
	Data *Merchant `json:"data"`
}

func InitPaymentClient() *PaymentsApiClient {
	client := PaymentsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_PAYMENTS_API_URL"))
	return &client
}

func (p *PaymentsApiClient) GetVoucherBalances(id string) ([]Balance, error) {
	endpoint := "/vouchers?account_id=" + id
	apiResponse := new(VoucherBalancesApiResponse)

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	return *apiResponse.Data, nil
}

func (p *PaymentsApiClient) GetFloatBalance(id string) (*Balance, error) {
	endpoint := "/float-accounts/" + id
	apiResponse := new(FLoatBalanceApiResponse)

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

func (p *PaymentsApiClient) GetWithdrawalCharges() ([]AmountCharge, error) {
	endpoint := "/charges/withdrawal"
	apiResponse := new(ChargesApiResponse)

	charges, err := cache.Get[[]AmountCharge](endpoint)
	if err == nil && len(*charges) > 0 {
		return *charges, nil
	}

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	cache.Set(endpoint, apiResponse.Data, 28*24*time.Hour)

	return *apiResponse.Data, nil
}

func (p *PaymentsApiClient) GetPayBillCharges() ([]AmountCharge, error) {
	endpoint := "/charges/pay-bill"
	apiResponse := new(ChargesApiResponse)

	charges, err := cache.Get[[]AmountCharge](endpoint)
	if err == nil && len(*charges) > 0 {
		return *charges, nil
	}

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	cache.Set(endpoint, apiResponse.Data, 28*24*time.Hour)

	return *apiResponse.Data, nil
}

func (p *PaymentsApiClient) GetBuyGoodsCharges() ([]AmountCharge, error) {
	endpoint := "/charges/buy-goods"
	apiResponse := new(ChargesApiResponse)

	charges, err := cache.Get[[]AmountCharge](endpoint)
	if err == nil && len(*charges) > 0 {
		return *charges, nil
	}

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	cache.Set(endpoint, apiResponse.Data, 28*24*time.Hour)

	return *apiResponse.Data, nil
}

func (p *PaymentsApiClient) GetMpesaWithdrawalCharges() ([]AmountCharge, error) {
	endpoint := "/charges/mpesa-withdrawal"
	apiResponse := new(ChargesApiResponse)

	charges, err := cache.Get[[]AmountCharge](endpoint)
	if err == nil && len(*charges) > 0 {
		return *charges, nil
	}

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	cache.Set(endpoint, apiResponse.Data, 28*24*time.Hour)

	return *apiResponse.Data, nil
}

func (p *PaymentsApiClient) SearchMerchant(code string) (*Merchant, error) {
	endpoint := "/merchants/search/" + code
	apiResponse := new(MerchantSearchApiResponse)

	merchant, err := cache.Get[Merchant](code)
	if merchant != nil {
		return merchant, err
	}

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	cache.Set(endpoint, apiResponse.Data, 180*24*time.Hour)

	return apiResponse.Data, err
}
