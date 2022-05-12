package models

type WebScraperRequest struct {
	URL string `json:"url"`
}

type WebScraperResponse struct {
	Data Data `json:"data"`
	Error string `json:"error"`
}

type Data struct {
	HTMLVersion   string  `json:"html_version"`
	Title         string  `json:"title"`
	Headers		  Headers `json:"headers"`
	Links         Links   `json:"links"`
	HasLoginForm  bool    `json:"has_login_form"`
}

type Headers struct{
	H1Count   int64  `json:"h1_Count"`
	H2Count   int64  `json:"h2_Count"`
	H3Count   int64  `json:"h3_Count"`
	H4Count   int64  `json:"h4_Count"`
	H5Count   int64  `json:"h5_Count"`
	H6Count   int64  `json:"h6_Count"`
}

type Links struct {
	InternalLinks     int64 `json:"internal_links"`
	ExternalLinks     int64 `json:"external_links"`
	InaccessibleLinks int64 `json:"inaccessible_links"`
}

