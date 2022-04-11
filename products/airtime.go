package products

import (
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"fmt"
	"strconv"
)

type Airtime struct {
	Product
	productRep string
}

func (a *Airtime) Process(input string) {
	fmt.Println("\t -- AIRTIME: process")
	a.productRep = "airtime"

	a.processScreen(input)
	a.finalize()
}

func (a *Airtime) processScreen(input string) {
	fmt.Println("\t -- AIRTIME: process screen", input)
	fmt.Println("\t --> selected: ", a.screen.Key)

	switch a.screen.Key {
	case utils.AIRTIME:
		a.vars["{product}"] = a.productRep
		a.vars["{number}"] = a.vars["{phone}"]
		break
	case utils.AIRTIME_OTHER_NUMBER_SELECT:
		break
	case utils.AIRTIME_OTHER_NUMBER:
		break
	case utils.AIRTIME_AMOUNT:
		a.vars["{amount}"] = input
		break
	case utils.PAYMENT_METHOD:
		switch input {
		case "1":
			a.vars["{payment_method}"] = utils.MPESA
			a.vars["{payment_method_text}"] = utils.MPESA + " " + a.vars["{phone}"]
			a.vars["{method_instruction}"] = "PLEASE ENTER MPESA PIN when prompted"
			break
		case "2":
			a.vars["{payment_method}"] = utils.VOUCHER
			a.vars["{payment_method_text}"] = utils.VOUCHER + "(" + a.vars["{voucher_balance}"] + ")"
			break
		}
		break
	}
}

func (a *Airtime) finalize() {
	fmt.Println("\t -- AIRTIME: finalize")

	//	Final checks
	fmt.Println("====", a.screen.Next.Type)

	if a.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(a.vars["{account_id}"])
		amount, _ := strconv.Atoi(a.vars["{amount}"])
		method := a.vars["{payment_method}"]

		request := client.AirtimePurchaseRequest{
			Initiator: client.CONSUMER,
			Amount:    amount,
			Method:    method,
			AccountId: accountId,
		}

		service.PurchaseAirtime(request)
	}
}
