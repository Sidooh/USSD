package data

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/utils"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Screen struct {
	Key         string          `json:"key"`
	Title       string          `json:"title"`
	Type        string          `json:"type"`
	Options     map[int]*Option `json:"options"`
	NextKey     string          `json:"next"`
	Next        *Screen         `json:"-"`
	Validations string          `json:"validations"`
}

type ScreenPath struct {
	Screen
	Previous *ScreenPath `json:"previous"`
}

var nextExceptionScreens = map[string]bool{
	utils.ABOUT:               true,
	utils.CANCEL:              true,
	utils.COMING_SOON:         true,
	utils.INVITE_END:          true,
	utils.NOT_TRANSACTED:      true,
	utils.SUBSCRIPTION_ACTIVE: true,
	utils.PIN_NOT_SET:         true,
	utils.PROFILE_UPDATE_END:  true,
}

var dynamicOptionScreens = map[string]bool{
	utils.AIRTIME_OTHER_NUMBER_SELECT: true,
	utils.UTILITY_ACCOUNT_SELECT:      true,
}

var optionsExceptionScreens = map[string]bool{
	utils.PROFILE_SECURITY_QUESTIONS_FIRST_OPTIONS: true,
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
		// exceptions about,
		if screen.Next != nil {
			if recursive {
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
	// TODO: Sanitize input to be logged? e.g. pins, etc. just replace with **** and use same length as input
	logger.UssdLog.Println("    Validating ", input, " against ", screen.Validations)

	validations := strings.Split(screen.Validations, ",")

	var currentValidationCheck bool = false
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
	case utils.PHONE:
		return isValidPhone(input)
	case utils.DISALLOW_CURRENT:
		return screen.isNotCurrentPhone(input, vars["{phone}"])
	case utils.SAFARICOM:
		return screen.isValidPhoneAndProvider(input, utils.SAFARICOM)
	case utils.EXISTING_ACCOUNT:
		return screen.isSidoohAccount(input)
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
	case utils.UTILITY_AMOUNTS:
		return isValidUtilityAmount(input, vars["{selected_utility}"])
	case utils.NAME:
		return isValidName(input)
	}

	return false
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
