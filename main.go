package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

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
	if absEnvPath != "" {
		log.Debug("[MAIN] Env loaded from: ", absEnvPath, " and merged with config from: ", absConfigPath)
	} else {
		log.Debug("[MAIN] Config loaded from: ", absConfigPath)
	}
	log.Info("[MAIN] Config initialized")
	app.InitLogger()
	log.Info("[MAIN] Logger initialized")

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.MongoDB.TimeoutSecs))
	defer cancel()
	app.InitDB(dbCtx)

	pokt.Client.ValidateNetwork()
	eth.Client.ValidateNetwork()

	var wg sync.WaitGroup

	healthcheck := app.NewHealthCheck(&wg)

	lastHealth, err := healthcheck.FindLastHealth()
	serviceHealthMap := make(map[string]models.ServiceHealth)
	if err != nil {
		log.Warn("[MAIN] Error getting last health: ", err)
	} else {
		for _, serviceHealth := range lastHealth.ServiceHealths {
			serviceHealthMap[serviceHealth.Name] = serviceHealth
		}
	}

	services := []models.Service{healthcheck}

	for serviceName, service := range GetServiceFactories() {
		services = append(services, CreateService(&wg, serviceName, serviceHealthMap, service.CreateService, service.CreateServiceWithLastHealth))
	}

	healthcheck.SetServices(services)

	wg.Add(len(services))

	for _, service := range services {
		go service.Start()
	}

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
