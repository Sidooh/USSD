package client

import (
	"github.com/spf13/viper"
	"net/http"
)

type PaymentsApiClient struct {
	ApiClient
}

func InitPaymentClient() *PaymentsApiClient {
	client := PaymentsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_PAYMENTS_API_URL"))
	return &client
}

func (p *PaymentsApiClient) GetVoucherBalances(id string) error {
	endpoint := "/vouchers?account_id=" + id
	apiResponse := new(ApiResponse)

	err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse)
	if err != nil {
		return err
	}

	return nil
}
