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
	pocket.PoktMonitorName: {
		serviceWithoutHealth: pocket.NewMonitor,
		serviceWithHealth:    pocket.NewMonitorWithLastHealth,
	},
	pocket.PoktSignerName: {
		serviceWithoutHealth: pocket.NewSigner,
		serviceWithHealth:    pocket.NewSignerWithLastHealth,
	},
	pocket.PoktExecutorName: {
		serviceWithoutHealth: pocket.NewExecutor,
		serviceWithHealth:    pocket.NewExecutorWithLastHealth,
	},
	ethereum.WPoktMonitorName: {
		serviceWithoutHealth: ethereum.NewMonitor,
		serviceWithHealth:    ethereum.NewMonitorWithLastHealth,
	},
	ethereum.WPoktSignerName: {
		serviceWithoutHealth: ethereum.NewSigner,
		serviceWithHealth:    ethereum.NewSignerWithLastHealth,
	},
	ethereum.WPoktExecutorName: {
		serviceWithoutHealth: ethereum.NewExecutor,
		serviceWithHealth:    ethereum.NewExecutorWithLastHealth,
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

	healthcheck := app.NewHealthCheck(&wg)

	lastHealth, err := healthcheck.FindLastHealth()
	serviceHealthMap := make(map[string]models.ServiceHealth)
	if err != nil {
		log.Error("[MAIN] Error getting last health: ", err)
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
