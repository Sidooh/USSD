package data

import (
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utils.DATA_DIRECTORY = ""
	utils.DATA_FILE = "test-data.json"

	_ = os.Remove(utils.DATA_FILE)

	os.Exit(m.Run())
}

func createDataTestFile(data interface{}) {
	_ = datastore.WriteFile(data, utils.DATA_FILE)
}

func TestLoadData(t *testing.T) {
	//1. Test non-existent file
	data, err := LoadData()
	if err == nil {
		t.Errorf("LoadData() = %v; want err", data)
	}

	//2. Test invalid content file
	createDataTestFile("test")

	data, err = LoadData()
	if err == nil {
		t.Errorf("LoadData() = %v; want err", data)
	}

	//3. Test empty json file
	createDataTestFile(map[string]*Screen{})

	data, err = LoadData()
	if err == nil {
		t.Errorf("LoadData() = %v; want err", data)
	}

	//4. Test valid file
	createDataTestFile(map[string]*Screen{
		"main_menu": {
			Key:   "main_menu",
			Title: "Main Menu Title",
			Options: map[int]*Option{
				1: {
					Label:   "end",
					NextKey: "end",
				},
			},
		},
		"end": {
			Key:   "end",
			Type:  utils.END,
			Title: "End Title",
		},
	})

	loadScreenKeys = []string{utils.MAIN_MENU}
	data, err = LoadData()
	if err != nil {
		t.Errorf("LoadData() = %v; want data", err)
	}
}
