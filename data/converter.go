package data

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func ReadFile(filename string) ([]byte, error) {
	wd, err := os.Getwd()
	file, err := os.ReadFile(filepath.Join(wd, DATA_DIRECTORY, filename))
	if err != nil {
		return nil, err
	}

	return file, err
}

func UnmarshalFromFile(to interface{}, filename string) error {
	file, err := ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &to)
	if err != nil {
		return nil
	}
	return nil
}

func WriteFile(data interface{}, filename string) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	err = os.WriteFile(filepath.Join(wd, DATA_DIRECTORY, filename), marshal, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadData() (map[string]*Screen, error) {
	file, err := ReadFile(DATA_FILE)
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

	setNextScreens(screens, screens[MAIN_MENU])

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
			setNextScreens(nil, next)
		}

	} else {

		for _, option := range current.Options {
			next, ok = screens[option.NextKey]

			if ok {

				if option.Next == nil {
					option.setNext(next)
					setNextScreens(nil, next)
				}

			}

		}

	}
}

func validateScreens(screens map[string]*Screen) error {
	for _, d := range screens {
		err := d.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
