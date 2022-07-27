package datastore

import (
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
	_ = WriteFile("test", filename)
}

func TestReadFile(t *testing.T) {
	file, err := ReadFile(filename)
	if err == nil {
		t.Errorf("ReadFile(%s) = %s; want err", filename, file)
	}

	createTestFile()

	file, err = ReadFile(filename)
	if file == nil {
		t.Errorf("ReadFile(%s) = %s; want file", filename, err)
	}
}

func TestUnmarshalFromFile(t *testing.T) {
	_ = os.Remove(filename)
	data := ""

	err := UnmarshalFromFile(data, filename)
	if err == nil {
		t.Errorf("UnmarshalFromFile(string, %s) = %s; want err", filename, err)
	}

	createTestFile()

	err = UnmarshalFromFile(data, filename)
	if err != nil {
		t.Errorf("UnmarshalFromFile(string, %s) = %s; want nil", filename, err)
	}
}

func TestWriteFile(t *testing.T) {
	data := "test"

	err := WriteFile(data, filename)
	if err != nil {
		t.Errorf("WriteFile(string, %s) = %s; want nil", filename, err)
	}

	data = "testLong"

	err = WriteFile(data, filename)
	if err != nil {
		t.Errorf("WriteFile(string, %s) = %s; want nil", filename, err)
	}
}
