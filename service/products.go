package service

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
	"strings"
)

var productsClient = client.InitProductClient()
var notifyClient = client.InitNotifyClient()

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

// TODO: Implement global caching procedure (Local memory, TTL3, network ...order?)
//var subscriptions = map[string]client.Subscription{}

func FetchSubscription(id string) (client.Subscription, error) {
	var subscription client.Subscription

	// TODO: Return once we have implemented global caching or can autoremove after duration
	//if cachedSub, ok := subscriptions[id]; ok {
	//	subscription = cachedSub
	//	return subscription, nil
	//}

	err := productsClient.GetSubscription(id, &subscription)
	if err != nil {
		logger.ServiceLog.Error("Failed to fetch subscription: ", err)
		return client.Subscription{}, err
	}
	// TODO: Should we only cache when ACTIVE?
	//if subscription.Id != 0 {
	//	subscriptions[id] = subscription
	//}

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
	// TODO: Test if this actually works
	//delete(subscriptions, strconv.Itoa(request.AccountId))
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
