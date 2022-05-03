package data

import (
	"fmt"
	"strconv"
)

type Option struct {
	Label   string  `json:"label"`
	Value   int     `json:"value,string"`
	NextKey string  `json:"next"`
	Next    *Screen `json:"-"`
	Acyclic bool    `json:"acyclic,omitempty"`
}

func (option *Option) setNext(s *Screen) {
	option.Next = s
}

func (option *Option) GetStringRep() string {
	return fmt.Sprintf("%v. %v", option.Value, option.Label)
}

func (option *Option) Validate() error {
	if option.Next == nil {
		return fmt.Errorf("next is not set for option " + option.Label + " with value " + strconv.Itoa(option.Value))
	} else if option.Acyclic {
		return nil
	} else {
		err := option.Next.Validate(true, true)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
