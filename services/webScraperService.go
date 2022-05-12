package services

import (
	"io/ioutil"
	"net/http"
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

	defer resp.Body.Close()

	dataInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return models.WebScraperResponse{
			Data:  models.Data{},
			Error: "Error while reading the content",
		}
	}
	pageContent := string(dataInBytes)

	title := findTitle(pageContent)
	elements := findElements(pageContent)

	response := models.WebScraperResponse{
		Data: models.Data{
			Title: title,
			Headers: models.Headers{
				H1Count: int64(elements["h1"]),
				H2Count: int64(elements["h2"]),
				H3Count: int64(elements["h3"]),
				H4Count: int64(elements["h4"]),
				H5Count: int64(elements["h5"]),
				H6Count: int64(elements["h6"]),
			},
		},
		Error: "",
	}

	return response
}

func findTitle(pageContent string) string {

	startIndex := strings.Index(pageContent, "<title>")
	if startIndex == -1 {
		log.Error("No title element found")
		return ""
	}

	startIndex += 7

	endIndex := strings.Index(pageContent, "</title>")
	if endIndex == -1 {
		log.Error("No closing tag for title found.")
		return ""
	}

	pageTitle := []byte(pageContent[startIndex:endIndex])

	return string(pageTitle)
}

func findElements(pageContent string) map[string]int {

	elementCount := make(map[string]int)
	elementList := [6]string{"h1", "h2", "h3", "h4", "h5", "h6"}

	for _, elem := range elementList {

		// re := regexp.MustCompile("<!--(.|\n)*?-->")
		// comments := re.FindAllString(string(body), -1)
		// if comments == nil {
		// 	fmt.Println("No matches.")
		// } else {
		// 	for _, comment := range comments {
		// 		fmt.Println(comment)
		// 	}
		// }





		titleStartIndex := strings.Index(pageContent, "<"+elem+">")
		if titleStartIndex == -1 {
			elementCount[elem] = 0
			continue
		} else {
			elementCount[elem] = elementCount[elem] + 1
		}
	}
	return elementCount
}
