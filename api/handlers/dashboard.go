package handlers

import (
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"encoding/json"
	"net/http"
)

func GetChartData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessions, err := datastore.ReadTimeSeriesCount(700)
		if err != nil {
			logger.ServiceLog.Error(err)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(marshal); err != nil {
			logger.ServiceLog.Error(err)
		}
	})
}
