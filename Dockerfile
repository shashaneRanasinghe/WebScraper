# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /WebScraper

##
## Deploy
##

FROM alpine:3.14.0
COPY --from=build /WebScraper /WebScraper

ENTRYPOINT ["./Webscraper"]