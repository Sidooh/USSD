package client

import (
	"github.com/spf13/viper"
	"net/http"
)

type PaymentsApiClient struct {
	ApiClient
}

type VouchersApiResponse struct {
	ApiResponse
	Data *[]Voucher `json:"data"`
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
