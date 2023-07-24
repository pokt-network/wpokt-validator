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
	EthBlockNumber() string
	PoktHeight() string
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

func (e *EmptyService) EthBlockNumber() string {
	return ""
}

func (e *EmptyService) PoktHeight() string {
	return ""
}

func NewEmptyService(wg *sync.WaitGroup, name string) *EmptyService {
	return &EmptyService{
		wg:   wg,
		name: name,
	}
}
