package main

import (
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/server"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type Data struct {
	SessionId   string `json:"sessionId"`
	ServiceCode string `json:"serviceCode"`
	PhoneNumber string `json:"phoneNumber"`
	NetworkCode string `json:"networkCode"`
	Text        string `json:"text"`
}

func decodeData(r *http.Request) *Data {
	content := r.Header.Get("Content-Type")

	if content == "application/json" {
		decoder := json.NewDecoder(r.Body)
		var t Data
		err := decoder.Decode(&t)
		if err != nil {
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

func ussdHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/ussd" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	data := decodeData(r)

	fmt.Fprintln(w, server.ProcessAndRespond(data.NetworkCode, data.PhoneNumber, data.SessionId, data.Text))
}

func Recovery() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)

				panic(err) //TODO: Maybe log this error? Send to sentry?
			}

		}()

		ussdHandler(w, r)
	})
}

func Logs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessions, err := datastore.FetchSessionLogs()
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
			w.Write(jsonBody)

			panic(err)
		}

		(w).Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(marshal)
	})
}

func main() {
	utils.SetupConfig(".")

	port := viper.GetString("PORT")
	if port == "" {
		port = "8004"
	}

	sentrySampleRate := viper.GetFloat64("SENTRY_TRACES_SAMPLE_RATE")

	err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: sentrySampleRate,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	server.InitUssd()

	fmt.Printf("Starting USSD server at port %v\n", port)

	http.Handle("/api/v1/ussd", Recovery())
	http.Handle("/api/v1/sessions/logs", Logs())

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
