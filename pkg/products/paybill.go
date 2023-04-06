package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/utils"
	"fmt"
)

type PayBill struct {
	Merchant
}

func (m *PayBill) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MPESA_PAY_BILL: process", m.screen.Key, input)
	m.productRep = "pay_mpesa_pay_bill"

	m.Merchant.Process(input)
	m.processScreen(input)
}

func (m *PayBill) processScreen(input string) {
	switch m.screen.Key {
	case utils.MERCHANT_PAY_BILL:
		m.vars["{merchant_type}"] = utils.MPESA_PAY_BILL
		m.vars["{merchant_number}"] = input
		m.vars["{product}"] = "to Paybill " + m.vars["{merchant_number}"]

		m.setMerchant(input)

	case utils.MERCHANT_PAY_BILL_ACCOUNT:
		m.vars["{merchant_account}"] = input

	case utils.MERCHANT_AMOUNT:
		m.vars["{number}"] = fmt.Sprintf("for Account %s (plus KES%s Savings)", m.vars["{merchant_account}"], m.vars["{merchant_fee}"])
	}
}

func (m *PayBill) setMerchant(input string) {
	m.searchMerchant(input)

	if name, ok := m.vars["{merchant_name}"]; ok {
		m.vars["{product}"] = fmt.Sprintf("to %s", name)
	}
}
