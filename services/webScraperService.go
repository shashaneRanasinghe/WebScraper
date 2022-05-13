package services

import (
	"fmt"
	HTMLVersion "github.com/lestoni/html-version"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"github.com/shashaneRanasinghe/WebScraper/models"
	"github.com/tryfix/log"
)

type WebScraper struct {
}

func NewWebScraper() interfaces.WebScraper {
	return &WebScraper{}
}

func (w *WebScraper) Scrape(URL string) models.WebScraperResponse {

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

	dataInBytes, err := ioutil.ReadAll(resp.Body)
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
			Error: "Error while getting the HTML Version",
		}
	}

	regexForTitle := "<title.*?>(.*)</title>"
	titles := w.SearchElements(pageContent, regexForTitle)

	elementList := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	elementCount := w.FindElementCount(pageContent, elementList)

	links := w.GetLinkCount(pageContent, URL)

	logins := w.SearchElements(pageContent, "<input.*type=\"?(password)\"?")
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

	return response
}

func (w *WebScraper) FindElementCount(pageContent string, elementList []string) map[string]int {

	elementCount := make(map[string]int)

	for _, elem := range elementList {

		re := regexp.MustCompile("</" + elem + ">")
		tags := re.FindAllString(pageContent, -1)
		if tags == nil {
			elementCount[elem] = 0
		} else {
			for _ = range tags {
				elementCount[elem] = elementCount[elem] + 1
			}
		}
	}
	return elementCount
}

func (w *WebScraper) GetLinkCount(pageContent string, currentURL string) map[string]int {
	linkCountMap := make(map[string]int)

	currentlink, err := url.Parse(currentURL)
	if err != nil {
		log.Error(err)
	}

	parts := strings.Split(currentlink.Hostname(), ".")
	currentDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	links := w.SearchElements(pageContent, "<a.*href=\"(.*?)\"")
	for _, link := range links {
		fmt.Println(link)
		linkURL, err := url.Parse(link)
		if err != nil {
			log.Error(err)
		}

		parts = strings.Split(linkURL.Hostname(), ".")
		if len(parts) < 2 {
			continue
		}
		linkDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]

		resp, err := http.Head(link)

		if err != nil || resp.StatusCode != 200 {
			linkCountMap["inaccessibleURL"] = linkCountMap["inaccessibleURL"] + 1
		} else if linkDomain != currentDomain {
			linkCountMap["externalURL"] = linkCountMap["externalURL"] + 1
		} else {
			linkCountMap["internalURL"] = linkCountMap["internalURL"] + 1
		}
	}
	return linkCountMap
}

func (w *WebScraper) SearchElements(pageContent string, regex string) []string {

	elements := make([]string, 0)
	re := regexp.MustCompile(regex)
	matches := re.FindAllStringSubmatch(pageContent, -1)

	for i, _ := range matches {
		elements = append(elements, matches[i][1])
	}

	return elements
}
