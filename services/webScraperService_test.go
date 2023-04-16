package services_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/shashaneRanasinghe/WebScraper/helpers"
	mock_interfaces "github.com/shashaneRanasinghe/WebScraper/mocks"
	"github.com/shashaneRanasinghe/WebScraper/models"
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

func TestScrape_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pageContent := getPageContent()

	url := "https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password"
	expected := models.WebScraperResponse{Error: ""}

	type test struct {
		URL      string
		expected models.WebScraperResponse
	}

	tests := []test{
		{
			URL:      url,
			expected: expected,
		},
	}

	searchElementResponse := []string{"title"}
	elementList := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	elementCount := make(map[string]int)
	elementCount["h1"] = 2
	urlCount := make(map[string]int)
	urlCount["internalURL"] = 2

	mockHelper := mock_interfaces.NewMockHelper(ctrl)
	service := services.NewWebScraper()

	mockHelper.EXPECT().SearchElements(pageContent, "<title.*?>(.*)</title>").Return(searchElementResponse)
	mockHelper.EXPECT().SearchElements(pageContent, "<input.*type=\"?(password)\"?").Return(searchElementResponse)
	mockHelper.EXPECT().FindElementCount(pageContent, elementList).Return(elementCount).AnyTimes()
	mockHelper.EXPECT().GetLinkCount(pageContent, url).Return(urlCount, nil).AnyTimes()

	for _, test := range tests {
		actual := service.Scrape(test.URL, mockHelper)
		if test.expected.Error != actual.Error {
			log.Info("Expected : %v, Got : %v ", expected.Error, actual.Error)
			t.Fail()
		}
	}
}

func TestScrape_ErrorPath_GetLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pageContent := getPageContent()

	url := "https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password"
	getLinkErr := errors.New("invalid URL provided")
	expected := models.WebScraperResponse{Error: "Error while getting the Link Count"}

	type test struct {
		URL      string
		expected models.WebScraperResponse
	}

	tests := []test{
		{
			URL:      url,
			expected: expected,
		},
	}

	searchElementResponse := []string{"title"}
	elementList := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	elementCount := make(map[string]int)
	elementCount["h1"] = 2
	urlCount := make(map[string]int)
	urlCount["internalURL"] = 2

	mockHelper := mock_interfaces.NewMockHelper(ctrl)
	service := services.NewWebScraper()

	mockHelper.EXPECT().SearchElements(pageContent, "<title.*?>(.*)</title>").Return(searchElementResponse)
	mockHelper.EXPECT().FindElementCount(pageContent, elementList).Return(elementCount).AnyTimes()
	mockHelper.EXPECT().GetLinkCount(pageContent, url).Return(urlCount, getLinkErr).AnyTimes()

	for _, test := range tests {
		actual := service.Scrape(test.URL, mockHelper)
		if test.expected.Error != actual.Error {
			log.Info("Expected : %v, Got : %v \n", expected.Error, actual.Error)
			t.Fail()
		}
	}
}

func TestScrape_ErrorURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	url := "https://www.w3schools.c"
	expected := models.WebScraperResponse{Error: "Could not get the data from the requested URL, please check the URL"}

	type test struct {
		URL      string
		expected models.WebScraperResponse
	}

	tests := []test{
		{
			URL:      url,
			expected: expected,
		},
	}

	mockHelper := mock_interfaces.NewMockHelper(ctrl)
	service := services.NewWebScraper()

	for _, test := range tests {
		actual := service.Scrape(test.URL, mockHelper)
		if test.expected.Error != actual.Error {
			log.Info("Expected : %v, Got : %v ", expected.Error, actual.Error)
			t.Fail()
		}
	}
}

func TestScrape_NoLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pageContent := getPageContent()

	url := "https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password"
	expected := models.WebScraperResponse{Error: ""}

	type test struct {
		URL      string
		expected models.WebScraperResponse
	}

	tests := []test{
		{
			URL:      url,
			expected: expected,
		},
	}

	searchElementResponse := []string{"title"}
	searchLoginResponse := []string{}
	elementList := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	elementCount := make(map[string]int)
	elementCount["h1"] = 2
	urlCount := make(map[string]int)
	urlCount["internalURL"] = 2

	mockHelper := mock_interfaces.NewMockHelper(ctrl)
	service := services.NewWebScraper()

	mockHelper.EXPECT().SearchElements(pageContent, "<title.*?>(.*)</title>").Return(searchElementResponse).AnyTimes()
	mockHelper.EXPECT().FindElementCount(pageContent, elementList).Return(elementCount).AnyTimes()
	mockHelper.EXPECT().GetLinkCount(pageContent, url).Return(urlCount, nil).AnyTimes()
	mockHelper.EXPECT().SearchElements(pageContent, "<input.*type=\"?(password)\"?").Return(searchLoginResponse).AnyTimes()

	for _, test := range tests {
		actual := service.Scrape(test.URL, mockHelper)
		if test.expected.Error != actual.Error {
			log.Info("Expected : %v, Got : %v ", expected.Error, actual.Error)
			t.Fail()
		}
	}
}

func BenchmarkWebScraper_Scrape(b *testing.B) {
	service := services.NewWebScraper()
	helper := helpers.NewWebServiceHelper()
	url := "https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password"

	for i := 0; i < b.N; i++ {
		service.Scrape(url, helper)
	}
}
