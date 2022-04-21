package main

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
	"USSD.sidooh/state"
	"USSD.sidooh/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var screens = map[string]*data.Screen{}

func Process(code, phone, session, input string) *state.State {
	stateData := state.RetrieveState(code, phone, session)

	// User is starting
	if stateData.ScreenPath.Key == "" {
		log.Println(" - User (" + phone + ") starting... - ")
		stateData.Init(screens)
		stateData.SaveState()

		log.Println(" - Return GENESIS response.")
		return stateData
	}

	// Check for global Navigation
	if input == "0" || input == "00" {
		stateData.NavigateBackOrHome(screens, input)
		stateData.SaveState()

		log.Println(" - Return BACK/HOME response.")
		return stateData
	}

	if stateData.ScreenPath.Type == utils.OPEN {
		if stateData.ScreenPath.ValidateInput(input, stateData.Phone) {
			stateData.ProcessOpenInput(screens, input)

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

func processAndRespond(code, phone, session, input string) string {
	log.Println("START ========================", phone, session, input)

	start := time.Now()
	response := Process(code, phone, session, input)

	log.Println("END ==========================", phone, time.Since(start))

	return response.GetStringResponse()
}

func LoadScreens() {
	loadedScreens, err := data.LoadData()
	if err != nil {
		log.Error(err)
	}
	log.Printf("Validated %v screens successfully", len(loadedScreens))

	screens = loadedScreens
}

func main() {
	logger.Init()

	LoadScreens()

	paths := map[string][]string{
		//"about": {"", "1"},
		//
		"airtime_self_20_mpesa_accept": {"", "2", "1", "20", "1", "1"},
		//	//"airtime_self_20_mpesa_cancel": {"", "2", "1", "20", "1", "2"},
		"airtime_self_20_mpesa_other_254714611696_accept": {"", "2", "1", "20", "1", "3", "254715611696", "1"},
		//
		//	//"airtime_self_20_voucher_valid-pin-accept": {"", "2", "1", "20", "2", "1234", "1"},
		//	//"airtime_self_20_voucher_invalid-pin-accept": {"", "2", "1", "20", "2", "123123", "1"},
		//
		//	//"airtime_other_new-acc_20_mpesa_accept": {"", "2", "1", "20", "1", "1"},
		//	//"airtime_other_20_mpesa_other_254714611696_accept": {"", "2", "1", "20", "1", "3", "254715611696", "1"},
		//
		//	//"airtime_self_20_voucher_accept": {"", "2", "1", "20", "1", "1"},
		//	//"airtime_self_20_voucher_cancel": {"", "2", "1", "20", "1", "2"},
		//	//"airtime_self_20_voucher_other_254714611696_accept": {"", "2", "1", "20", "1", "3", "254715611696", "1"},
		//	//"airtime_self_20_voucher_other_254714611696_cancel": {"", "2", "1", "20", "1", "3", "254715611696", "2"},
		//
		//	//"pay_utility_tokens_existing-acc_100_mpesa_accept": {"", "3", "4", "1", "100", "1", "1"},
		//	//"pay_utility_tokens_new-acc_100_mpesa_accept": {"", "3", "4", "1", "100", "1", "1"},
	}
	x := time.Now()
	for path, inputs := range paths {
		for _, input := range inputs {
			fmt.Println(processAndRespond("*384*99#", "254714611696", "254714611696"+path, input))
			//time.Sleep(300 * time.Millisecond)
			fmt.Println(processAndRespond("*384*99#", "254764611696", "254764611696"+path, input))
			//time.Sleep(200 * time.Millisecond)

		}
	}

	log.Println(time.Since(x))
}
