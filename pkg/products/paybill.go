package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/utils"
	"fmt"
)

type PayBill struct {
	Pay
}

func (p *PayBill) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MPESA_PAY_BILL: process", p.screen.Key, input)
	p.productRep = "pay_mpesa_pay_bill"

	p.Pay.Process(input)
	p.processScreen(input)
}

func (p *PayBill) processScreen(input string) {
	switch p.screen.Key {
	case utils.MERCHANT_PAY_BILL:
		p.vars["{merchant_type}"] = utils.MPESA_PAY_BILL
		p.vars["{merchant_number}"] = input
		p.vars["{product}"] = "to Paybill " + p.vars["{merchant_number}"]

		p.setMerchant(input)

	case utils.MERCHANT_PAY_BILL_ACCOUNT:
		p.vars["{merchant_account}"] = input

	case utils.MERCHANT_AMOUNT:
		p.vars["{number}"] = fmt.Sprintf("for Account %s (plus KES%s Savings)", p.vars["{merchant_account}"], p.vars["{merchant_fee}"])
	}
}

func (p *PayBill) setMerchant(input string) {
	p.searchMerchant(input)

	if name, ok := p.vars["{merchant_name}"]; ok {
		p.vars["{product}"] = fmt.Sprintf("to %s", name)
	}
}

func (p *PayBill) searchMerchant(input string) {
	merchant, _ := service.SearchMerchant(input)
	if merchant != nil {
		p.vars["{merchant_name}"] = merchant.Name
	} else {
		delete(p.vars, "{merchant_name}")
	}
}
