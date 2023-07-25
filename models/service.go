package models

import (
	"sync"
	"time"
)

type Service interface {
	Start()
	Health() ServiceHealth
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

func (e *EmptyService) Health() ServiceHealth {
	return ServiceHealth{
		Name:           EmptyServiceName,
		LastSyncTime:   time.Now(),
		NextSyncTime:   time.Now(),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func NewEmptyService(wg *sync.WaitGroup) *EmptyService {
	return &EmptyService{
		wg: wg,
	}
}
