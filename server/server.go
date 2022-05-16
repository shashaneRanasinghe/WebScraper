package server

import (
	"context"
	"github.com/shashaneRanasinghe/WebScraper/interfaces"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tryfix/log"
)

// The Serve function creates the server
func Serve(r interfaces.Router) chan string {

	server := http.Server{
		Addr:         ":8001",
		Handler:      r.Route(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	closeChannel := make(chan string)

	//This goroutine will make sure that the service is stopped gracefully
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		signal.Notify(sig, syscall.SIGTERM)
		signal.Notify(sig, syscall.SIGQUIT)
		<-sig

		log.Info("service interruption received")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Error("Server shutdown error : %v", err)
		}

		log.Info("HTTP server stopped")
		close(closeChannel)
	}()

	log.Info("server is starting on port " + server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}

	return closeChannel
}
