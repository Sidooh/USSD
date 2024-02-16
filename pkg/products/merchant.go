package products

import (
	"USSD.sidooh/pkg/data"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"encoding/json"
	"strconv"
)

type Merchant struct {
	Product
}

func (m *Merchant) Process(input string) {
	logger.UssdLog.Println(" -- MERCHANT: process", m.screen.Key, input)
	m.productRep = "merchant"

	m.Product.Process(input)
	m.processScreen(input)
	m.finalize()
}

func (m *Merchant) processScreen(input string) {
	switch m.screen.Key {
	case utils.MERCHANT_FIRST_NAME:
		m.vars["{first_name}"] = input
	case utils.MERCHANT_LAST_NAME:
		m.vars["{last_name}"] = input
		m.vars["{name}"] = m.vars["{first_name}"] + " " + input
	case utils.MERCHANT_ID_NUMBER:
		m.vars["{id_number}"] = input
		//if m.checkHasPin() {
		//	m.screen.Next.Options[1].NextKey = utils.MERCHANT_SECURITY_QUESTIONS_OPTIONS
		//m.fetchSecurityQuestionOptions()
		//}
	case utils.MERCHANT_NEW_PIN:
		m.vars["{pin}"] = input
	case utils.MERCHANT_NEW_PIN_CONFIRM:
		m.vars["{confirm_pin}"] = input
		m.fetchSecurityQuestionOptions()

	case utils.MERCHANT_SECURITY_QUESTIONS_OPTIONS:
		m.processQuestionSelection(input)

	case utils.MERCHANT_SECURITY_QUESTIONS_ANSWER:
		m.processAnswer(input)

	case utils.MERCHANT_KYB:
		m.vars["{business_name}"] = input
		m.fetchCountyOptions()
	case utils.MERCHANT_COUNTY:
		m.processCountySelection(input)
	case utils.MERCHANT_SUB_COUNTY:
		m.processSubCountySelection(input)
	case utils.MERCHANT_WARD:
		m.processWardSelection(input)
	case utils.MERCHANT_LANDMARK:
		m.processLandmarkSelection(input)
	case utils.MERCHANT_LANDMARK_OTHER:
		m.vars["{landmark_id}"] = input
		m.vars["{landmark}"] = input
	}
}

func (m *Merchant) finalize() {
	logger.UssdLog.Println(" -- MERCHANT: finalize", m.screen.Next.Type)

	if m.screen.Key == utils.MERCHANT_TERMS {
		accountId, _ := strconv.Atoi(m.vars["{account_id}"])
		if accountId == 0 {
			logger.UssdLog.Println(" -- MERCHANT: creating acc")

			account, err := service.CreateAccount(m.vars["{phone}"], m.vars["{invite_code_string}"])
			if err != nil {
				// TODO: Send message to user
				logger.UssdLog.Error("Failed to create account: ", err)
			}

			accountId = account.Id
			m.vars["{account_id}"] = strconv.Itoa(accountId)
		}

		//set pin asynchronously
		go func() {
			if !m.checkHasPin() && m.vars["{confirm_pin}"] != "" {
				service.SetPin(strconv.Itoa(accountId), m.vars["{confirm_pin}"])
			}
		}()

		//set sec qns asynchronously
		go func() {
			if _, ok := m.vars["{question_answers}"]; ok {
				questionAnswerVars := map[string]string{}

				_ = json.Unmarshal([]byte(m.vars["{question_answers}"]), &questionAnswerVars)

				// TODO: Make into goroutine if applicable
				// TODO: Should we check returned value? Or should we make it a void function?
				_ = service.SetSecurityQuestions(strconv.Itoa(accountId), questionAnswerVars)
			}
		}()

		request := client.MerchantKYCDetails{
			FirstName: m.vars["{first_name}"],
			LastName:  m.vars["{last_name}"],
			IdNumber:  m.vars["{id_number}"],
			AccountId: accountId,
		}
		service.CreateMerchant(request)
	}

	if m.screen.Key == utils.MERCHANT_KYB_CONSENT {
		request := client.MerchantKYBDetails{
			BusinessName: m.vars["{business_name}"],
			Landmark:     m.vars["{landmark_id}"],
		}
		service.UpdateMerchant(m.vars["{merchant_id}"], request)
	}
}

func (m *Merchant) fetchCountyOptions() {
	logger.UssdLog.Println("   ++ MERCHANT: fetch county options")

	counties, _ := service.FetchCounties()

	m.screen.Next.Options = map[int]*data.Option{}

	for i, county := range *counties {
		if county.Id == 24 || county.Id == 25 {
			continue
		}
		m.screen.Next.Options[i+1] = &data.Option{
			Label:      county.County,
			Value:      i + 1,
			NextKey:    utils.MERCHANT_SUB_COUNTY,
			NoFullStop: true,
		}
	}
}

func (m *Merchant) processCountySelection(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process county selection", input)

	counties, _ := service.FetchCounties()

	selectedOption, _ := strconv.Atoi(input)

	for _, county := range *counties {
		if county.Id == selectedOption {
			m.vars["{county_id}"] = input
			m.vars["{county}"] = county.County
		}
	}

	m.fetchSubCountyOptions(input)
}

