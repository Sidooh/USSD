package state

import (
	"USSD.sidooh/data"
	"USSD.sidooh/datastore"
	"USSD.sidooh/logger"
	"USSD.sidooh/products"
	"USSD.sidooh/service"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
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

func (s *State) ToSession() *datastore.Session {
	vars, err := json.Marshal(s.Vars)
	if err != nil {
		panic(err)
	}

	return &datastore.Session{
		SessionId:  s.Session,
		Phone:      s.Phone,
		Code:       s.Code,
		Status:     s.Status,
		Product:    s.ProductKey,
		ScreenPath: s.ScreenPath.Encode(),
		Vars:       vars,
	}
}

var screens = map[string]*data.Screen{}

func (s *State) Init(sc map[string]*data.Screen) {
	s.Vars = map[string]string{}
	screens = sc
	s.ScreenPath.Screen = *screens[utils.MAIN_MENU]

	s.Vars["{name}"] = ""
	s.Vars["{voucher_balance}"] = "0"
	s.Vars["{customer_support_email}"] = "customersupport@sidooh.co.ke"

	// Can we use go defer/concurrency to fetch other details like voucher balances?
	account, err := service.FetchAccount(s.Phone)
	//Possible Use-cases
	//1. Error is thrown -> phone "", name "", balances 0
	//2. Account has no user -> phone, name "", balances,
	//3. Account and User -> phone, name, balances
	if err != nil {
		logger.UssdLog.Error("FetchAccount: ", err)
		s.Vars["{voucher_balance}"] = "0"
		s.Vars["{phone}"] = s.Phone

		_, err := service.FetchInvite(s.Phone)
		if err != nil {
			logger.UssdLog.Error("FetchInvite: ", err)

			s.ScreenPath.Screen = *screens[utils.INVITE_CODE]
		}
	} else {
		s.Vars["{account_id}"] = strconv.Itoa(account.Id)
		s.Vars["{phone}"] = account.Phone

		if len(account.Balances) != 0 {
			s.Vars["{voucher_balance}"] = fmt.Sprintf("%.0f", account.Balances[0].Balance)
		}

		if account.User.Name != "" {
			s.Vars["{name}"] = " " + strings.Split(account.User.Name, " ")[0]
			s.Vars["{full_name}"] = account.User.Name
		}

		if account.Subscription.Id != 0 {
			s.Vars["{subscription_status}"] = account.Subscription.Status
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
	case products.PAY_VOUCHER:
		s.product = &products.Voucher{}
		s.ProductKey = products.PAY_VOUCHER
	//case products.PAY_MERCHANT:
	//	s.product = &products.Merchant{}
	//	s.ProductKey = products.PAY_MERCHANT
	case products.SAVE:
		s.product = &products.Save{}
		s.ProductKey = products.SAVE
	case products.INVITE:
		s.product = &products.Invite{}
		s.ProductKey = products.INVITE
	case products.SUBSCRIPTION:
		s.product = &products.Subscription{}
		s.ProductKey = products.SUBSCRIPTION
	case products.ACCOUNT:
		s.product = &products.Account{}
		s.ProductKey = products.ACCOUNT
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
	case &products.Voucher{}:
		return products.PAY_VOUCHER
	case &products.Save{}:
		return products.SAVE
	case &products.Invite{}:
		return products.INVITE
	case &products.Subscription{}:
		return products.SUBSCRIPTION
	case &products.Account{}:
		return products.ACCOUNT
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

	sessionData := new(datastore.Session)
	err := datastore.UnmarshalFromDatabase(session, sessionData)
	if err != nil {
		// TODO: Get the actual error and decide whether log is info or error
		logger.UssdLog.Error(err)
	} else {
		stateData.FromSession(sessionData)
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

	err := datastore.MarshalToDatabase(*s.ToSession())
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *State) unsetState() {
	// Unnecessary since we are using DB and not file store
	//_ = datastore.RemoveFile(s.Session + utils.STATE_FILE)
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

		// Unauthorized screens for first time users
		keys := []int{5, 6, 7}
		_, ok := s.Vars["{account_id}"]
		if !ok {
			for _, k := range keys {
				s.ScreenPath.Options[k].NextKey = utils.NOT_TRANSACTED
			}
		}
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
		if s.ScreenPath.Previous.Validations == utils.PIN {
			s.ScreenPath.Screen = s.ScreenPath.Previous.Previous.Screen
			s.ScreenPath.Previous = s.ScreenPath.Previous.Previous.Previous
		} else {
			s.ScreenPath.Screen = s.ScreenPath.Previous.Screen
			s.ScreenPath.Previous = s.ScreenPath.Previous.Previous
		}
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
		if !s.ScreenPath.Paginated {
			response += "0.Back "
			response += " "
		}
		response += "00.Home"
	}

	return response
}

func getScreen(screens map[string]*data.Screen, screenKey string) *data.Screen {
	// here we need a value and not reference since it will be translated and we don't want to change the original
	return screens[screenKey]
}

func (s *State) FromSession(session *datastore.Session) {
	s.ProductKey = session.Product

	err := json.Unmarshal([]byte(session.ScreenPath), &s.ScreenPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(session.Vars, &s.Vars)
	if err != nil {
		panic(err)
	}
}
