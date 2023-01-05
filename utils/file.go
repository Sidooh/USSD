package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func GetFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Error(err)
		return nil
	}

	return file
}

func GetLogFile(filename string) *os.File {
	pwd, err := os.Getwd()

	file := GetFile(filepath.Join(pwd, "logs/", filename))
	if err != nil || file == nil {
		file = os.Stdout
	}

	return file
}

func ReadFile(filename string) ([]byte, error) {
	wd, err := os.Getwd()
	file, err := os.ReadFile(filepath.Join(wd, DATA_DIRECTORY, filename))
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
	err = os.WriteFile(filepath.Join(wd, DATA_DIRECTORY, filename), marshal, 0644)
	if err != nil {
		return err
	}

	return nil
}
