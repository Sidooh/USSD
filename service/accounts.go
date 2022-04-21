package service

import (
	"USSD.sidooh/service/client"
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
	Type   string
	Amount string
}

var accountsClient = client.InitAccountClient()

func FetchAccount(phone string) (*Account, error) {
	var account = new(Account)

	err := accountsClient.GetAccount(phone, &account)
	if err != nil {
		return nil, err
	}

	account.Balances = []Balance{
		{"VOUCHER", "10"},
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
