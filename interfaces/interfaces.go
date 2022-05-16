package interfaces

import (
	"github.com/shashaneRanasinghe/WebScraper/models"
	"net/http"
)

type WebScraper interface {
	Scrape(URL string, h Helper) models.WebScraperResponse
}

type Router interface {
	Route() http.Handler
}

type Helper interface {
	FindElementCount(pageContent string, elementList []string) map[string]int
	GetLinkCount(pageContent string, currentURL string) (map[string]int, error)
	SearchElements(pageContent string, regex string) []string
}
