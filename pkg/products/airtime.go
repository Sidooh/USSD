package products

import (
	"USSD.sidooh/pkg/data"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

type Airtime struct {
	Product
}

func (a *Airtime) Process(input string) {
	logger.UssdLog.Println(" -- AIRTIME: process", a.screen.Key, input)
	a.productRep = "airtime"

	a.Product.Process(input)
	a.processScreen(input)
	a.finalize()
}

func (a *Airtime) processScreen(input string) {
	switch a.screen.Key {
	case utils.AIRTIME:
		a.vars["{product}"] = a.productRep + " for"
		a.vars["{number}"] = a.vars["{phone}"]
		a.vars["{payment_charge_text}"] = ""

		a.setOtherNumberOptions(input)
		break
	case utils.AIRTIME_OTHER_NUMBER_SELECT:
		a.processOtherNumberSelection(input)
		break
	case utils.AIRTIME_OTHER_NUMBER:
		a.vars["{number}"], _ = utils.FormatPhone(input)
		break
	case utils.AIRTIME_AMOUNT:
		a.vars["{amount}"] = input
		a.setPaymentMethods(input)

		amount, _ := strconv.Atoi(input)
		subscription, _ := a.vars["{subscription_status}"]
		provider, _ := utils.GetPhoneProvider(a.vars["{number}"])
		a.vars["{product}"] = fmt.Sprintf(
			"%s (which will earn you %.2f points) for",
			a.productRep,
			service.GetPotentialEarnings(provider, amount, subscription == "ACTIVE"),
		)
		break
	}
}

func (a *Airtime) finalize() {
	logger.UssdLog.Println(" -- AIRTIME: finalize", a.screen.Next.Type)

	if a.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(a.vars["{account_id}"])
		amount, _ := strconv.Atoi(a.vars["{amount}"])
		method := a.vars["{payment_method}"]

		if accountId == 0 {
			logger.UssdLog.Println(" -- AIRTIME: creating acc")

			account, err := service.CreateAccount(a.vars["{phone}"], a.vars["{invite_code_string}"])
			if err != nil {
				// TODO: Send message to user
				logger.UssdLog.Error("Failed to create account: ", err)
			}

			accountId = account.Id
		}

		request := client.PurchaseRequest{
			Initiator: utils.CONSUMER,
			Amount:    amount,
			Method:    method,
			AccountId: accountId,
		}

		if a.vars["{number}"] != a.vars["{phone}"] {
			request.TargetNumber = a.vars["{number}"]
		}

		if acc, ok := a.vars["{mpesa_number}"]; ok {
			request.DebitAccount = acc
		}

		logger.UssdLog.Println(" -- AIRTIME: purchase", request)

		// TODO: Make into goroutine if applicable
		go service.PurchaseAirtime(&request)
	}
}

func (a *Airtime) setOtherNumberOptions(input string) {
	logger.UssdLog.Println("   ++ AIRTIME: set other number options", input)

	if input == "2" {
		accountId := a.vars["{account_id}"]
		accounts, _ := service.FetchAirtimeAccounts(accountId)

		if accounts != nil {
			airtimeAccountOptionVars := map[int]string{}

			maxAccounts := accounts
			if len(accounts) > 7 {
				maxAccounts = accounts[:7]
			}

			for i, account := range maxAccounts {
				a.screen.Next.Options[i+1] = &data.Option{
					Label:   account.AccountNumber,
					Value:   i + 1,
					NextKey: utils.AIRTIME_AMOUNT,
				}

				airtimeAccountOptionVars[i+1] = account.AccountNumber
			}
			stringVars, _ := json.Marshal(airtimeAccountOptionVars)
			a.vars["{airtime_account_options}"] = string(stringVars)
		} else {
			a.screen.Options[2].NextKey = utils.AIRTIME_OTHER_NUMBER
		}
	}
}

func (a *Airtime) processOtherNumberSelection(input string) {
	logger.UssdLog.Println("   ++ AIRTIME: process other number selection", input)

	selectedAirtimeAccount, _ := strconv.Atoi(input)
	airtimeAccountOptionVars := map[int]string{}

	_ = json.Unmarshal([]byte(a.vars["{airtime_account_options}"]), &airtimeAccountOptionVars)

	a.vars["{number}"] = airtimeAccountOptionVars[selectedAirtimeAccount]
}
