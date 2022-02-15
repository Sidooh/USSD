package products

import (
	"USSD/data"
	"fmt"
)

type Airtime struct {
	Product
	productRep string
}

func (a *Airtime) Process(input string) {
	fmt.Println("\t -- AIRTIME: process")
	a.productRep = "airtime"

	a.processScreen(input)
}

func (a *Airtime) processScreen(input string) {
	fmt.Println("\t -- AIRTIME: process screen", input)
	fmt.Println("\t --> selected: ", a.screen.Key)

	switch a.screen.Key {
	case data.AIRTIME:
		a.vars["{product}"] = a.productRep
		a.vars["{number}"] = a.vars["{phone}"]
		break
	case data.AIRTIME_OTHER_NUMBER_SELECT:
		break
	case data.AIRTIME_OTHER_NUMBER:
		break
	case data.AIRTIME_AMOUNT:
		a.vars["{amount}"] = input
		break

	case data.PAYMENT_METHOD:
		switch input {
		case "1":
			a.vars["{payment_method}"] = data.MPESA
			a.vars["{payment_method_text}"] = data.MPESA + " " + a.vars["{phone}"]
			a.vars["{method_instruction}"] = "PLEASE ENTER MPESA PIN when prompted"
			break
		case "2":
			a.vars["{payment_method}"] = data.VOUCHER
			a.vars["{payment_method_text}"] = data.VOUCHER + "{voucher_balance}"
			break
		}
		break
	}
}
