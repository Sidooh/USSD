package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"fmt"
	"strconv"
)

type MerchantVoucher struct {
	Merchant
}

func (w *MerchantVoucher) Process(input string) {
	logger.UssdLog.Println(" -- MERCHANT_VOUCHER: process", w.screen.Key, input)
	w.productRep = "merchant_voucher"

	w.Product.Process(input)
	w.processScreen(input)
	w.finalize()
}

func (w *MerchantVoucher) processScreen(input string) {
	switch w.screen.Key {

	case utils.MERCHANT:
		//f.vars["{pay_buy}"] = "Buy"
		w.vars["{product}"] = "Voucher for "
		w.vars["{number}"] = w.vars["{phone}"]

		w.setPaymentMethodText("1")

	case utils.MERCHANT_VOUCHER_AMOUNT:
		w.vars["{amount}"] = input
		w.getTopUpCharge(input)
	}
}

func (w *MerchantVoucher) finalize() {
	logger.UssdLog.Println(" -- MERCHANT_VOUCHER: finalize", w.screen.Next.Type)

	if w.screen.Key == utils.MERCHANT_VOUCHER_CONFIRMATION {
		amount, _ := strconv.Atoi(w.vars["{amount}"])

		request := client.MerchantMpesaWithdrawalRequest{
			Amount: amount,
			Phone:  w.vars["{number}"],
		}
		service.VoucherPurchase(w.vars["{merchant_id}"], request)
	}
}

func (w *MerchantVoucher) getTopUpCharge(input string) {
	amount, _ := strconv.Atoi(input)

	charge := service.GetMpesaCollectionCharge(amount)

	//w.vars["{withdrawal_charge}"] = strconv.Itoa(charge)
	w.vars["{payment_charge_text}"] = fmt.Sprintf("\n\nCost: KES %v", charge)

}
