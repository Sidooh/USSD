package client

import (
	"USSD.sidooh/pkg/cache"
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

type AccountsApiClient struct {
	ApiClient
}

type SecurityQuestion struct {
	Id       uint
	Question string
}

type UserSecurityQuestion struct {
	Id       uint
	Question SecurityQuestion
}

type EarningAccount struct {
	Id     uint `json:"id,string"`
	Type   string
	Self   float64 `json:"self_amount,string"`
	Invite float64 `json:"invite_amount,string"`
}

type SavingAccount struct {
	Id       uint `json:"id,string"`
	Type     string
	Balance  float64 `json:"balance"`
	Interest float64 `json:"interest"`
}

type AccountApiResponse struct {
	ApiResponse
	Data *Account `json:"data"`
}

func InitAccountClient() *AccountsApiClient {
	client := AccountsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_ACCOUNTS_API_URL"))
	return &client
}

func (a *AccountsApiClient) GetAccount(phone string) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/phone/"+phone, nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) GetAccountByIdOrPhone(search string) error {
	var apiResponse = new(ApiResponse)

	// TODO: Test cache

	err := a.newRequest(http.MethodGet, "/accounts/search/id_or_phone?search="+search, nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) GetAccountWithUser(phone string) (*Account, error) {
	var apiResponse = new(AccountApiResponse)
	var account Account

	err := cache.Get("account_"+phone, account)
	if err == nil {
		return &account, nil
	}

	err = a.newRequest(http.MethodGet, "/accounts/phone/"+phone+"?with_user=true", nil).send(apiResponse)
	if err != nil {
		return &Account{}, err
	}

	return apiResponse.Data, nil
}

func (a *AccountsApiClient) CheckInvite(phone string) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/invites/phone/"+phone, nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) CheckPin(id string, pin string) error {
	var apiResponse = new(ApiResponse)

	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/check-pin", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) CheckHasPin(id string) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/has-pin", nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) CheckHasSecurityQuestions(id string) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/has-security-questions", nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) SetPin(id string, pin string) error {
	var apiResponse = new(ApiResponse)

	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/set-pin", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) CreateAccount(phone string, inviteCode interface{}) error {
	var apiResponse = new(ApiResponse)

	values := map[string]interface{}{"phone": phone, "invite_code": inviteCode}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) CreateInvite(id string, phone string) error {
	var apiResponse = new(ApiResponse)

	values := map[string]string{"inviter_id": id, "phone": phone}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/invites", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) UpdateProfile(id string, request ProfileDetails) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/update-profile", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) FetchSecurityQuestions() error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/security-questions", nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) SetSecurityQuestion(id string, request SecurityQuestionRequest) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/security-questions/answers", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) FetchUserSecurityQuestions(id string) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/security-questions", nil).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountsApiClient) CheckSecurityQuestionAnswers(id string, request SecurityQuestionRequest) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/security-questions/check", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	return nil
}
