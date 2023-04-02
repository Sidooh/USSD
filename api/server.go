package api

import (
	"USSD.sidooh/api/handlers"
	"USSD.sidooh/pkg/cache"
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"fmt"
	"net/http"
)

func Setup() {
	fmt.Println("==== Starting Server ====")

	logger.Init()
	cache.Init()
	datastore.Init()
	service.Init()

	handlers.LoadScreens()

	http.Handle("/api/v1/ussd", handlers.Recovery())
	http.Handle("/api/v1/sessions/logs", handlers.LogSession())
	http.Handle("/api/v1/dashboard/chart", handlers.GetChartData())
	http.Handle("/api/v1/dashboard/recent-sessions", handlers.GetRecentSessions())
	http.Handle("/api/v1/dashboard/summaries", handlers.GetSummaries())
}
