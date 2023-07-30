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
	HealthServiceName = "health"
)

type HealthService struct {
	wg               *sync.WaitGroup
	stop             chan bool
	poktSigners      []string
	poktPublicKey    string
	poktAddress      string
	poktVaultAddress string
	ethValidators    []string
	ethAddress       string
	wpoktAddress     string
	hostname         string
	validatorId      string
	interval         time.Duration
	services         []models.Service
}

type HealthServiceInterface interface {
	models.Service
	FindLastHealth() (models.Health, error)
	SetServices(services []models.Service)
}

func (b *HealthService) Health() models.ServiceHealth {
	lastSyncTime := time.Now()
	return models.ServiceHealth{
		Name:           HealthServiceName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(b.interval),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}
func (b *HealthService) Start() {
	log.Info("[HEALTH] Starting service")
	stop := false
	for !stop {
		log.Info("[HEALTH] Starting sync")

		b.PostHealth()

		log.Info("[HEALTH] Finished sync, Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			b.PostHealth()
			log.Info("[HEALTH] Stopped service")
			b.wg.Done()
		case <-time.After(b.interval):
		}
	}
}

func (b *HealthService) Stop() {
	log.Debug("[HEALTH] Stopping health")
	b.stop <- true
}

func (b *HealthService) FindLastHealth() (models.Health, error) {
	var health models.Health
	filter := bson.M{"hostname": b.hostname}
	err := DB.FindOne(models.CollectionHealthChecks, filter, &health)
	return health, err
}

func (b *HealthService) ServiceHealths() []models.ServiceHealth {
	var serviceHealths []models.ServiceHealth
	for _, service := range b.services {
		serviceHealth := service.Health()
		if serviceHealth.Name == models.EmptyServiceName || serviceHealth.Name == "" {
			continue
		}
		serviceHealths = append(serviceHealths, serviceHealth)
	}
	return serviceHealths
}

func (b *HealthService) PostHealth() bool {
	log.Debug("[HEALTH] Posting health")

	filter := bson.M{
		"validator_id": b.validatorId,
		"hostname":     b.hostname,
	}

	onInsert := bson.M{
		"pokt_vault_address": b.poktVaultAddress,
		"pokt_signers":       b.poktSigners,
		"pokt_public_key":    b.poktPublicKey,
		"pokt_address":       b.poktAddress,
		"eth_validators":     b.ethValidators,
		"eth_address":        b.ethAddress,
		"wpokt_address":      b.wpoktAddress,
		"hostname":           b.hostname,
		"validator_id":       b.validatorId,
		"created_at":         time.Now(),
	}

	onUpdate := bson.M{
		"healthy":         true,
		"service_healths": b.ServiceHealths(),
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

func (b *HealthService) SetServices(services []models.Service) {
	b.services = services
}

func NewHealthCheck(wg *sync.WaitGroup) HealthServiceInterface {
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

	b := &HealthService{
		stop:             make(chan bool),
		interval:         time.Duration(Config.HealthCheck.IntervalSecs) * time.Second,
		poktVaultAddress: multisigPkAddress,
		poktSigners:      Config.Pocket.MultisigPublicKeys,
		poktPublicKey:    pk.PublicKey().RawString(),
		poktAddress:      poktAddress,
		ethValidators:    Config.Ethereum.ValidatorAddresses,
		ethAddress:       ethAddress,
		wpoktAddress:     Config.Ethereum.WrappedPocketAddress,
		hostname:         hostname,
		validatorId:      validatorId,
		wg:               wg,
	}

	log.Info("[HEALTH] Initialized health")

	return b
}
