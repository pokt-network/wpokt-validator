package app

import (
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
	log "github.com/sirupsen/logrus"
)

type Runner interface {
	Run()
	Status() models.RunnerStatus
}

type RunnerService struct {
	wg       *sync.WaitGroup
	name     string
	runner   Runner
	interval time.Duration

	stop chan struct{}

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *RunnerService) Start() {
	log.Infof("[%s] Service started", x.name)
	stop := false
	for !stop {
		log.Infof("[%s] Run started", x.name)

		x.runner.Run()

		x.updateHealth(x.runner.Status())

		log.Infof("[%s] Run complete, next run in %s", x.name, x.interval)

		select {
		case <-x.stop:
			log.Infof("[%s] Service stopped", x.name)
			x.wg.Done()
			stop = true
		case <-time.After(x.interval):
		}
	}
}

func (x *RunnerService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *RunnerService) updateHealth(status models.RunnerStatus) {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           x.name,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     status.PoktHeight,
		EthBlockNumber: status.EthBlockNumber,
		Healthy:        true,
	}
}

func (x *RunnerService) Stop() {
	log.Debugf("[%s] Stopping", x.name)
	close(x.stop)
}

func NewRunnerService(
	name string,
	runner Runner,
	wg *sync.WaitGroup,
	interval time.Duration,

) Service {
	if (name == "") || (runner == nil) || (wg == nil) || (interval == 0) {
		log.Debug("[RUNNER] Invalid parameters")
		return nil
	}

	return &RunnerService{
		name:     name,
		runner:   runner,
		wg:       wg,
		interval: interval,
		stop:     make(chan struct{}),
		health: models.ServiceHealth{
			Name: name,
		},
	}
}
