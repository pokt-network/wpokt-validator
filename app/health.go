package app

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	poktCrypto "github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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

	err := DB.UpsertOne(models.CollectionHealthChecks, filter, update)

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

	pk, err := poktCrypto.NewPrivateKey(Config.Pocket.PrivateKey)
	if err != nil {
		log.Fatal("[HEALTH] Error initializing pokt signer: ", err)
	}
	log.Debug("[HEALTH] Initialized pokt signer private key")
	log.Debug("[HEALTH] Pokt signer public key: ", pk.PublicKey().RawString())
	log.Debug("[HEALTH] Pokt signer address: ", pk.PublicKey().Address().String())

	ethPK, err := ethCrypto.HexToECDSA(Config.Ethereum.PrivateKey)
	log.Debug("[HEALTH] Initialized private key")
	log.Debug("[HEALTH] ETH Address: ", ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex())

	ethAddress := ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex()
	poktAddress := pk.PublicKey().Address().String()

	var pks []poktCrypto.PublicKey
	var signerIndex int
	for _, pk := range Config.Pocket.MultisigPublicKeys {
		p, err := poktCrypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[HEALTH] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
		if p.Address().String() == poktAddress {
			signerIndex = len(pks)
		}
	}

	validatorId := "wpokt-validator-" + fmt.Sprintf("%02d", signerIndex)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("[HEALTH] Error getting hostname: ", err)
	}

	multisigPkAddress := poktCrypto.PublicKeyMultiSignature{PublicKeys: pks}.Address().String()
	log.Debug("[HEALTH] Multisig address: ", multisigPkAddress)
	if multisigPkAddress != Config.Pocket.VaultAddress {
		log.Fatal("[HEALTH] Multisig address does not match vault address")
	}

	x := &HealthCheckRunner{
		poktVaultAddress: multisigPkAddress,
		poktSigners:      Config.Pocket.MultisigPublicKeys,
		poktPublicKey:    pk.PublicKey().RawString(),
		poktAddress:      poktAddress,
		ethValidators:    Config.Ethereum.ValidatorAddresses,
		ethAddress:       ethAddress,
		wpoktAddress:     Config.Ethereum.WrappedPocketAddress,
		hostname:         hostname,
		validatorId:      validatorId,
	}

	log.Info("[HEALTH] Initialized health")

	return x
}

func NewHealthService(x *HealthCheckRunner, wg *sync.WaitGroup) Service {
	service := NewRunnerService(HealthCheckName, x, wg, time.Duration(Config.HealthCheck.IntervalSecs)*time.Second)
	return service
}
