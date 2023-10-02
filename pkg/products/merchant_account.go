package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
)

type MerchantAccount struct {
	Merchant
}

func (a *MerchantAccount) Process(input string) {
	logger.UssdLog.Println(" -- MERCH_ACC: process", a.screen.Key, input)
	a.productRep = "merch_acc"

	a.Product.Process(input)
	a.processScreen(input)
	a.finalize()
}

func (a *MerchantAccount) processScreen(input string) {
	switch a.screen.Key {
	case utils.MERCHANT_ACCOUNT:
		a.fetchUserSecurityQuestionOptions()

	case utils.MERCHANT_PROFILE_CHANGE_PIN_QUESTION:
		a.processUserAnswer(input)

	case utils.MERCHANT_PROFILE_NEW_PIN:
		a.vars["{pin}"] = input
	case utils.MERCHANT_PROFILE_NEW_PIN_CONFIRM:
		a.vars["{confirm_pin}"] = input

	case utils.MERCHANT_BALANCES:
		a.setbalances()
	}
}

func (a *MerchantAccount) finalize() {
	logger.UssdLog.Println(" -- MERCH_ACC: finalize", a.screen.Next.Type)

	// User has just input their security question answers which need verification
	if a.screen.NextKey == utils.MERCHANT_PROFILE_NEW_PIN && a.screen.Key == utils.MERCHANT_PROFILE_CHANGE_PIN_QUESTION {
		accountId, _ := a.vars["{account_id}"]

		questionAnswerVars := map[string]string{}

		_ = json.Unmarshal([]byte(a.vars["{question_answers}"]), &questionAnswerVars)

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		valid := service.CheckSecurityQuestionAnswers(accountId, questionAnswerVars)

		if !valid {
			a.screen.NextKey = utils.PROFILE_SECURITY_QUESTIONS_END
			a.vars["{profile_security_questions_end_title}"] = "Sorry. We failed to verify your security questions, please try again later."
		}
	}

	// User has just updated pin
	if a.screen.Key == utils.MERCHANT_PROFILE_NEW_PIN_CONFIRM {
		accountId, _ := a.vars["{account_id}"]
		pin := a.vars["{confirm_pin}"]

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		status := service.SetPin(accountId, pin)
		if !status {
			a.screen.Next.Title = "Sorry. We failed to set your pin, please try again later."
			return
		} else {
			//	TODO: Notify user of new pin set and also ask to set id and security questions
		}

	}
}

func (a *MerchantAccount) setbalances() {
	earnings, err := service.FetchEarningAccounts(a.vars["{merchant_id}"])
	if err != nil {
		a.vars["{cashback_balance}"] = "0"
		a.vars["{withdrawable_cashback}"] = "0"

		a.vars["{commission_balance}"] = "0"
		a.vars["{withdrawable_commission}"] = "0"
	}

	var cashback client.MerchantEarningAccount
	var commission client.MerchantEarningAccount
	for _, earning := range earnings {
		if earning.Type == "CASHBACK" {
			cashback = earning
		}
		if earning.Type == "COMMISSION" {
			commission = earning
		}
	}

	a.vars["{cashback_balance}"] = formatAmount(cashback.Amount, "%.0f")
	a.vars["{withdrawable_cashback}"] = formatAmount(cashback.Amount*.8, "%.0f")
	a.vars["{saved_cashback}"] = formatAmount(cashback.Amount*.2, "%.0f")

	a.vars["{commission_balance}"] = formatAmount(commission.Amount, "%.0f")
	a.vars["{withdrawable_commission}"] = formatAmount(commission.Amount*.8, "%.0f")
	a.vars["{saved_commission}"] = formatAmount(commission.Amount*.2, "%.0f")

}
