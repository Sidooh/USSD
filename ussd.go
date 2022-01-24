package main

import (
	"USSD/data"
	"USSD/products"
	"fmt"
	"strconv"
)

type state struct {
	session string
	product int
	path    *data.Screen
}

type screens struct {
	all     map[string]*data.Screen
	current *data.Screen
	status  string
	product products.ProductI
}

func (s *screens) moveNext() {
	s.product.Process(s.current)

	s.current = s.current.Next
}

func (s *screens) selectOption(o int) {
	s.current = s.current.Options[o].Next
}

func (s *screens) validate(input string) bool {
	//TODO: Add validations
	return true
}

func (s *screens) setProduct(option int) {
	switch option {
	case products.AIRTIME:
		s.product = &products.Airtime{}
	default:
		s.product = &products.Product{}
	}
}

func (s *screens) getProduct() int {
	switch s.product {
	case &products.Airtime{}:
		return products.AIRTIME
	default:
		return 0
	}
}

func (s *screens) process(input string) string {
	s.retrieveState()

	if s.current.Type == data.OPEN {
		if s.validate(input) {
			s.moveNext()
		}
	} else {
		if v, e := strconv.Atoi(input); e == nil {
			if _, ok := s.current.Options[v]; ok {
				s.selectOption(v)
			}
		}
	}

	if s.current.Type == data.GENESIS {
		s.status = data.GENESIS
	}
	if s.current.Type == data.END {
		s.status = data.END
	}

	s.saveState()

	return s.current.GetStringRep()
}

func (s *screens) retrieveState() error {
	var stateData = state{}
	err := data.UnmarshalFromFile(stateData, "state")
	if err != nil {
		return err
	}

	s.setProduct(stateData.product)

	return nil
}

func (s *screens) saveState() error {
	var stateData = state{
		"", s.getProduct(), s.current,
	}

	err := data.WriteFile(stateData, "state")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	screens := screens{}

	loadedScreens, err := data.LoadData()
	if err != nil {
		fmt.Println(err)
	}

	screens.all = loadedScreens
	screens.current = screens.all["main_menu"]

	fmt.Println(screens.current.GetStringRep())

	screens.process("2")

	fmt.Println(screens.current.GetStringRep())
	fmt.Println(screens.product)

	screens.process("1")

	fmt.Println(screens.current.GetStringRep())

	screens.process("20")

	fmt.Println(screens.current.GetStringRep())

	screens.process("1")

	fmt.Println(screens.current.GetStringRep())

	screens.process("1")

	fmt.Println(screens.current.GetStringRep())

	screens.process("1")

	fmt.Println(screens.current.GetStringRep())

}
