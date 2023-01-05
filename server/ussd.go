package server

import (
	"USSD.sidooh/data"
	"USSD.sidooh/pkg/cache"
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/state"
	"USSD.sidooh/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var screens = map[string]*data.Screen{}

func Process(code, phone, session, input string) *state.State {
	stateData := state.RetrieveState(code, phone, session)

	// TODO: Add pagination capability

	// User is starting
	if stateData.ScreenPath.Key == "" {
		logger.UssdLog.Println("\nSTART ========================", session)
		stateData.Init(screens)
		stateData.SaveState()

		logger.UssdLog.Println(" - Return GENESIS response.")
		return stateData
	}

	// Check for global Navigation
	if input == "0" || input == "00" {
		stateData.NavigateBackOrHome(screens, input)
		stateData.SaveState()

		logger.UssdLog.Println(" - Return BACK/HOME response.")
		return stateData
	}

	if stateData.ScreenPath.Type == utils.OPEN {
		if stateData.ScreenPath.ValidateInput(input, stateData.Vars) {
			stateData.ProcessOpenInput(screens, input)

			stateData.SaveState()
		}

		// Cater for pin screen tries which are updated in the validator
		// Alternatively, check pin in the products themselves
		// TODO: Check which is more performant
		if stateData.ScreenPath.Validations == "PIN" {
			if stateData.ScreenPath.NextKey == utils.PIN_BLOCKED {
				stateData.MoveNext(stateData.ScreenPath.NextKey)
			}
			stateData.SaveState()
		}
	} else {
		if v, e := strconv.Atoi(input); e == nil {
			if option, ok := stateData.ScreenPath.Options[v]; ok {
				stateData.ProcessOptionInput(screens, option)

				stateData.SaveState()
			}
		}
	}

	return stateData
}

func ProcessAndRespond(code, phone, session, input string) string {
	start := time.Now()
	//TODO: Update ussd session in table here

	text := strings.Split(input, "*")
	phone, _ = utils.FormatPhone(phone)
	response := Process(code, phone, session, text[len(text)-1])

	if response.Status == utils.END {
		logger.UssdLog.Println(utils.END+" ==========================", session, time.Since(start))
	} else {
		logger.UssdLog.Println(response.Status+" --------------------------", session, time.Since(start))
	}

	return response.GetStringResponse()
}

func LoadScreens() {
	loadedScreens, err := data.LoadData()
	if err != nil {
		logger.UssdLog.Fatal(err)
	}
	logger.UssdLog.Printf("Validated %v screens successfully", len(loadedScreens))

	screens = loadedScreens
}

func InitUssd() {
	fmt.Println("Initializing USSD subsystem")
	utils.SetupConfig(".")

	logger.Init()
	cache.Init()
	datastore.Init()
	service.Init()

	LoadScreens()
}
