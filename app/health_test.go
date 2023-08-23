package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-validator/mocks"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestHealthCheck() *HealthCheckRunner {
	x := &HealthCheckRunner{
		validatorId: "validatorId",
		hostname:    "hostname",
	}
	return x
}

func TestHealthStatus(t *testing.T) {
	x := NewTestHealthCheck()

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "")
}

func TestFindLastHealth(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockDB := mocks.NewMockDatabase(t)
		DB = mockDB

		x := NewTestHealthCheck()
		filter := bson.M{
			"validator_id": x.validatorId,
			"hostname":     x.hostname,
		}
		var health models.Health
		mockDB.EXPECT().FindOne(models.CollectionHealthChecks, filter, &health).Return(nil)

		_, err := x.FindLastHealth()

		assert.Nil(t, err)
	})

	t.Run("With Error", func(t *testing.T) {
		mockDB := mocks.NewMockDatabase(t)
		DB = mockDB

		x := NewTestHealthCheck()
		filter := bson.M{
			"validator_id": x.validatorId,
			"hostname":     x.hostname,
		}
		var health models.Health
		mockDB.EXPECT().FindOne(models.CollectionHealthChecks, filter, &health).Return(errors.New("error"))

		_, err := x.FindLastHealth()

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "error")
	})

}

type MockService struct {
}

func (e *MockService) Start() {}

func (e *MockService) Stop() {
}

const MockServiceName = "mock"

func (e *MockService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           MockServiceName,
		LastSyncTime:   time.Now(),
		NextSyncTime:   time.Now(),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func NewMockService() Service {
	return &MockService{}
}

func TestServices(t *testing.T) {
	x := NewTestHealthCheck()
	wg := &sync.WaitGroup{}
	x.SetServices([]Service{
		NewEmptyService(wg),
		NewEmptyService(wg),
		NewMockService(),
	})

	assert.Equal(t, len(x.services), 3)

	assert.Equal(t, x.services[0].Health().Name, EmptyServiceName)
	assert.Equal(t, x.services[1].Health().Name, EmptyServiceName)
	assert.Equal(t, x.services[2].Health().Name, MockServiceName)
}

func TestServiceHealths(t *testing.T) {
	x := NewTestHealthCheck()
	wg := &sync.WaitGroup{}
	x.SetServices([]Service{
		NewEmptyService(wg),
		NewEmptyService(wg),
		NewMockService(),
	})

	healths := x.ServiceHealths()

	assert.Equal(t, len(healths), 1)

	assert.Equal(t, healths[0].Name, MockServiceName)

}

func TestPostHealth(t *testing.T) {
	t.Run("No Error", func(t *testing.T) {
		x := NewTestHealthCheck()
		wg := &sync.WaitGroup{}
		x.SetServices([]Service{
			NewEmptyService(wg),
			NewEmptyService(wg),
			NewMockService(),
		})

		mockDB := mocks.NewMockDatabase(t)
		DB = mockDB

		filter := bson.M{
			"validator_id": x.validatorId,
			"hostname":     x.hostname,
		}

		onInsert := bson.M{
			"pokt_vault_address": x.poktVaultAddress,
			"pokt_signers":       x.poktSigners,
			"pokt_public_key":    x.poktPublicKey,
			"pokt_address":       x.poktAddress,
			"eth_validators":     x.ethValidators,
			"eth_address":        x.ethAddress,
			"wpokt_address":      x.wpoktAddress,
			"hostname":           x.hostname,
			"validator_id":       x.validatorId,
			"created_at":         nil,
		}

		onUpdate := bson.M{
			"healthy":         true,
			"service_healths": []models.ServiceHealth{},
			"updated_at":      nil,
		}

		update := bson.M{"$set": onUpdate, "$setOnInsert": onInsert}

		call := mockDB.EXPECT().UpsertOne(models.CollectionHealthChecks, filter, mock.Anything)
		call.Run(func(_ string, _ interface{}, arg interface{}) {

			updateArg := arg.(bson.M)

			updateArg["$setOnInsert"].(bson.M)["created_at"] = nil
			updateArg["$set"].(bson.M)["updated_at"] = nil
			updateArg["$set"].(bson.M)["service_healths"] = []models.ServiceHealth{}

			assert.Equal(t, updateArg, update)
		})
		call.Return(nil)

		success := x.PostHealth()
		assert.True(t, success)
	})

	t.Run("With Error", func(t *testing.T) {
		x := NewTestHealthCheck()
		wg := &sync.WaitGroup{}
		x.SetServices([]Service{
			NewEmptyService(wg),
			NewEmptyService(wg),
			NewMockService(),
		})

		mockDB := mocks.NewMockDatabase(t)
		DB = mockDB

		call := mockDB.EXPECT().UpsertOne(mock.Anything, mock.Anything, mock.Anything)
		call.Return(errors.New("error"))

		success := x.PostHealth()
		assert.False(t, success)
	})

	t.Run("Via Run", func(t *testing.T) {
		x := NewTestHealthCheck()
		wg := &sync.WaitGroup{}
		x.SetServices([]Service{
			NewEmptyService(wg),
			NewEmptyService(wg),
			NewMockService(),
		})

		mockDB := mocks.NewMockDatabase(t)
		DB = mockDB

		call := mockDB.EXPECT().UpsertOne(mock.Anything, mock.Anything, mock.Anything)
		call.Return(errors.New("error"))

		x.Run()
	})

}

/*
  validator_addresses:
    - "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
    - "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
    - "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
  vault_address: "E3BB46007E9BF127FD69B02DD5538848A80CADCE"
  multisig_public_keys:
    - "eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
    - "ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2"
    - "abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283"

*/

func TestNewHealthCheck(t *testing.T) {
	t.Run("With Empty Pocket Private Key", func(t *testing.T) {
		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Empty Eth Private Key", func(t *testing.T) {
		Config.Pocket.PrivateKey = "5efedbbc3d3d6f82d78eaf21258c81f462f3a25268be0018d4d75e1a4787bd14eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Empty MultiSig Keys", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.PrivateKey = "5efedbbc3d3d6f82d78eaf21258c81f462f3a25268be0018d4d75e1a4787bd14eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Invalid MultiSig Keys", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.PrivateKey = "5efedbbc3d3d6f82d78eaf21258c81f462f3a25268be0018d4d75e1a4787bd14eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
		Config.Pocket.MultisigPublicKeys = []string{"0x1234"}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Valid MultiSig Keys but Without Signer", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.PrivateKey = "5efedbbc3d3d6f82d78eaf21258c81f462f3a25268be0018d4d75e1a4787bd14eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
		Config.Pocket.MultisigPublicKeys = []string{
			// "eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Valid MultiSig Keys but Empty Vault Address", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.PrivateKey = "5efedbbc3d3d6f82d78eaf21258c81f462f3a25268be0018d4d75e1a4787bd14eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
		Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Valid Config", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.PrivateKey = "5efedbbc3d3d6f82d78eaf21258c81f462f3a25268be0018d4d75e1a4787bd14eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743"
		Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}
		Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

		x := NewHealthCheck()

		hostname, _ := os.Hostname()

		assert.NotNil(t, x)
		assert.Equal(t, strings.ToLower(Config.Pocket.VaultAddress), x.poktVaultAddress)
		assert.Equal(t, Config.Pocket.MultisigPublicKeys, x.poktSigners)
		assert.Equal(t, "wpokt-validator-01", x.validatorId)
		assert.Equal(t, hostname, x.hostname)

	})
}
