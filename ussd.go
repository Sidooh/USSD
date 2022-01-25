package main

import (
	"USSD/data"
	"USSD/state"
	"fmt"
	"strconv"
)

var screens = map[string]*data.Screen{}

//var stateData = state.State{}

func getNextScreen(screenKey string) data.Screen {
	return *screens[screenKey]
}

func process(session string, input string) string {
	stateData, err := state.RetrieveState(session)
	if err != nil {
		return ""
	} else {
		stateData.Screen = *screens["main_menu"]
	}

	if stateData.Screen.Type == data.OPEN {
		if stateData.Screen.ValidateInput(input) {
			stateData.ProcessInput(input)
		}
	} else {
		if v, e := strconv.Atoi(input); e == nil {
			if _, ok := stateData.Screen.Options[v]; ok {
				stateData.ProcessInput(input)
			}
		}
	}

	if stateData.Screen.Type == data.GENESIS {
		stateData.Status = data.GENESIS
	}
	if stateData.Screen.Type == data.END {
		stateData.Status = data.END
	}

	stateData.SaveState()

	return stateData.Screen.GetStringRep()
}

func main() {
	//screens := screens{}

	loadedScreens, err := data.LoadData()
	if err != nil {
		fmt.Println(err)
	}

	screens = loadedScreens

	fmt.Println(process("a", "2"))

	fmt.Println(process("a", "1"))

	fmt.Println(process("a", "20"))

	fmt.Println(process("a", "1"))

	fmt.Println(process("a", "1"))

	fmt.Println(process("a", "1"))

}
