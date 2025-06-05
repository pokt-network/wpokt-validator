package app

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

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

	signer, err := common.NewMnemonicSigner(Config.Pocket.Mnemonic)
	if err != nil {
		log.Fatal("[HEALTH] Error initializing pokt signer: ", err)
	}

	cosmosPubKey := signer.CosmosPublicKey()
	poktPubKeyHex := hex.EncodeToString(cosmosPubKey.Bytes())
	poktAddress, err := common.Bech32FromBytes(Config.Pocket.Bech32Prefix, cosmosPubKey.Address().Bytes())

	log.Debug("[HEALTH] Initialized pokt signer private key")
	log.Debug("[HEALTH] Pokt signer public key: ", poktPubKeyHex)
	log.Debug("[HEALTH] Pokt signer address: ", poktAddress)

	ethPK, err := ethCrypto.HexToECDSA(Config.Ethereum.PrivateKey)
	if err != nil {
		log.Fatal("[HEALTH] Error initializing ethereum signer: ", err)
	}
	log.Debug("[HEALTH] Initialized private key")
	log.Debug("[HEALTH] ETH Address: ", ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex())

	ethAddress := ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex()

	var pks []crypto.PubKey
	signerIndex := -1
	for _, pk := range Config.Pocket.MultisigPublicKeys {
		pKey, err := common.CosmosPublicKeyFromHex(pk)
		if err != nil {
			log.Fatalf("Error parsing public key: %s", err)
		}
		pks = append(pks, pKey)
		if strings.EqualFold(hex.EncodeToString(pKey.Bytes()), poktPubKeyHex) {
			signerIndex = len(pks)
		}
	}

	if signerIndex == -1 {
		log.Fatal("[HEALTH] Multisig public keys do not contain signer")
	}

	multisigPk := multisig.NewLegacyAminoPubKey(int(Config.Pocket.MultisigThreshold), pks)
	multisigAddressBytes := multisigPk.Address().Bytes()
	multisigAddress, _ := common.Bech32FromBytes(Config.Pocket.Bech32Prefix, multisigAddressBytes)

	if !strings.EqualFold(multisigAddress, Config.Pocket.MultisigAddress) {
		log.Fatal("[HEALTH] Multisig address does not match vault address")
	}

	validatorId := "wpokt-validator-" + fmt.Sprintf("%02d", signerIndex)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("[HEALTH] Error getting hostname: ", err)
	}

	log.Debug("[HEALTH] Multisig address: ", multisigAddress)

	x := &HealthCheckRunner{
		poktVaultAddress: strings.ToLower(multisigAddress),
		poktSigners:      Config.Pocket.MultisigPublicKeys,
		poktPublicKey:    poktPubKeyHex,
		poktAddress:      strings.ToLower(poktAddress),
		ethValidators:    Config.Ethereum.ValidatorAddresses,
		ethAddress:       strings.ToLower(ethAddress),
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
