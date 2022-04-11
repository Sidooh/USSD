package service

import (
	"USSD.sidooh/service/client"
	"fmt"
)

//TODO: Flesh out services i.e. Accounts - fetch an account; Products - perform purchases;

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
	//
	//values := map[string]string{"email": "aa@a.a", "password": "12345678"}
	//jsonData, err := json.Marshal(values)

	var account = new(Account)

	err := accountsClient.GetAccount(phone, &account)
	if err != nil {
		return nil, err
	}

	account.Balances = []Balance{
		{"VOUCHER", "10"},
	}

	fmt.Println(account, account)
	return account, nil
}

func CheckPin(phone string, pin string) bool {
	return pin == "1234"
}
