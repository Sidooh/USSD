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
		w.vars["{voucher_option}"] = "1"
		w.vars["{amount}"] = input
		w.getTopUpCharge(input)

	case utils.MERCHANT_VOUCHER_TRANSFER_OTHER_ACCOUNT:
		w.vars["{number}"] = input

	case utils.MERCHANT_VOUCHER_TRANSFER_AMOUNT:
		w.vars["{voucher_option}"] = "2"
		w.vars["{amount}"] = input
		//delete(w.screen.Next.Options, 1)
		//delete(w.screen.Next.Options, 2)
		//amount, _ := strconv.Atoi(input)
		//w.addMerchantFloatPaymentMethod(amount)
		w.vars["{payment_charge_text}"] = ""

	case utils.MERCHANT_VOUCHER_TRANSFER_PAYMENT_METHOD:
		w.setPaymentMethodText(input)

	}
}

func (w *MerchantVoucher) finalize() {
	logger.UssdLog.Println(" -- MERCHANT_VOUCHER: finalize", w.screen.Next.Type)

	if w.screen.Next.Type == utils.END {
		if w.vars["{voucher_option}"] == "1" {
			amount, _ := strconv.Atoi(w.vars["{amount}"])

			request := client.MerchantMpesaWithdrawalRequest{
				Amount: amount,
				Phone:  w.vars["{number}"],
			}

			if phone, ok := w.vars["{mpesa_number}"]; ok {
				request.Phone = phone
			}

			go service.VoucherPurchase(w.vars["{merchant_id}"], request)
		}

		if w.vars["{voucher_option}"] == "2" {
			amount, _ := strconv.Atoi(w.vars["{amount}"])

			request := client.MerchantFloatTransferRequest{
				Amount:  amount,
				Account: w.vars["{merchant_account_validated}"],
			}
			//
			//if phone, ok := w.vars["{mpesa_number}"]; ok {
			//	request.Phone = phone
			//}

			go service.VoucherTransfer(w.vars["{merchant_id}"], request)
		}
	}
}

func (w *MerchantVoucher) getTopUpCharge(input string) {
	amount, _ := strconv.Atoi(input)

	charge := service.GetPayBillCharge(amount)

	if amount < 11000 {
		charge = service.GetBuyGoodsCharge(amount)
	}

	w.vars["{payment_charge_text}"] = fmt.Sprintf("\n\nCost: KES %v", charge)

}
