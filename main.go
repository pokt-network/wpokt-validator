package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/pocket"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if len(os.Args) < 2 {
		log.Fatal("Please provide config file as parameter")
	}
	absConfigPath, _ := filepath.Abs(os.Args[1])

	app.InitConfig(absConfigPath)
	app.InitLogger()

	m := pocket.NewMintMonitor()

	m.Start()

	// Gracefully shut down server
	gracefulStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go waitForExitSignals(gracefulStop, done)
	<-done
	log.Info("Server shutting down")
	m.Cancel()
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Debug("Got signal:", sig)
	log.Debug("Sending done signal to main")
	done <- true
}
