package client

import (
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

type UserApiResponse struct {
	ApiResponse
	Data *User `json:"data"`
}

type AccountApiResponse struct {
	ApiResponse
	Data *Account `json:"data"`
}

type InviteApiResponse struct {
	ApiResponse
	Data *Invite `json:"data"`
}

type SecurityQuestionsApiResponse struct {
	ApiResponse
	Data []SecurityQuestion `json:"data"`
}

type UserSecurityQuestionsApiResponse struct {
	ApiResponse
	Data []UserSecurityQuestion `json:"data"`
}

type CheckSecurityQuestionAnswersApiResponse struct {
	ApiResponse
	Data bool `json:"data"`
}

type SetPinApiResponse struct {
	ApiResponse
	Data bool `json:"data"`
}

type CheckHasPinApiResponse struct {
	ApiResponse
	Data bool `json:"data"`
}

type CheckPinApiResponse struct {
	ApiResponse
	Data bool `json:"data"`
}

type HasSecurityQuestionsApiResponse struct {
	ApiResponse
	Data bool `json:"data"`
}

func InitAccountClient() *AccountsApiClient {
	client := AccountsApiClient{}
	client.ApiClient.init(viper.GetString("SIDOOH_ACCOUNTS_API_URL"))
	return &client
}

func (a *AccountsApiClient) GetAccount(phone string) (*Account, error) {
	var res = new(AccountApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/phone/"+phone, nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) GetAccountByIdOrPhone(search string) (*Account, error) {
	var res = new(AccountApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/search/id_or_phone?search="+search, nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) GetAccountWithUser(phone string) (*Account, error) {
	var res = new(AccountApiResponse)

	// TODO: Only add cache once we can invalidate cache from UI/endpoint
	err := a.newRequest(http.MethodGet, "/accounts/phone/"+phone+"?with_user=true", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CheckInvite(phone string) (*Invite, error) {
	var res = new(InviteApiResponse)

	err := a.newRequest(http.MethodGet, "/invites/phone/"+phone, nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CheckPin(id, pin string) (bool, error) {
	var res = new(CheckPinApiResponse)

	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/check-pin", dataBytes).send(res)
	if err != nil {
		return false, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CheckHasPin(id string) (bool, error) {
	var res = new(CheckHasPinApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/has-pin", nil).send(res)
	if err != nil {
		return false, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CheckHasSecurityQuestions(id string) (bool, error) {
	var res = new(HasSecurityQuestionsApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/has-security-questions", nil).send(res)
	if err != nil {
		return false, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) SetPin(id string, pin string) (bool, error) {
	var res = new(SetPinApiResponse)

	values := map[string]string{"pin": pin}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/set-pin", dataBytes).send(res)
	if err != nil {
		return false, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CreateAccount(phone string, inviteCode interface{}) (*Account, error) {
	var res = new(AccountApiResponse)

	values := map[string]interface{}{"phone": phone, "invite_code": inviteCode}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts", dataBytes).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CreateInvite(id, phone, inviteType string) (*Invite, error) {
	var res = new(InviteApiResponse)

	values := map[string]string{"inviter_id": id, "phone": phone, "type": inviteType}
	jsonData, err := json.Marshal(values)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/invites", dataBytes).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) UpdateProfile(id string, request ProfileDetails) (*User, error) {
	var res = new(UserApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/update-profile", dataBytes).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) FetchSecurityQuestions() ([]SecurityQuestion, error) {
	var res = new(SecurityQuestionsApiResponse)

	err := a.newRequest(http.MethodGet, "/security-questions", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) SetSecurityQuestion(id string, request SecurityQuestionRequest) (interface{}, error) {
	var res = new(ApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/security-question-answers", dataBytes).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) FetchUserSecurityQuestions(id string) ([]UserSecurityQuestion, error) {
	var res = new(UserSecurityQuestionsApiResponse)

	err := a.newRequest(http.MethodGet, "/accounts/"+id+"/security-question-answers", nil).send(res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (a *AccountsApiClient) CheckSecurityQuestionAnswers(id string, request SecurityQuestionRequest) (bool, error) {
	var res = new(CheckSecurityQuestionAnswersApiResponse)

	jsonData, err := json.Marshal(request)
	dataBytes := bytes.NewBuffer(jsonData)

	err = a.newRequest(http.MethodPost, "/accounts/"+id+"/security-question-answers/check", dataBytes).send(res)
	if err != nil {
		return false, err
	}

	return res.Data, nil
}
