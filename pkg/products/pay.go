package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/utils"
)

type Pay struct {
	Product
}

func (p *Pay) Process(input string) {
	logger.UssdLog.Println(" -- PAY: process", p.screen.Key, input)
	p.productRep = "pay"

	p.Product.Process(input)
	p.processScreen(input)
	p.finalize()
}

func (p *Pay) processScreen(input string) {
	switch p.screen.Key {
	case utils.PAY:
		p.vars["{payment_charge_text}"] = ""
	}
}
