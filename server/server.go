package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shashaneRanasinghe/WebScraper/handlers"
	"github.com/tryfix/log"
)

//the RequestHandler function creates the router and handles the requests
func Serve() {
	router := mux.NewRouter()

	server := http.Server{
		Addr:         ":8001",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	router.HandleFunc("/scrape", handlers.Scrape).
		Methods("GET")

	log.Info("server is starting on port "+server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}