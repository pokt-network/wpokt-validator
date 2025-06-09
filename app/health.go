package app

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dan13ram/wpokt-validator/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"encoding/hex"
)

const (
	HealthCheckName = "HEALTH"
)

type HealthCheckRunner struct {
	poktSigners      []string
	poktPublicKey    string
	poktAddress      string
	poktVaultAddress string
	ethValidators    []string
	ethAddress       string
	wpoktAddress     string
	hostname         string
	validatorId      string
	services         []Service
}

func (x *HealthCheckRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{}
}

func (x *HealthCheckRunner) Run() {
	x.PostHealth()
}

func (x *HealthCheckRunner) FindLastHealth() (models.Health, error) {
	var health models.Health
	filter := bson.M{
		"validator_id": x.validatorId,
		"hostname":     x.hostname,
	}
	err := DB.FindOne(models.CollectionHealthChecks, filter, &health)
	return health, err
}

func (x *HealthCheckRunner) ServiceHealths() []models.ServiceHealth {
	var serviceHealths []models.ServiceHealth
	for _, service := range x.services {
		serviceHealth := service.Health()
		if serviceHealth.Name == EmptyServiceName || serviceHealth.Name == "" {
			continue
		}
		serviceHealths = append(serviceHealths, serviceHealth)
	}
	return serviceHealths
}

func (x *HealthCheckRunner) PostHealth() bool {
	log.Debug("[HEALTH] Posting health")

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
		"created_at":         time.Now(),
	}

	onUpdate := bson.M{
		"healthy":         true,
		"service_healths": x.ServiceHealths(),
		"updated_at":      time.Now(),
	}

	update := bson.M{"$set": onUpdate, "$setOnInsert": onInsert}

	_, err := DB.UpsertOne(models.CollectionHealthChecks, filter, update)

	if err != nil {
		log.Error("[HEALTH] Error posting health: ", err)
		return false
	}

	log.Info("[HEALTH] Posted health")
	return true
}

func (x *HealthCheckRunner) SetServices(services []Service) {
	x.services = services
}

func NewHealthCheck() *HealthCheckRunner {
	log.Debug("[HEALTH] Initializing health")

	poktSigner, err := GetPocketSignerAndMultisig()
	if err != nil {
		log.Fatal("[HEALTH] Error getting pokt signer and multisig: ", err)
	}

	log.Debug("[HEALTH] POKT address: ", poktSigner.Address)

	ethSigner, err := GetEthereumSigner()
	if err != nil {
		log.Fatal("[HEALTH] Error getting ethereum signer: ", err)
	}

	log.Debug("[HEALTH] ETH Address: ", ethSigner.Address)

	validatorId := "wpokt-validator-" + fmt.Sprintf("%02d", poktSigner.SignerIndex+1)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("[HEALTH] Error getting hostname: ", err)
	}

	x := &HealthCheckRunner{
		poktVaultAddress: poktSigner.MultisigAddress,
		poktSigners:      Config.Pocket.MultisigPublicKeys,
		poktPublicKey:    hex.EncodeToString(poktSigner.Signer.CosmosPublicKey().Bytes()),
		poktAddress:      poktSigner.Address,
		ethValidators:    Config.Ethereum.ValidatorAddresses,
		ethAddress:       ethSigner.Address,
		wpoktAddress:     strings.ToLower(Config.Ethereum.WrappedPocketAddress),
		hostname:         hostname,
		validatorId:      validatorId,
	}

	log.Info("[HEALTH] Initialized health")

	return x
}

func NewHealthService(x *HealthCheckRunner, wg *sync.WaitGroup) Service {
	service := NewRunnerService(HealthCheckName, x, wg, time.Duration(Config.HealthCheck.IntervalMillis)*time.Millisecond)
	return service
}
