package api

import (
	"USSD.sidooh/api/handlers"
	"USSD.sidooh/pkg/cache"
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"fmt"
	"github.com/gorilla/mux"
)

func Setup() *mux.Router {
	fmt.Println("==== Starting Server ====")

	logger.Init()
	cache.Init()
	datastore.Init()
	service.Init()

	handlers.LoadScreens()

	router := mux.NewRouter()

	router.Handle("/api/v1/ussd", handlers.Recovery())
	router.Handle("/api/v1/sessions", handlers.GetSessions())
	router.Handle("/api/v1/sessions/{id:[0-9]+}", handlers.GetSession())

	router.Handle("/api/v1/dashboard/chart", handlers.GetChartData())
	router.Handle("/api/v1/dashboard/recent-sessions", handlers.GetRecentSessions())
	//router.Handle("/api/v1/dashboard/sessions/{id}", handlers.GetSession())
	router.Handle("/api/v1/dashboard/summaries", handlers.GetSummaries())

	return router

}
