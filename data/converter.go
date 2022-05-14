package data

import (
	"USSD.sidooh/utils"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func ReadFile(filename string) ([]byte, error) {
	wd, err := os.Getwd()
	file, err := os.ReadFile(filepath.Join(wd, utils.DATA_DIRECTORY, filename))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func UnmarshalFromFile(to interface{}, filename string) error {
	file, err := ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &to)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(data interface{}, filename string) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	err = os.WriteFile(filepath.Join(wd, utils.DATA_DIRECTORY, filename), marshal, 0644)
	if err != nil {
		return err
	}

	return nil
}

func RemoveFile(filename string) error {
	wd, err := os.Getwd()
	err = os.Remove(filepath.Join(wd, utils.DATA_DIRECTORY, filename))
	if err != nil {
		return err
	}

	return nil
}

func LoadData() (map[string]*Screen, error) {
	file, err := ReadFile(utils.DATA_FILE)
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

	setNextScreens(screens, screens[utils.MAIN_MENU])
	// The below screens are hanging screens, i.e. have no parent
	setNextScreens(screens, screens[utils.SUBSCRIPTION_RENEW])
	setNextScreens(screens, screens[utils.PROFILE_SECURITY_QUESTIONS_FIRST_CHOICE])

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
