package data

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Screen struct {
	Key         string          `json:"key"`
	Title       string          `json:"title"`
	Type        string          `json:"type"`
	Options     map[int]*Option `json:"options,omitempty"`
	NextKey     string          `json:"next"`
	Next        *Screen         `json:"-"`
	Validations string          `json:"validations,omitempty"`
	Acyclic     bool            `json:"acyclic,omitempty"`
	Paginated   bool            `json:"paginated,omitempty"`
}

type ScreenPath struct {
	Screen
	Previous *ScreenPath `json:"previous,omitempty"`
}

func (path *ScreenPath) Encode() string {
	screenPath, err := json.Marshal(path)
	if err != nil {
		return ""
	}

	return string(screenPath)
}

var nextExceptionScreens = map[string]bool{
	utils.ABOUT:                      true,
	utils.CANCEL:                     true,
	utils.COMING_SOON:                true,
	utils.INVITE_END:                 true,
	utils.NOT_TRANSACTED:             true,
	utils.SUBSCRIPTION_ACTIVE:        true,
	utils.PROFILE_UPDATE_END:         true,
	utils.HAS_SECURITY_QUESTIONS:     true,
	utils.SECURITY_QUESTIONS_NOT_SET: true,
	utils.FLOAT_BALANCE_INSUFFICIENT: true,
}

var dynamicOptionScreens = map[string]bool{
	utils.AIRTIME_OTHER_NUMBER_SELECT: true,
	utils.UTILITY_ACCOUNT_SELECT:      true,
}

var optionsExceptionScreens = map[string]bool{
	utils.PROFILE_SECURITY_QUESTIONS_OPTIONS:  true,
	utils.MERCHANT_SECURITY_QUESTIONS_OPTIONS: true,
	utils.MERCHANT_COUNTY:                     true,
	utils.MERCHANT_SUB_COUNTY:                 true,
	utils.MERCHANT_WARD:                       true,
	utils.MERCHANT_LANDMARK:                   true,
}

func (screen *Screen) setNext(s *Screen) {
	screen.Next = s
}

