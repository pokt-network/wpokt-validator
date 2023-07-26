package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth"
	ethClient "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/dan13ram/wpokt-validator/pokt"
	poktClient "github.com/dan13ram/wpokt-validator/pokt/client"
	log "github.com/sirupsen/logrus"
)

func createService(
	wg *sync.WaitGroup,
	serviceName string,
	serviceHealthMap map[string]models.ServiceHealth,
	serviceWithoutHealth func(*sync.WaitGroup) models.Service,
	serviceWithHealth func(*sync.WaitGroup, models.ServiceHealth) models.Service,
) models.Service {
	serviceHealth, ok := serviceHealthMap[serviceName]
	if ok {
		return serviceWithHealth(wg, serviceHealth)
	} else {
		return serviceWithoutHealth(wg)
	}
}

var SERVICES = map[string]struct {
	serviceWithoutHealth func(*sync.WaitGroup) models.Service
	serviceWithHealth    func(*sync.WaitGroup, models.ServiceHealth) models.Service
}{
	pokt.MintMonitorName: {
		serviceWithoutHealth: pokt.NewMonitor,
		serviceWithHealth:    pokt.NewMonitorWithLastHealth,
	},
	pokt.BurnSignerName: {
		serviceWithoutHealth: pokt.NewSigner,
		serviceWithHealth:    pokt.NewSignerWithLastHealth,
	},
	pokt.BurnExecutorName: {
		serviceWithoutHealth: pokt.NewExecutor,
		serviceWithHealth:    pokt.NewExecutorWithLastHealth,
	},
	eth.BurnMonitorName: {
		serviceWithoutHealth: eth.NewMonitor,
		serviceWithHealth:    eth.NewMonitorWithLastHealth,
	},
	eth.MintSignerName: {
		serviceWithoutHealth: eth.NewSigner,
		serviceWithHealth:    eth.NewSignerWithLastHealth,
	},
	eth.MintExecutorName: {
		serviceWithoutHealth: eth.NewExecutor,
		serviceWithHealth:    eth.NewExecutorWithLastHealth,
	},
}

func main() {

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if len(os.Args) < 2 {
		log.Fatal("[MAIN] Missing config file path argument")
	}
	absConfigPath, _ := filepath.Abs(os.Args[1])

	absEnvPath := ""
	if len(os.Args) == 3 {
		absEnvPath, _ = filepath.Abs(os.Args[2])
	}

	log.Debug("[MAIN] Starting server")
	app.InitConfig(absConfigPath, absEnvPath)
	if absEnvPath != "" {
		log.Debug("[MAIN] Env loaded from: ", absEnvPath, " and merged with config from: ", absConfigPath)
	} else {
		log.Debug("[MAIN] Config loaded from: ", absConfigPath)
	}
	log.Info("[MAIN] Config initialized")
	app.InitLogger()
	log.Info("[MAIN] Logger initialized")

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.MongoDB.TimeOutSecs))
	defer cancel()
	app.InitDB(dbCtx)

	poktClient.Client.ValidateNetwork()
	ethClient.Client.ValidateNetwork()

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

	for serviceName, service := range SERVICES {
		services = append(services, createService(&wg, serviceName, serviceHealthMap, service.serviceWithoutHealth, service.serviceWithHealth))
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
