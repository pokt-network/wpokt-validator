package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/dan13ram/wpokt-validator/pokt"
	log "github.com/sirupsen/logrus"
)

type ServiceFactory = func(*sync.WaitGroup, models.ServiceHealth) app.Service

var ServiceFactoryMap map[string]ServiceFactory = map[string]ServiceFactory{
	pokt.MintMonitorName:  pokt.NewMintMonitor,
	pokt.BurnSignerName:   pokt.NewSigner,
	pokt.BurnExecutorName: pokt.NewExecutor,
	eth.BurnMonitorName:   eth.NewBurnMonitor,
	eth.MintSignerName:    eth.NewMintSigner,
	eth.MintExecutorName:  eth.NewMintExecutor,
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var configPath string
	var envPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&envPath, "env", "", "path to env file")
	flag.Parse()

	var absConfigPath string = ""
	var err error
	if configPath != "" {
		absConfigPath, err = filepath.Abs(configPath)
		if err != nil {
			log.Fatal("[MAIN] Error getting absolute path for config file: ", err)
		}
	}

	var absEnvPath string = ""
	if envPath != "" {
		absEnvPath, err = filepath.Abs(envPath)
		if err != nil {
			log.Fatal("[MAIN] Error getting absolute path for env file: ", err)
		}
	}

	app.InitConfig(absConfigPath, absEnvPath)
	app.InitLogger()
	app.InitDB()

	pokt.ValidateNetwork()
	eth.ValidateNetwork()

	healthcheck := app.NewHealthCheck()

	serviceHealthMap := make(map[string]models.ServiceHealth)

	if app.Config.HealthCheck.ReadLastHealth {
		if lastHealth, err := healthcheck.FindLastHealth(); err == nil {
			for _, serviceHealth := range lastHealth.ServiceHealths {
				serviceHealthMap[serviceHealth.Name] = serviceHealth
			}
		}
	}

	services := []app.Service{}
	var wg sync.WaitGroup

	for serviceName, NewService := range ServiceFactoryMap {
		health := models.ServiceHealth{}
		if lastHealth, ok := serviceHealthMap[serviceName]; ok {
			health = lastHealth
		}
		services = append(services, NewService(&wg, health))
	}

	services = append(services, app.NewHealthService(healthcheck, &wg))

	healthcheck.SetServices(services)

	wg.Add(len(services))

	for _, service := range services {
		go service.Start()
	}

	log.Info("[MAIN] Server started")

	gracefulStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go waitForExitSignals(gracefulStop, done)
	<-done

	log.Debug("[MAIN] Stopping server gracefully")

	for _, service := range services {
		service.Stop()
	}

	wg.Wait()

	app.DB.Disconnect()
	log.Info("[MAIN] Server stopped")
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Debug("[MAIN] Caught signal: ", sig)
	done <- true
}
