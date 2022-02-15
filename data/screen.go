package data

import (
	"fmt"
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
	"about": true, "cancel": true, "coming_soon": true,
	"refer_end": true,
}

func (screen *Screen) setNext(s *Screen) {
	screen.Next = s
}

func (screen *Screen) GetStringRep() string {
	// The below is needed in order to iterate over options in order.
	// Can be updated when go 1.18 is stable using generics
	var keys []int
	for k := range screen.Options {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	optionsString := ""
	for _, k := range keys {
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

	if screen.Type == GENESIS || screen.Type == CLOSED {
		// Validate that Next is not set
		if screen.Next != nil {
			return fmt.Errorf("next should not be set for screen " + screen.Key + " of type " + screen.Type)
		}

		// Validate that options exist
		if len(screen.Options) == 0 {
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
					return fmt.Errorf("screen options contains duplicates of " + option.Label + " with value " + strconv.Itoa(option.Value))
				} else {
					existingOptions[option.Label] = struct{}{}
				}
			}
		}
	}

	if screen.Type == OPEN {
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

	if screen.Type == END {
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

func (screen *Screen) ValidateInput(input string) bool {
	fmt.Println("\tValidating ", input, " against ", screen.Validations)

	validations := strings.Split(screen.Validations, ",")

	var currentValidationCheck bool = false
	for _, s := range validations {
		validation := strings.Split(s, ":")

		if !screen.checkValidation(validation, input) {
			fmt.Println("\t*** ", input, " failed on ", validation)
			currentValidationCheck = false
			break
		}

		currentValidationCheck = true
	}

	fmt.Println("\t*** Validation result ", currentValidationCheck)

	return currentValidationCheck
}

func (screen *Screen) checkValidation(v []string, input string) bool {
	var validateAgainst = 0
	if len(v) > 1 {
		if v, e := strconv.Atoi(v[1]); e == nil {
			validateAgainst = v
		}
	}

	switch v[0] {
	case INT:
		return getIntVal(input) > 0
	case MIN:
		return getIntVal(input) >= validateAgainst
	case MAX:
		return getIntVal(input) <= validateAgainst
	case PHONE:
		return isValidPhone(input)
	case DISALLOW_CURRENT:
		return isValidPhone(input)
	case MPESA_NUMBER:
		return getIntVal(input) <= validateAgainst
	}

	return false
}

func isValidPhone(input string) bool {
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
