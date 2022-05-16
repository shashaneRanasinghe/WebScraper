# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine3.14 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /WebScraper

##
## Deploy
##

FROM alpine:3.14.0
WORKDIR /
COPY --from=build /WebScraper /WebScraper
ENTRYPOINT ["/WebScraper"]