package data

import (
	"USSD.sidooh/datastore"
	"USSD.sidooh/utils"
	"os"
	"testing"
)

var filename = "test.file"

func TestMain(m *testing.M) {
	utils.DATA_DIRECTORY = ""
	utils.DATA_FILE = "test-data.json"

	_ = os.Remove(filename)
	_ = os.Remove(utils.DATA_FILE)

	os.Exit(m.Run())
}

func createTestFile() {
	_ = datastore.WriteFile("test", filename)
}

func createDataTestFile(data interface{}) {
	_ = datastore.WriteFile(data, utils.DATA_FILE)
}

func TestReadFile(t *testing.T) {
	file, err := datastore.ReadFile(filename)
	if err == nil {
		t.Errorf("ReadFile(%s) = %s; want err", filename, file)
	}

	createTestFile()

	file, err = datastore.ReadFile(filename)
	if file == nil {
		t.Errorf("ReadFile(%s) = %s; want file", filename, err)
	}
}

func TestUnmarshalFromFile(t *testing.T) {
	_ = os.Remove(filename)
	data := ""

	err := datastore.UnmarshalFromFile(data, filename)
	if err == nil {
		t.Errorf("UnmarshalFromFile(string, %s) = %s; want err", filename, err)
	}

	createTestFile()

	err = datastore.UnmarshalFromFile(data, filename)
	if err != nil {
		t.Errorf("UnmarshalFromFile(string, %s) = %s; want nil", filename, err)
	}
}

func TestWriteFile(t *testing.T) {
	data := "test"

	err := datastore.WriteFile(data, filename)
	if err != nil {
		t.Errorf("WriteFile(string, %s) = %s; want nil", filename, err)
	}

	data = "testLong"

	err = datastore.WriteFile(data, filename)
	if err != nil {
		t.Errorf("WriteFile(string, %s) = %s; want nil", filename, err)
	}
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
