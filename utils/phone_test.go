package utils

import (
	"os"
	"testing"
)

var (
	validPhones   = map[string]string{}
	invalidPhones = map[string]string{}
)

func TestMain(m *testing.M) {
	validPhones["700000000"] = SAFARICOM
	validPhones["748000000"] = SAFARICOM
	validPhones["110000000"] = SAFARICOM
	validPhones["730000000"] = AIRTEL
	validPhones["762000000"] = AIRTEL
	validPhones["106000000"] = AIRTEL
	validPhones["779000000"] = TELKOM
	validPhones["764000000"] = EQUITEL
	validPhones["747000000"] = FAIBA

	validPhones["+254110000000"] = SAFARICOM
	validPhones["254762000000"] = AIRTEL
	validPhones["0779000000"] = TELKOM
	validPhones["0106000000"] = AIRTEL
	validPhones["254110000000"] = SAFARICOM
	validPhones["0110000000"] = SAFARICOM

	invalidPhones["107000000"] = "err"
	invalidPhones["116000000"] = "err"

	os.Exit(m.Run())
}

func TestGetPhoneProvider(t *testing.T) {
	for phone, provider := range validPhones {
		got, err := GetPhoneProvider(phone)
		if err != nil {
			t.Errorf("GetPhoneProvider(%s): got %s; expect %s", phone, err, provider)
		}

		if got != provider {
			t.Errorf("GetPhoneProvider(%s): got %s; expect %s", phone, got, provider)
		}
	}

	for phone, provider := range invalidPhones {
		got, err := GetPhoneProvider(phone)
		if err == nil {
			t.Errorf("GetPhoneProvider(%s): got %s; expect %v", phone, got, provider)
		}
	}
}

func TestFormatPhone(t *testing.T) {
	for phone, _ := range validPhones {
		formattedPhone, err := FormatPhone(phone)
		if err != nil {
			t.Errorf("FormatPhone(%s): got %v; expect %v", phone, err, "254"+phone)
		}

		if formattedPhone != "254"+phone {

		}
	}
}
