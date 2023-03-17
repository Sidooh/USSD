package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/utils"
)

type BuyGoods struct {
	Merchant
}

func (m *BuyGoods) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MPESA_BUY_GOODS: process", m.screen.Key, input)
	m.productRep = "pay_mpesa_buy_goods"

	m.Merchant.Process(input)
	m.processScreen(input)
	m.finalize()
}

func (m *BuyGoods) processScreen(input string) {
	switch m.screen.Key {
	case utils.MERCHANT_BUY_GOODS:
		m.vars["{merchant_type}"] = utils.MPESA_BUY_GOODS
		m.vars["{merchant_number}"] = input
		m.vars["{product}"] = "to Till Number " + m.vars["{merchant_number}"]
		m.vars["{number}"] = ""

	}
}
