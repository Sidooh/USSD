package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/utils"
	"fmt"
)

type BuyGoods struct {
	Pay
}

func (b *BuyGoods) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MPESA_BUY_GOODS: process", b.screen.Key, input)
	b.productRep = "pay_mpesa_buy_goods"

	b.Pay.Process(input)
	b.processScreen(input)
}

func (b *BuyGoods) processScreen(input string) {
	switch b.screen.Key {
	case utils.MERCHANT_BUY_GOODS:
		b.vars["{merchant_type}"] = utils.MPESA_BUY_GOODS
		b.vars["{merchant_number}"] = input
		b.vars["{product}"] = "to Till Number " + b.vars["{merchant_number}"]
		b.vars["{number}"] = ""

	case utils.MERCHANT_AMOUNT:
		b.vars["{number}"] = fmt.Sprintf("(plus KES%s Savings)", b.vars["{merchant_fee}"])
	}
}
