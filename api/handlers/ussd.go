package handlers

import (
	"USSD.sidooh/pkg/data"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/state"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	SessionId   string `json:"sessionId"`
	ServiceCode string `json:"serviceCode"`
	PhoneNumber string `json:"phoneNumber"`
	NetworkCode string `json:"networkCode"`
	Text        string `json:"text"`
}

var screens = map[string]*data.Screen{}

func LoadScreens() {
	loadedScreens, err := data.LoadData()
	if err != nil {
		logger.UssdLog.Fatal(err)
	}
	logger.UssdLog.Printf("Validated %v screens successfully", len(loadedScreens))

	screens = loadedScreens
}

func decodeData(r *http.Request) *Data {
	content := r.Header.Get("Content-Type")

	if content == "application/json" {
		decoder := json.NewDecoder(r.Body)
		var t Data

		if err := decoder.Decode(&t); err != nil {
			panic(err)
		}

		return &t
	} else if content == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		return &Data{
			SessionId:   r.FormValue("sessionId"),
			ServiceCode: r.FormValue("serviceCode"),
			PhoneNumber: r.FormValue("phoneNumber"),
			NetworkCode: r.FormValue("networkCode"),
			Text:        r.FormValue("text"),
		}
	}

	return nil
}

func saveState(stateData *state.State) {
	if err := stateData.SaveState(); err != nil {
		logger.UssdLog.Error("\nUnable to save state", err)
	}
}

func process(code, phone, session, input string) *state.State {
	stateData := state.RetrieveState(code, phone, session)

	// TODO: Add pagination capability

	// User is starting
	if stateData.ScreenPath.Key == "" {
		logger.UssdLog.Println("\nSTART ========================", session)
		stateData.Init(screens)

		saveState(stateData)

		logger.UssdLog.Println(" - Return GENESIS response.")
		return stateData
	}

	// On server restart, session may still be viable, we need to ensure screens are set
	stateData.EnsureScreensAreSet(screens)

	// Check for global Navigation
	if stateData.ScreenPath.Type != utils.GENESIS && (input == "0" || input == "00") {
		stateData.NavigateBackOrHome(input)

		saveState(stateData)

		logger.UssdLog.Println(" - Return BACK/HOME response.")
		return stateData
	}

	if stateData.ScreenPath.Type == utils.OPEN {
		if stateData.ScreenPath.ValidateInput(input, stateData.Vars) {
			stateData.ProcessOpenInput(input)

			saveState(stateData)
		}

		// Cater for pin screen tries which are updated in the validator
		// Alternatively, check pin in the products themselves
		// TODO: Check which is more performant
		if stateData.ScreenPath.Validations == "PIN" {
			if stateData.ScreenPath.NextKey == utils.PIN_BLOCKED {
				stateData.MoveNext(stateData.ScreenPath.NextKey)
			}
			saveState(stateData)
		}
	} else {
		if v, e := strconv.Atoi(input); e == nil {
			if option, ok := stateData.ScreenPath.Options[v]; ok {
				stateData.ProcessOptionInput(option)

				saveState(stateData)
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
	response := process(code, phone, session, text[len(text)-1])

	if response.Status == utils.END {
		logger.UssdLog.Println(utils.END+" ==========================", session, time.Since(start))
	} else {
		logger.UssdLog.Println(response.Status+" --------------------------", session, time.Since(start))
	}

	return response.GetStringResponse()
}

func ussdHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/ussd" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	decodedData := decodeData(r)

	_, err := fmt.Fprintln(w, processAndRespond(decodedData.NetworkCode, decodedData.PhoneNumber, decodedData.SessionId, decodedData.Text))
	if err != nil {
		return
	}
}

func Recovery() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if _, err := w.Write(jsonBody); err != nil {
					logger.ServiceLog.Error(err)
				}

				panic(err) //TODO: Maybe log this error? Send to sentry?
			}
		}()

		ussdHandler(w, r)
	})
}
