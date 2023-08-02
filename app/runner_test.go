package app

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
)

// MockRunner is a mock implementation of the Runner interface for testing purposes.
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
	// Create a mock Runner and set the interval to a short duration for testing.
	mockRunner := &MockRunner{}
	interval := 100 * time.Millisecond
	wg := &sync.WaitGroup{}
	service := NewRunnerService("TestService", mockRunner, wg, interval)
	wg.Add(1)

	// Run the service asynchronously in a goroutine.
	go service.Start()

	// Wait for a short duration to allow the service to run.
	time.Sleep(600 * time.Millisecond)

	// Stop the service.
	service.Stop()

	// Wait for the service to stop.
	wg.Wait()

	// Check if the health status has been updated correctly.
	health := service.Health()
	assert.True(t, health.Healthy)
	assert.Equal(t, "TestService", health.Name)
	runs, err := strconv.Atoi(health.PoktHeight)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, runs, 5)
	assert.Equal(t, "456", health.EthBlockNumber)
}

func TestNewRunnerServiceInvalidParameters(t *testing.T) {
	// Test NewRunnerService with invalid parameters.
	wg := &sync.WaitGroup{}
	invalidService := NewRunnerService("", nil, wg, 0)
	assert.Nil(t, invalidService)
}

func TestRunnerServiceStop(t *testing.T) {
	// Test stopping the service without starting it.
	wg := &sync.WaitGroup{}
	mockRunner := &MockRunner{}
	service := NewRunnerService("TestService", mockRunner, wg, 100*time.Millisecond)
	service.Stop()

	// The service should exit immediately without any error, so there's nothing to assert.
}
