package service_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/shashaneRanasinghe/WebScraper/services"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"testing"
)

func getPageContent() string {

	resp, err := http.Get("https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password")
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
	expected["h1"] = 1
	expected["h2"] = 0

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
			fmt.Printf("Expected : %v, Got : %v ", expected, actual)
			t.Fail()
		}
	}

}

func TestGetLinkCount_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regex := "<title.*?>(.*)</title>"
	expected := []string{"Tryit Editor v3.7"}

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
			fmt.Printf("Expected : %v, Got : %v ", expected, actual)
			t.Fail()
		}
	}
}

func TestSearchElements_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regex := "<a.*href=\"(.*?)\""
	expected := 10 //[]string{"https://www.w3schools.com/"}

	type test struct {
		pageContent string
		regex       string
		expected    int //[]string
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
		if test.expected != len(actual) {
			fmt.Printf("Expected : %v, Got : %v ", expected, len(actual))
			t.Fail()
		}
	}
}
