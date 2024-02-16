package main

import (
	"USSD.sidooh/api"
	"USSD.sidooh/utils"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	utils.SetupConfig(".")

	port := viper.GetString("PORT")
	if port == "" {
		port = "8004"
	}

	sentrySampleRate := viper.GetFloat64("SENTRY_TRACES_SAMPLE_RATE")

	err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100% of transactions for performance monitoring.
		// We recommend adjusting this value in production.
		TracesSampleRate: sentrySampleRate,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	apiRouter := api.Setup()

	if err := http.ListenAndServe(":"+port, apiRouter); err != nil {
		log.Fatal(err)
	}
}
