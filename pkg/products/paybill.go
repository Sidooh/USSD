package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/utils"
)

type PayBill struct {
	Merchant
}

func (m *PayBill) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MPESA_PAY_BILL: process", m.screen.Key, input)
	m.productRep = "pay_mpesa_pay_bill"

	m.Merchant.Process(input)
	m.processScreen(input)
	m.finalize()
}

func (m *PayBill) processScreen(input string) {
	switch m.screen.Key {
	case utils.MERCHANT_PAY_BILL:
		m.vars["{merchant_type}"] = utils.MPESA_PAY_BILL
		m.vars["{merchant_number}"] = input
		m.vars["{product}"] = "to Paybill " + m.vars["{merchant_number}"]

	case utils.MERCHANT_PAY_BILL_ACCOUNT:
		m.vars["{merchant_account}"] = input
		m.vars["{number}"] = "for Account " + m.vars["{merchant_account}"]

	}
}
