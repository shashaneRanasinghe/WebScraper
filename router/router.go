package router

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shashaneRanasinghe/WebScraper/handlers"
	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"net/http"
)

type Router struct {
}

func NewRouter() interfaces.Router {
	return &Router{}
}

//Route creates the router and contains all the endpoints
func (r *Router) Route() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/scrape", handlers.Scrape).
		Methods("GET")
	router.Handle("/metrics", promhttp.Handler())

	return router
}
