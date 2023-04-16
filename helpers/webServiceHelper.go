package helpers

import (
	"errors"
	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"github.com/tryfix/log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type webServiceHelper struct {
}

func NewWebServiceHelper() interfaces.Helper {
	return &webServiceHelper{}
}

//FindElementCount counts the number of times an element is in the webpage
func (w *webServiceHelper) FindElementCount(pageContent string, elementList []string) map[string]int {

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

//SearchElements finds the values of the given html elements according to the
//given regex
func (w *webServiceHelper) SearchElements(pageContent string, regex string) []string {

	elements := make([]string, 0)
	re := regexp.MustCompile(regex)
	matches := re.FindAllStringSubmatch(pageContent, -1)

	for i := range matches {
		elements = append(elements, matches[i][1])
	}

	return elements
}

//GetLinkCount counts the number of internal, external and inaccessible links in the webpage
func (w *webServiceHelper) GetLinkCount(pageContent string, currentURL string) (map[string]int, error) {

	linkCountMap := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	currentLink, err := url.Parse(currentURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	parts := strings.Split(currentLink.Hostname(), ".")
	if len(parts) < 2 {
		err := errors.New("invalid URL provided")
		log.Error(err)
		return nil, err
	}
	currentDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	links := w.SearchElements(pageContent, "<a.*href=\"(.*?)\"")

	for _, link := range links {
		linkURL, err := url.Parse(link)
		if err != nil {
			log.Error(err)
			continue
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

	return linkCountMap, nil
}
