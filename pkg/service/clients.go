package service

import (
	"USSD.sidooh/pkg/service/client"
)

var accountsClient *client.AccountsApiClient
var paymentsClient *client.PaymentsApiClient
var savingsClient *client.SavingsApiClient
var productsClient *client.ProductsApiClient
var notifyClient *client.NotifyApiClient

func Init() {
	accountsClient = client.InitAccountClient()
	productsClient = client.InitProductClient()
	paymentsClient = client.InitPaymentClient()
	savingsClient = client.InitSavingsClient()
	notifyClient = client.InitNotifyClient()
}
