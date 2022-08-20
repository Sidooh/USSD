package client

import (
	"bytes"
	"encoding/json"
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

func (s *SavingsApiClient) FetchAccountSavings(id string, response interface{}) error {
	endpoint := "/accounts/" + id + "/earnings"

	apiResponse := new(Response)
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

	var response = Response{}
	err = s.newRequest(http.MethodPost, "/personal-accounts/withdraw", dataBytes).send(&response)
	if err != nil {
		return err
	}

	return nil
}
