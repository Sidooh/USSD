package client

import (
	"USSD.sidooh/cache"
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"time"
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

// TODO: Understand and use custom unmarshalers
//type Money struct {
//	float64
//}
//
//func (m Money) UnmarshalJSON(b []byte) error {
//	strippedBytes := b[1 : len(b)-1]
//
//	val, err := strconv.ParseFloat(string(strippedBytes), 64)
//	fmt.Println(val, err)
//
//	if err != nil {
//		return err
//	}
//
//	m = Money{val}
//	return nil
//}

func InitAccountClient() *AccountsApiClient {
	client := AccountsApiClient{}
	client.ApiClient.init(viper.GetString("ACCOUNTS_URL"))
	return &client
}

func (a *AccountsApiClient) GetAccount(phone string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/phone/"+phone, nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) GetAccountByIdOrPhone(search string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	// TODO: Test cache

	err := a.newRequest(http.MethodGet, "/accounts/search/id_or_phone?search="+search, nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) GetAccountWithUser(phone string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := cache.Get("account_"+phone, response)
	if err == nil {
		return nil
	}

	err = a.newRequest(http.MethodGet, "/accounts/phone/"+phone+"?with_user=true", nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	cache.Set("account_"+phone, response, 24*time.Hour)

	return nil
}

func (a *AccountsApiClient) CheckInvite(phone string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/invites/phone/"+phone, nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) CheckPin(id string, pin string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/check-pin", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) CheckHasPin(id string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/has-pin", nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) CheckHasSecurityQuestions(id string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/has-security-questions", nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) SetPin(id string, pin string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/set-pin", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) CreateAccount(phone string, inviteCode interface{}, response interface{}) error {
	var apiResponse = new(ApiResponse)

	values := map[string]interface{}{"phone": phone, "invite_code": inviteCode}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) CreateInvite(id string, phone string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	values := map[string]string{"inviter_id": id, "phone": phone}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/invites", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) UpdateProfile(id string, request ProfileDetails, response interface{}) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/update-profile", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) FetchSecurityQuestions(response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/security-questions", nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) SetSecurityQuestion(id string, request SecurityQuestionRequest, response interface{}) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/security-questions/answers", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) FetchUserSecurityQuestions(id string, response interface{}) error {
	var apiResponse = new(ApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/security-questions", nil).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}

func (a *AccountsApiClient) CheckSecurityQuestionAnswers(id string, request SecurityQuestionRequest, response interface{}) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/security-questions/check", dataBytes).send(apiResponse)
	if err != nil {
		return err
	}

	ConvertStruct(apiResponse.Data, response)

	return nil
}
