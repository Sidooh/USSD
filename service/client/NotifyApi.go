package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type NotifyApiClient struct {
	ApiClient
}

func InitNotifyClient() *NotifyApiClient {
	client := NotifyApiClient{}
	client.ApiClient.init(os.Getenv("NOTIFY_URL"))
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
