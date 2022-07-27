package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/utils"
)

type Save struct {
	Product
}

func (s *Save) Process(input string) {
	logger.UssdLog.Println(" -- SAVE: process", s.screen.Key, input)
	s.productRep = "save"

	s.Product.Process(input)
	s.processScreen(input)
	s.finalize()
}

func (s *Save) processScreen(input string) {
	switch s.screen.Key {
	case utils.MAIN_MENU:
		//s.vars["{product}"] = s.productRep
		//s.vars["{number}"] = s.vars["{phone}"]

	}
}
