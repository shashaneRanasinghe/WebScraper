package interfaces

import "github.com/shashaneRanasinghe/WebScraper/models"

type WebScraper interface {
	Scrape(URL string) models.WebScraperResponse
}
