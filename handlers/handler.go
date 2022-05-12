package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tryfix/log"

	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"github.com/shashaneRanasinghe/WebScraper/services"
)

var WebScraper interfaces.WebScraper = services.NewWebScraper()

//the classify method gets an image and returns the labels for
//the given image
func Scrape(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")

	if url == "" {
		err := errors.New("URL missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}
	response := WebScraper.Scrape(url)

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("%v", err)
		_, _ = fmt.Fprint(w, err)
		return
	}
	w.Write(res)

}
