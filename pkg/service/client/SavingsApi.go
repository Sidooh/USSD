package client

import (
	"bytes"
	"encoding/json"
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

func (s *SavingsApiClient) FetchAccountSavings(id string, response interface{}) ([]SavingAccount, error) {
	endpoint := "/accounts/" + id + "/earnings"

	res := new(SavingsAccountApiResponse)
	err := s.newRequest(http.MethodGet, endpoint, nil).send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (s *SavingsApiClient) WithdrawEarnings(request *EarningsWithdrawalRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = s.newRequest(http.MethodPost, "/personal-accounts/withdraw", dataBytes).send(nil)
	if err != nil {
		return err
	}

	return nil
}
