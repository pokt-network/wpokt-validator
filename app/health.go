package app

import (
	"os"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	poktCrypto "github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
	wpoktAddress     string
	hostname         string
	lastSyncTime     time.Time
	interval         time.Duration
	services         []models.Service
}

type HealthServiceInterface interface {
	models.Service
	FindLastHealth() (models.Health, error)
	SetServices(services []models.Service)
}

func (b *HealthService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           b.name,
		LastSyncTime:   b.lastSyncTime,
		NextSyncTime:   b.lastSyncTime.Add(b.interval),
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
		b.lastSyncTime = time.Now()

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
		health := service.Health()
		if health.Name == models.EmptyServiceName {
			continue
		}
		serviceHealths = append(serviceHealths, health)
	}
	return serviceHealths
}

func (b *HealthService) PostHealth() bool {
	log.Debug("[HEALTH] Posting health")

	filter := bson.M{"hostname": b.hostname}

	health := models.Health{
		PoktVaultAddress: b.poktVaultAddress,
		PoktSigners:      b.poktSigners,
		PoktPublicKey:    b.poktPublicKey,
		PoktAddress:      b.poktAddress,
		EthValidators:    b.ethValidators,
		EthAddress:       b.ethAddress,
		WPoktAddress:     b.wpoktAddress,
		Hostname:         b.hostname,
		Healthy:          true,
		CreatedAt:        time.Now(),
		ServiceHealths:   b.ServiceHealths(),
	}

	update := bson.M{"$set": health}

	err := DB.UpsertOne(models.CollectionHealthChecks, filter, update)

	if err != nil {
		log.Error("[HEALTH] Error posting health: ", err)
		return false
	}
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

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("[HEALTH] Error getting hostname: ", err)
	}

	var pks []poktCrypto.PublicKey
	for _, pk := range Config.Pocket.MultisigPublicKeys {
		p, err := poktCrypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[HEALTH] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPkAddress := poktCrypto.PublicKeyMultiSignature{PublicKeys: pks}.Address().String()
	log.Debug("[HEALTH] Multisig address: ", multisigPkAddress)
	if multisigPkAddress != Config.Pocket.VaultAddress {
		log.Fatal("[HEALTH] Multisig address does not match vault address")
	}

	b := &HealthService{
		name:             "health",
		stop:             make(chan bool),
		interval:         time.Duration(Config.Health.IntervalSecs) * time.Second,
		poktVaultAddress: multisigPkAddress,
		poktSigners:      Config.Pocket.MultisigPublicKeys,
		poktPublicKey:    pk.PublicKey().RawString(),
		poktAddress:      pk.PublicKey().Address().String(),
		ethValidators:    Config.Ethereum.ValidatorAddresses,
		ethAddress:       ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex(),
		wpoktAddress:     Config.Ethereum.WPOKTAddress,
		hostname:         hostname,
		wg:               wg,
	}

	log.Info("[HEALTH] Initialized health")

	return b
}
