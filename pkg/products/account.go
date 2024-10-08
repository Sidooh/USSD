package products

import (
	"USSD.sidooh/pkg/data"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

type Account struct {
	Product
}

func (a *Account) Process(input string) {
	// TODO: Add check for sensitive inputs and ignore. Maybe add sensitive/sanitize/hide key to screen?
	logger.UssdLog.Println(" -- ACCOUNT: process", a.screen.Key, input)
	a.productRep = "account"

	a.Product.Process(input)
	a.processScreen(input)
	a.finalize()
}

func (a *Account) processScreen(input string) {
	switch a.screen.Key {
	case utils.MAIN_MENU:
		// Set to '' for display purposes
		if _, ok := a.vars["{full_name}"]; !ok {
			a.vars["{full_name}"] = ""
		}

	case utils.ACCOUNT:
		a.setAccountOptions()

	case utils.ACCOUNT_PROFILE:
		a.setAccountProfileOptions()

	case utils.PROFILE_SECURITY:
		//if name, ok := a.vars["{full_name}"]; ok && len(name) > 0 {
		//	if _, exists := a.screen.Options[1]; exists {
		//		a.screen.Options[1].NextKey = utils.PROFILE_NEW_PIN
		//	}
		//}

	case utils.PROFILE_NAME, utils.PROFILE_UPDATE_NAME:
		a.vars["{full_name}"] = input

	case utils.PROFILE_NEW_PIN:
		a.vars["{pin}"] = input

	case utils.PROFILE_NEW_PIN_CONFIRM:
		a.vars["{confirm_pin}"] = input

	case utils.PROFILE_CHANGE_PIN_METHODS:
		a.fetchUserSecurityQuestionOptions()

	case utils.PROFILE_CHANGE_PIN_QUESTION:
		a.processUserAnswer(input)

	case utils.PROFILE_SECURITY_QUESTIONS_PIN:
		a.fetchSecurityQuestionOptions()

	case utils.PROFILE_SECURITY_QUESTIONS_OPTIONS:
		a.processQuestionSelection(input)

	case utils.PROFILE_SECURITY_QUESTIONS_ANSWER:
		a.processAnswer(input)

	case utils.ACCOUNT_BALANCES:
		a.getAccountBalances(input)

	case utils.ACCOUNT_WITHDRAW_PIN:
		a.fetchSavings()

	case utils.ACCOUNT_WITHDRAW:
		a.vars["{points}"] = input
		a.vars["{amount}"] = input
		a.setWithdrawalCharge(input)

	case utils.WITHDRAW_DESTINATION:
		switch input {
		case "1":
			a.vars["{account_type}"] = utils.MPESA
		case "2":
			a.vars["{account_type}"] = utils.VOUCHER
		case "3":
			a.vars["{account_type}"] = utils.BANK
		}

	case utils.WITHDRAW_MPESA:
		a.vars["{account_number}"] = a.vars["{phone}"]

	case utils.WITHDRAW_OTHER_NUMBER_MPESA:
		a.vars["{account_number}"], _ = utils.FormatPhone(input)

	}
}

func (a *Account) finalize() {
	logger.UssdLog.Println(" -- ACCOUNT: finalize", a.screen.Next.Type)

	// User has just created a new pin
	if a.screen.Key == utils.PROFILE_NEW_PIN_CONFIRM {
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

		name := a.vars["{full_name}"]

		request := client.ProfileDetails{
			Name: name,
		}

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		_, err := service.UpdateProfile(accountId, request)
		if err != nil {
			// TODO: Handle this on SMS
			name, _ := a.vars["{name}"]
			email, _ := a.vars["{customer_support_email}"]
			phone, _ := a.vars["{phone}"]

			message := fmt.Sprintf("Dear %s"+
				"we successfully updated your pin but were unable to update your name. "+
				"Please try again later or contact support at %s", name, email)
			request := client.NotificationRequest{
				Channel:     "SMS",
				Destination: []string{phone},
				EventType:   "ERROR_ALERT",
				Content:     message,
			}
			service.Notify(&request)
		}

	}

	// User has just updated their name/profile
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

	// User has just created security questions
	if a.screen.NextKey == utils.PROFILE_SECURITY_QUESTIONS_END {
		accountId, _ := a.vars["{account_id}"]

		questionAnswerVars := map[string]string{}

		_ = json.Unmarshal([]byte(a.vars["{question_answers}"]), &questionAnswerVars)

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		err := service.SetSecurityQuestions(accountId, questionAnswerVars)
		if err != nil {
			a.vars["{profile_security_questions_end_title}"] = "Sorry. We failed to set your security questions, please try again later."
		} else {
			a.vars["{profile_security_questions_end_title}"] = "Your security questions and answers have been recorded. Please remember them as you will need them if and when resetting your Sidooh PIN."
		}
	}

	// User has just input their security question answers which need verification
	if a.screen.NextKey == utils.PROFILE_NEW_PIN && a.screen.Key == utils.PROFILE_CHANGE_PIN_QUESTION {
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

	// User has just requested a withdrawal
	if a.screen.Key == utils.WITHDRAW_CONFIRM_PIN {
		accountId, _ := strconv.Atoi(a.vars["{account_id}"])
		amount, _ := strconv.Atoi(a.vars["{amount}"])

		request := &client.EarningsWithdrawalRequest{
			PurchaseRequest: client.PurchaseRequest{
				Initiator: utils.CONSUMER,
				AccountId: accountId,
				Amount:    amount,
			},
		}

		if a.vars["{account_number}"] != a.vars["{phone}"] {
			request.TargetNumber = a.vars["{account_number}"]
		}

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		err := service.RequestEarningsWithdrawal(request)
		if err != nil {
			a.screen.Next.Title = "Sorry. We failed to process your withdrawal request, please try again later."
		}
	}

}

func (a *Account) fetchUserSubscription() {
	logger.UssdLog.Println("   ++ ACCOUNT: fetch user subscription")

	if accountId, ok := a.vars["{account_id}"]; ok {

		subscription, _ := service.FetchSubscription(accountId)

		// TODO: Get subscription type as well and use the name
		if subscription.Id != 0 && subscription.Status == utils.ACTIVE {
			a.vars["{subscription_type}"] = "Earn More"
		} else {
			a.vars["{subscription_type}"] = "None"
		}
	}
}

func (a *Account) checkHasSecurityQuestions() bool {
	accountId := a.vars["{account_id}"]

	// Check if user already has_pin in state else fetch from service
	var hasSecurityQuestions bool
	err := json.Unmarshal([]byte(a.vars["{has_security_questions}"]), &hasSecurityQuestions)
	if err != nil {
		hasSecurityQuestions = service.CheckHasSecurityQuestions(accountId)
		stringVars, _ := json.Marshal(hasSecurityQuestions)
		a.vars["{has_security_questions}"] = string(stringVars)
	}

	return hasSecurityQuestions
}

func (a *Account) fetchSecurityQuestionOptions() {
	logger.UssdLog.Println("   ++ ACCOUNT: fetch security question options")

	questions, _ := service.FetchSecurityQuestions()

	questionAnswerVars := map[uint]string{}
	_ = json.Unmarshal([]byte(a.vars["{question_answers}"]), &questionAnswerVars)

	var unansweredQuestions []client.SecurityQuestion
	for _, question := range questions {
		if _, ok := questionAnswerVars[question.Id]; !ok {
			unansweredQuestions = append(unansweredQuestions, question)
		}
	}

	if unansweredQuestions != nil {
		questionOptionVars := map[int]client.SecurityQuestion{}

		maxQuestions := unansweredQuestions
		if len(questions) > 5 {
			maxQuestions = unansweredQuestions[:5]
		}

		a.screen.Next.Options = map[int]*data.Option{}

		for i, question := range maxQuestions {
			a.screen.Next.Options[i+1] = &data.Option{
				Label:   question.Question,
				Value:   i + 1,
				NextKey: utils.PROFILE_SECURITY_QUESTIONS_ANSWER,
			}

			questionOptionVars[i+1] = question
		}
		stringVars, _ := json.Marshal(questionOptionVars)
		a.vars["{question_options}"] = string(stringVars)
	} else {
		if _, exists := a.screen.Options[3]; exists {
			a.screen.Options[3].NextKey = utils.COMING_SOON
		}
	}
}

func (a *Account) processQuestionSelection(input string) {
	logger.UssdLog.Println("   ++ ACCOUNT: process question selection", input)

	selectedQuestion, _ := strconv.Atoi(input)
	questionOptionVars := map[int]client.SecurityQuestion{}

	_ = json.Unmarshal([]byte(a.vars["{question_options}"]), &questionOptionVars)

	a.vars["{security_question}"] = questionOptionVars[selectedQuestion].Question
	a.vars["{security_question_id}"] = strconv.Itoa(int(questionOptionVars[selectedQuestion].Id))

}

func (a *Account) processAnswer(input string) {
	logger.UssdLog.Println("   ++ ACCOUNT: process answer", input)

	questionAnswerVars := map[string]string{}

	_ = json.Unmarshal([]byte(a.vars["{question_answers}"]), &questionAnswerVars)

	questionAnswerVars[a.vars["{security_question_id}"]] = input

	stringVars, _ := json.Marshal(questionAnswerVars)

	a.vars["{question_answers}"] = string(stringVars)

	if len(questionAnswerVars) == 3 {
		a.screen.NextKey = utils.PROFILE_SECURITY_QUESTIONS_END
	} else {
		a.fetchSecurityQuestionOptions()
	}
}

func (a *Account) fetchUserSecurityQuestionOptions() {
	logger.UssdLog.Println("   ++ ACCOUNT: fetch user security question options")

	// Check if user already has questions in state
	var userQuestions []client.UserSecurityQuestion
	_ = json.Unmarshal([]byte(a.vars["{user_questions}"]), &userQuestions)

	// Fetch from API otherwise
	if len(userQuestions) == 0 {
		accountId := a.vars["{account_id}"]
		userQuestions, _ = service.FetchUserSecurityQuestions(accountId)

		stringVars, _ := json.Marshal(userQuestions)
		a.vars["{user_questions}"] = string(stringVars)
	}

	// Get the answered questions
	questionAnswerVars := map[uint]string{}
	_ = json.Unmarshal([]byte(a.vars["{question_answers}"]), &questionAnswerVars)

	// Filter only unanswered questions so we pick from them
	var unansweredQuestions []client.UserSecurityQuestion
	for _, question := range userQuestions {
		if _, ok := questionAnswerVars[question.Question.Id]; !ok {
			unansweredQuestions = append(unansweredQuestions, question)
		}
	}

	// Ensure there are still unanswered questions, otherwise proceed to verify them
	if len(unansweredQuestions) != 0 {
		a.vars["{security_question}"] = unansweredQuestions[0].Question.Question
		a.vars["{security_question_id}"] = strconv.Itoa(int(unansweredQuestions[0].Question.Id))
	} else {
		a.screen.Options[3].NextKey = utils.COMING_SOON
	}

}

func (a *Account) processUserAnswer(input string) {
	logger.UssdLog.Println("   ++ ACCOUNT: process user answer", input)

	questionAnswerVars := map[string]string{}

	_ = json.Unmarshal([]byte(a.vars["{question_answers}"]), &questionAnswerVars)

	questionAnswerVars[a.vars["{security_question_id}"]] = input

	stringVars, _ := json.Marshal(questionAnswerVars)

	a.vars["{question_answers}"] = string(stringVars)

	if len(questionAnswerVars) == 3 {
		a.screen.NextKey = utils.PROFILE_NEW_PIN
	} else {
		a.fetchUserSecurityQuestionOptions()
	}
}

func (a *Account) setAccountOptions() {
	logger.UssdLog.Println("   ++ ACCOUNT: set account options")

	hasPin := a.checkHasPin()

	if !hasPin {
		// Account Profile option
		if option, ok := a.screen.Options[1]; ok {
			option.NextKey = utils.ACCOUNT_PROFILE
		}
		// Account Balances option
		if option, ok := a.screen.Options[2]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}
		// Account Withdrawal option
		if option, ok := a.screen.Options[3]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}
	}

	a.fetchUserSubscription()
}

