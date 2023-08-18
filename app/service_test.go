package app

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestEmptyService(t *testing.T) {
	t.Run("Empty Service", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		service := NewEmptyService(wg)

		assert.NotNil(t, service)

		wg.Add(1)

		service.Start()

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, EmptyServiceName)
		assert.WithinDuration(t, health.LastSyncTime, time.Now(), 1*time.Second)
		assert.WithinDuration(t, health.NextSyncTime, time.Now(), 1*time.Second)
		assert.Equal(t, health.PoktHeight, "")
		assert.Equal(t, health.EthBlockNumber, "")
		assert.Equal(t, health.Healthy, true)

		service.Stop()

		wg.Wait()
	})

}
