package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"strconv"
	"strings"
	"time"
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
	case utils.MAIN_MENU:
		s.vars["{product}"] = s.productRep
		s.vars["{number}"] = s.vars["{phone}"]

		// TODO: Move to different screen after selection
		// TODO: Make dynamic with fetch from api
		s.vars["{subscription_type}"] = "Sidooh Agent"
		s.vars["{subscription_amount}"] = "KES365"
		s.vars["{duration}"] = "month"
		s.vars["{amount}"] = "365"

		s.fetchUserSubscription()

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

		// TODO: Make subscription_type dynamic with api fetch
		request := client.SubscriptionPurchaseRequest{
			PurchaseRequest: client.PurchaseRequest{
				Initiator: utils.CONSUMER,
				Amount:    amount,
				Method:    method,
				AccountId: accountId,
			},
			SubscriptionTypeId: 1,
		}

		logger.UssdLog.Println(" -- SUBSCRIPTION: purchase", request)

		// TODO: Make into goroutine if applicable
		service.PurchaseSubscription(&request)
	}
}

func (s *Subscription) fetchUserSubscription() {
	logger.UssdLog.Println("   ++ SUBSCRIPTION: fetch user subscription")

	if accountId, ok := s.vars["{account_id}"]; ok {
		subscription, _ := service.FetchSubscription(accountId)

		if subscription.Id != 0 {

			endDate := strings.Split(subscription.EndDate, " ")[0]
			s.vars["{subscription_end_date}"] = endDate

			if subscription.Status == utils.ACTIVE {
				s.screen.Options[6].NextKey = utils.SUBSCRIPTION_ACTIVE
			}

			expiryTime, err := time.Parse(`2006-01-02 15:04:05`, subscription.EndDate)

			if subscription.Status == utils.EXPIRED || (time.Until(expiryTime) < 3*24*time.Hour && err == nil) {
				s.screen.Options[6].NextKey = utils.SUBSCRIPTION_RENEW
			}

		}
	}
}
