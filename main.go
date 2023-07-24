package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum"
	ethClient "github.com/dan13ram/wpokt-backend/ethereum/client"
	"github.com/dan13ram/wpokt-backend/models"

	"github.com/dan13ram/wpokt-backend/pocket"
	poktClient "github.com/dan13ram/wpokt-backend/pocket/client"
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

	// check 2 argument, if it exists set it as absEnvPath
	absEnvPath := ""
	if len(os.Args) == 3 {
		absEnvPath, _ = filepath.Abs(os.Args[2])
	}

	log.Info("[MAIN] Starting server")
	app.InitConfig(absConfigPath, absEnvPath)
	if absEnvPath != "" {
		log.Info("[MAIN] Env loaded from: ", absEnvPath, " and merged with config from: ", absConfigPath)
	} else {
		log.Info("[MAIN] Config loaded from: ", absConfigPath)
	}
	app.InitLogger()
	log.Info("[MAIN] Logger initialized")

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.MongoDB.TimeOutSecs))
	defer cancel()
	app.InitDB(dbCtx)

	poktClient.Client.ValidateNetwork()
	ethClient.Client.ValidateNetwork()

	var wg sync.WaitGroup

	poktMonitor := pocket.NewMonitor(&wg)
	poktSigner := pocket.NewSigner(&wg)
	poktExecutor := pocket.NewExecutor(&wg)

	wpoktMonitor := ethereum.NewMonitor(&wg)
	wpoktSigner := ethereum.NewSigner(&wg)
	wpoktExecutor := ethereum.NewExecutor(&wg)

	services := []models.Service{poktMonitor, poktSigner, poktExecutor, wpoktMonitor, wpoktSigner, wpoktExecutor}

	healthcheck := app.NewHealthCheck(services, &wg)

	services = append(services, healthcheck)

	wg.Add(len(services))

	for _, service := range services {
		go service.Start()
	}

	// Gracefully shut down server
	gracefulStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go waitForExitSignals(gracefulStop, done)
	<-done

	log.Info("[MAIN] Stopping server gracefully")

	for _, service := range services {
		service.Stop()
	}

	wg.Wait()

	app.DB.Disconnect()
	log.Info("[MAIN] Server stopped")
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Info("[MAIN] Caught signal: ", sig)
	done <- true
}
