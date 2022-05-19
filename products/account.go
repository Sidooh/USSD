package products

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"strconv"
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
		a.checkHasSecurityQuestions()

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
		err := service.SetPin(accountId, pin)
		if !err {
			a.screen.Next.Title = "Sorry. We failed to set your pin, please try again later."
		} else {
			//	TODO: Notify user of new pin set and also ask to set id and security questions
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
		if option, ok := a.screen.Next.Options[2]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}

		if option, ok := a.screen.Next.Options[3]; ok {
			option.NextKey = utils.PIN_NOT_SET
		}
	}
}

func (a *Account) checkHasSecurityQuestions() {
	accountId := a.vars["{account_id}"]

	hasSecurityQuestions := service.CheckHasSecurityQuestions(accountId)

	if hasSecurityQuestions {
		a.screen.Next.Options[3].NextKey = utils.HAS_SECURITY_QUESTIONS
	} else {
		if option, ok := a.screen.Next.Options[2]; ok {
			option.NextKey = utils.SECURITY_QUESTIONS_NOT_SET
		}
	}
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
		a.screen.Options[3].NextKey = utils.COMING_SOON
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
