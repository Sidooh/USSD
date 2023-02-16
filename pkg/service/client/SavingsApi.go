package client

import (
	"github.com/spf13/viper"
	"net/http"
)

type SavingsApiClient struct {
	ApiClient
}

type SavingsAccountApiResponse struct {
	ApiResponse
	Data []SavingAccount `json:"data"`
}

func InitSavingsClient() *SavingsApiClient {
	client := SavingsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_SAVINGS_API_URL"))
	return &client
}

func (s *SavingsApiClient) FetchAccountSavings(id string) ([]SavingAccount, error) {
	res := new(SavingsAccountApiResponse)
	err := s.newRequest(http.MethodGet, "/accounts/"+id+"/earnings", nil).send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
