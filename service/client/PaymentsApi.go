package client

import (
	"net/http"
	"os"
)

type PaymentsApiClient struct {
	ApiClient
}

func InitPaymentClient() *PaymentsApiClient {
	client := PaymentsApiClient{}
	client.ApiClient.init(os.Getenv("PAYMENTS_URL"))
	return &client
}

func (p *PaymentsApiClient) GetVoucherBalances(id string, response interface{}) error {
	endpoint := "/accounts/" + id + "/vouchers"

	err := p.newRequest(http.MethodGet, endpoint, nil).send(&response)
	if err != nil {
		return err
	}

	return nil
}
