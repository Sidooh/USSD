package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
	"strconv"
)

type Account struct {
	Id     int    `json:"id,omitempty"`
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
	User   *struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
	Balances []Balance
}

type Balance struct {
	Type    string
	Balance string
}

type Invite struct {
	Id      int      `json:"id"`
	Phone   string   `json:"phone"`
	Status  string   `json:"status"`
	Inviter *Account `json:"inviter"`
}

var accountsClient = client.InitAccountClient()
var paymentsClient = client.InitPaymentClient()

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

func CheckAccount(phone string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.GetAccount(phone, &account)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch account", err)
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
