package state

import (
	"USSD.sidooh/pkg/data"
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/products"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
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

func (s *State) EnsureScreensAreSet(sc map[string]*data.Screen) {
	if len(screens) == 0 {
		screens = sc
	}
}

func (s *State) Init(sc map[string]*data.Screen) {
	s.Vars = map[string]string{}

	// TODO: Test efficiency/ time comp of this
	//tempScreens, err := json.Marshal(sc)
	//err = json.Unmarshal(tempScreens, &screens)
	//if err != nil {
	//	screens = sc
	//}

	screens = sc
	//s.ScreenPath.Screen = *screens[utils.MAIN_MENU]

	//make copy of main screen to prevent modifications on master screens
	var mainScreen = data.Screen{}

	byteScreen, err := json.Marshal(screens[utils.MAIN_MENU])
	err = json.Unmarshal(byteScreen, &mainScreen)
	if err != nil {
		mainScreen = *screens[utils.MAIN_MENU]
	}
	s.ScreenPath.Screen = mainScreen

	s.Vars["{name}"] = ""
	s.Vars["{voucher_balance}"] = "0"
	s.Vars["{float_balance}"] = "0"
	s.Vars["{customer_support_email}"] = "customersupport@sidooh.co.ke"

	// Can we use go defer/concurrency to fetch other details like voucher balances?
	account, err := service.FetchAccount(s.Phone)
	//Possible Use-cases
	//1. Error is thrown -> phone "", name "", balances 0
	//2. Account has no user -> phone, name "", balances,
	//3. Account and User -> phone, name, balances
	if err != nil {
		logger.UssdLog.Error("FetchAccount: ", err)
		s.Vars["{phone}"] = s.Phone

		_, err := service.FetchInvite(s.Phone)
		if err != nil {
			logger.UssdLog.Error("FetchInvite: ", err)

			s.ScreenPath.Screen = *screens[utils.INVITE_CODE]
		}
	} else {
		if account.Id != 0 {
			s.Vars["{account_id}"] = strconv.Itoa(account.Id)
			s.Vars["{phone}"] = account.Phone
		}

		merchantBetaAccounts := strings.Split(viper.GetString("MERCHANT_BETA_ACCOUNTS"), ",")

		if !slices.Contains(merchantBetaAccounts, s.Vars["{account_id}"]) {
			delete(s.ScreenPath.Options, 0)
		}

		if account.User.Name != "" {
			s.Vars["{name}"] = " " + strings.Split(account.User.Name, " ")[0]
			s.Vars["{full_name}"] = account.User.Name
		}

		if !account.Active {
			s.ScreenPath.Screen = *screens[utils.INACTIVE_ACCOUNT]
		} else {
			if len(account.Vouchers) != 0 {
				s.Vars["{voucher_balance}"] = fmt.Sprintf("%.0f", account.Vouchers[0].Balance)
			}

			if account.Float != nil {
				s.Vars["{float_balance}"] = fmt.Sprintf("%.0f", account.Float.Balance)
			}

			if account.Subscription.Id != 0 {
				s.Vars["{subscription_status}"] = account.Subscription.Status
			}

			if account.HasPin {
				s.Vars["{has_pin}"] = "true"
			}
		}

		if account.Merchant != nil {
			s.Vars["{merchant_id}"] = strconv.Itoa(int(account.Merchant.Id))
			if account.Merchant.BusinessName != "" {
				s.Vars["{merchant_business_name}"] = account.Merchant.BusinessName
				s.Vars["{merchant_code}"] = account.Merchant.Code
				s.Vars["{merchant_float}"] = strconv.Itoa(int(account.Merchant.FloatAccountId))
			}
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
	case products.PAY_MPESA_PAY_BILL:
		s.product = &products.PayBill{}
		s.ProductKey = products.PAY_MPESA_PAY_BILL
	case products.PAY_MPESA_BUY_GOODS:
		s.product = &products.BuyGoods{}
		s.ProductKey = products.PAY_MPESA_BUY_GOODS
	case products.PAY_UTILITY:
		s.product = &products.Utility{}
		s.ProductKey = products.PAY_UTILITY
	case products.PAY_VOUCHER:
		s.product = &products.Voucher{}
		s.ProductKey = products.PAY_VOUCHER
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
	case products.MERCHANT:
		s.product = &products.Merchant{}
		s.ProductKey = products.MERCHANT
	case products.MERCHANT_FLOAT:
		s.product = &products.MerchantFloat{}
		s.ProductKey = products.MERCHANT_FLOAT
	case products.MERCHANT_ACCOUNT:
		s.product = &products.MerchantAccount{}
		s.ProductKey = products.MERCHANT_ACCOUNT
	default:
		s.product = &products.Product{}
		s.ProductKey = products.DEFAULT
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

func (s *State) ProcessOpenInput(input string) {
	logger.UssdLog.Println("Processing open input: ", input)

	s.ScreenPath.Screen.Next = getScreen(s.ScreenPath.Screen.NextKey)

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(input)

	//s.MoveNext("")
	// TODO: Confirm this works across board
	if s.ScreenPath.NextKey != s.ScreenPath.Next.Key {
		s.MoveNext(s.ScreenPath.NextKey)
	} else {
		s.MoveNext("")
	}
}

func (s *State) ProcessOptionInput(option *data.Option) {
	logger.UssdLog.Println("Processing option input: ", option.Value)

	for i, o := range s.ScreenPath.Options {
		if o.Rules == "BETA" {
			betaAccounts := strings.Split(viper.GetString("BETA_ACCOUNTS"), ",")

			if !slices.Contains(betaAccounts, s.Vars["{account_id}"]) {
				s.ScreenPath.Options[i].NextKey = utils.COMING_SOON
			}
		}
	}

	if s.ScreenPath.Type == utils.GENESIS {
		s.setProduct(option.Value)
		if option.Value == 0 {
			s.setProduct(products.MERCHANT)
		}

		// Unauthorized screens for first time users
		keys := []int{5, 6, 7}
		_, ok := s.Vars["{account_id}"]
		if !ok {
			for _, k := range keys {
				s.ScreenPath.Options[k].NextKey = utils.NOT_TRANSACTED
			}
			s.ScreenPath.Options[0].NextKey = utils.MERCHANT_CONSENT
		} else {
			if hasPin, ok := s.Vars["{has_pin}"]; !ok || hasPin != "true" {
				s.ScreenPath.Options[5].NextKey = utils.PIN_NOT_SET
			}
		}

		if _, exists := s.ScreenPath.Options[0]; exists {
			_, ok = s.Vars["{merchant_id}"]
			if !ok {
				s.ScreenPath.Options[0].NextKey = utils.MERCHANT_CONSENT
			} else {
				_, ok = s.Vars["{merchant_business_name}"]
				if !ok {
					s.ScreenPath.Options[0].NextKey = utils.MERCHANT_KYB
				}
			}
		}
	}

	if s.ScreenPath.Key == utils.PAY {
		s.setProduct(products.PAY*10 + option.Value)
	}

	if s.ScreenPath.Key == utils.MERCHANT {
		s.setProduct(products.MERCHANT*10 + option.Value)
	}

	if (s.ScreenPath.Key == utils.PIN_NOT_SET || s.ScreenPath.Key == utils.MERCHANT_PIN_NOT_SET) && option.Value == 1 {
		s.setProduct(products.ACCOUNT)
	}

	if s.ScreenPath.Key == utils.VOUCHER_BALANCE_INSUFFICIENT && option.Value == 1 {
		s.setProduct(products.PAY_VOUCHER)
	}

	s.ScreenPath.Screen.Next = getScreen(option.NextKey)

	s.product.Initialize(s.Vars, &s.ScreenPath.Screen)
	s.product.Process(strconv.Itoa(option.Value))

	if option.NextKey != s.ScreenPath.Next.Key {
		s.MoveNext(option.NextKey)
	} else {
		s.MoveNext("")
	}
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

	if screenKey != "" {
		s.ScreenPath.Screen = *getScreen(screenKey)
	} else {
		s.ScreenPath.Screen = *s.ScreenPath.Next
	}
}

func (s *State) NavigateBackOrHome(input string) {
	if s.ScreenPath.Key == utils.INVITE_CODE || s.ScreenPath.Type == utils.GENESIS {
		return
	}
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
		s.ScreenPath.Screen = *getScreen(utils.MAIN_MENU)
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

	if s.ScreenPath.Type != utils.GENESIS && s.ScreenPath.Key != utils.INVITE_CODE && s.ScreenPath.Type != utils.END {
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

func getScreen(screenKey string) *data.Screen {
	// here we need a value and not reference since it will be translated, and we don't want to change the original
	screen := *screens[screenKey]

	options := make(map[int]*data.Option)
	for i, v := range screen.Options {
		options[i] = &data.Option{
			Label:   v.Label,
			Value:   v.Value,
			NextKey: v.NextKey,
			Next:    v.Next,
			Acyclic: v.Acyclic,
			Rules:   v.Rules,
		}
	}
	screen.Options = options

	return &screen
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