func (m *Merchant) fetchSubCountyOptions(county string) {
	logger.UssdLog.Println("   ++ MERCHANT: fetch sub county options")

	subCounties, _ := service.FetchSubCounties(county)

	m.screen.Next.Options = map[int]*data.Option{}
	subCountiesOptionsMapper := map[int]int{}

	for i, subCounty := range *subCounties {
		m.screen.Next.Options[i+1] = &data.Option{
			Label:   subCounty.SubCounty,
			Value:   i + 1,
			NextKey: utils.MERCHANT_WARD,
		}

		subCountiesOptionsMapper[i+1] = subCounty.Id
	}

	stringVars, _ := json.Marshal(subCountiesOptionsMapper)
	m.vars["{sub_county_options}"] = string(stringVars)
}

func (m *Merchant) processSubCountySelection(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process sub county selection", input)

	subCountyOptionVars := map[string]int{}

	_ = json.Unmarshal([]byte(m.vars["{sub_county_options}"]), &subCountyOptionVars)
	delete(m.vars, "{sub_county_options}")

	m.vars["{sub_county_id}"] = strconv.Itoa(subCountyOptionVars[input])

	subCounties, _ := service.FetchSubCounties(m.vars["{county_id}"])

	for _, subCounty := range *subCounties {
		if subCounty.Id == subCountyOptionVars[input] {
			m.vars["{sub_county}"] = subCounty.SubCounty
		}
	}

	m.fetchWardOptions(m.vars["{sub_county_id}"])
}

func (m *Merchant) fetchWardOptions(subCounty string) {
	logger.UssdLog.Println("   ++ MERCHANT: fetch ward options")

	wards, _ := service.FetchWards(m.vars["{county_id}"], subCounty)

	m.screen.Next.Options = map[int]*data.Option{}

	//wardOptionVars := map[int]client.Ward{}

	wardsOptionsMapper := map[int]int{}

	for i, ward := range *wards {
		m.screen.Next.Options[i+1] = &data.Option{
			Label:   ward.Ward,
			Value:   i + 1,
			NextKey: utils.MERCHANT_LANDMARK,
		}

		wardsOptionsMapper[i+1] = ward.Id

		//wardOptionVars[i+1] = subCounty
	}

	stringVars, _ := json.Marshal(wardsOptionsMapper)
	m.vars["{ward_options}"] = string(stringVars)
}

func (m *Merchant) processWardSelection(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process ward selection", input)

	//TODO: is using string more efficient than int conversions?
	//selectedWard, _ := strconv.Atoi(input)
	wardOptionVars := map[string]int{}

	_ = json.Unmarshal([]byte(m.vars["{ward_options}"]), &wardOptionVars)
	delete(m.vars, "{ward_options}")

	m.vars["{ward_id}"] = strconv.Itoa(wardOptionVars[input])

	wards, _ := service.FetchWards(m.vars["{county_id}"], m.vars["{sub_county_id}"])

	for _, ward := range *wards {
		if ward.Id == wardOptionVars[input] {
			m.vars["{ward}"] = ward.Ward
		}
	}

	m.fetchLandmarkOptions(m.vars["{ward_id}"])
}

func (m *Merchant) fetchLandmarkOptions(ward string) {
	logger.UssdLog.Println("   ++ MERCHANT: fetch landmark options")

	landmarks, _ := service.FetchLandmarks(m.vars["{county_id}"], m.vars["{sub_county_id}"], ward)

	m.screen.Next.Options = map[int]*data.Option{}

	landmarksOptionsMapper := map[int]int{}

	for i, landmark := range *landmarks {
		landmarksOptionsMapper[i+1] = landmark.Id

		if len(*landmarks) > 10 && i == 9 {
			m.screen.Next.Options[i+1] = &data.Option{
				Label:   "Other",
				Value:   i + 1,
				NextKey: utils.MERCHANT_LANDMARK_OTHER,
			}

			break
		}

		m.screen.Next.Options[i+1] = &data.Option{
			Label:   landmark.Landmark,
			Value:   i + 1,
			NextKey: utils.MERCHANT_KYB_CONFIRMATION,
		}
	}

	stringVars, _ := json.Marshal(landmarksOptionsMapper)
	m.vars["{landmark_options}"] = string(stringVars)
}

func (m *Merchant) processLandmarkSelection(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process landmark selection", input)

	if input == "10" {
		return
	}

	selectedLandmark, _ := strconv.Atoi(input)
	landmarkOptionVars := map[int]int{}

	_ = json.Unmarshal([]byte(m.vars["{landmark_options}"]), &landmarkOptionVars)
	delete(m.vars, "{landmark_options}")

	m.vars["{landmark_id}"] = strconv.Itoa(landmarkOptionVars[selectedLandmark])

	landmarks, _ := service.FetchLandmarks(m.vars["{county_id}"], m.vars["{sub_county_id}"], m.vars["{ward_id}"])

	for _, landmark := range *landmarks {
		if landmark.Id == landmarkOptionVars[selectedLandmark] {
			m.vars["{landmark}"] = landmark.Landmark
		}
	}
}