func (screen *Screen) GetStringRep() string {
	// The below is needed in order to iterate over options in order.
	// Can be updated when go 1.18 is stable using generics
	// TODO: Check on this now that 1.18 is GA
	var keys []int
	for k := range screen.Options {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	// TODO: optimize by using strings.builder
	optionsString := ""
	for _, k := range keys {
		// Add new line between dynamic options and static one
		if _, ok := dynamicOptionScreens[screen.Key]; ok && k == 9 {
			optionsString += "\n"
		}

		optionsString += screen.Options[k].GetStringRep() + "\n"
	}

	return fmt.Sprintf("%v\n\n%v", screen.Title, optionsString)
}

func (screen *Screen) Validate(withOptions bool, recursive bool) error {
	if screen.Key == "" {
		return fmt.Errorf("key should be set for screen")
	}

	if screen.Title == "" {
		return fmt.Errorf("title should be set for screen " + screen.Key)
	}

	if screen.Type == "" {
		return fmt.Errorf("type should be set for screen " + screen.Key)
	}

	if screen.Type == utils.GENESIS || screen.Type == utils.CLOSED {
		// Validate that Next is not set
		if screen.Next != nil {
			return fmt.Errorf("next should not be set for screen " + screen.Key + " of type " + screen.Type)
		}

		// Validate that options exist
		if _, ok := optionsExceptionScreens[screen.Key]; !ok && len(screen.Options) == 0 {
			return fmt.Errorf("screen options are not set for screen " + screen.Key + " of type " + screen.Type)
		}

		// Validate that options are valid
		if withOptions {
			existingOptions := map[string]struct{}{}
			for _, option := range screen.Options {
				// Validate option
				err := option.Validate()
				if err != nil {
					panic(err)
				}

				//	Check if option already exists in list
				_, ok := existingOptions[option.Label]
				if ok {
					return fmt.Errorf("screen options for screen " + screen.Key + " contains duplicates of " + option.Label + " with value " + strconv.Itoa(option.Value))
				} else {
					existingOptions[option.Label] = struct{}{}
				}
			}
		}
	}

	if screen.Type == utils.OPEN {
		// Validate that Next must be set
		// exceptions about, acyclic screens...
		if screen.Next != nil {
			if recursive && !screen.Acyclic {
				err := screen.Next.Validate(withOptions, recursive)
				if err != nil {
					panic(err)
				}
			}
		} else {
			if _, ok := nextExceptionScreens[screen.Key]; ok != true {
				return fmt.Errorf("next is not set for screen " + screen.Key + " of type " + screen.Type)
			}
		}

		//	 Validate has no options
		if screen.Options != nil {
			return fmt.Errorf("screen options should not be set for screen " + screen.Key + " of type " + screen.Type)
		}
	}

	if screen.Type == utils.END {
		// Validate that Next must not set
		if screen.Next != nil {
			return fmt.Errorf("next should not be set for screen " + screen.Key + " of type " + screen.Type)
		}

		//	Validate has no options
		if screen.Options != nil {
			return fmt.Errorf("screen options should not be set for screen " + screen.Key + " of type " + screen.Type)
		}
	}

	return nil
}

func (screen *Screen) ValidateInput(input string, vars map[string]string) bool {
	// TODO: Sanitize input to be logged? e.g. pins, sec qns, etc. just replace with **** and use same length as input
	logger.UssdLog.Println("    Validating ", input, " against ", screen.Validations)

	if screen.Validations == "" {
		return true
	}

	validations := strings.Split(screen.Validations, ",")

	var currentValidationCheck = false
	for _, s := range validations {
		validation := strings.Split(s, ":")

		if !screen.checkValidation(validation, input, vars) {
			logger.UssdLog.Println("    *** ", input, " failed on ", validation)
			currentValidationCheck = false
			break
		}

		currentValidationCheck = true
	}

	logger.UssdLog.Println("    *** Validation result ", currentValidationCheck)

	return currentValidationCheck
}

func (screen *Screen) checkValidation(v []string, input string, vars map[string]string) bool {
	var validateAgainst = 0
	if len(v) > 1 {
		if v, e := strconv.Atoi(v[1]); e == nil {
			validateAgainst = v
		}
	}

	switch v[0] {
	case utils.INT:
		return getIntVal(input) > 0
	case utils.MIN:
		return getIntVal(input) >= validateAgainst
	case utils.MAX:
		return getIntVal(input) <= validateAgainst
	case utils.ALPHANUM:
		return isAlphaNumeric(input, validateAgainst)
	case utils.LENGTH:
		return len(input) == validateAgainst
	case utils.MAX_CHARS:
		return len(input) <= validateAgainst
	case utils.PHONE:
		return isValidPhone(input)
	case utils.DISALLOW_CURRENT:
		return screen.isNotCurrentPhone(input, vars["{phone}"])
	case utils.SAFARICOM:
		return screen.isValidPhoneAndProvider(input, utils.SAFARICOM)
	case utils.EXISTING_ACCOUNT:
		return screen.isSidoohAccount(input)
	case utils.EXISTING_MERCHANT_ACCOUNT:
		return screen.isSidoohMerchantAccount(input, vars)
	case utils.NOT_INVITED_OR_EXISTING_ACCOUNT:
		return screen.isUninvitedAndNonExistent(input)
	case utils.PIN_LENGTH:
		return screen.checkPinLength(input)
	case utils.PIN_CONFIRMED:
		return screen.confirmPin(input, vars)
	case utils.PIN:
		// TODO: Handle both -no pin set- and -invalid pin-
		// 	Also note, one may not have an account. maybe it is best if voucher isn't shown for first time user
		return screen.checkPin(input, vars)
	case utils.EXISTING_MERCHANT:
		return screen.checkMerchantIdNumber(input)
	case utils.UTILITY_AMOUNTS:
		return isValidUtilityAmount(input, vars["{selected_utility}"])
	case utils.NAME:
		return isValidName(input)
	case utils.SECURITY_QUESTION:
		return isValidSecurityQuestionAnswer(input)
	case utils.WITHDRAW_LIMITS:
		return isValidWithdrawalAmount(input, vars["{withdrawable_savings}"])

	case utils.MERCHANT_WITHDRAW_LIMITS:
		return screen.isValidMerchantWithdrawalAmount(input, vars["{withdrawable_earnings}"])

	case utils.INVITE_CODE_VALIDATION:
		return screen.isSocialInvite(input, vars) || screen.isSidoohAccountIdOrPhone(input, vars)
	}

	return false
}

func isValidWithdrawalAmount(input string, points string) bool {
	min := 25
	max := 1000

	val := getIntVal(input)
	pointsVal := getIntVal(points)

	return val <= max && val >= min && val <= pointsVal
}

func (screen *Screen) isValidMerchantWithdrawalAmount(input string, points string) bool {
	min := 10
	max := 10000

	val := getIntVal(input)
	pointsVal := getIntVal(points)

	valid := val <= max && val >= min && val <= pointsVal

	if val > max {
		screen.Title = "Amount is greater than threshold"
	}
	if val < min {
		screen.Title = "Amount is below the minimum required. Please enter amount above KES" + strconv.Itoa(min)
	}
	if val > pointsVal {
		screen.Title = "Amount is more than available balance. Your balance is KES" + points + ".\nEnter amount you wish to withdraw (MIN. KES10)"
	}

	return valid
}

func isValidSecurityQuestionAnswer(input string) bool {
	ansRegx := regexp.MustCompile(`^[A-z]{3,}$`)

	return ansRegx.MatchString(input)
}

func isValidName(input string) bool {
	nameRegx := regexp.MustCompile(`^[A-z .'-]{3,}$`)

	return nameRegx.MatchString(input)
}

func isValidUtilityAmount(input string, utility string) bool {
	min := 200
	max := 10000

	switch utility {
	case utils.KPLC_POSTPAID, utils.KPLC_PREPAID:
		max = 35000
	}

	val := getIntVal(input)

	return val <= max && val >= min
}

func (screen *Screen) isSocialInvite(input string, vars map[string]string) bool {
	inviteCodes := viper.GetString("INVITE_CODES")
	if inviteCodes == "" {
		return false
	}

	codes := strings.Split(inviteCodes, ",")

	for _, code := range codes {
		if strings.ToLower(code) == strings.ToLower(input) {
			vars["{invite_code_string}"] = strings.ToUpper(code)
			return true
		}
	}

	return false
}

func (screen *Screen) isSidoohAccount(input string) bool {
	account, err := service.CheckAccount(input)
	if err != nil {
		screen.Title = "Account not found, please try again."
		return false
	}

	if account != nil {
		return true
	}

	return false
}

func (screen *Screen) isSidoohMerchantAccount(input string, vars map[string]string) bool {
	account, err := service.CheckAccount(input)
	if err != nil {
		screen.Title = "Account not found, please try again."
		return false
	}

	if account != nil {
		merchant, err := service.GetMerchantByAccount(strconv.Itoa(account.Id))
		if err != nil {
			screen.Title = "Account not found, please try again."
			return false
		}

		if merchant != nil {
			vars["{merchant_account_validated}"] = strconv.Itoa(int(merchant.Id))
			vars["{merchant_account_validated_name}"] = merchant.BusinessName
			return true
		}
	}

	return false
}

func (screen *Screen) isSidoohAccountIdOrPhone(input string, vars map[string]string) bool {
	account, err := service.CheckAccountByIdOrPhone(input)
	if err != nil {
		screen.Title = "Invalid Code, please try again."
		return false
	}

	if account != nil {
		vars["{invite_code}"] = strconv.Itoa(account.Id)
		_, _ = service.CreateInvite(vars["{invite_code}"], vars["{phone}"], "CODE")
		return true
	}

	return false
}

func (screen *Screen) isUninvitedAndNonExistent(input string) bool {
	exists := service.InviteOrAccountExists(input)
	if exists {
		screen.Title = "Sorry, this number is not eligible for invite at the moment."
	}

	return !exists
}

func (screen *Screen) checkPinLength(input string) bool {
	// TODO: Should we also check for consecutive e.g. 1234, 5678
	actualVal := strconv.Itoa(getIntVal(input))
	if len(actualVal) == 4 {
		return true
	}

	screen.Title = "Please enter a valid pin (4 digits; Can't start with 0)"

	return false
}

func (screen *Screen) confirmPin(input string, vars map[string]string) bool {
	pin := vars["{pin}"]

	if pin == input {
		return true
	}
	screen.Title = "The PIN entered does not seem to match. Please try again."

	return false
}

func (screen *Screen) checkPin(input string, vars map[string]string) bool {
	if screen.checkPinLength(input) {
		id := vars["{account_id}"]
		pinTries, _ := strconv.Atoi(vars["{pin_tries}"])

		if pinTries > 2 {
			screen.NextKey = utils.PIN_BLOCKED
			//	TODO: Inform accounts, or use accounts to determine blockage
			// Notify relevant parties e.g. support...
		}

		isValid := service.CheckPin(id, input)

		if !isValid {
			vars["{pin_tries}"] = strconv.Itoa(pinTries + 1)
			screen.Title = "Invalid Pin!\nPlease try again."
		}

		return isValid
	}

	return false
}

func (screen *Screen) checkMerchantIdNumber(input string) bool {
	exists := service.MerchantIdNumberExists(input)
	if exists {
		screen.Title = "Hi, we are unable to process your signup request as the national Id you have provided is already registered on Sidooh and belongs to another customer\n\nEnter National ID Number"
	}

	return !exists
}

func (screen *Screen) isNotCurrentPhone(input string, phone string) bool {
	valid := !screen.isCurrentPhone(input, phone)
	if !valid {
		screen.Title = "Invalid phone.\nPlease use a number different from yours"
	}

	return valid
}

func (screen *Screen) isCurrentPhone(input string, phone string) bool {
	s, err := utils.FormatPhone(input)
	if err != nil {
		return false
	}

	return s == phone
}

func (screen *Screen) isValidPhoneAndProvider(input string, requiredProvider string) bool {
	provider, err := utils.GetPhoneProvider(input)
	if err != nil {
		return false
	}

	valid := provider == requiredProvider
	if !valid {
		screen.Title = "Invalid phone.\nPlease try again with a valid " + requiredProvider + " number."
	}

	return provider == requiredProvider
}

func isValidPhone(input string) bool {
	_, err := utils.GetPhoneProvider(input)
	if err != nil {
		logger.UssdLog.Println("isValidPhoneError", err)
		return false
	}

	return true
}

func getIntVal(str string) int {
	if v, e := strconv.Atoi(str); e == nil {
		return v
	}
	return 0
}

func isAlphaNumeric(str string, count int) bool {
	if count == 0 {
		count = 64
	}
	alphaNumRegx := regexp.MustCompile(fmt.Sprintf(`^[A-z0-9 ]{0,%v}$`, count))

	return alphaNumRegx.MatchString(str)
}

func (screen *Screen) SubstituteVars(vars map[string]string) {
	screen.Title = Strtr(screen.Title, vars)

	for _, v := range screen.Options {
		v.Label = Strtr(v.Label, vars)
	}
}

// Strtr SOURCE: https://github.com/syyongx/php2go/blob/master/php.go
// Strtr strtr()
//
// If the parameter length is 1, type is: map[string]string
// Strtr("baab", map[string]string{"ab": "01"}) will return "ba01"
// If the parameter length is 2, type is: string, string
// Strtr("baab", "ab", "01") will return "1001", a => 0; b => 1.
func Strtr(haystack string, params ...interface{}) string {
	ac := len(params)

	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)

		if length == 0 {
			return haystack
		}

		oldnew := make([]string, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}

		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)

		if trlen > lt {
			trlen = lt
		}

		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8

		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}

			return string(str)
		}

		// trlen != 1
		for {
			xlat[j] = j

			if j++; j == 0 {
				break
			}
		}
		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}
		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}

		return string(str)
	}

	return haystack
}
