package app

import (
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
)

type Service interface {
	Start()
	Health() models.ServiceHealth
	Stop()
}

type EmptyService struct {
	wg *sync.WaitGroup
}

func (e *EmptyService) Start() {}

func (e *EmptyService) Stop() {
	e.wg.Done()
}

const EmptyServiceName = "empty"

func (e *EmptyService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           EmptyServiceName,
		LastSyncTime:   time.Now(),
		NextSyncTime:   time.Now(),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func NewEmptyService(wg *sync.WaitGroup) Service {
	return &EmptyService{
		wg: wg,
	}
}
