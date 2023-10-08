package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
		switch input {
		case "2":
			a.setInvites()
		case "3", "4":
			a.setbalances()
		}

	case utils.MERCHANT_PROFILE_CHANGE_PIN_QUESTION:
		a.processUserAnswer(input)

	case utils.MERCHANT_PROFILE_NEW_PIN:
		a.vars["{pin}"] = input
	case utils.MERCHANT_PROFILE_NEW_PIN_CONFIRM:
		a.vars["{confirm_pin}"] = input

		//case utils.MERCHANT_INVITES:
		//	a.setInvites()
		//
		//case utils.MERCHANT_BALANCES:
		//	a.setbalances()

	case utils.MERCHANT_WITHDRAW:
		switch input {
		case "1":
			a.vars["{source}"] = "Cashback"
			a.vars["{withdrawable_earnings}"] = a.vars["{cashback_balance}"]
		case "2":
			a.vars["{source}"] = "Commission"
			a.vars["{withdrawable_earnings}"] = a.vars["{commission_balance}"]

		}

		// TODO: get actual charges
		a.vars["{withdrawal_charge}"] = "15"
	case utils.MERCHANT_WITHDRAW_AMOUNT:
		a.vars["{amount}"] = input

	case utils.MERCHANT_WITHDRAW_DESTINATION:
		switch input {
		case "3":
			a.vars["{destination}"] = utils.FLOAT
			a.vars["{destination_text}"] = utils.VOUCHER
			a.vars["{account}"] = a.vars["{merchant_float}"]

		default:
			a.vars["{destination}"] = utils.MPESA
			a.vars["{destination_text}"] = utils.MPESA

		}

	case utils.MERCHANT_WITHDRAW_MPESA:
		a.vars["{account}"] = a.vars["{phone}"]
		a.vars["{destination_text}"] = utils.MPESA + ": " + a.vars["{phone}"]

	case utils.MERCHANT_WITHDRAW_MPESA_OTHER:
		a.vars["{account}"], _ = utils.FormatPhone(input)
		a.vars["{destination_text}"] += utils.MPESA + ": " + a.vars["{account}"]

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

	// User has just requested a withdrawal
	if a.screen.Key == utils.MERCHANT_WITHDRAW_CONFIRM_PIN {
		amount, _ := strconv.Atoi(a.vars["{amount}"])

		request := client.MerchantWithdrawalRequest{
			Destination: a.vars["{destination}"],
			Account:     a.vars["{account}"],
			Amount:      amount,
		}

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		err := service.WithdrawEarnings(a.vars["{merchant_id}"], request)
		if err != nil {
			a.screen.Next.Title = "Sorry. We failed to process your withdrawal request, please try again later."
		}
	}
}

func (a *MerchantAccount) setInvites() {
	a.vars["{subagents}"] = "0"
	a.vars["{subagents_active}"] = "0"
	a.vars["{subagent_referrals}"] = "0"
	a.vars["{subagent_referrals_active}"] = "0"

	descendants, err := service.FetchDescendants(a.vars["{account_id}"], "2")
	if err != nil {
		return
	}

	var subAgents, subAgentReferrals []string

	// Group/sum by level
	for _, descendant := range descendants {
		if descendant.Level == 1 {
			subAgents = append(subAgents, strconv.Itoa(descendant.Id))
		}

		if descendant.Level == 2 {
			subAgentReferrals = append(subAgentReferrals, strconv.Itoa(descendant.Id))
		}
	}

	// TODO: should we only count those who are merchants?

	a.vars["{subagents}"] = strconv.Itoa(len(subAgents))
	a.vars["{subagent_referrals}"] = strconv.Itoa(len(subAgentReferrals))

	// Get tx done by each acc
	if (len(subAgents) + len(subAgentReferrals)) == 0 {
		return
	}

	acc := strings.Join(subAgents, ",")
	acc += strings.Join(subAgentReferrals, ",")
	transactions, err := service.FetchTransactions(acc, "7")
	fmt.Println(transactions, err)
	if err != nil {
		return
	}

	// TODO: finalize this

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
