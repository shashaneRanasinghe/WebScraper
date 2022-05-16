package handlers

import (
	"encoding/json"
	"errors"
	"github.com/shashaneRanasinghe/WebScraper/helpers"
	"net/http"

	"github.com/tryfix/log"

	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"github.com/shashaneRanasinghe/WebScraper/services"
)

var WebScraper interfaces.WebScraper = services.NewWebScraper()

//The Scrape function validates the input sent by the request
func Scrape(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")

	if url == "" {
		err := errors.New("URL missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	helper := helpers.NewWebServiceHelper()
	response := WebScraper.Scrape(url, helper)

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("%v", err)
		return
	}
	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("%v", err)
		return
	}
}
