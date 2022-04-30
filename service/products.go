package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
)

var productsClient = client.InitProductClient()

func PurchaseAirtime(request *client.AirtimePurchaseRequest) {
	err := productsClient.BuyAirtime(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to buy airtime: ", err)
	}
}

func PurchaseUtility(request client.UtilityPurchaseRequest) {
	err := productsClient.PayUtility(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to pay utility: ", err)
	}
}

func FetchAirtimeAccounts(id string) ([]client.UtilityAccount, error) {
	var accounts []client.UtilityAccount

	err := productsClient.GetAirtimeAccounts(id, &accounts)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch airtime accounts: ", err)
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	return accounts, nil
}

func FetchUtilityAccounts(id string, provider string) ([]client.UtilityAccount, error) {
	var accounts []client.UtilityAccount

	// TODO: Add stack traces for easier log tracing
	err := productsClient.GetUtilityAccounts(id, &accounts)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch utility accounts: ", err)
		return nil, err
	}

	// TODO: Can we cache the data so we don't have to re-fetch full accounts
	var providerAccounts []client.UtilityAccount
	for _, account := range accounts {

		if account.Provider == provider {
			providerAccounts = append(providerAccounts, account)
		}
	}

	if len(providerAccounts) == 0 {
		return nil, nil
	}

	return providerAccounts, nil
}

func PurchaseVoucher(request *client.VoucherPurchaseRequest) {
	err := productsClient.PurchaseVoucher(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to purchase voucher: ", err)
	}
}
