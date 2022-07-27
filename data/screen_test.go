package data

import (
	"USSD.sidooh/utils"
	"errors"
	"testing"
)

func TestScreen_GetStringRep(t *testing.T) {
	// Test bare-bones
	screen := Screen{
		Key:   "test",
		Title: "Test",
	}

	got := screen.GetStringRep()
	expect := "Test\n\n"
	if got != expect {
		t.Errorf("Screen_GetStringRep(): got %s; expect %s", got, expect)
	}

	// Test with option
	screen.Options = map[int]*Option{
		0: {
			Label:   "Test",
			Value:   0,
			NextKey: "next",
			Next:    nil,
		},
	}

	got = screen.GetStringRep()
	expect += "0. Test\n"
	if got != expect {
		t.Errorf("Screen_GetStringRep(): got %s; expect %s", got, expect)
	}
}

func TestScreen_Validate(t *testing.T) {
	// Test bare-bones
	screen := Screen{
		Key:   "",
		Title: "",
	}
	//###################################################################################
	/* Test Basic Items*/

	// Check key is set
	got := screen.Validate(false, false)
	expect := errors.New("key should be set for screen")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Key = "test_key"

	// Check title is set
	got = screen.Validate(false, false)
	expect = errors.New("title should be set for screen test_key")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Title = "Test Title"

	// Check type is set
	got = screen.Validate(false, false)
	expect = errors.New("type should be set for screen test_key")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	//###################################################################################

	/* Test CLOSED screen*/

	screen.Type = utils.CLOSED // GENESIS is a type of closed screen
	screen.Next = &Screen{Type: utils.END}

	// Check next is NOT set
	got = screen.Validate(false, false)
	expect = errors.New("next should not be set for screen test_key of type CLOSED")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Next = nil

	// Check options are set
	got = screen.Validate(false, false)
	expect = errors.New("screen options are not set for screen test_key of type CLOSED")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Options = map[int]*Option{
		0: {
			Label:   "Test",
			Value:   0,
			NextKey: "next",
			Next:    nil,
		},
	}

	got = screen.Validate(false, false)

	if got != nil {
		t.Errorf("Screen_Validate(): got %s; expect %v", got, nil)
	}

	//###################################################################################

	/* Test OPEN screen*/

	screen = Screen{
		Key:   "test_open",
		Title: "Test OPEN",
		Type:  utils.OPEN,
	}

	// Check next is set
	got = screen.Validate(false, false)
	expect = errors.New("next is not set for screen test_open of type OPEN")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Next = &Screen{Type: utils.END}

	// Check options are not set
	screen.Options = map[int]*Option{
		0: {
			Label:   "Test",
			Value:   0,
			NextKey: "next",
			Next:    nil,
		},
	}
	got = screen.Validate(false, false)
	expect = errors.New("screen options should not be set for screen test_open of type OPEN")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Options = nil

	got = screen.Validate(false, false)

	if got != nil {
		t.Errorf("Screen_Validate(): got %s; expect %v", got, nil)
	}

	//###################################################################################

	/* Test END screen*/

	screen = Screen{
		Key:   "test_end",
		Title: "Test END",
		Type:  utils.END,
	}

	// Check options are not set
	screen.Options = map[int]*Option{
		0: {
			Label:   "Test",
			Value:   0,
			NextKey: "next",
			Next:    nil,
		},
	}
	got = screen.Validate(false, false)
	expect = errors.New("screen options should not be set for screen test_end of type END")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Options = nil

	// Check next is NOT set
	screen.Next = &Screen{Type: utils.END}
	got = screen.Validate(false, false)
	expect = errors.New("next should not be set for screen test_end of type END")
	if got.Error() != expect.Error() {
		t.Errorf("Screen_Validate(): got %s; expect %s", got, expect)
	}

	screen.Next = nil

	got = screen.Validate(false, false)

	if got != nil {
		t.Errorf("Screen_Validate(): got %s; expect %v", got, nil)
	}
}
