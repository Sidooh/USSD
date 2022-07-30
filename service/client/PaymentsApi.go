package client

import (
	"encoding/json"
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
	apiResponse := new(Response)

	err := p.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse)
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
