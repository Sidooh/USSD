package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
)

type AirtimeAccount struct {
	Id            int
	Provider      string
	AccountNumber string
}

var productsClient = client.InitProductClient()

func PurchaseAirtime(request client.AirtimePurchaseRequest) {
	err := productsClient.BuyAirtime(request)
	if err != nil {
		logger.ServiceLog.Panic(err)
	}
}

func FetchAirtimeAccounts(id string) ([]AirtimeAccount, error) {
	var accounts []AirtimeAccount

	account := AirtimeAccount{
		Id:            1,
		Provider:      "AIRTEL",
		AccountNumber: "254780611696",
	}

	accounts = append(accounts,
		account,
		AirtimeAccount{
			Id:            1,
			Provider:      "AIRTEL",
			AccountNumber: "254781611696",
		},
		AirtimeAccount{
			Id:            1,
			Provider:      "AIRTEL",
			AccountNumber: "254782611696",
		},
		AirtimeAccount{
			Id:            1,
			Provider:      "AIRTEL",
			AccountNumber: "254783611696",
		},
		AirtimeAccount{
			Id:            1,
			Provider:      "AIRTEL",
			AccountNumber: "254784611696",
		},
	)

	//err := productsClient.GetAirtimeAccounts(id, &accounts)
	//if err != nil {
	//	return nil, err
	//}

	return accounts, nil
}

func FetchUtilityAccounts(id string) ([]AirtimeAccount, error) {
	var accounts []AirtimeAccount

	account := AirtimeAccount{
		Id:            1,
		Provider:      "KPLC_POSTPAID",
		AccountNumber: "2390423904",
	}

	accounts = append(accounts,
		account,
		AirtimeAccount{
			Id:            2,
			Provider:      "KPLC_PREPAID",
			AccountNumber: "0123912",
		})

	//err := productsClient.GetUtilityAccounts(id, &accounts)
	//if err != nil {
	//	return nil, err
	//}

	return accounts, nil
}
