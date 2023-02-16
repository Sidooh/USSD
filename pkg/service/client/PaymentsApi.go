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

type VouchersApiResponse struct {
	ApiResponse
	Data *[]Voucher `json:"data"`
}

type ChargesApiResponse struct {
	ApiResponse
	Data *[]AmountCharge `json:"data"`
}

func InitPaymentClient() *PaymentsApiClient {
	client := PaymentsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_PAYMENTS_API_URL"))
	return &client
}

func (p *PaymentsApiClient) GetVouchers(id string) ([]Voucher, error) {
	endpoint := "/vouchers?account_id=" + id
	apiResponse := new(VouchersApiResponse)

	if err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse); err != nil {
		return nil, err
	}

	return *apiResponse.Data, nil
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

func (p *PaymentsApiClient) GetPaybillCharges() ([]AmountCharge, error) {
	endpoint := "/charges/paybill"
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
