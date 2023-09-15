# WebScraper

WebScraper is a simple application which allows you to get the
title, the header count and the number of a links of a
particular web page

### How to run the application

- Clone the repository
- run `go build` to build the binary
- run `.\{binary name}` to run the binary (eg .\WebScraper)

### Testing
- run `go test -v ./...` to run the unit tests
- run `go test -v ./... -coverprofile=c.out` to create the cover profile
- run `go tool cover -html c` to see the 
coverage of the unit tests

## Endpoints
This application has a single endpoint called scrape which
will take an url as a query parameter 

A request to the endpoint would look like this

    curl --location --request GET 'http://localhost:8001/scrape?url=https://www.w3schools.com/tags/tryit.asp?filename=tryhtml5_input_type_password'

The response would look like this

    {
    "data": {
        "html_version": "HTML 5",
        "title": "Tryit Editor v3.7",
        "headers": {
            "h1_Count": 1,
            "h2_Count": 0,
            "h3_Count": 0,
            "h4_Count": 0,
            "h5_Count": 0,
            "h6_Count": 0
        },
        "links": {
            "internal_links": 2,
            "external_links": 0,
            "inaccessible_links": 0
        },
        "has_login_form": true
    },
    "error": ""
    }

## Assumptions

1. A login form can be detected by checking if an input
   password field exists, since typically we would only use
   a password input for a login.

2. Internal links are links that are redirecting to the same
   domain as the provided url

3. External links are links that have a domain which is 
    different from the domain of the provided url

