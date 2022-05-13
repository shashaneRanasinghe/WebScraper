package service_test

import (
	"github.com/golang/mock/gomock"
	"github.com/shashaneRanasinghe/WebScraper/services"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"testing"
)

func getPageContent() string {

	resp, err := http.Get("https://www.google.com/")
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	dataInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	pageContent := string(dataInBytes)
	return pageContent
}

func TestFindElementCount_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	elementList := []string{"h1", "h2"}
	expected := make(map[string]int)
	expected["h1"] = 10
	expected["h2"] = 2

	type test struct {
		pageContent string
		elementList []string
		expected    map[string]int
	}

	tests := []test{
		{
			pageContent: getPageContent(),
			elementList: elementList,
			expected:    expected,
		},
	}

	service := services.NewWebScraper()

	for _, test := range tests {
		actual := service.FindElementCount(test.pageContent, test.elementList)
		if test.expected["h1"] != actual["h1"] {
			t.Fail()
		}
	}

}

func TestGetLinkCount_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regex := "<title.*?>(.*)</title>"
	expected := []string{"Google"}

	type test struct {
		pageContent string
		regex       string
		expected    []string
	}

	tests := []test{
		{
			pageContent: getPageContent(),
			regex:       regex,
			expected:    expected,
		},
	}

	service := services.NewWebScraper()

	for _, test := range tests {
		actual := service.SearchElements(test.pageContent, test.regex)
		if len(test.expected) != len(actual) {
			t.Fail()
		}
	}
}

func TestSearchElements_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regex := "<a.*href=\"(.*?)\""
	expected := []string{"https://www.google.com/"}

	type test struct {
		pageContent string
		regex       string
		expected    []string
	}

	tests := []test{
		{
			pageContent: getPageContent(),
			regex:       regex,
			expected:    expected,
		},
	}

	service := services.NewWebScraper()

	for _, test := range tests {
		actual := service.SearchElements(test.pageContent, test.regex)
		if len(test.expected) != len(actual) {
			t.Fail()
		}
	}
}
