package datastore

import (
	"USSD.sidooh/utils"
	"encoding/json"
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
