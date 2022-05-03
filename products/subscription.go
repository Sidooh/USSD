package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"strconv"
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

func (s *Subscription) finalize() {
	logger.UssdLog.Println(" -- SUBSCRIPTION: finalize", s.screen.Next.Type)

	if s.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(s.vars["{account_id}"])
		amount, _ := strconv.Atoi(s.vars["{amount}"])
		method := s.vars["{payment_method}"]

		request := client.AirtimePurchaseRequest{
			Initiator: client.CONSUMER,
			Amount:    amount,
			Method:    method,
			AccountId: accountId,
		}

		logger.UssdLog.Println(" -- SUBSCRIPTION: purchase", request)

		// TODO: Make into goroutine if applicable
		//service.PurchaseSubscription(&request)
	}
}
