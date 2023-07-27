package main

import (
	"sync"

	"github.com/dan13ram/wpokt-validator/eth"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/dan13ram/wpokt-validator/pokt"
)

func CreateService(
	wg *sync.WaitGroup,
	serviceName string,
	serviceHealthMap map[string]models.ServiceHealth,
	createService func(*sync.WaitGroup) models.Service,
	createServiceWithLastHealth func(*sync.WaitGroup, models.ServiceHealth) models.Service,
) models.Service {
	serviceHealth, ok := serviceHealthMap[serviceName]
	if ok {
		return createServiceWithLastHealth(wg, serviceHealth)
	} else {
		return createService(wg)
	}
}

type ServiceFactory struct {
	CreateService               func(*sync.WaitGroup) models.Service
	CreateServiceWithLastHealth func(*sync.WaitGroup, models.ServiceHealth) models.Service
}

func GetServiceFactories() map[string]ServiceFactory {
	services := map[string]ServiceFactory{
		pokt.MintMonitorName: {
			CreateService:               pokt.NewMonitor,
			CreateServiceWithLastHealth: pokt.NewMonitorWithLastHealth,
		},
		pokt.BurnSignerName: {
			CreateService:               pokt.NewSigner,
			CreateServiceWithLastHealth: pokt.NewSignerWithLastHealth,
		},
		pokt.BurnExecutorName: {
			CreateService:               pokt.NewExecutor,
			CreateServiceWithLastHealth: pokt.NewExecutorWithLastHealth,
		},
		eth.BurnMonitorName: {
			CreateService:               eth.NewMonitor,
			CreateServiceWithLastHealth: eth.NewMonitorWithLastHealth,
		},
		eth.MintSignerName: {
			CreateService:               eth.NewSigner,
			CreateServiceWithLastHealth: eth.NewSignerWithLastHealth,
		},
		eth.MintExecutorName: {
			CreateService:               eth.NewExecutor,
			CreateServiceWithLastHealth: eth.NewExecutorWithLastHealth,
		},
	}

	return services
}
