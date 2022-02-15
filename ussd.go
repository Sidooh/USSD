package main

import (
	"USSD/data"
	"USSD/state"
	"fmt"
	"strconv"
	"time"
)

var screens = map[string]*data.Screen{}

func Process(code, phone, session, input string) *state.State {
	fmt.Println("START ========================")
	stateData := state.RetrieveState(code, phone, session)

	// User is starting
	if stateData.ScreenPath.Key == "" {
		fmt.Println("User starting...")
		stateData.Init(screens)
		stateData.SaveState()

		fmt.Println("- Return GENESIS response.")
		return stateData
	}

	// Check for global Navigation
	if input == "0" || input == "00" {
		stateData.NavigateBackOrHome(screens, input)
		stateData.SaveState()

		fmt.Println("- Return BACK/HOME response.")
		return stateData
	}

	if stateData.ScreenPath.Type == data.OPEN {
		if stateData.ScreenPath.ValidateInput(input) {
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

	fmt.Println("- Return response.")
	return stateData
}

func processAndRespond(code, phone, session, input string) string {
	response := Process(code, phone, session, input)
	return response.GetStringResponse()
}

func LoadScreens() {
	loadedScreens, err := data.LoadData()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Validated %v screens successfully\n", len(loadedScreens))

	screens = loadedScreens
}

func main() {
	LoadScreens()

	paths := map[string][]string{
		//"about": {"", "1"},
		"airtime_self_20_mpesa_accept": {"", "2", "1", "20", "1", "1"},
		//"pay_utility_tokens_existing-acc_100_mpesa_accept": {"", "3", "4", "1", "100", "1", "1"},
	}

	for path, inputs := range paths {
		for _, input := range inputs {
			start := time.Now()
			fmt.Println(processAndRespond("*384*99#", "254714611696", path, input))
			fmt.Println("LATENCY == ", time.Since(start))
			fmt.Println("========================== END")
		}
	}

}
