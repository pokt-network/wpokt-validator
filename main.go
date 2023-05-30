package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum"
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

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.MongoDB.TimeOutSecs))
	defer cancel()
	app.InitDB(dbCtx)

	pocket.ValidateNetwork()
	ethereum.ValidateNetwork()

	m := pocket.NewMintMonitor()
	b := ethereum.NewBurnMonitor()

	go m.Start()
	go b.Start()

	// Gracefully shut down server
	gracefulStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go waitForExitSignals(gracefulStop, done)
	<-done

	log.Debug("Gracefully shutting down server...")
	b.Stop()
	m.Stop()
	app.DB.Disconnect()
	log.Debug("Server gracefully stopped")
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Debug("Got signal:", sig)
	done <- true
}
