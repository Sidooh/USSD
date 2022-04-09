package client

import (
	"net/http"
	"os"
)

type AccountApiClient struct {
	ApiClient
}

var apiUrl = os.Getenv("ACCOUNTS_URL")

func getUrl(endpoint string) string {
	return apiUrl + endpoint
}

func InitAccountClient() *AccountApiClient {
	client := AccountApiClient{}
	client.ApiClient.init()
	return &client
}

func (a *AccountApiClient) GetAccount(phone string, response interface{}) error {
	err := a.newRequest(http.MethodGet, getUrl("/accounts/phone/"+phone+"?with_user=true"), nil).send(response)
	if err != nil {
		return err
	}

	return nil
}
