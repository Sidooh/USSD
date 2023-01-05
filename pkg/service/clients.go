package service

import (
	client2 "USSD.sidooh/pkg/service/client"
)

var accountsClient *client2.AccountsApiClient
var paymentsClient *client2.PaymentsApiClient
var savingsClient *client2.SavingsApiClient
var productsClient *client2.ProductsApiClient
var notifyClient *client2.NotifyApiClient

func Init() {
	accountsClient = client2.InitAccountClient()
	productsClient = client2.InitProductClient()
	paymentsClient = client2.InitPaymentClient()
	savingsClient = client2.InitSavingsClient()
	notifyClient = client2.InitNotifyClient()
}
