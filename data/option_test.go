package data

import (
	"errors"
	"testing"
)

func TestOption_GetStringRep(t *testing.T) {
	option := Option{
		Label:   "Test",
		Value:   0,
		NextKey: "next",
		Next:    nil,
	}

	got := option.GetStringRep()
	expect := "0. Test"

	if got != expect {
		t.Errorf("Option_GetStringRep(): got %s; expect %s", got, expect)
	}
}

func TestOption_Validate(t *testing.T) {
	option := Option{
		Label:   "Test",
		Value:   0,
		NextKey: "next",
		Next:    nil,
	}

	got := option.Validate()
	expect := errors.New("next is not set for option Test with value 0")
	if got.Error() != expect.Error() {
		t.Errorf("Option_Validate(): got %s; expect %s", got, expect)
	}
}
