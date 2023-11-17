package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"fmt"
	"strconv"
)

type MerchantMpesaWithdraw struct {
	Merchant
}

func (w *MerchantMpesaWithdraw) Process(input string) {
	logger.UssdLog.Println(" -- MPESA_WITHDRAW: process", w.screen.Key, input)
	w.productRep = "mpesa withdraw"

	w.Product.Process(input)
	w.processScreen(input)
	w.finalize()
}

func (w *MerchantMpesaWithdraw) processScreen(input string) {
	switch w.screen.Key {

	case utils.MERCHANT:
		//f.vars["{pay_buy}"] = "Buy"
		w.vars["{product}"] = "Mpesa withdrawal for"

		w.vars["{payment_method}"] = utils.MPESA
		w.vars["{payment_method_instruction}"] = "Ask customer to enter MPESA PIN in the screen that asks \"Do you want to pay to Sidooh International?\" to complete the transaction."

	case utils.MERCHANT_MPESA_WITHDRAW:
		w.vars["{number}"] = input
		w.vars["{payment_method_text}"] = utils.MPESA + " " + w.vars["{number}"]

	case utils.MERCHANT_MPESA_WITHDRAW_AMOUNT:
		w.vars["{amount}"] = input
		//w.vars["{payment_charge_text}"] = "\n\nCost: KES 30"
		w.getWithdrawalCharge(input)
	}
}

func (w *MerchantMpesaWithdraw) finalize() {
	logger.UssdLog.Println(" -- MPESA_WITHDRAW: finalize", w.screen.Next.Type)

	if w.screen.Key == utils.MERCHANT_MPESA_WITHDRAW_CONFIRMATION {
		amount, _ := strconv.Atoi(w.vars["{amount}"])

		request := client.MerchantMpesaWithdrawalRequest{
			Amount: amount,
			Phone:  w.vars["{number}"],
		}
		service.MpesaWithdrawal(w.vars["{merchant_id}"], request)
	}
}

func (w *MerchantMpesaWithdraw) getWithdrawalCharge(input string) {
	amount, _ := strconv.Atoi(input)

	charge := service.GetMpesaWithdrawalCharge(amount)

	//w.vars["{withdrawal_charge}"] = strconv.Itoa(charge)
	w.vars["{payment_charge_text}"] = fmt.Sprintf("\n\nCost: KES %v", charge)

}
