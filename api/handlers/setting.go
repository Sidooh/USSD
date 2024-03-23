package handlers

import (
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetSettings() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		settings, err := datastore.FetchSettings()

		marshal, err := json.Marshal(settings)

		if len(settings) == 0 {
			marshal, err = json.Marshal([]interface{}{})
		}

		if err != nil {
			jsonBody, _ := json.Marshal(map[string]string{
				"error": "There was an internal server error",
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write(jsonBody); err != nil {
				logger.ServiceLog.Error(err)
			}

			panic(err)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(marshal); err != nil {
			logger.ServiceLog.Error(err)
		}
	})
}

func SetSetting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust origin as needed

		decoder := json.NewDecoder(r.Body)
		var data datastore.Setting
		err := decoder.Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error decoding request body: %v", err)
			return
		}

		err = datastore.SetSettingByName(mux.Vars(r)["name"], data.Value)
		marshal, err := json.Marshal(map[string]string{
			"status": "success",
		})
		if err != nil {
			jsonBody, _ := json.Marshal(map[string]string{
				"error": "There was an internal server error",
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write(jsonBody); err != nil {
				logger.ServiceLog.Error(err)
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(marshal); err != nil {
			logger.ServiceLog.Error(err)
		}
	})
}
