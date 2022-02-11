package service

import "errors"

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