func (a *Account) setAccountProfileOptions() {
	logger.UssdLog.Println("   ++ ACCOUNT: set account profile options")

	hasPin := a.checkHasPin()

	if hasPin {
		delete(a.screen.Next.Options, 1)

		hasSecurityQuestions := a.checkHasSecurityQuestions()

		if hasSecurityQuestions {
			if option, ok := a.screen.Next.Options[3]; ok {
				option.NextKey = utils.HAS_SECURITY_QUESTIONS
			}
		} else {
			if option, ok := a.screen.Next.Options[2]; ok {
				option.NextKey = utils.SECURITY_QUESTIONS_NOT_SET
			}
		}
	} else {
		// Update Profile option
		if option, ok := a.screen.Options[2]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}

		// Change Pin option
		if option, ok := a.screen.Next.Options[2]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}
		// Security Questions option
		if option, ok := a.screen.Next.Options[3]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}
	}

}

func (a *Account) getAccountBalances(input string) {
	logger.UssdLog.Println("   ++ ACCOUNT: get account balances")

	if input == "2" {
		a.fetchEarnings()
	} else if input == "3" {
		a.fetchSavings()
	} else if input == "4" {

		nextWithdrawal := "July 1st"
		if time.Now().Month() >= time.July {
			nextWithdrawal = "January 1st"
		}

		a.vars["{merchant_withdrawal_text}"] = "Locked for 6 months. Withdraw on " + nextWithdrawal

		a.fetchEarnings()
		a.fetchSavings()

	}
}

