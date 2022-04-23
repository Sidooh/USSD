package client

import (
	"net/http"
	"os"
)

type PaymentApiClient struct {
	ApiClient
}

func InitPaymentClient() *PaymentApiClient {
	client := PaymentApiClient{}
	client.ApiClient.init(os.Getenv("PAYMENTS_URL"))
	return &client
}

func (p *PaymentApiClient) GetVoucherBalances(id string, response interface{}) error {
	endpoint := "/accounts/" + id + "/vouchers"

	err := p.newRequest(http.MethodGet, endpoint, nil).send(&response)
	if err != nil {
		return err
	}

	return nil
}
