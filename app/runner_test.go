package app

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
)

type MockRunner struct {
	runs int
}

func (m *MockRunner) Run() {
	m.runs += 1
}

func (m *MockRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{
		PoktHeight:     strconv.Itoa(m.runs),
		EthBlockNumber: "456",
	}
}

func TestRunnerService(t *testing.T) {
	mockRunner := &MockRunner{}
	interval := 100 * time.Millisecond
	wg := &sync.WaitGroup{}
	service := NewRunnerService("TestService", mockRunner, wg, interval)
	wg.Add(1)

	go service.Start()

	time.Sleep(600 * time.Millisecond)

	service.Stop()

	wg.Wait()

	health := service.Health()
	assert.True(t, health.Healthy)
	assert.Equal(t, "TestService", health.Name)
	runs, err := strconv.Atoi(health.PoktHeight)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, runs, 5)
	assert.Equal(t, "456", health.EthBlockNumber)
}

func TestNewRunnerServiceInvalidParameters(t *testing.T) {
	wg := &sync.WaitGroup{}
	invalidService := NewRunnerService("", nil, wg, 0)
	assert.Nil(t, invalidService)
}

func TestRunnerServiceStop(t *testing.T) {
	wg := &sync.WaitGroup{}
	mockRunner := &MockRunner{}
	service := NewRunnerService("TestService", mockRunner, wg, 100*time.Millisecond)
	service.Stop()
}
