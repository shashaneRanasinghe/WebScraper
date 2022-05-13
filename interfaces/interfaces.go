package interfaces

import "github.com/shashaneRanasinghe/WebScraper/models"

type WebScraper interface {
	Scrape(URL string) models.WebScraperResponse
	FindElementCount(pageContent string, elementList []string) map[string]int
	GetLinkCount(pageContent string, currentURL string) map[string]int
	SearchElements(pageContent string, regex string) []string
}
