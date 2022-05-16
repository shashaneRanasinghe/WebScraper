package router

import (
	"github.com/gorilla/mux"
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

	return router
}
