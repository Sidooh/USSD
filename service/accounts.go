package service

import "errors"

//TODO: Flesh out services i.e. Accounts - fetch an account; Products - perform purchases;

type Account struct {
	Phone string
	Name  string
}

func FetchAccount(phone string) (*Account, error) {
	return &Account{
		Phone: phone,
		Name:  "Customer",
	}, errors.New("no Account")
}
