package client

import (
	"net/http"
	"os"
)

type SavingsApiClient struct {
	ApiClient
}

func InitSavingsClient() *SavingsApiClient {
	client := SavingsApiClient{}
	client.ApiClient.init(os.Getenv("SAVINGS_URL"))
	return &client
}

func (p *SavingsApiClient) FetchAccountEarnings(id string, response interface{}) error {
	endpoint := "/accounts/" + id + "/earnings"

	err := p.newRequest(http.MethodGet, endpoint, nil).send(&response)
	if err != nil {
		return err
	}

	return nil
}
