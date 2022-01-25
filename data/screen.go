package data

import (
	"fmt"
	"strconv"
)

type Screen struct {
	Key     string          `json:"key"`
	Title   string          `json:"title"`
	Type    string          `json:"type"`
	Options map[int]*Option `json:"options"`
	NextKey string          `json:"next"`
	Next    *Screen         `json:"-"`
}

// types of screens
const (
	GENESIS = "GENESIS"
	OPEN    = "OPEN"
	CLOSED  = "CLOSED"
	END     = "END"
)

var nextExceptionScreens = map[string]bool{
	"about": true, "cancel": true, "coming_soon": true,
	"refer_end": true,
}

func (screen *Screen) setNext(s *Screen) {
	screen.Next = s
}

func (screen *Screen) GetStringRep() string {
	optionsString := ""
	for _, v := range screen.Options {
		optionsString += v.GetStringRep() + "\n"
	}
	return fmt.Sprintf("%v\n\n%v", screen.Title, optionsString)
}

func (screen *Screen) Validate() error {
	if screen.Type == "" {
		screen.Type = CLOSED
		//return fmt.Errorf("type should be set for " + screen.Key)
	}

	if screen.Type == GENESIS || screen.Type == CLOSED {
		// Validate that Next is not set
		if screen.Next != nil {
			return fmt.Errorf("next should not be set for " + screen.Key + " of type " + screen.Type)
		}

		// Validate that options exist
		if len(screen.Options) == 0 {
			return fmt.Errorf("screen options are not set for " + screen.Key + " of type " + screen.Type)
		}

		// Validate that options are valid
		existingOptions := map[string]struct{}{}
		for _, option := range screen.Options {
			// Validate option
			err := option.Validate()
			if err != nil {
				panic(err)
			}

			//	Check if option already exists in list
			_, ok := existingOptions[option.Label]
			if ok {
				return fmt.Errorf("screen options contains duplicates of " + option.Label + " with value " + strconv.Itoa(option.Value))
			} else {
				existingOptions[option.Label] = struct{}{}
			}
		}
	}

	if screen.Type == OPEN {
		// Validate that Next must be set
		// exceptions about,
		if screen.Next != nil {
			err := screen.Next.Validate()
			if err != nil {
				panic(err)
			}
		} else {

			if _, ok := nextExceptionScreens[screen.Key]; ok != true {
				return fmt.Errorf("next is not set for " + screen.Key + " of type " + screen.Type)
			}
		}

		//	 Validate has no options
		if screen.Options != nil {
			return fmt.Errorf("screen options should not be set for " + screen.Key + " of type " + screen.Type)
		}
	}

	if screen.Type == END {
		// Validate that Next must not set
		if screen.Next != nil {
			return fmt.Errorf("next should not be set for " + screen.Key + " of type " + screen.Type)
		}

		//	Validate has no options
		if screen.Options != nil {
			return fmt.Errorf("screen options should not be set for " + screen.Key + " of type " + screen.Type)
		}
	}

	return nil
}

func (screen Screen) ValidateInput(input string) bool {
	//TODO: Add validations
	return true
}
