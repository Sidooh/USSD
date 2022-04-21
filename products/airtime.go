package products

import (
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Airtime struct {
	Product
	productRep string
}

func (a *Airtime) Process(input string) {
	log.Println(" -- AIRTIME: process", a.screen.Key, input)
	a.productRep = "airtime"

	a.processScreen(input)
	a.finalize()
}

func (a *Airtime) processScreen(input string) {
	switch a.screen.Key {
	case utils.AIRTIME:
		a.vars["{product}"] = a.productRep
		a.vars["{number}"] = a.vars["{phone}"]
		break
	case utils.AIRTIME_OTHER_NUMBER_SELECT:
		// TODO: Get and Set other numbers
		break
	case utils.AIRTIME_OTHER_NUMBER:
		// TODO: Get other number input
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
	log.Println(" -- AIRTIME: finalize", a.screen.Next.Type)

	if a.screen.Next.Type == utils.END {
		account_id, _ := strconv.Atoi(a.vars["{account_id}"])
		amount, _ := strconv.Atoi(a.vars["{amount}"])
		method := a.vars["{payment_method}"]
		//mpesa_number := a.vars["{payment_method}"]

		if account_id == 0 {
			log.Println(" -- AIRTIME: creating acc")

			account, err := service.CreateAccount(a.vars["{phone}"])
			if err != nil {
				// TODO: Send message to user
				log.Error(err)
			}

			account_id = account.Id
		}

		request := client.AirtimePurchaseRequest{
			Initiator: client.CONSUMER,
			Amount:    amount,
			Method:    method,
			AccountId: account_id,
		}

		log.Println(" -- AIRTIME: purchase", request)

		service.PurchaseAirtime(request)
	}
}
