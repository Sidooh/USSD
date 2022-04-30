package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
	"strconv"
)

type Account struct {
	Id     int    `json:"id"`
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

var accountsClient = client.InitAccountClient()
var paymentsClient = client.InitPaymentClient()

func FetchAccount(phone string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.GetAccount(phone, &account)
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
