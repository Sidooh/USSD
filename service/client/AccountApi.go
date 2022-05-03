package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type AccountApiClient struct {
	ApiClient
}

func InitAccountClient() *AccountApiClient {
	client := AccountApiClient{}
	client.ApiClient.init(os.Getenv("ACCOUNTS_URL"))
	return &client
}

func (a *AccountApiClient) GetAccount(phone string, response interface{}) error {
	err := a.newRequest(http.MethodGet, "/accounts/phone/"+phone, nil).send(response)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountApiClient) GetAccountWithUser(phone string, response interface{}) error {
	err := a.newRequest(http.MethodGet, "/accounts/phone/"+phone+"?with_user=true", nil).send(response)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountApiClient) CheckInvite(phone string, response interface{}) error {
	err := a.newRequest(http.MethodGet, "/invites/phone/"+phone, nil).send(response)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountApiClient) CheckPin(id string, pin string, response interface{}) error {
	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/check-pin", dataBytes).send(response)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountApiClient) CreateAccount(phone string, response interface{}) error {
	values := map[string]string{"phone": phone}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts", dataBytes).send(response)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountApiClient) CreateInvite(id string, phone string, response interface{}) error {
	values := map[string]string{"inviter_id": id, "phone": phone}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/invites", dataBytes).send(response)
	if err != nil {
		return err
	}

	return nil
}
