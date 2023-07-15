package app

import (
	"crypto/ecdsa"
	"time"

	"github.com/dan13ram/wpokt-backend/models"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	poktCrypto "github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
)

type HealthService struct {
	stop          chan bool
	poktPublicKey string
	ethAddress    string
	interval      time.Duration
}

func (b *HealthService) Stop() {
	log.Debug("[HEALTH] Stopping health")
	b.stop <- true
}

func (b *HealthService) PostHealth() bool {
	log.Debug("[HEALTH] Posting health")
	health := models.Health{
		PoktPublicKey: b.poktPublicKey,
		EthAddress:    b.ethAddress,
		CreatedAt:     time.Now(),
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

		b.PostHealth()

		log.Debug("[HEALTH] Finished health sync")
		log.Debug("[HEALTH] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[HEALTH] Stopped health")
		case <-time.After(b.interval):
		}
	}
}

func NewHealthCheck() models.Service {
	log.Debug("[HEALTH] Initializing health")

	pk, err := poktCrypto.NewPrivateKey(Config.PoktSigner.PrivateKey)
	if err != nil {
		log.Fatal("[POKT SIGNER] Error initializing pokt signer: ", err)
	}
	log.Debug("[POKT SIGNER] Initialized pokt signer private key")
	log.Debug("[POKT SIGNER] Pokt signer public key: ", pk.PublicKey().RawString())
	log.Debug("[POKT SIGNER] Pokt signer address: ", pk.PublicKey().Address().String())

	privateKey, err := ethCrypto.HexToECDSA(Config.WPOKTSigner.PrivateKey)
	if err != nil {
		log.Fatal("[WPOKT SIGNER] Error loading private key: ", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("[WPOKT SIGNER] Error casting public key to ECDSA")
	}
	address := ethCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	b := &HealthService{
		stop:          make(chan bool),
		interval:      time.Duration(Config.Health.IntervalSecs) * time.Second,
		poktPublicKey: pk.PublicKey().RawString(),
		ethAddress:    address,
	}

	log.Debug("[HEALTH] Initialized health")

	return b
}
