package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
	"fmt"
)

var productsClient = client.InitProductClient()

func PurchaseAirtime(request client.AirtimePurchaseRequest) {
	err := productsClient.BuyAirtime(request)
	if err != nil {
		logger.ServiceLog.Error(err)
	}
}

func PurchaseUtility(request client.UtilityPurchaseRequest) {
	err := productsClient.PayUtility(request)
	if err != nil {
		logger.ServiceLog.Error(err)
	}
}

func FetchAirtimeAccounts(id string) ([]client.UtilityAccount, error) {
	var accounts []client.UtilityAccount

	err := productsClient.GetAirtimeAccounts(id, &accounts)
	if err != nil {
		logger.ServiceLog.Error(err)
		return nil, err
	}

	fmt.Println(accounts)

	return accounts, nil
}

func FetchUtilityAccounts(id string, provider string) ([]client.UtilityAccount, error) {
	var accounts []client.UtilityAccount

	err := productsClient.GetUtilityAccounts(id, &accounts)
	if err != nil {
		logger.ServiceLog.Error(err)
		return nil, err
	}

	// TODO: Can we cache the data so we don't have to re-fetch full accounts
	var providerAccounts []client.UtilityAccount
	for _, account := range accounts {

		if account.Provider == provider {
			providerAccounts = append(providerAccounts, account)
		}
	}

	return providerAccounts, nil
}
