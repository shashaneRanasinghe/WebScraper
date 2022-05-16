package main

import (
	"github.com/shashaneRanasinghe/WebScraper/router"
	"github.com/shashaneRanasinghe/WebScraper/server"
	"github.com/tryfix/log"
)

func main() {

	r := router.NewRouter()
	closeChannel := server.Serve(r)
	<-closeChannel

	log.Info("Service Stopped")
}
