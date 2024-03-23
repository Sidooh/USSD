package api

import (
	"USSD.sidooh/api/handlers"
	"USSD.sidooh/pkg/cache"
	"USSD.sidooh/pkg/datastore"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust origin as needed
			//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Setup() http.Handler {
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

	router.Handle("/api/v1/settings", handlers.GetSettings())
	router.Handle("/api/v1/settings/{name:[a-zA-Z]+}", handlers.SetSetting()).Methods("POST")

	router.Handle("/api/v1/dashboard/chart", handlers.GetChartData())
	router.Handle("/api/v1/dashboard/recent-sessions", handlers.GetRecentSessions())
	router.Handle("/api/v1/dashboard/summaries", handlers.GetSummaries())

	return corsMiddleware(router)

}
