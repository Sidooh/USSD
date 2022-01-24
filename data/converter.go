package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var screens = map[string]*Screen{}

func ReadFile(filename string) ([]byte, error) {
	wd, err := os.Getwd()
	file, err := os.ReadFile(filepath.Join(wd, "data/", filename))
	if err != nil {
		return nil, err
	}

	return file, err
}

func UnmarshalFromFile(to interface{}, filename string) error {
	file, err := ReadFile("state")
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
	err = os.WriteFile(filepath.Join(wd, filename), marshal, 0777)
	if err != nil {
		return err
	}

	return nil
}

func LoadData() (map[string]*Screen, error) {
	file, err := ReadFile("data.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &screens)
	if err != nil {
		return nil, err
	}

	setNextScreens(screens["main_menu"])

	validateScreens()

	return screens, nil
}

func setNextScreens(current *Screen) {
	next, ok := screens[current.NextKey]
	if ok {

		if current.Next == nil {
			fmt.Println(current.Key, "Setting next", next.Key)

			current.setNext(next)
			setNextScreens(next)
		}

	} else {
		for _, option := range current.Options {
			next, ok = screens[option.NextKey]
			if ok {

				if option.Next == nil {
					fmt.Println(current.Key, option.Value, "Setting next", next.Key)

					option.setNext(next)
					setNextScreens(next)
				}
			}
		}
	}
}

func validateScreens() {
	for _, d := range screens {
		err := d.Validate()
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Validated %v screens successfully\n", len(screens))
}
