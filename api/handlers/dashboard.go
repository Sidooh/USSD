package handlers

import (
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/utils"
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

		marshal, err := json.Marshal(utils.SuccessResponse(sessions))
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

func GetRecentSessions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessions, err := datastore.FetchSessionLogs(20)
		if err != nil {
			return
		}

		marshal, err := json.Marshal(utils.SuccessResponse(sessions))
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

func GetSummaries() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ussdBalance, err := service.GetUSSDBalance()
		if err != nil {
			return
		}

		sessionsCount, err := datastore.ReadSummaries()
		if err != nil {
			return
		}

		var summaries = struct {
			UssdBalance   float64     `json:"ussd_balance"`
			SessionsCount interface{} `json:"sessions_count"`
		}{
			ussdBalance,
			sessionsCount,
		}

		marshal, err := json.Marshal(utils.SuccessResponse(summaries))
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
