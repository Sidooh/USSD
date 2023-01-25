package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
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
		s.vars["{subscription_type}"] = "Earn More"
		s.vars["{subscription_amount}"] = "KES395"
		s.vars["{duration}"] = "month"
		s.vars["{amount}"] = "395"

		s.FetchSubscriptionType()

		s.fetchUserSubscription()

	case utils.SUBSCRIPTION_REGISTER:
		if name, ok := s.vars["{full_name}"]; ok && len(name) > 0 {
			// TODO: Fix pin not set checking in setPaymentMethods
			s.screen.Options[1].NextKey = utils.PAYMENT_METHOD
		}

	case utils.SUBSCRIPTION_SUBSCRIBER_NAME:
		s.vars["{subscriber_name}"] = input

	case utils.SUBSCRIPTION_SUBSCRIBER_CONFIRM:
		s.setPaymentMethods(s.vars["{amount}"])

	case utils.SUBSCRIPTION_RENEW:
		s.setPaymentMethods(s.vars["{amount}"])
	}
}

func (s *Subscription) finalize() {
	logger.UssdLog.Println(" -- SUBSCRIPTION: finalize", s.screen.Next.Type)

	if s.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(s.vars["{account_id}"])
		method := s.vars["{payment_method}"]

		// Update name if necessary
		if name, ok := s.vars["{subscriber_name}"]; ok {
			profileRequest := client.ProfileDetails{
				Name: name,
			}

			// TODO: Make into goroutine if applicable
			// TODO: Should we check returned value? Or should we make it a void function?
			_, _ = service.UpdateProfile(s.vars["{account_id}"], profileRequest)
		}

		// TODO: Make subscription_type dynamic with api fetch
		request := client.SubscriptionPurchaseRequest{
			PurchaseRequest: client.PurchaseRequest{
				Initiator: utils.CONSUMER,
				Method:    method,
				AccountId: accountId,
			},
		}

		if mpesa, ok := s.vars["{mpesa_number}"]; ok {
			request.DebitAccount = mpesa
		}

		if subType, ok := s.vars["{subscription_type_id}"]; ok {
			request.SubscriptionTypeId, _ = strconv.Atoi(subType)
		}

		logger.UssdLog.Println(" -- SUBSCRIPTION: purchase", request)

		// TODO: Make into goroutine if applicable
		service.PurchaseSubscription(&request)
	}
}

func (s *Subscription) FetchSubscriptionType() {
	logger.UssdLog.Println("   ++ SUBSCRIPTION: fetch default subscription type")

	subscriptionType, err := service.FetchSubscriptionType()
	if err != nil {
		return
	}

	s.vars["{subscription_type_id}"] = strconv.Itoa(subscriptionType.Id)
	//s.vars["{subscription_type}"] = subscriptionType.Title
	s.vars["{subscription_amount}"] = "KES" + strconv.Itoa(subscriptionType.Price)
	s.vars["{amount}"] = strconv.Itoa(subscriptionType.Price)

}

func (s *Subscription) fetchUserSubscription() {
	logger.UssdLog.Println("   ++ SUBSCRIPTION: fetch user subscription")

	if accountId, ok := s.vars["{account_id}"]; ok {
		subscription, _ := service.FetchSubscription(accountId)

		if subscription.Id != 0 {

			endDate := strings.Split(subscription.EndDate, " ")[0]
			s.vars["{subscription_end_date}"] = "valid until " + endDate

			if subscription.Status == utils.ACTIVE {
				s.screen.Options[6].NextKey = utils.SUBSCRIPTION_ACTIVE
			} else {
				s.vars["{subscription_end_date}"] = "expired"
			}

			expiryTime, err := time.Parse(`2006-01-02 15:04:05`, subscription.EndDate)

			isPast := expiryTime.Before(time.Now())
			isIn3Days := expiryTime.Before(time.Now().Add(3*24*time.Hour)) && !isPast
			isToday := time.Now().YearDay() == expiryTime.YearDay()

			if isPast {
				s.vars["{subscription_end_date}"] = "expired on " + endDate
			} else if isToday {
				s.vars["{subscription_end_date}"] = "expires today"
			} else if isIn3Days {
				s.vars["{subscription_end_date}"] = "expires on " + endDate
			}

			if subscription.Status == utils.EXPIRED || isPast || (isIn3Days && err == nil) {
				s.screen.Options[6].NextKey = utils.SUBSCRIPTION_RENEW
			}

		}
	}
}