// TODO: Test Pay bill and goods works accordingly
//func (m *Merchant) getCharge(input string) {
//	amount, _ := strconv.Atoi(input)
//	fee := 0
//
//	if m.vars["{merchant_type}"] == utils.MPESA_PAY_BILL {
//		fee = service.GetPayBillCharge(amount)
//	} else {
//		fee = service.GetBuyGoodsCharge(amount)
//	}
//
//	m.vars["{merchant_fee}"] = strconv.Itoa(fee)
//}

func (m *Merchant) fetchSecurityQuestionOptions() {
	logger.UssdLog.Println("   ++ MERCHANT: fetch security question options")

	questions, _ := service.FetchSecurityQuestions()

	questionAnswerVars := map[uint]string{}
	_ = json.Unmarshal([]byte(m.vars["{question_answers}"]), &questionAnswerVars)

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

		m.screen.Next.Options = map[int]*data.Option{}

		for i, question := range maxQuestions {
			m.screen.Next.Options[i+1] = &data.Option{
				Label:   question.Question,
				Value:   i + 1,
				NextKey: utils.MERCHANT_SECURITY_QUESTIONS_ANSWER,
			}

			questionOptionVars[i+1] = question
		}
		stringVars, _ := json.Marshal(questionOptionVars)
		m.vars["{question_options}"] = string(stringVars)
	} else {
		//if _, exists := m.screen.Options[3]; exists {
		//	m.screen.Options[3].NextKey = utils.COMING_SOON
		//}
	}
}

func (m *Merchant) processQuestionSelection(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process question selection", input)

	selectedQuestion, _ := strconv.Atoi(input)
	questionOptionVars := map[int]client.SecurityQuestion{}

	_ = json.Unmarshal([]byte(m.vars["{question_options}"]), &questionOptionVars)

	m.vars["{security_question}"] = questionOptionVars[selectedQuestion].Question
	m.vars["{security_question_id}"] = strconv.Itoa(int(questionOptionVars[selectedQuestion].Id))

}

func (m *Merchant) processAnswer(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process answer", input)

	questionAnswerVars := map[string]string{}

	_ = json.Unmarshal([]byte(m.vars["{question_answers}"]), &questionAnswerVars)

	questionAnswerVars[m.vars["{security_question_id}"]] = input

	stringVars, _ := json.Marshal(questionAnswerVars)

	m.vars["{question_answers}"] = string(stringVars)

	//fmt.Println(questionAnswerVars, m.screen)
	if len(questionAnswerVars) == 3 {
		m.screen.NextKey = utils.MERCHANT_TERMS
	} else {
		m.fetchSecurityQuestionOptions()
	}
	//fmt.Println(questionAnswerVars, m.screen)

}

func (m *Merchant) fetchUserSecurityQuestionOptions() {
	logger.UssdLog.Println("   ++ MERCHANT: fetch user security question options")

	// Check if user already has questions in state
	var userQuestions []client.UserSecurityQuestion
	_ = json.Unmarshal([]byte(m.vars["{user_questions}"]), &userQuestions)

	// Fetch from API otherwise
	if len(userQuestions) == 0 {
		accountId := m.vars["{account_id}"]
		userQuestions, _ = service.FetchUserSecurityQuestions(accountId)

		stringVars, _ := json.Marshal(userQuestions)
		m.vars["{user_questions}"] = string(stringVars)
	}

	// Get the answered questions
	questionAnswerVars := map[uint]string{}
	_ = json.Unmarshal([]byte(m.vars["{question_answers}"]), &questionAnswerVars)

	// Filter only unanswered questions so we pick from them
	var unansweredQuestions []client.UserSecurityQuestion
	for _, question := range userQuestions {
		if _, ok := questionAnswerVars[question.Question.Id]; !ok {
			unansweredQuestions = append(unansweredQuestions, question)
		}
	}

	// Ensure there are still unanswered questions, otherwise proceed to verify them
	if len(unansweredQuestions) != 0 {
		m.vars["{security_question}"] = unansweredQuestions[0].Question.Question
		m.vars["{security_question_id}"] = strconv.Itoa(int(unansweredQuestions[0].Question.Id))
	} else {
		//m.screen.Options[3].NextKey = utils.COMING_SOON
	}

}

func (m *Merchant) processUserAnswer(input string) {
	logger.UssdLog.Println("   ++ MERCHANT: process user answer", input)

	questionAnswerVars := map[string]string{}

	_ = json.Unmarshal([]byte(m.vars["{question_answers}"]), &questionAnswerVars)

	questionAnswerVars[m.vars["{security_question_id}"]] = input

	stringVars, _ := json.Marshal(questionAnswerVars)

	m.vars["{question_answers}"] = string(stringVars)

	if len(questionAnswerVars) == 3 {
		m.screen.NextKey = utils.MERCHANT_PROFILE_NEW_PIN
	} else {
		m.fetchUserSecurityQuestionOptions()
	}
}
