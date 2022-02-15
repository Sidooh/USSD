package service

import "errors"

//TODO: Flesh out services i.e. Accounts - fetch an account; Products - perform purchases;

type Account struct {
	Phone    string
	Name     string
	Balances []Balance
}

type Balance struct {
	Type   string
	Amount string
}

func FetchAccount(phone string) (*Account, error) {
	return &Account{
		Phone: phone,
		Name:  "Customer",
		Balances: []Balance{
			{"VOUCHER", "10"},
		},
	}, errors.New("no Account")
}
