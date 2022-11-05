package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"strconv"
)

type Voucher struct {
	Pay
}

func (v *Voucher) Process(input string) {
	logger.UssdLog.Println(" -- VOUCHER: process", v.screen.Key, input)
	v.productRep = "voucher"

	v.Product.Process(input)
	v.processScreen(input)
	v.finalize()
}

func (v *Voucher) processScreen(input string) {
	switch v.screen.Key {
	case utils.PAY_VOUCHER, utils.VOUCHER_BALANCE_INSUFFICIENT:
		v.vars["{product}"] = v.productRep + " for"
		v.vars["{number}"] = v.vars["{phone}"]
	case utils.VOUCHER_OTHER_ACCOUNT:
		v.vars["{number}"], _ = utils.FormatPhone(input)
		break
	case utils.VOUCHER_AMOUNT:
		v.vars["{amount}"] = input
		v.setPaymentMethods(input)
		break
	}
}

func (v *Voucher) finalize() {
	logger.UssdLog.Println(" -- VOUCHER: finalize", v.screen.Next.Type)

	if v.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(v.vars["{account_id}"])
		amount, _ := strconv.Atoi(v.vars["{amount}"])
		method := v.vars["{payment_method}"]

		if accountId == 0 {
			logger.UssdLog.Println(" -- AIRTIME: creating acc")

			account, err := service.CreateAccount(v.vars["{phone}"])
			if err != nil {
				// TODO: Send message to user
				logger.UssdLog.Error("Failed to create account: ", err)
			}

			accountId = account.Id
		}

		request := client.VoucherPurchaseRequest{
			PurchaseRequest: client.PurchaseRequest{
				Initiator: utils.CONSUMER,
				Amount:    amount,
				Method:    method,
				AccountId: accountId,
			},
		}

		if v.vars["{number}"] != v.vars["{phone}"] {
			request.TargetNumber = v.vars["{number}"]
		}

		if _, ok := v.vars["{mpesa_number}"]; ok {
			request.DebitAccount = v.vars["{mpesa_number}"]
		}

		logger.UssdLog.Println(" -- VOUCHER: purchase", request)

		// TODO: Make into goroutine if applicable
		service.PurchaseVoucher(&request)
	}
}