func (a *Account) fetchEarnings() {
	accountId := a.vars["{account_id}"]

	earnings, err := service.FetchEarningBalances(accountId)
	if err != nil {
		a.screen.Next.Type = "Sorry, we failed to fetch your earnings. Please try again later."
		logger.UssdLog.Error(err)
		//return
	}

	var purchasesAccount client.EarningAccount
	var subscriptionsAccount client.EarningAccount
	var withdrawalAccount client.EarningAccount
	var merchantAccount client.EarningAccount
	for _, earning := range earnings {
		if earning.Type == "PURCHASES" {
			purchasesAccount = earning
		}
		if earning.Type == "SUBSCRIPTIONS" {
			subscriptionsAccount = earning
		}
		if earning.Type == "WITHDRAWALS" {
			withdrawalAccount = earning
		}
		if earning.Type == "MERCHANT" {
			merchantAccount = earning
		}
	}

	pE := purchasesAccount.Self + purchasesAccount.Invite
	sE := subscriptionsAccount.Self + subscriptionsAccount.Invite

	total := pE + sE

	balance := total - withdrawalAccount.Self

	a.vars["{purchase_earnings}"] = formatAmount(pE, "")
	a.vars["{self_purchase_earnings}"] = formatAmount(purchasesAccount.Self, "")
	a.vars["{invite_purchase_earnings}"] = formatAmount(purchasesAccount.Invite, "")
	a.vars["{subscriptions_earnings}"] = formatAmount(sE, "")
	a.vars["{self_subscriptions_earnings}"] = formatAmount(subscriptionsAccount.Self, "")
	a.vars["{invite_subscriptions_earnings}"] = formatAmount(subscriptionsAccount.Invite, "")
	a.vars["{withdrawn_earnings}"] = formatAmount(withdrawalAccount.Self, "")
	a.vars["{earnings_balance}"] = formatAmount(balance, "")

	a.vars["{merchant_earnings}"] = formatAmount(merchantAccount.Self, "")

}

