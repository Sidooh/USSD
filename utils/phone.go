package utils

import (
	"errors"
	"fmt"
	"regexp"
)

func GetPhoneProvider(input string) (string, error) {
	sReg := regexp.MustCompile(SAFARICOM_REGEX)
	aReg := regexp.MustCompile(AIRTEL_REGEX)
	tReg := regexp.MustCompile(TELKOM_REGEX)
	eReg := regexp.MustCompile(EQUITEL_REGEX)
	fReg := regexp.MustCompile(FAIBA_REGEX)

	switch true {
	case sReg.MatchString(input):
		return SAFARICOM, nil
	case aReg.MatchString(input):
		return AIRTEL, nil
	case tReg.MatchString(input):
		return TELKOM, nil
	case eReg.MatchString(input):
		return EQUITEL, nil
	case fReg.MatchString(input):
		return FAIBA, nil
	}

	return "", errors.New(fmt.Sprintf("phone %s does not seem to be supported", input))
}

func FormatPhone(input string) (string, error) {
	_, err := GetPhoneProvider(input)
	if err != nil {
		return "", err
	}

	return "254" + string(input[len(input)-9:]), nil
}
