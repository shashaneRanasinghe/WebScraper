package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shashaneRanasinghe/WebScraper/handlers"
	"github.com/tryfix/log"
)

// The Serve function creates the router and handles the requests
func Serve() {
	router := mux.NewRouter()

	server := http.Server{
		Addr:         ":8001",
		Handler:      router,
		ReadTimeout:  500 * time.Second,
		WriteTimeout: 500 * time.Second,
	}

	router.HandleFunc("/scrape", handlers.Scrape).
		Methods("GET")

	log.Info("server is starting on port " + server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
