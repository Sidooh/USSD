package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/utils"
)

type Subscription struct {
	Product
}

func (s *Subscription) Process(input string) {
	logger.UssdLog.Println(" -- SUBSCRIPTION: process", s.screen.Key, input)
	s.productRep = "subscription"

	s.Product.Process(input)
	s.processScreen(input)
	s.finalize()
}

func (s *Subscription) processScreen(input string) {
	switch s.screen.Key {
	case utils.SUBSCRIPTION:
		s.vars["{product}"] = s.productRep
		s.vars["{number}"] = s.vars["{phone}"]

		// TODO: Move to different screen after selection
		s.vars["{subscription_type}"] = "Sidooh Agent"
		s.vars["{subscription_amount}"] = "KES365"
		s.vars["{duration}"] = "month"
		s.vars["{amount}"] = "365"

	case utils.SUBSCRIPTION_AGENT_CONFIRM:
		s.setPaymentMethods(input)
	}
}
