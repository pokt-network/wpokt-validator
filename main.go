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

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if len(os.Args) < 2 {
		log.Fatal("[MAIN] Missing config file path argument")
	}
	absConfigPath, _ := filepath.Abs(os.Args[1])

	log.Info("[MAIN] Starting server")
	app.InitConfig(absConfigPath)
	log.Info("[MAIN] Config loaded from: ", absConfigPath)
	app.InitLogger()
	log.Info("[MAIN] Logger initialized")

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.MongoDB.TimeOutSecs))
	defer cancel()
	app.InitDB(dbCtx)

	pocket.Client.ValidateNetwork()
	ethereum.Client.ValidateNetwork()

	healthcheck := app.NewHealthCheck()

	poktMonitor := pocket.NewMonitor()
	poktSigner := pocket.NewSigner()
	poktExecutor := pocket.NewExecutor()

	wpoktMonitor := ethereum.NewMonitor()
	wpoktSigner := ethereum.NewSigner()
	wpoktExecutor := ethereum.NewExecutor()

	go healthcheck.Start()

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

	log.Info("[MAIN] Stopping server gracefully")

	wpoktExecutor.Stop()
	wpoktSigner.Stop()
	wpoktMonitor.Stop()

	poktMonitor.Stop()
	poktSigner.Stop()
	poktExecutor.Stop()

	healthcheck.Stop()

	app.DB.Disconnect()
	log.Info("[MAIN] Server stopped")
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Info("[MAIN] Caught signal: ", sig)
	done <- true
}
