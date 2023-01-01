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

func InitSavingsClient() *SavingsApiClient {
	client := SavingsApiClient{}
	client.ApiClient.init(viper.GetString("SAVINGS_URL"))
	return &client
}

func (s *SavingsApiClient) FetchAccountSavings(id string, response interface{}) error {
	endpoint := "/accounts/" + id + "/earnings"

	apiResponse := new(ApiResponse)
	err := s.newRequest(http.MethodGet, endpoint, nil).send(&apiResponse)
	if err != nil {
		return err
	}

	// TODO: Can we get rid of this round trip?
	dbByte, err := json.Marshal(apiResponse.Data)
	err = json.Unmarshal(dbByte, response)
	if err != nil {
		return err
	}

	return nil
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
