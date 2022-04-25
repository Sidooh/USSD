package state

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
	"USSD.sidooh/products"
	"USSD.sidooh/service"
	"USSD.sidooh/utils"
	"strconv"
	"strings"
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
	s.ScreenPath.Screen = *screens[utils.MAIN_MENU]

	s.Vars["{name}"] = ""
	s.Vars["{voucher_balance}"] = "0"

	// Can we use go defer/concurrency to fetch other details like voucher balances?
	account, err := service.FetchAccount(s.Phone)
	//Possible Use-cases
	//1. Error is thrown -> phone "", name "", balances 0
	//2. Account has no user -> phone, name "", balances,
	//3. Account and User -> phone, name, balances
	if err != nil {
		logger.UssdLog.Error(err)
		s.Vars["{voucher_balance}"] = "0"
		s.Vars["{phone}"] = s.Phone
	} else {
		s.Vars["{account_id}"] = strconv.Itoa(account.Id)
		s.Vars["{phone}"] = account.Phone

		if len(account.Balances) != 0 {
			s.Vars["{voucher_balance}"] = account.Balances[0].Balance
		}

		if account.User.Name != "" {
			s.Vars["{name}"] = " " + strings.Split(account.User.Name, " ")[0]
		}
	}
}

func (s *State) setProduct(option int) {
	switch option {
	case products.AIRTIME:
		s.product = &products.Airtime{}
		s.ProductKey = products.AIRTIME
	case products.PAY:
		s.product = &products.Pay{}
		s.ProductKey = products.PAY
	case products.PAY_UTILITY:
		s.product = &products.Utility{}
		s.ProductKey = products.PAY_UTILITY
	//case products.PAY_VOUCHER:
	//	s.product = &products.Voucher{}
	//	s.ProductKey = products.PAY_VOUCHER
	//case products.PAY_MERCHANT:
	//	s.product = &products.Merhcant{}
	//	s.ProductKey = products.PAY_MERCHANT
	default:
		s.product = &products.Product{}
		s.ProductKey = products.DEFAULT
	}
}

func (s *State) getProduct() int {
	switch s.product {
	case &products.Airtime{}:
		return products.AIRTIME
	case &products.Pay{}:
		return products.PAY
	case &products.Utility{}:
		return products.PAY_UTILITY
	//case &products.Voucher{}:
	//	return products.PAY_VOUCHER
	default:
		return 0
	}
}

func RetrieveState(code, phone, session string) *State {
	stateData := State{
		Code:    code,
		Session: session,
		Status:  utils.GENESIS,
		Phone:   phone,
	}
	err := data.UnmarshalFromFile(&stateData, session+utils.STATE_FILE)
	if err != nil {
		// TODO: Get the actual error and decide whether log is info or error
		logger.UssdLog.Error(err)
	}

	stateData.setProduct(stateData.ProductKey)

	return &stateData
}

func (s *State) SaveState() error {
	if s.ScreenPath.Type == utils.END {
		s.Status = utils.END
	} else if s.ScreenPath.Type != utils.GENESIS {
		s.Status = utils.OPEN
	}

	s.ScreenPath.SubstituteVars(s.Vars)

	err := data.WriteFile(s, s.Session+utils.STATE_FILE)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *State) unsetState() {
	_ = data.RemoveFile(s.Session + utils.STATE_FILE)
}

func (s *State) ProcessOpenInput(m map[string]*data.Screen, input string) {
	logger.UssdLog.Println("Processing open input: ", input)
	screens = m

	s.ScreenPath.Screen.Next = getScreen(screens, s.ScreenPath.Screen.NextKey)

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(input)

	s.MoveNext(s.ScreenPath.Screen.NextKey)
}

func (s *State) ProcessOptionInput(m map[string]*data.Screen, option *data.Option) {
	logger.UssdLog.Println("Processing option input: ", option.Value)
	screens = m

	if s.ScreenPath.Type == utils.GENESIS {
		s.setProduct(option.Value)
	}

	if s.ScreenPath.Key == utils.PAY {
		s.setProduct(products.PAY*10 + option.Value)
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
	i := s.ensurePathDepth(s.ScreenPath.Previous, 7)
	logger.UssdLog.Println(" -- Depth", 7-i+1)
}

func (s *State) ensurePathDepth(previous *data.ScreenPath, i int) int {
	if previous.Previous != nil && i > 0 {
		return s.ensurePathDepth(previous.Previous, i-1)
	} else {
		previous.Previous = nil
	}

	return i
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
		s.ScreenPath.Screen = *getScreen(screens, utils.MAIN_MENU)
		s.ScreenPath.Previous = nil
		s.setProduct(0)
	}
}

func (s *State) GetStringResponse() string {
	response := ""

	if s.ScreenPath.Type == utils.END {
		response += "END "
		s.unsetState()
	} else {
		response += "CON "
	}

	response += s.ScreenPath.GetStringRep()

	if s.ScreenPath.Type != utils.GENESIS && s.ScreenPath.Type != utils.END {
		if s.ScreenPath.Type == utils.CLOSED {
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
