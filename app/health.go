package app

import (
	"os"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-backend/models"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	poktCrypto "github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
)

type HealthService struct {
	wg               *sync.WaitGroup
	name             string
	stop             chan bool
	poktSigners      []string
	poktPublicKey    string
	poktAddress      string
	poktVaultAddress string
	ethValidators    []string
	ethAddress       string
	hostname         string
	lastSyncTime     time.Time
	interval         time.Duration
	services         []models.Service
}

func (b *HealthService) EthBlockNumber() string {
	return ""
}

func (b *HealthService) PoktHeight() string {
	return ""
}

func (b *HealthService) LastSyncTime() time.Time {
	return b.lastSyncTime
}

func (b *HealthService) Interval() time.Duration {
	return b.interval
}

func (b *HealthService) Name() string {
	return b.name
}

func (b *HealthService) Stop() {
	log.Debug("[HEALTH] Stopping health")
	b.stop <- true
}

func (b *HealthService) ServiceHealths() []models.ServiceHealth {
	var serviceHealths []models.ServiceHealth
	for _, service := range b.services {
		var serviceHealth models.ServiceHealth
		serviceHealth.Name = service.Name()
		serviceHealth.Healthy = true
		serviceHealth.LastSyncTime = service.LastSyncTime()
		serviceHealth.NextSyncTime = service.LastSyncTime().Add(service.Interval())
		serviceHealth.EthBlockNumber = service.EthBlockNumber()
		serviceHealth.PoktHeight = service.PoktHeight()
		serviceHealths = append(serviceHealths, serviceHealth)
	}
	serviceHealths = append(serviceHealths, models.ServiceHealth{
		Name:           "health",
		Healthy:        true,
		LastSyncTime:   b.LastSyncTime(),
		NextSyncTime:   b.LastSyncTime().Add(b.Interval()),
		EthBlockNumber: b.EthBlockNumber(),
		PoktHeight:     b.PoktHeight(),
	})
	return serviceHealths
}

func (b *HealthService) PostHealth() bool {
	log.Debug("[HEALTH] Posting health")
	health := models.Health{
		PoktVaultAddress: b.poktVaultAddress,
		PoktSigners:      b.poktSigners,
		PoktPublicKey:    b.poktPublicKey,
		PoktAddress:      b.poktAddress,
		EthValidators:    b.ethValidators,
		EthAddress:       b.ethAddress,
		Hostname:         b.hostname,
		Healthy:          true,
		CreatedAt:        time.Now(),
		ServiceHealths:   b.ServiceHealths(),
	}

	err := DB.InsertOne(models.CollectionHealthChecks, health)

	if err != nil {
		log.Error("[HEALTH] Error posting health: ", err)
		return false
	}
	return true
}

func (b *HealthService) Start() {
	log.Debug("[HEALTH] Starting health")
	stop := false
	for !stop {
		log.Debug("[HEALTH] Starting health sync")
		b.lastSyncTime = time.Now()

		b.PostHealth()

		log.Debug("[HEALTH] Finished health sync")
		log.Debug("[HEALTH] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			b.PostHealth()
			log.Debug("[HEALTH] Stopped health")
			b.wg.Done()
		case <-time.After(b.interval):
		}
	}
}

func NewHealthCheck(services []models.Service, wg *sync.WaitGroup) models.Service {
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

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("[HEALTH] Error getting hostname: ", err)
	}

	var pks []poktCrypto.PublicKey
	for _, pk := range Config.Pocket.MultisigPublicKeys {
		p, err := poktCrypto.NewPublicKey(pk)
		if err != nil {
			log.Debug("[HEALTH] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := poktCrypto.PublicKeyMultiSignature{PublicKeys: pks}
	log.Debug("[HEALTH] Multisig address: ", multisigPk.Address().String())

	b := &HealthService{
		name:             "health",
		stop:             make(chan bool),
		interval:         time.Duration(Config.Health.IntervalSecs) * time.Second,
		poktVaultAddress: multisigPk.Address().String(),
		poktSigners:      Config.Pocket.MultisigPublicKeys,
		poktPublicKey:    pk.PublicKey().RawString(),
		poktAddress:      pk.PublicKey().Address().String(),
		ethValidators:    Config.Ethereum.ValidatorAddresses,
		ethAddress:       ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex(),
		hostname:         hostname,
		services:         services,
		wg:               wg,
	}

	log.Debug("[HEALTH] Initialized health")

	return b
}
