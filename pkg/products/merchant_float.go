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
	"strings"
)

type MerchantFloat struct {
	Merchant
}

func (f *MerchantFloat) Process(input string) {
	logger.UssdLog.Println(" -- FLOAT: process", f.screen.Key, input)
	f.productRep = "float"

	f.Product.Process(input)
	f.processScreen(input)
	f.finalize()
}

func (f *MerchantFloat) processScreen(input string) {
	switch f.screen.Key {

	case utils.MERCHANT:
		f.setStoreOptions(input)
		f.vars["{pay_buy}"] = "Buy"
		f.vars["{product}"] = "float purchase for"

	case utils.MERCHANT_FLOAT:
		f.processStoreSelection(input)

	case utils.MERCHANT_FLOAT_AGENT:
		f.vars["{agent}"] = input

	case utils.MERCHANT_FLOAT_STORE:
		f.vars["{store}"] = input
		f.vars["{number}"] = f.vars["{agent}"] + " - " + f.vars["{store}"]

	case utils.MERCHANT_FLOAT_AMOUNT:

		f.vars["{amount}"] = input
		f.setPaymentMethods(input)
		f.vars["{payment_charge_text}"] = "\n\nCost: KES 30"

	case utils.PAYMENT_METHOD:
		f.getFloatCharge(f.vars["{amount}"])

	}
}
func (f *MerchantFloat) finalize() {
	logger.UssdLog.Println(" -- FLOAT: finalize", f.screen.Next.Type)

	if f.screen.Key == utils.PAYMENT_CONFIRMATION {
		account, _ := strconv.Atoi(f.vars["{account_id}"])
		amount, _ := strconv.Atoi(f.vars["{amount}"])
		method := f.vars["{payment_method}"]

		request := client.MerchantMpesaFloatPurchaseRequest{
			PurchaseRequest: client.PurchaseRequest{
				AccountId: account,
				Amount:    amount,
				Method:    method,
			},
			Agent: f.vars["{agent}"],
			Store: f.vars["{store}"],
		}

		if method == utils.MPESA {
			request.DebitAccount = f.vars["{phone}"]
		}

		if acc, ok := f.vars["{mpesa_number}"]; ok {
			request.DebitAccount = acc
		}

		go service.BuyFloat(f.vars["{merchant_id}"], request)
	}
}

func (f *MerchantFloat) setStoreOptions(input string) {
	logger.UssdLog.Println("   ++ FLOAT: set store options", input)

	merchantId := f.vars["{merchant_id}"]
	accounts, _ := service.FetchMpesaStoreAccounts(merchantId)

	if accounts != nil && len(accounts) != 0 {
		storeAccountOptionVars := map[int]string{}

		for i, account := range accounts {
			f.screen.Next.Options[i+2] = &data.Option{
				Label:   strings.Join(strings.Split(account.Name, " ")[0:3], " "),
				Value:   i + 2,
				NextKey: utils.MERCHANT_FLOAT_AMOUNT,
			}

			storeAccountOptionVars[i+2] = account.Agent + " __ " + account.Name
		}
		stringVars, _ := json.Marshal(storeAccountOptionVars)
		f.vars["{mpesa_store_account_options}"] = string(stringVars)
	} else {
		f.screen.Options[1].NextKey = utils.MERCHANT_FLOAT_AGENT
	}
}

func (f *MerchantFloat) processStoreSelection(input string) {
	logger.UssdLog.Println("   ++ FLOAT: process store selection", input)

	selectedStoreAccount, _ := strconv.Atoi(input)
	storeAccountOptionVars := map[int]string{}

	_ = json.Unmarshal([]byte(f.vars["{mpesa_store_account_options}"]), &storeAccountOptionVars)

	agentName := strings.Split(storeAccountOptionVars[selectedStoreAccount], " __ ")

	if len(agentName) > 1 {
		f.vars["{agent}"] = agentName[0]
		f.vars["{store}"] = strings.Split(agentName[1], " - ")[0]
		f.vars["{number}"] = agentName[1]
	}
}

func (f *MerchantFloat) getFloatCharge(input string) {
	amount, _ := strconv.Atoi(input)

	charge := service.GetMpesaFloatCharge(amount)

	if f.vars["{payment_method}"] == utils.MPESA {
		charge += service.GetMpesaCollectionCharge(amount)
	}

	f.vars["{payment_charge_text}"] = fmt.Sprintf("\n\nCost: KES %v", charge)

}
