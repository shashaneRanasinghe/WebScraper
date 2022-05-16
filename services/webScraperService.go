package services

import (
	HTMLVersion "github.com/lestoni/html-version"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

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

	//links := w.GetLinkCount(pageContent, URL) //36 6 11 1m3.30sec
	links := w.GetLinkCount(pageContent, URL) //36 6 11 14sec

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

	log.Debug(response)
	return response
}

//TODO should these functions be methods of the Webscraper struct or should it be independent from the struct
func (w *WebScraper) FindElementCount(pageContent string, elementList []string) map[string]int {

	elementCount := make(map[string]int)

	for _, elem := range elementList {

		re := regexp.MustCompile("</" + elem + ">")
		tags := re.FindAllString(pageContent, -1)
		if tags == nil {
			elementCount[elem] = 0
		} else {
			for range tags {
				elementCount[elem] = elementCount[elem] + 1
			}
		}
	}
	return elementCount
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

func (w *WebScraper) GetLinkCount(pageContent string, currentURL string) map[string]int {

	linkCountMap := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	currentLink, err := url.Parse(currentURL)
	if err != nil {
		log.Error(err)
	}

	parts := strings.Split(currentLink.Hostname(), ".")
	currentDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	links := w.SearchElements(pageContent, "<a.*href=\"(.*?)\"")

	for _, link := range links {
		linkURL, err := url.Parse(link)
		if err != nil {
			log.Error(err)
		}

		parts = strings.Split(linkURL.Hostname(), ".")
		if len(parts) < 2 {
			continue
		}
		linkDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]

		wg.Add(1)

		go func(link string, linkDomain string, linkCountMap map[string]int) {

			resp, err := http.Head(link)

			mu.Lock()
			if err != nil || resp.StatusCode != 200 {
				linkCountMap["inaccessibleURL"] = linkCountMap["inaccessibleURL"] + 1
			} else if linkDomain != currentDomain {
				linkCountMap["externalURL"] = linkCountMap["externalURL"] + 1
			} else {
				linkCountMap["internalURL"] = linkCountMap["internalURL"] + 1
			}
			mu.Unlock()
			wg.Done()
		}(link, linkDomain, linkCountMap)
	}
	wg.Wait()

	return linkCountMap
}
