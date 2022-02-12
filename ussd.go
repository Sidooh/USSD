package main

import (
	"USSD/data"
	"USSD/state"
	"fmt"
	"strconv"
)

var screens = map[string]*data.Screen{}

func process(session string, input string) *state.State {
	fmt.Println("==========================")
	fmt.Println("Processing...")
	stateData := state.RetrieveState(session)

	// User is starting
	if stateData.ScreenPath.Key == "" {
		fmt.Println("User starting...")
		stateData.Init(screens)
		stateData.SaveState()

		fmt.Println("==========================")
		fmt.Println("Return GENESIS response.")
		return stateData
	}

	// Check for global Navigation
	if input == "0" || input == "00" {
		stateData.NavigateBackOrHome(screens, input)
		stateData.SaveState()

		fmt.Println("==========================")
		fmt.Println("Return BACK/HOME response.")
		return stateData
	}

	if stateData.ScreenPath.Type == data.OPEN {
		if stateData.ScreenPath.ValidateInput(input) {
			stateData.ProcessOpenInput(input)

			stateData.MoveNext(screens, stateData.ScreenPath.Screen.NextKey)

			stateData.SaveState()
		}
	} else {
		if v, e := strconv.Atoi(input); e == nil {
			if o, ok := stateData.ScreenPath.Options[v]; ok {
				stateData.ProcessOptionInput(v)

				stateData.MoveNext(screens, o.NextKey)

				stateData.SaveState()
			}
		}
	}

	fmt.Println("==========================")
	fmt.Println("Return response.")
	return stateData
}

func processAndRespond(session string, input string) string {
	response := process(session, input)
	return response.GetStringResponse()
}

func main() {
	//screens := screens{}

	loadedScreens, err := data.LoadData()
	if err != nil {
		fmt.Println(err)
	}

	screens = loadedScreens

	fmt.Println(processAndRespond("a", ""))

	fmt.Println(processAndRespond("a", "2"))
	//
	fmt.Println(processAndRespond("a", "1"))
	//
	//fmt.Println(processAndRespond("a", "00"))
	//
	//fmt.Println(processAndRespond("a", "0"))

	fmt.Println(processAndRespond("a", "+20"))
	//
	//fmt.Println(processAndRespond("a", "1"))
	//
	//fmt.Println(processAndRespond("a", "1"))

}
