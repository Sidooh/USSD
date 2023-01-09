package handlers

import (
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"encoding/json"
	"net/http"
)

func LogSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessions, err := datastore.FetchSessionLogs(300)
		if err != nil {
			return
		}

		marshal, err := json.Marshal(sessions)
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

		(w).Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(marshal); err != nil {
			logger.ServiceLog.Error(err)
		}
	})
}
