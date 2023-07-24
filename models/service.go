package models

import (
	"sync"
	"time"
)

type Service interface {
	Name() string
	Start()
	LastSyncTime() time.Time
	Interval() time.Duration
	Health() ServiceHealth
	Stop()
}

type EmptyService struct {
	wg   *sync.WaitGroup
	name string
}

func (e *EmptyService) Name() string {
	return e.name
}

func (e *EmptyService) Start() {}

func (e *EmptyService) Stop() {
	e.wg.Done()
}

func (e *EmptyService) LastSyncTime() time.Time {
	return time.Now()
}

func (e *EmptyService) Interval() time.Duration {
	return time.Second * 0
}

func (e *EmptyService) Health() ServiceHealth {
	return ServiceHealth{
		Name:           e.Name(),
		LastSyncTime:   e.LastSyncTime(),
		NextSyncTime:   e.LastSyncTime().Add(e.Interval()),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func NewEmptyService(wg *sync.WaitGroup, name string) *EmptyService {
	return &EmptyService{
		wg:   wg,
		name: name,
	}
}
