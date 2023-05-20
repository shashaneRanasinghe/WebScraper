package services

import (
	HTMLVersion "github.com/lestoni/html-version"
	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"github.com/shashaneRanasinghe/WebScraper/models"
	"github.com/tryfix/log"
	"io"
	"net/http"
)

type webScraper struct {
}

func NewWebScraper() interfaces.WebScraper {
	return &webScraper{}
}

// The Scrape function get the content of the webpage of the url and
// create the webScraper response
func (w *webScraper) Scrape(URL string, h interfaces.Helper) models.WebScraperResponse {

	resp, err := http.Get(URL)
	if err != nil {
		log.Error(err)
		return models.WebScraperResponse{
			Data:  models.Data{},
			Error: "Could not get the data from the requested URL, please check the URL",
		}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return models.WebScraperResponse{
			Data:  models.Data{},
			Error: "Error while reading the content",
		}
	}
	pageContent := string(dataInBytes)

	htmlVersion, err := HTMLVersion.DetectFromURL(URL)
	if err != nil {
		log.Error(err)
		return models.WebScraperResponse{
			Data:  models.Data{},
			Error: "error while getting the HTML Version",
		}
	}

	regexForTitle := "<title.*?>(.*)</title>"
	titles := h.SearchElements(pageContent, regexForTitle)

	elementList := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	elementCount := h.FindElementCount(pageContent, elementList)

	links, err := h.GetLinkCount(pageContent, URL)
	if err != nil {
		return models.WebScraperResponse{
			Data:  models.Data{},
			Error: "Error while getting the Link Count",
		}
	}

	logins := h.SearchElements(pageContent, "<input.*type=\"?(password)\"?")
	hasLogin := false
	if len(logins) > 0 {
		hasLogin = true
	} else {
		hasLogin = false
	}

	response := models.WebScraperResponse{
		Data: models.Data{
			HTMLVersion: htmlVersion,
			Title:       titles[0],
			Headers: models.Headers{
				H1Count: int16(elementCount["h1"]),
				H2Count: int16(elementCount["h2"]),
				H3Count: int16(elementCount["h3"]),
				H4Count: int16(elementCount["h4"]),
				H5Count: int16(elementCount["h5"]),
				H6Count: int16(elementCount["h6"]),
			},
			Links: models.Links{
				InternalLinks:     int64(links["internalURL"]),
				ExternalLinks:     int64(links["externalURL"]),
				InaccessibleLinks: int64(links["inaccessibleURL"]),
			},
			HasLoginForm: hasLogin,
		},
		Error: "",
	}

	log.Debug(response)
	return response
}
