package interfaces

import (
	"github.com/shashaneRanasinghe/WebScraper/models"
	"net/http"
)

type WebScraper interface {
	Scrape(URL string) models.WebScraperResponse
	FindElementCount(pageContent string, elementList []string) map[string]int
	GetLinkCount(pageContent string, currentURL string) map[string]int
	SearchElements(pageContent string, regex string) []string
}

type Router interface {
	Route() http.Handler
}
