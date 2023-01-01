package main

import (
	"USSD.sidooh/cache"
	"USSD.sidooh/data"
	"USSD.sidooh/datastore"
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
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

func processAndRespond(code, phone, session, input string) string {
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

func initUssd() {
	fmt.Println("Initializing USSD subsystem")

	logger.Init()
	cache.Init()
	datastore.Init()
	service.Init()

	LoadScreens()
}

//func main() {
//	initUssd()
//
//	paths := map[string][]string{
//		// 1. ########## ABOUT
//		// ... > About
//		"about": {"", "1"}, // --- valid
//		//
//		//############## ABOUT END
//
//		// 2. ########## AIRTIME
//		// ... > Airtime > self > amount > mpesa > final
//		//"airtime_self_20_mpesa_accept": {"", "2", "1", "20", "1", "1"}, // --- valid
//		//
//		// ... > Airtime > self > amount > other mpesa > final
//		//"airtime_self_20_other-mpesa_254714611696_accept": {"", "2", "1", "20", "1", "3", "254714611696", "1"}, // --- valid
//		//
//		// ... > Airtime > self > amount > voucher > final
//		//"airtime_self_20_voucher_pin_accept": {"", "2", "1", "31", "2", "1234", "1"}, // --- valid
//		//
//		// ... > Airtime > other > new phone > amount > payment > final
//		//"airtime_other_new-phone_20_mpesa_accept": {"", "2", "2", "780611696", "20", "1", "1"}, // --- valid
//		//
//		// ... > Airtime > other > phone > amount > payment > final
//		//"airtime_other_phone_20_mpesa_accept": {"", "2", "2", "1", "20", "1", "1"}, // --- valid
//		//
//		//	... > Extra paths
//		//"airtime_self_20_mpesa_cancel": {"", "2", "1", "20", "1", "2"}, // --- valid
//		//"airtime_self_20_voucher_invalid-pin_blocked": {"", "2", "1", "20", "2", "123123", "1231", "1232", "7667", "3245"}, // --- valid
//		//"airtime_other_existing-new-phone_20_mpesa_accept": {"", "2", "2", "9", "254780611696", "20", "1", "1"}, // --- valid
//		//"airtime_other_existing_20_mpesa_other_254714611696_accept": {"", "2", "2", "1", "20", "1", "3", "254110039317", "1"}, // --- valid
//		//
//		//"airtime_self_20_voucher_back": {"", "2", "1", "20", "2", "1234", "3", "0"}, // --- valid
//		//
//		// ############## AIRTIME END
//
//		// 3.1 ########## UTILITY
//		// ... > Pay > Utility > provider > select account > amount > payment > final
//		//"pay_utility_kplc_existing-acc_200_mpesa_accept": {"", "3", "1", "2", "1", "200", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Utility > provider > no account > amount > payment > final
//		//"pay_utility_dstv_new-acc_200_mpesa_accept": {"", "3", "1", "4", "1234567", "200", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Utility > provider > existing but new account > amount > payment > final
//		//"pay_utility_kplc_new-acc_200_mpesa_accept": {"", "3", "1", "2", "9", "1234567", "200", "1", "1"},
//		//
//		// ############## UTILITY END
//
//		// 3.2 ########## VOUCHER
//		// ... > Pay > Voucher > self > amount > payment > final
//		//"voucher_self_100_mpesa_accept": {"", "3", "2", "1", "100", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Voucher > other > account > amount > mpesa > final
//		//"voucher_other_phone_100_mpesa_accept": {"", "3", "2", "2", "110039317", "100", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Voucher > other > account > amount > voucher > final
//		//"voucher_other_phone_100_voucher_accept": {"", "3", "2", "2", "110039317", "100", "2", "1234", "1"}, // --- valid
//		//
//		// ############## VOUCHER END
//
//		// 4 ########## SAVE
//		// ... > Save > Voucher > self > amount > payment > final
//		//"voucher_self_100_mpesa_accept": {"", "3", "2", "1", "100", "1", "1"}, // --- valid
//		//
//		// ... > Save > Voucher > other > account > amount > mpesa > final
//		//"voucher_other_phone_100_mpesa_accept": {"", "3", "2", "2", "110039317", "100", "1", "1"}, // --- valid
//		//
//		// ... > Save > Voucher > other > account > amount > voucher > final
//		//"voucher_other_phone_100_voucher_accept": {"", "3", "2", "2", "110039317", "100", "2", "1234", "1"}, // --- valid
//		//
//		// ############## SAVE END
//
//		// 5 ########## INVITE
//		// ... > Invite > Pin > phone > final
//		//"invite_pin_716611696_end": {"", "5", "1234", "716611696"},
//		//
//		// ... > Invite > Pin > phone [existing invite] > final
//		//"invite_pin_718611696_end": {"", "5", "1234", "718611696"}, // --- valid
//		//
//		// ... > Invite > Pin > phone [existing account] > final
//		//"invite_pin_110039317_end": {"", "5", "1234", "110039317"}, // --- valid
//		//
//		// ############## INVITE END
//
//		// 6 ########## SUBSCRIPTION
//		// ... > Subscription > info > name > confirm > payment > final
//		//"subscription_info_Dr-H_confirm_payment_end": {"", "6", "1", "1", "Dr H", "1", "2", "1234" /*, "1"*/},
//		//
//		// ... > Subscription > renew > payment > final
//		//"subscription_renew_payment_end": {"", "6", "1", "1", "1"},
//		//
//		// ############## SUBSCRIPTION END
//
//		// 7 ########## ACCOUNT
//		// ... > Account > Profile > view
//		//"account_profile_end": {"", "7", "1"},
//		//
//		// ... > Account > Profile > pin > set > name > new > confirm
//		//"account_profile_pin_new": {"", "7", "1", "1", "1", "Dr H", "1000", "1000"},
//		//
//		// ... > Account > Profile > pin > change > use pin > new > confirm
//		//"account_profile_pin_change_use_pin": {"", "7", "1", "1", "2", "1", "1234", "1001", "1001"},
//		//
//		// ... > Account > Profile > pin > change > use sec qns > new > confirm
//		//"account_profile_pin_change_use_qns": {"", "7", "1", "1", "2", "2", "Blue", "1001", "1001"},
//		//
//		// ... > Account > Profile > sec qn > option 1 > choice 1 > ...
//		//"account_profile_security_questions": {"", "7", "1", "1", "3", "1234", "2", "Blue", "1", "Jack", "1", "Dabber"},
//		//
//		// ... > Account > Profile > update > pin > name > end
//		//"account_profile_update_name": {"", "7", "1", "2", "1234", "Jack Dabbs"},
//		//
//		// ############## ACCOUNT END
//	}
//	x := time.Now()
//	for path, inputs := range paths {
//		for _, input := range inputs {
//			//254110039317
//			// TODO: Test with 7, 07, 2547, +2547... determine if mpesa validation will work for different scenarios
//			fmt.Println(processAndRespond("*384*99#", "254714611696", "254714611696"+path, input))
//			//time.Sleep(300 * time.Millisecond)
//			//fmt.Println(processAndRespond("*384*99#", "254110039317", "254110039317"+path, input))
//			//time.Sleep(200 * time.Millisecond)
//
//		}
//	}
//
//	logger.UssdLog.Println(time.Since(x))
//}
