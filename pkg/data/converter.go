package data

import (
	"USSD.sidooh/utils"
	"encoding/json"
	"errors"
)

var loadScreenKeys = []string{
	// TODO: Reset to MAIN_MENU when done with invite_code beta
	utils.INVITE_CODE,

	// The below screens are hanging screens, i.e. have no parent
	utils.SUBSCRIPTION_RENEW,
	utils.PROFILE_SECURITY_QUESTIONS_ANSWER,
	utils.MERCHANT_SECURITY_QUESTIONS_ANSWER,
	utils.PIN_NOT_SET,
	utils.MERCHANT_PIN_NOT_SET,
	utils.VOUCHER_BALANCE_INSUFFICIENT,
	utils.FLOAT_BALANCE_INSUFFICIENT,
	utils.MERCHANT_TERMS,
	utils.MERCHANT_PROFILE_NEW_PIN,
	utils.MERCHANT_SUB_COUNTY,
	utils.MERCHANT_WARD,
	utils.MERCHANT_LANDMARK,
	utils.MERCHANT_KYB_CONFIRMATION,

	//TODO: Remove once county screen is decided
	utils.MERCHANT_LOCATION,
	utils.MERCHANT_LANDMARK_OTHER,

	//TODO: Remove once save is added
	utils.SAVE,
	utils.PROFILE_SECURITY_QUESTIONS_PIN,

	//TODO: Remove once merchant is added
	utils.MERCHANT_PAY_BILL,
	utils.MERCHANT_BUY_GOODS,
	utils.MERCHANT_CONSENT,
}

func LoadData() (map[string]*Screen, error) {
	file, err := utils.ReadFile(utils.DATA_FILE)
	if err != nil {
		return nil, err
	}

	var screens = map[string]*Screen{}
	err = json.Unmarshal(file, &screens)
	if err != nil {
		return nil, err
	}

	if len(screens) == 0 {
		return nil, errors.New("data file is empty")
	}

	for _, screenKey := range loadScreenKeys {
		setNextScreens(screens, screens[screenKey])
	}
	//// TODO: Reset to MAIN_MENU when done with invite_code beta
	//setNextScreens(screens, screens[utils.INVITE_CODE])
	//// The below screens are hanging screens, i.e. have no parent
	//setNextScreens(screens, screens[utils.SUBSCRIPTION_RENEW])
	//setNextScreens(screens, screens[utils.PROFILE_SECURITY_QUESTIONS_ANSWER])
	//setNextScreens(screens, screens[utils.PIN_NOT_SET])

	err = validateScreens(screens)
	if err != nil {
		return nil, err
	}

	return screens, nil
}

func setNextScreens(screens map[string]*Screen, current *Screen) {
	next, ok := screens[current.NextKey]
	if ok {
		if current.Next == nil {
			current.setNext(next)
			setNextScreens(screens, next)
		}
	} else {
		// Set default type
		if current.Type == "" {
			current.Type = utils.CLOSED
		}

		for _, option := range current.Options {
			next, ok = screens[option.NextKey]

			if ok {
				if option.Next == nil {
					option.setNext(next)
					setNextScreens(screens, next)
				}
			}
		}
	}
}

func validateScreens(screens map[string]*Screen) error {
	for _, d := range screens {
		err := d.Validate(true, true)
		if err != nil {
			return err
		}
	}

	return nil
}
