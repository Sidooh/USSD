package data

import (
	"fmt"
	"strconv"
)

type Option struct {
	Label      string  `json:"label"`
	Value      int     `json:"value,string"`
	NextKey    string  `json:"next"`
	Next       *Screen `json:"-"`
	Acyclic    bool    `json:"acyclic,omitempty"`
	Rules      string  `json:"rules"`
	NoFullStop bool    `json:"no_full_stop"`
}

func (option *Option) setNext(s *Screen) {
	option.Next = s
}

func (option *Option) GetStringRep() string {
	format := "%v. %v"
	if option.NoFullStop {
		format = "%v %v"
	}

	return fmt.Sprintf(format, option.Value, option.Label)
}

func (option *Option) Validate() error {
	if option.Label == "" {
		return fmt.Errorf("label is not set for option with value " + strconv.Itoa(option.Value))
	}

	if option.Next == nil {
		return fmt.Errorf("next is not set for option '" + option.Label + "' with value " + strconv.Itoa(option.Value))
	} else if option.Acyclic {
		return nil
	} else {
		if err := option.Next.Validate(true, true); err != nil {
			panic(err)
		}
	}

	return nil
}
