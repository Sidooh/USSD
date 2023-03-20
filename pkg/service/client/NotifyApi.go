package client

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

type NotifyApiClient struct {
	ApiClient
}
type USSDBalanceApiResponse struct {
	ApiResponse
	Data *int `json:"data"`
}

func InitNotifyClient() *NotifyApiClient {
	client := NotifyApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_NOTIFY_API_URL"))
	return &client
}

func (n *NotifyApiClient) SendNotification(request *NotificationRequest) error {
	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = n.newRequest(http.MethodPost, "/notifications", dataBytes).send(nil)
	if err != nil {
		return err
	}

	return nil
}

func (n *NotifyApiClient) GetUSSDBalance() (int, error) {
	res := new(USSDBalanceApiResponse)

	if err := n.newRequest(http.MethodGet, "/dashboard/ussd-balance", nil).send(&res); err != nil {
		return 0, err
	}

	return *res.Data, nil
}
