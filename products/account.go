package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
)

type Account struct {
	Product
}

func (a *Account) Process(input string) {
	logger.UssdLog.Println(" -- ACCOUNT: process", a.screen.Key, input)
	a.productRep = "account"

	a.Product.Process(input)
	a.processScreen(input)
	a.finalize()
}

func (a *Account) processScreen(input string) {
	switch a.screen.Key {
	case utils.MAIN_MENU:
		//a.vars["{product}"] = a.productRep
		//a.vars["{number}"] = a.vars["{phone}"]
		if _, ok := a.vars["{full_name}"]; !ok {
			a.vars["{full_name}"] = ""
		}

		a.fetchUserSubscription()
		break
	case utils.ACCOUNT_PROFILE:
		a.checkHasPin()

	case utils.PROFILE_NAME, utils.PROFILE_UPDATE_NAME:
		a.vars["{full_name}"] = input

	case utils.PROFILE_NEW_PIN:
		a.vars["{pin}"] = input

	case utils.PROFILE_NEW_PIN_CONFIRM:
		a.vars["{confirm_pin}"] = input

	}
}

func (a *Account) finalize() {
	logger.UssdLog.Println(" -- ACCOUNT: finalize", a.screen.Next.Type)

	if a.screen.Key == utils.PROFILE_NEW_PIN_CONFIRM {
		accountId, _ := a.vars["{account_id}"]
		pin := a.vars["{confirm_pin}"]

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		err := service.SetPin(accountId, pin)
		if !err {
			a.screen.Next.Title = "Sorry. We failed to set your pin, please try again later."
		} else {
			//	TODO: Notify user of new pin set and also ask to set id and security questions
		}
	}

	if a.screen.Key == utils.PROFILE_UPDATE_NAME {
		accountId, _ := a.vars["{account_id}"]
		name := a.vars["{full_name}"]

		request := client.ProfileDetails{
			Name: name,
		}

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		_, err := service.UpdateProfile(accountId, request)
		if err != nil {
			a.screen.Next.Title = "Sorry. We failed to update your details, please try again later."
		}
	}
}

func (a *Account) fetchUserSubscription() {
	logger.UssdLog.Println("   ++ ACCOUNT: fetch user subscription")

	accountId := a.vars["{account_id}"]

	subscription, _ := service.FetchSubscription(accountId)

	if subscription.Id != 0 && subscription.Status == utils.ACTIVE {
		a.vars["{subscription_type}"] = "Sidooh Agent"
	} else {
		a.vars["{subscription_type}"] = "None"
	}
}

func (a *Account) checkHasPin() {
	accountId := a.vars["{account_id}"]

	hasPin := service.CheckHasPin(accountId)

	if hasPin {
		delete(a.screen.Next.Options, 1)
	} else {
		a.screen.Next.Options[2].NextKey = utils.PIN_NOT_SET

		if option, ok := a.screen.Next.Options[3]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}
	}
}
