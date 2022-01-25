package state

import (
	"USSD/data"
	"USSD/products"
)

type State struct {
	Session    string `json:"session"`
	ProductKey int    `json:"product_key"`
	product    products.ProductI
	Screen     data.Screen `json:"screen_path"`
	Status     string      `json:"status"`
}

func (s *State) setProduct(option int) {
	switch option {
	case products.AIRTIME:
		s.product = &products.Airtime{}
	default:
		s.product = &products.Product{}
	}
}

func (s *State) getProduct() int {
	switch s.product {
	case &products.Airtime{}:
		return products.AIRTIME
	default:
		return 0
	}
}

func RetrieveState(session string) (*State, error) {
	stateData := State{
		Session: session,
		Status:  data.GENESIS,
	}
	err := data.UnmarshalFromFile(stateData, session+"_state.json")
	if err != nil {
		panic(err)
	}

	stateData.setProduct(stateData.ProductKey)

	return &stateData, nil
}

func (s *State) SaveState() error {
	err := data.WriteFile(s, s.Session+"_state.json")
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func (s *State) ProcessInput(input string) {

}
