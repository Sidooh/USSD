package main

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
	"USSD.sidooh/state"
	"USSD.sidooh/utils"
	"fmt"
	"strconv"
	"time"
)

var screens = map[string]*data.Screen{}

func Process(code, phone, session, input string) *state.State {
	stateData := state.RetrieveState(code, phone, session)

	// User is starting
	if stateData.ScreenPath.Key == "" {
		logger.UssdLog.Println("START ========================", session)
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

func processAndRespond(code, phone, session, input string) string {
	start := time.Now()
	response := Process(code, phone, session, input)

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
		logger.UssdLog.Error(err)
	}
	logger.UssdLog.Printf("Validated %v screens successfully", len(loadedScreens))

	screens = loadedScreens
}

func main() {
	logger.Init()

	LoadScreens()

	paths := map[string][]string{
		// 1. ########## ABOUT
		// ... > About
		//"about": {"", "1"}, // --- valid
		//
		//############## ABOUT END

		// 2. ########## AIRTIME
		// ... > Airtime > self > amount > mpesa > final
		//"airtime_self_20_mpesa_accept": {"", "2", "1", "20", "1", "1"}, // --- valid
		//
		// ... > Airtime > self > amount > other mpesa > final
		//"airtime_self_20_other-mpesa_254714611696_accept": {"", "2", "1", "20", "1", "3", "254714611696", "1"}, // --- valid
		//
		// ... > Airtime > self > amount > voucher > final
		//"airtime_self_20_voucher_pin_accept": {"", "2", "1", "31", "2", "1234", "1"}, // --- valid
		//
		// ... > Airtime > other > new phone > amount > payment > final
		//"airtime_other_new-phone_20_mpesa_accept": {"", "2", "2", "780611696", "20", "1", "1"}, // --- valid
		//
		// ... > Airtime > other > phone > amount > payment > final
		//"airtime_other_phone_20_mpesa_accept": {"", "2", "2", "1", "20", "1", "1"}, // --- valid
		//
		//	... > Extra paths
		//"airtime_self_20_mpesa_cancel": {"", "2", "1", "20", "1", "2"}, // --- valid
		//"airtime_self_20_voucher_invalid-pin_blocked": {"", "2", "1", "20", "2", "123123", "1231", "1232", "7667", "3245"}, // --- valid
		//"airtime_other_existing-new-phone_20_mpesa_accept": {"", "2", "2", "9", "254780611696", "20", "1", "1"}, // --- valid
		//"airtime_other_existing_20_mpesa_other_254714611696_accept": {"", "2", "2", "1", "20", "1", "3", "254110039317", "1"}, // --- valid
		//
		//"airtime_self_20_voucher_back": {"", "2", "1", "20", "2", "1234", "3", "0"}, // --- valid
		//
		// ############## AIRTIME END

		// 3.1 ########## UTILITY
		// ... > Pay > Utility > provider > select account > amount > payment > final
		//"pay_utility_kplc_existing-acc_200_mpesa_accept": {"", "3", "1", "2", "1", "200", "1", "1"}, // --- valid
		//
		// ... > Pay > Utility > provider > no account > amount > payment > final
		//"pay_utility_dstv_new-acc_200_mpesa_accept": {"", "3", "1", "4", "1234567", "200", "1", "1"}, // --- valid
		//
		// ... > Pay > Utility > provider > existing but new account > amount > payment > final
		//"pay_utility_kplc_new-acc_200_mpesa_accept": {"", "3", "1", "2", "9", "1234567", "200", "1", "1"},
		//
		// ############## UTILITY END

		// 3.2 ########## VOUCHER
		// ... > Pay > Voucher > self > amount > payment > final
		//"voucher_self_100_mpesa_accept": {"", "3", "2", "1", "100", "1", "1"}, // --- valid
		//
		// ... > Pay > Voucher > other > account > amount > mpesa > final
		//"voucher_other_phone_100_mpesa_accept": {"", "3", "2", "2", "110039317", "100", "1", "1"}, // --- valid
		//
		// ... > Pay > Voucher > other > account > amount > voucher > final
		//"voucher_other_phone_100_voucher_accept": {"", "3", "2", "2", "110039317", "100", "2", "1234", "1"}, // --- valid
		//
		// ############## VOUCHER END
	}
	x := time.Now()
	for path, inputs := range paths {
		for _, input := range inputs {
			//254110039317
			fmt.Println(processAndRespond("*384*99#", "254714611696", "254714611696"+path, input))
			//time.Sleep(300 * time.Millisecond)
			//fmt.Println(processAndRespond("*384*99#", "254110039317", "254110039317"+path, input))
			//time.Sleep(200 * time.Millisecond)

		}
	}

	logger.UssdLog.Println(time.Since(x))
}
