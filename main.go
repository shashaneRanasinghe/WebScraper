package main

import (
	"github.com/joho/godotenv"
	"github.com/shashaneRanasinghe/WebScraper/server"
	"github.com/tryfix/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("%v", err)
	}
	server.Serve()
}
