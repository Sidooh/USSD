package products

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"strconv"
)

type Utility struct {
	Pay
}

func (u *Utility) Process(input string) {
	logger.UssdLog.Println(" -- PAY_UTILITY: process", u.screen.Key, input)
	u.productRep = "pay_utility"

	u.Product.Process(input)
	u.processScreen(input)
	u.finalize()
}

func (u *Utility) processScreen(input string) {
	switch u.screen.Key {
	case utils.UTILITY:
		u.setUtilityAccountOptions(input)
		break
	case utils.UTILITY_ACCOUNT_SELECT:
		u.processUtilityAccountSelection(input)
		break
	case utils.UTILITY_OTHER_ACCOUNT:
		u.vars["{number}"] = input
		break
	case utils.UTILITY_AMOUNT:
		u.vars["{amount}"] = input
		break
	}
}

func (u *Utility) setUtilityAccountOptions(input string) {
	logger.UssdLog.Println("   ++ PAY_UTILITY: set utility account options", input)

	integerInput, _ := strconv.Atoi(input)
	u.vars["{product}"] = "to " + u.screen.Options[integerInput].Label

	provider := ""
	switch integerInput {
	case 1:
		provider = utils.KPLC_PREPAID
		break
	case 2:
		provider = utils.KPLC_POSTPAID
		break
	case 3:
		provider = utils.NAIROBI_WTR
		break
	case 4:
		provider = utils.DSTV
		break
	case 5:
		provider = utils.ZUKU
		break
	case 6:
		provider = utils.GOTV
		break
	case 7:
		provider = utils.STARTIMES
		break
	}

	u.vars["{selected_utility}"] = provider

	accountId := u.vars["{account_id}"]
	accounts, _ := service.FetchUtilityAccounts(accountId, provider)

	if accounts != nil {
		utilityAccountOptionVars := map[int]string{}

		maxAccounts := accounts
		if len(accounts) > 5 {
			maxAccounts = accounts[:5]
		}

		for i, account := range maxAccounts {
			u.screen.Next.Options[i+1] = &data.Option{
				Label:   account.AccountNumber,
				Value:   i + 1,
				NextKey: utils.UTILITY_AMOUNT,
			}

			utilityAccountOptionVars[i+1] = account.AccountNumber
		}
		stringVars, _ := json.Marshal(utilityAccountOptionVars)
		u.vars["{utility_account_options}"] = string(stringVars)
	} else {
		u.screen.Options[integerInput].NextKey = utils.UTILITY_OTHER_ACCOUNT
	}

}

func (u *Utility) processUtilityAccountSelection(input string) {
	logger.UssdLog.Println("   ++ PAY_UTILITY: process utility account selection", input)

	selectedUtilityAccount, _ := strconv.Atoi(input)
	utilityAccountOptionVars := map[int]string{}

	_ = json.Unmarshal([]byte(u.vars["{utility_account_options}"]), &utilityAccountOptionVars)

	u.vars["{number}"] = utilityAccountOptionVars[selectedUtilityAccount]
}

func (u *Utility) finalize() {
	logger.UssdLog.Println(" -- PAY_UTILITY: finalize", u.screen.Next.Type)

	if u.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(u.vars["{account_id}"])
		amount, _ := strconv.Atoi(u.vars["{amount}"])
		method := u.vars["{payment_method}"]
		provider := u.vars["{selected_utility}"]
		accountNumber := u.vars["{number}"]

		if accountId == 0 {
			logger.UssdLog.Println(" -- AIRTIME: creating acc")

			account, err := service.CreateAccount(u.vars["{phone}"])
			if err != nil {
				// TODO: Send message to user
				logger.UssdLog.Error(err)
			}

			accountId = account.Id
		}

		request := client.UtilityPurchaseRequest{
			PurchaseRequest: client.PurchaseRequest{
				Initiator: utils.CONSUMER,
				Amount:    amount,
				Method:    method,
				AccountId: accountId,
			},
			Provider:      provider,
			AccountNumber: accountNumber,
		}

		if _, ok := u.vars["{mpesa_number}"]; ok {
			request.DebitAccount = u.vars["{mpesa_number}"]
		}

		logger.UssdLog.Println(" -- PAY_UTILITY: purchase", request)

		service.PurchaseUtility(request)
	}
}
