package main

import (
	"USSD/data"
	"USSD/state"
	"fmt"
	"strconv"
)

var screens = map[string]*data.Screen{}

func process(session string, input string) string {
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
		return stateData.ScreenPath.GetStringRep()
	}

	// Check for global Navigation
	if input == "0" || input == "00" {
		stateData.NavigateBackOrHome(screens, input)
		stateData.SaveState()

		fmt.Println("==========================")
		fmt.Println("Return BACK/HOME response.")
		return stateData.ScreenPath.GetStringRep()
	}

	if stateData.ScreenPath.Type == data.OPEN {
		if stateData.ScreenPath.ValidateInput(input) {
			stateData.ProcessOpenInput(input)

			stateData.MoveNext(screens, stateData.ScreenPath.Screen.NextKey)
		}
	} else {
		if v, e := strconv.Atoi(input); e == nil {
			if o, ok := stateData.ScreenPath.Options[v]; ok {
				stateData.ProcessOptionInput(v)

				stateData.MoveNext(screens, o.NextKey)
			}
		}
	}

	if stateData.ScreenPath.Type == data.GENESIS {
		stateData.Status = data.GENESIS
	}
	if stateData.ScreenPath.Type == data.END {
		stateData.Status = data.END
	}

	stateData.SaveState()

	fmt.Println("==========================")
	fmt.Println("Return response.")
	return stateData.ScreenPath.GetStringRep()
}

func main() {
	//screens := screens{}

	loadedScreens, err := data.LoadData()
	if err != nil {
		fmt.Println(err)
	}

	screens = loadedScreens

	fmt.Println(process("a", ""))

	fmt.Println(process("a", "2"))
	//
	//fmt.Println(process("a", "1"))
	//
	//fmt.Println(process("a", "00"))
	//
	//fmt.Println(process("a", "0"))

	//fmt.Println(process("a", "20"))
	//
	//fmt.Println(process("a", "1"))
	//
	//fmt.Println(process("a", "1"))

}
