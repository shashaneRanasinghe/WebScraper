package helpers_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/shashaneRanasinghe/WebScraper/helpers"
	"github.com/tryfix/log"
	"io"
	"net/http"
	"testing"
)

func getPageContent() string {

	resp, err := http.Get("https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password")
	if err != nil {
		log.Error(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Error when closing ReadCloser")
		}
	}(resp.Body)

	dataInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	pageContent := string(dataInBytes)
	return pageContent
}

func TestFindElementCount(t *testing.T) {
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

	helper := helpers.NewWebServiceHelper()

	for _, test := range tests {
		actual := helper.FindElementCount(test.pageContent, test.elementList)
		if test.expected["h1"] != actual["h1"] {
			log.Info("Expected : %v, Got : %v ", expected, actual)
			t.Fail()
		}
	}

}

func TestGetLinkCount_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	url := "https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password"
	expected := 2

	type test struct {
		pageContent string
		URL         string
		expected    int
	}

	tests := []test{
		{
			pageContent: getPageContent(),
			URL:         url,
			expected:    expected,
		},
	}

	helper := helpers.NewWebServiceHelper()

	for _, test := range tests {
		actual, err := helper.GetLinkCount(test.pageContent, test.URL)
		if err != nil {
			log.Error(err)
			t.Fail()
		}
		if test.expected != actual["internalURL"] {
			log.Info("Expected : %v, Got : %v ", expected, actual)
			t.Fail()
		}
	}
}

func TestGetLinkCount_ErrorPath_InvalidURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	url := "errorURL"
	expected := errors.New("invalid URL provided")

	type test struct {
		pageContent string
		URL         string
		expected    error
	}

	tests := []test{
		{
			pageContent: getPageContent(),
			URL:         url,
			expected:    expected,
		},
	}

	helper := helpers.NewWebServiceHelper()

	for _, test := range tests {
		_, actualErr := helper.GetLinkCount(test.pageContent, test.URL)
		if test.expected.Error() != actualErr.Error() {
			log.Info("Expected : %v, Got : %v \n", expected, actualErr)
			t.Fail()
		}
	}
}

func TestGetLinkCount_ErrorPath_URLParse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	url := "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode"

	type test struct {
		pageContent string
		URL         string
		expected    bool
	}

	tests := []test{
		{
			pageContent: getPageContent(),
			URL:         url,
			expected:    true,
		},
	}

	helper := helpers.NewWebServiceHelper()

	for _, test := range tests {
		_, ok := helper.GetLinkCount(test.pageContent, test.URL)
		if ok == nil {
			log.Info("Expected : %v, Got : %v \n", true, ok)
			t.Fail()
		}
	}
}

func TestSearchElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regex := "<title.*?>(.*)</title>"
	expected := []string{"W3Schools Tryit Editor"}

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

	helper := helpers.NewWebServiceHelper()

	for _, test := range tests {
		actual := helper.SearchElements(test.pageContent, test.regex)
		if test.expected[0] != actual[0] {
			log.Info("Expected : %v, Got : %v \n", expected[0], actual[0])
			t.Fail()
		}
	}
}

func BenchmarkWebServiceHelper_FindElementCount(b *testing.B) {
	helper := helpers.NewWebServiceHelper()
	pageContent := getPageContent()
	elementList := []string{"h1", "h2"}

	for i := 0; i < b.N; i++ {
		helper.FindElementCount(pageContent, elementList)
	}
}

func BenchmarkWebServiceHelper_GetLinkCount(b *testing.B) {
	helper := helpers.NewWebServiceHelper()
	url := "https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password"
	pageContent := getPageContent()

	for i := 0; i < b.N; i++ {
		helper.GetLinkCount(pageContent, url)
	}
}

func BenchmarkWebServiceHelper_SearchElements(b *testing.B) {
	helper := helpers.NewWebServiceHelper()
	regex := "<title.*?>(.*)</title>"
	pageContent := getPageContent()

	for i := 0; i < b.N; i++ {
		helper.SearchElements(pageContent, regex)
	}
}
