package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
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
			a.setEarnings()
		}

	case utils.MERCHANT_PROFILE_CHANGE_PIN_QUESTION:
		a.processUserAnswer(input)

	case utils.MERCHANT_PROFILE_NEW_PIN:
		a.vars["{pin}"] = input
	case utils.MERCHANT_PROFILE_NEW_PIN_CONFIRM:
		a.vars["{confirm_pin}"] = input

	case utils.MERCHANT_CASHBACK:
		a.vars["{source}"] = "CASHBACK"
		a.vars["{withdrawable_earnings}"] = a.vars["{cashback_balance}"]

	case utils.MERCHANT_COMMISSION:
		a.vars["{source}"] = "COMMISSION"
		a.vars["{withdrawable_earnings}"] = a.vars["{commission_balance}"]

	case utils.MERCHANT_WITHDRAW:
		switch input {
		case "1":
			a.vars["{source}"] = "CASHBACK"
			a.vars["{withdrawable_earnings}"] = a.vars["{cashback_balance}"]
		case "2":
			a.vars["{source}"] = "COMMISSION"
			a.vars["{withdrawable_earnings}"] = a.vars["{commission_balance}"]

		}

		// TODO: get actual charges
	case utils.MERCHANT_WITHDRAW_AMOUNT:
		a.vars["{amount}"] = input
		a.setWithdrawalCharge(input)

	case utils.MERCHANT_WITHDRAW_DESTINATION:
		switch input {
		case "3":
			a.vars["{destination}"] = utils.FLOAT
			a.vars["{destination_text}"] = utils.VOUCHER
			a.vars["{account}"] = a.vars["{merchant_float}"]
			a.vars["{withdrawal_charge}"] = "0"

		default:
			a.vars["{destination}"] = utils.MPESA
			a.vars["{destination_text}"] = utils.MPESA

		}

	case utils.MERCHANT_WITHDRAW_MPESA:
		a.vars["{account}"] = a.vars["{phone}"]
		a.vars["{destination_text}"] = utils.MPESA + ": " + a.vars["{phone}"]

	case utils.MERCHANT_WITHDRAW_MPESA_OTHER:
		a.vars["{account}"], _ = utils.FormatPhone(input)
		a.vars["{destination_text}"] = utils.MPESA + ": " + a.vars["{account}"]

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
			Source:      a.vars["{source}"],
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

func (a *MerchantAccount) setEarnings() {
	earnings, err := service.FetchEarningAccounts(a.vars["{merchant_id}"])
	if err != nil {
		a.vars["{cashback_balance}"] = "0"

		a.vars["{commission_balance}"] = "0"
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

	a.vars["{cashback_earnings}"] = formatAmount(cashback.Amount, "%.2f")

	a.vars["{commission_earnings}"] = formatAmount(commission.Amount, "%.2f")

	balances := a.fetchSavings()

	a.vars["{cashback_balance}"] = formatAmount(cashback.Amount+balances[1], "%.0f")
	a.vars["{commission_balance}"] = formatAmount(commission.Amount+balances[2], "%.0f")

}

func (a *MerchantAccount) fetchSavings() map[int]float64 {
	a.vars["{cashback_savings}"] = "0"
	a.vars["{cashback_interest}"] = "0"
	a.vars["{commission_savings}"] = "0"
	a.vars["{commission_interest}"] = "0"
	a.vars["{total_savings}"] = "0"

	a.vars["{withdrawable_savings}"] = "0"

	nextWithdrawal := "July 1st"
	if time.Now().Month() >= time.July {
		nextWithdrawal = "January 1st"
	}

	a.vars["{merchant_withdrawal_text}"] = "locked for 6 months. Withdraw next on " + nextWithdrawal

	accountId := a.vars["{account_id}"]

	savings, err := service.FetchSavingBalances(accountId)
	if err != nil {
		return nil
	}

	var cashbackAccount client.SavingAccount
	var commissionAccount client.SavingAccount
	for _, earning := range savings {
		if earning.Type == "MERCHANT_CASHBACK" {
			cashbackAccount = earning
		}
		if earning.Type == "MERCHANT_COMMISSION" {
			commissionAccount = earning
		}
	}

	cashS := cashbackAccount.Balance
	commS := commissionAccount.Balance

	total := cashS + commS

	wS := 0.0
	likelyCharge := float64(service.GetWithdrawalCharge(int(cashS)))
	if cashS > 30 {
		wS = math.Floor(cashS - likelyCharge)
	}

	a.vars["{cashback_savings}"] = formatAmount(cashS, "%.2f")
	a.vars["{cashback_interest}"] = formatAmount(cashbackAccount.Interest, "%.2f")
	a.vars["{commission_savings}"] = formatAmount(commS, "%.2f")
	a.vars["{commission_interest}"] = formatAmount(cashbackAccount.Interest, "%.2f")
	a.vars["{total_savings}"] = formatAmount(total, "%.2f")

	a.vars["{withdrawable_savings}"] = formatAmount(wS, "%.0f")

	return map[int]float64{1: cashS, 2: commS}
}
