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
	client.ApiClient.init(viper.GetString("PAYMENTS_URL"))
	return &client
}

func (p *PaymentsApiClient) GetVoucherBalances(id string, response interface{}) error {
	endpoint := "/vouchers?account_id=" + id
	apiResponse := new(ApiResponse)

	err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}
