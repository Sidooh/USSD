package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
	"strconv"
)

type Account struct {
	Id       int    `json:"id,omitempty"`
	Phone    string `json:"phone"`
	Active   bool   `json:"active"`
	User     `json:"user"`
	Balances []Balance
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Balance struct {
	Type    string
	Balance float64 `json:"balance,string"`
}

type Invite struct {
	Id      int      `json:"id"`
	Phone   string   `json:"phone"`
	Status  string   `json:"status"`
	Inviter *Account `json:"inviter"`
}

var accountsClient = client.InitAccountClient()
var paymentsClient = client.InitPaymentClient()
var savingsClient = client.InitSavingsClient()

func FetchAccount(phone string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.GetAccountWithUser(phone, &account)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch account", err)
		return nil, err
	}

	if account != nil {
		err = paymentsClient.GetVoucherBalances(strconv.Itoa(account.Id), &account.Balances)
		if err != nil {
			logger.ServiceLog.Error("Failed to fetch voucher balances: ", err)
		}
	}

	return account, nil
}

func FetchInvite(phone string) (*Invite, error) {
	var invite = new(Invite)

	// Check invite existence
	err := accountsClient.CheckInvite(phone, &invite)
	if err != nil {
		return &Invite{}, err
	}

	return invite, nil
}

func CheckAccount(phone string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.GetAccount(phone, &account)
	if err != nil {
		logger.ServiceLog.Error("Failed to check account: ", err)
		return nil, err
	}

	return account, nil
}

func CheckAccountByIdOrPhone(search string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.GetAccountByIdOrPhone(search, &account)
	if err != nil {
		logger.ServiceLog.Error("Failed to search account: ", err)
		return nil, err
	}

	return account, nil
}

func InviteOrAccountExists(phone string) bool {
	var account = new(Account)

	// Check account existence
	err := accountsClient.GetAccount(phone, &account)
	if err != nil && err.Error() != "record not found" {
		logger.ServiceLog.Error("Failed to check invite/account - account: ", err)
	}

	if account.Id != 0 {
		return true
	}

	var invite = new(Invite)

	// Check invite existence
	err = accountsClient.CheckInvite(phone, &invite)
	if err != nil && err.Error() != "record not found" {
		logger.ServiceLog.Error("Failed to check invite/account - invite: ", err)
	}

	if invite.Id != 0 {
		return true
	}

	return false
}

func CheckPin(id string, pin string) bool {
	var valid map[string]string

	err := accountsClient.CheckPin(id, pin, &valid)
	if err != nil {
		return false
	}

	return valid["message"] == "ok"
}

func CheckHasPin(id string) bool {
	var valid map[string]bool

	err := accountsClient.CheckHasPin(id, &valid)
	if err != nil {
		return false
	}

	return valid["message"]
}

func CheckHasSecurityQuestions(id string) bool {
	var valid map[string]bool

	err := accountsClient.CheckHasSecurityQuestions(id, &valid)
	if err != nil {
		return false
	}

	return valid["message"]
}

func CreateAccount(phone string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.CreateAccount(phone, &account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func CreateInvite(id string, phone string) (*Invite, error) {
	var invite = new(Invite)

	err := accountsClient.CreateInvite(id, phone, &invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func SetPin(id string, pin string) bool {
	var valid map[string]string

	err := accountsClient.SetPin(id, pin, &valid)
	if err != nil {
		return false
	}

	return valid["message"] == "ok"
}

func UpdateProfile(id string, details client.ProfileDetails) (User, error) {
	var user = User{}

	err := accountsClient.UpdateProfile(id, details, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

var securityQuestions []client.SecurityQuestion

func FetchSecurityQuestions() ([]client.SecurityQuestion, error) {
	var questions []client.SecurityQuestion

	if len(securityQuestions) > 0 {
		return securityQuestions, nil
	}

	err := accountsClient.FetchSecurityQuestions(&questions)
	if err != nil {
		return nil, err
	}

	securityQuestions = questions

	return questions, nil
}

func SetSecurityQuestions(id string, answers map[string]string) error {
	var results []interface{}

	for i, answer := range answers {
		var res interface{}
		err := accountsClient.SetSecurityQuestion(id, client.SecurityQuestionRequest{
			QuestionId: i,
			Answer:     answer,
		}, &res)
		if err != nil {
			return err
		} else {
			results = append(results, res)
		}
	}

	logger.ServiceLog.Println("Security Questions Request: ", results)
	return nil
}

func FetchUserSecurityQuestions(id string) ([]client.UserSecurityQuestion, error) {
	var questions []client.UserSecurityQuestion

	err := accountsClient.FetchUserSecurityQuestions(id, &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func CheckSecurityQuestionAnswers(id string, answers map[string]string) bool {
	var valid = false

	for i, answer := range answers {
		var res map[string]bool
		err := accountsClient.CheckSecurityQuestionAnswers(id, client.SecurityQuestionRequest{
			QuestionId: i,
			Answer:     answer,
		}, &res)
		if err != nil {
			return false
		} else {
			valid = res["message"]
		}
	}

	return valid
}

func FetchEarningBalances(id string) ([]client.EarningAccount, error) {
	var earnings []client.EarningAccount

	err := savingsClient.FetchAccountEarnings(id, &earnings)
	if err != nil {
		return nil, err
	}

	return earnings, nil
}

func RequestEarningsWithdrawal(request *client.EarningsWithdrawalRequest) error {
	err := savingsClient.WithdrawEarnings(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to withdraw earnings: ", err)
		return err
	}

	return nil
}
