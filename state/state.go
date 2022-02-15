package state

import (
	"USSD/data"
	"USSD/products"
	"USSD/service"
	"fmt"
	"strconv"
)

type State struct {
	Session    string `json:"session"`
	ProductKey int    `json:"product_key"`
	product    products.ProductI
	ScreenPath data.ScreenPath   `json:"screen_path"`
	Status     string            `json:"status"`
	Vars       map[string]string `json:"vars"`
}

func (s *State) Init(screens map[string]*data.Screen) {
	s.Vars = map[string]string{}
	s.ScreenPath.Screen = *screens[data.MAIN_MENU]

	account, err := service.FetchAccount(s.Session)
	if err != nil {
		s.Vars["{name}"] = ""
	} else {
		s.Vars["{name}"] = " " + account.Name
	}

	s.Vars["{phone}"] = account.Phone
}

func (s *State) SetProduct(option int) {
	switch option {
	case products.AIRTIME:
		s.product = &products.Airtime{}
		s.ProductKey = products.AIRTIME
	default:
		s.product = &products.Product{}
		s.ProductKey = products.DEFAULT
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

func RetrieveState(session string) *State {
	stateData := State{
		Session: session,
		Status:  data.GENESIS,
	}
	err := data.UnmarshalFromFile(&stateData, session+data.STATE_FILE)
	if err != nil {
		fmt.Println(err)
	}

	stateData.SetProduct(stateData.ProductKey)

	return &stateData
}

func (s *State) SaveState() error {
	if s.ScreenPath.Type == data.END {
		s.Status = data.END
	} else if s.ScreenPath.Type != data.GENESIS {
		s.Status = data.OPEN
	}

	s.ScreenPath.SubstituteVars(s.Vars)

	err := data.WriteFile(s, s.Session+data.STATE_FILE)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *State) ProcessOpenInput(input string) {
	fmt.Println("Processing open input: ", input)

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(input)
}

func (s *State) ProcessOptionInput(input int) {
	fmt.Println("Processing option input: ", input)
	if s.ScreenPath.Type == data.GENESIS {
		s.SetProduct(input)
	}

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(strconv.Itoa(input))
}

func (s *State) SetPrevious() {
	if s.ScreenPath.Previous == nil || s.ScreenPath.Previous.Key != s.ScreenPath.Key {
		s.ScreenPath.Previous = &data.ScreenPath{
			Screen: s.ScreenPath.Screen, Previous: s.ScreenPath.Previous,
		}
	}
}

func (s *State) MoveNext(screens map[string]*data.Screen, screenKey string) {
	s.SetPrevious()
	s.ScreenPath.Screen = getScreen(screens, screenKey)
}

func (s *State) NavigateBackOrHome(screens map[string]*data.Screen, input string) {
	if input == "0" && s.ScreenPath.Previous != nil {
		s.ScreenPath.Screen = s.ScreenPath.Previous.Screen
		s.ScreenPath.Previous = s.ScreenPath.Previous.Previous
	}

	if input == "00" {
		s.ScreenPath.Screen = getScreen(screens, data.MAIN_MENU)
		s.ScreenPath.Previous = nil
		s.SetProduct(0)
	}
}

func (s *State) GetStringResponse() string {
	response := s.ScreenPath.GetStringRep()

	if s.ScreenPath.Type != data.GENESIS {
		response += "\n"
		response += "0. Back"
		response += "\n"
		response += "00. Home"
	}

	return response
}

func getScreen(screens map[string]*data.Screen, screenKey string) data.Screen {
	return *screens[screenKey]
}
