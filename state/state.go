package state

import (
	"USSD/data"
	"USSD/products"
	"USSD/service"
	"fmt"
	"log"
	"strconv"
)

type State struct {
	Code       string `json:"code"`
	Session    string `json:"session"`
	Phone      string `json:"phone"`
	ProductKey int    `json:"product_key"`
	product    products.ProductI
	ScreenPath data.ScreenPath   `json:"screen_path"`
	Status     string            `json:"status"`
	Vars       map[string]string `json:"vars"`
}

var screens = map[string]*data.Screen{}

func (s *State) Init(sc map[string]*data.Screen) {
	s.Vars = map[string]string{}
	screens = sc
	s.ScreenPath.Screen = *screens[data.MAIN_MENU]

	account, err := service.FetchAccount(s.Phone)
	if err != nil {
		s.Vars["{name}"] = ""
		s.Vars["{voucher_balance}"] = "0"
	} else {
		s.Vars["{name}"] = " " + account.Name
		s.Vars["{voucher_balance}"] = account.Balances[0].Amount
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

func RetrieveState(code, phone, session string) *State {
	stateData := State{
		Code:    code,
		Session: session,
		Status:  data.GENESIS,
		Phone:   phone,
	}
	err := data.UnmarshalFromFile(&stateData, session+data.STATE_FILE)
	if err != nil {
		log.Default().Println(err)
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

func (s *State) unsetState() {
	_ = data.RemoveFile(s.Session + data.STATE_FILE)
}

func (s *State) ProcessOpenInput(m map[string]*data.Screen, input string) {
	fmt.Println("Processing open input: ", input)
	screens = m

	s.ScreenPath.Screen.Next = getScreen(screens, s.ScreenPath.Screen.NextKey)

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(input)

	s.MoveNext(s.ScreenPath.Screen.NextKey)
}

func (s *State) ProcessOptionInput(m map[string]*data.Screen, option *data.Option) {
	fmt.Println("Processing option input: ", option.Value)
	screens = m

	if s.ScreenPath.Type == data.GENESIS {
		s.SetProduct(option.Value)
	}

	s.ScreenPath.Screen.Next = getScreen(screens, option.NextKey)

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(strconv.Itoa(option.Value))

	s.MoveNext(option.NextKey)
}

func (s *State) SetPrevious() {
	if s.ScreenPath.Previous == nil || s.ScreenPath.Previous.Key != s.ScreenPath.Key {
		s.ScreenPath.Previous = &data.ScreenPath{
			Screen: s.ScreenPath.Screen, Previous: s.ScreenPath.Previous,
		}
	}

	// TODO: Fully check and test this: What will happen if I go back and no screen exists?
	s.ensurePathDepth(s.ScreenPath.Previous, 7)
}

func (s *State) ensurePathDepth(previous *data.ScreenPath, i int) {
	fmt.Println(previous.Key, i)
	if previous.Previous != nil && i > 0 {
		s.ensurePathDepth(previous.Previous, i-1)
	} else {
		previous.Previous = nil
	}
}

func (s *State) MoveNext(screenKey string) {
	s.SetPrevious()
	s.ScreenPath.Screen = *getScreen(screens, screenKey)
}

func (s *State) NavigateBackOrHome(screens map[string]*data.Screen, input string) {
	if input == "0" && s.ScreenPath.Previous != nil {
		s.ScreenPath.Screen = s.ScreenPath.Previous.Screen
		s.ScreenPath.Previous = s.ScreenPath.Previous.Previous
	}

	if input == "00" {
		s.ScreenPath.Screen = *getScreen(screens, data.MAIN_MENU)
		s.ScreenPath.Previous = nil
		s.SetProduct(0)
	}
}

func (s *State) GetStringResponse() string {
	response := ""

	if s.ScreenPath.Type == data.END {
		response += "END "
		s.unsetState()
	} else {
		response += "CON "
	}

	response += s.ScreenPath.GetStringRep()

	if s.ScreenPath.Type != data.GENESIS && s.ScreenPath.Type != data.END {
		if s.ScreenPath.Type == data.CLOSED {
			response += "\n"
		}
		response += "0. Back"
		response += "\n"
		response += "00. Home"
	}

	return response
}

func getScreen(screens map[string]*data.Screen, screenKey string) *data.Screen {
	return screens[screenKey]
}
