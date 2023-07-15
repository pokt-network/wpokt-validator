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
		log.Fatal("[MAIN] Missing config file path argument")
	}
	absConfigPath, _ := filepath.Abs(os.Args[1])

	app.InitConfig(absConfigPath)
	app.InitLogger()

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.MongoDB.TimeOutSecs))
	defer cancel()
	app.InitDB(dbCtx)

	pocket.Client.ValidateNetwork()
	ethereum.Client.ValidateNetwork()

	poktMonitor := pocket.NewMonitor()
	poktSigner := pocket.NewSigner()
	poktExecutor := pocket.NewExecutor()

	wpoktMonitor := ethereum.NewMonitor()
	wpoktSigner := ethereum.NewSigner()
	wpoktExecutor := ethereum.NewExecutor()

	go poktMonitor.Start()
	go poktSigner.Start()
	go poktExecutor.Start()

	go wpoktMonitor.Start()
	go wpoktSigner.Start()
	go wpoktExecutor.Start()

	// Gracefully shut down server
	gracefulStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go waitForExitSignals(gracefulStop, done)
	<-done

	log.Debug("[MAIN] Stopping server gracefully")

	wpoktExecutor.Stop()
	wpoktSigner.Stop()
	wpoktMonitor.Stop()

	poktMonitor.Stop()
	poktSigner.Stop()
	poktExecutor.Stop()

	app.DB.Disconnect()
	log.Debug("[MAIN] Server stopped")
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Debug("[MAIN] Caught signal: ", sig)
	done <- true
}
