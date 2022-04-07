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

func (a AccountApiClient) GetAccount(phone string, response interface{}) error {
	a.EnsureAuthenticated()

	err := a.newRequest(http.MethodGet, getUrl("/accounts/phone/"+phone+"?with_user=true"), nil).send(response)
	if err != nil {
		return err
	}

	return nil
}