func (a *Account) fetchSavings() {
	accountId := a.vars["{account_id}"]

	earnings, err := service.FetchSavingBalances(accountId)
	if err != nil {
		a.screen.Next.Title = "Sorry, we failed to fetch your saved points. Please try again later."
		logger.UssdLog.Error(err)
		//return
	}

	var currentAccount client.SavingAccount
	var lockedAccount client.SavingAccount
	var merchantAccount client.SavingAccount
	for _, earning := range earnings {
		if earning.Type == "CURRENT" {
			currentAccount = earning
		}
		if earning.Type == "LOCKED" {
			lockedAccount = earning
		}
		if earning.Type == "MERCHANT" {
			merchantAccount = earning
		}
	}

	cS := currentAccount.Balance
	lS := lockedAccount.Balance

	total := cS + lS

	interest := currentAccount.Interest + lockedAccount.Interest

	wS := 0.0
	likelyCharge := float64(service.GetWithdrawalCharge(int(cS)))
	if cS > 25 {
		wS = math.Floor(cS - likelyCharge)
	}

	a.vars["{current_savings}"] = formatAmount(cS, "%.2f")
	a.vars["{locked_savings}"] = formatAmount(lS, "%.2f")
	a.vars["{total_savings}"] = formatAmount(total, "%.2f")
	a.vars["{withdrawable_savings}"] = formatAmount(wS, "%.0f")
	a.vars["{interest_savings}"] = formatAmount(interest, "")

	mE, _ := strconv.ParseFloat(a.vars["{merchant_earnings}"], 64)
	mEToday := mE - merchantAccount.Balance
	if mEToday < 0 {
		mEToday = 0
	}
	a.vars["{merchant_earnings_today}"] = formatAmount(mEToday, "%.0f")

	a.vars["{merchant_savings}"] = formatAmount(merchantAccount.Balance, "%.2f")
	a.vars["{merchant_interest}"] = formatAmount(merchantAccount.Interest, "%.2f")
}

func formatAmount(amount float64, format string) string {
	if format == "" {
		format = "%.4f"
	}
	if amount == 0 {
		format = "%.0f"
	}

	return fmt.Sprintf(format, amount)
}
