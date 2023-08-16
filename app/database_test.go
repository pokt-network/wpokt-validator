package app

import (
	"fmt"
	"testing"
)

func TestInitDB(t *testing.T) {
	Config.MongoDB.TimeoutMillis = int64(1234)
	fmt.Println("TestInitDB")
	fmt.Println("TimeouMillis: ", Config.MongoDB.TimeoutMillis)

	// Check if the health status has been updated correctly.
	// health := service.Health()
	// assert.True(t, health.Healthy)
	// assert.Equal(t, "TestService", health.Name)
	// runs, err := strconv.Atoi(health.PoktHeight)
	// assert.NoError(t, err)
	// assert.GreaterOrEqual(t, runs, 5)
	// assert.Equal(t, "456", health.EthBlockNumber)
}
