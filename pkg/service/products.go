package service

import (
	"USSD.sidooh/pkg/cache"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service/client"
	"strconv"
	"strings"
)

func PurchaseAirtime(request *client.PurchaseRequest) {
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

func PayMerchant(request client.MerchantPurchaseRequest) {
	err := productsClient.PayMerchant(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to pay merchant: ", err)
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

func FetchSubscription(id string) (client.Subscription, error) {
	var subscription client.Subscription

	err := productsClient.GetSubscription(id, &subscription)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch subscription: ", err)
		return client.Subscription{}, err
	}

	parts := strings.Split(subscription.EndDate, "T")
	if len(parts) == 2 {
		timeParts := strings.Split(parts[1], ".")
		subscription.EndDate = parts[0] + " " + timeParts[0]
	}

	return subscription, nil
}

func PurchaseSubscription(request *client.SubscriptionPurchaseRequest) {
	err := productsClient.PurchaseSubscription(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to purchase subscription: ", err)
	}

	cache.Remove("subscription_" + strconv.Itoa(request.AccountId))

}

func FetchSubscriptionType() (client.SubscriptionType, error) {
	var subscriptionType client.SubscriptionType

	err := productsClient.GetSubscriptionType(&subscriptionType)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch subscription type: ", err)
		return client.SubscriptionType{}, err
	}

	return subscriptionType, nil
}

func Notify(request *client.NotificationRequest) {
	err := notifyClient.SendNotification(request)
	if err != nil {
		logger.ServiceLog.Error("Failed to send notification: ", err)
	}
}

// Move to cache client storage, and enable sync process if changed in products service
var earningRates = map[string]client.EarningRate{}

func GetEarningRate(provider string) (*client.EarningRate, error) {
	var rate client.EarningRate

	if cachedRate, ok := earningRates[provider]; ok {
		rate = cachedRate
		return &rate, nil
	}

	var rates map[string]client.EarningRate
	err := productsClient.FetchEarningRates(&rates)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch earning rates: ", err)
		return nil, err
	}

	earningRates = rates
	rate = earningRates[provider]

	return &rate, nil
}

func GetPotentialEarnings(provider string, amount int, subscribed bool) float64 {
	// Get users earning ratio
	ratio := .6

	// Get ripples
	ripples := 6

	if subscribed {
		// Get subscribed users earning ratio
		ratio = 1.0

		// Get ripples (Subscribed users earn 100% pass-thru at the moment)
		ripples = 1
	}

	return calculateEarnings(provider, float64(amount), ratio, float64(ripples))
}

func calculateEarnings(provider string, amount float64, ratio float64, ripples float64) float64 {
	var discountAmount float64

	rate, err := GetEarningRate(provider)
	if err == nil && rate.Value != 0 {
		switch rate.Type {
		case "%":
			discountAmount = rate.Value * amount
		case "$":
			discountAmount = rate.Value
		}
	} else {
		discountAmount = float64(amount) * 0
	}

	return discountAmount * ratio / ripples
}
