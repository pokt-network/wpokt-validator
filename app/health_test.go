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

	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"

	"github.com/dan13ram/wpokt-validator/app/mocks"
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
		call.Return(primitive.NewObjectID(), nil)

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
		call.Return(primitive.NewObjectID(), errors.New("error"))

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
		call.Return(primitive.NewObjectID(), errors.New("error"))

		x.Run()
	})

}

func TestNewHealthCheck(t *testing.T) {
	t.Run("With Empty Pocket Private Key", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Empty Eth Private Key", func(t *testing.T) {
		Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Empty MultiSig Keys", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		Config.Pocket.MultisigAddress = "pokt10r5n6x28p9qntchsmhxd4ftq9lk6vzcx3dv4gx"
		Config.Pocket.MultisigThreshold = 2
		Config.Pocket.Bech32Prefix = "pokt"
		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Invalid MultiSig Keys", func(t *testing.T) {
		Config.Pocket.MultisigPublicKeys = []string{"0x1234"}
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		Config.Pocket.MultisigAddress = "pokt10r5n6x28p9qntchsmhxd4ftq9lk6vzcx3dv4gx"
		Config.Pocket.MultisigThreshold = 2
		Config.Pocket.Bech32Prefix = "pokt"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Valid MultiSig Keys but Without Signer", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		Config.Pocket.MultisigPublicKeys = []string{
			// "0223aa679d6d5344e201e0df9f02ab15a84726eee0dfb4e953c46a9e2cb52349dc",
			"02faaaf0f385bb17381f36dcd86ab2486e8ff8d93440436496665ac007953076c2",
			"02cae233806460db75a941a269490ca5165a620b43241edb8bc72e169f4143a6df",
		}
		Config.Pocket.MultisigAddress = "pokt10r5n6x28p9qntchsmhxd4ftq9lk6vzcx3dv4gx"
		Config.Pocket.MultisigThreshold = 2
		Config.Pocket.Bech32Prefix = "pokt"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Valid MultiSig Keys but Empty Vault Address", func(t *testing.T) {
		Config.Pocket.MultisigAddress = ""
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		Config.Pocket.MultisigPublicKeys = []string{
			"0223aa679d6d5344e201e0df9f02ab15a84726eee0dfb4e953c46a9e2cb52349dc",
			"02faaaf0f385bb17381f36dcd86ab2486e8ff8d93440436496665ac007953076c2",
			"02cae233806460db75a941a269490ca5165a620b43241edb8bc72e169f4143a6df",
		}
		Config.Pocket.MultisigThreshold = 2
		Config.Pocket.Bech32Prefix = "pokt"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { NewHealthCheck() })
	})

	t.Run("With Valid Config", func(t *testing.T) {
		Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		Config.Pocket.MultisigPublicKeys = []string{
			"0223aa679d6d5344e201e0df9f02ab15a84726eee0dfb4e953c46a9e2cb52349dc",
			"02faaaf0f385bb17381f36dcd86ab2486e8ff8d93440436496665ac007953076c2",
			"02cae233806460db75a941a269490ca5165a620b43241edb8bc72e169f4143a6df",
		}
		Config.Pocket.MultisigAddress = "pokt10r5n6x28p9qntchsmhxd4ftq9lk6vzcx3dv4gx"
		Config.Pocket.MultisigThreshold = 2
		Config.Pocket.Bech32Prefix = "pokt"

		x := NewHealthCheck()

		hostname, _ := os.Hostname()

		assert.NotNil(t, x)
		assert.Equal(t, strings.ToLower(Config.Pocket.MultisigAddress), x.poktVaultAddress)
		assert.Equal(t, Config.Pocket.MultisigPublicKeys, x.poktSigners)
		assert.Equal(t, "wpokt-validator-01", x.validatorId)
		assert.Equal(t, hostname, x.hostname)

	})
}
