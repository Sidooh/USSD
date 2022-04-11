package service

import "USSD.sidooh/service/client"

var productsClient = client.InitProductClient()

type paymentMethod struct {
	method       string
	debitAccount string
}

func PurchaseAirtime(request client.AirtimePurchaseRequest) {
	err := productsClient.BuyAirtime(request)
	if err != nil {
		panic(err)
	}
}
