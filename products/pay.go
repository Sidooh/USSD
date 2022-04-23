package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/utils"
)

type Pay struct {
	Product
	productRep string
}

func (p *Pay) Process(input string) {
	logger.UssdLog.Println(" -- PAY: process", p.screen.Key, input)
	p.productRep = "pay"

	p.processScreen(input)
	p.finalize()
}

func (p *Pay) processScreen(input string) {
	switch p.screen.Key {
	case utils.PAY:

	}
}
