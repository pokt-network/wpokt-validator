package app

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"bytes"
	"sort"
	"time"

	"crypto/ecdsa"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	multisigtypes "github.com/cosmos/cosmos-sdk/crypto/types/multisig"
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

type PocketSigner struct {
	Signer          common.Signer
	Address         string
	SignerIndex     int
	Multisig        multisigtypes.PubKey
	MultisigAddress string
}

func CreatePocketSigner() (common.Signer, error) {
	return common.NewMnemonicSigner(Config.Pocket.Mnemonic)

	// // Mnemonic for both Ethereum and Cosmos networks
	// if config.Mnemonic == "" && config.GcpKmsKeyName == "" {
	// 	return nil, fmt.Errorf("Mnemonic or GcpKmsKeyName is required")
	// }
	// if config.Mnemonic != "" {
	// 	if !bip39.IsMnemonicValid(config.Mnemonic) {
	// 		return nil, fmt.Errorf("Mnemonic is invalid")
	// 	}
	//
	// 	return common.NewMnemonicSigner(config.Mnemonic)
	// }
	//
	// return common.NewGcpKmsSigner(config.GcpKmsKeyName)

}

func GetPocketSignerAndMultisig() (*PocketSigner, error) {
	signer, err := CreatePocketSigner()
	if err != nil {
		return nil, fmt.Errorf("error initializing pokt signer: %w", err)
	}

	cosmosPubKey := signer.CosmosPublicKey()
	poktAddress, err := common.Bech32FromBytes(Config.Pocket.Bech32Prefix, cosmosPubKey.Address().Bytes())
	if err != nil {
		return nil, fmt.Errorf("error getting pokt address: %w", err)
	}

	var pks []crypto.PubKey
	signerIndex := -1
	for index, pk := range Config.Pocket.MultisigPublicKeys {
		pKey, err := common.CosmosPublicKeyFromHex(pk)
		if err != nil {
			return nil, fmt.Errorf("error parsing multisig public key [%d]: %w", index, err)
		}
		pks = append(pks, pKey)
		if pKey.Equals(cosmosPubKey) {
			signerIndex = index
			log.Debugf("[HEALTH] Found current pocket signer at index %d", index)
		}
	}

	if signerIndex == -1 {
		// log.Fatal("[HEALTH] Multisig public keys do not contain signer")
		return nil, fmt.Errorf("multisig public keys do not contain signer")
	}

	if Config.Pocket.MultisigThreshold == 0 || Config.Pocket.MultisigThreshold > uint64(len(Config.Pocket.MultisigPublicKeys)) {
		return nil, fmt.Errorf("multisig threshold is invalid")
	}

	sort.Slice(pks, func(i, j int) bool {
		return bytes.Compare(pks[i].Address(), pks[j].Address()) < 0
	})

	multisigPk := multisig.NewLegacyAminoPubKey(int(Config.Pocket.MultisigThreshold), pks)
	multisigAddressBytes := multisigPk.Address().Bytes()
	multisigAddress, _ := common.Bech32FromBytes(Config.Pocket.Bech32Prefix, multisigAddressBytes)

	if !strings.EqualFold(multisigAddress, Config.Pocket.MultisigAddress) {
		return nil, fmt.Errorf("multisig address does not match vault address")
	}

	return &PocketSigner{
		Signer:          signer,
		SignerIndex:     signerIndex,
		Multisig:        multisigPk,
		MultisigAddress: multisigAddress,
		Address:         poktAddress,
	}, nil
}

type EthereumSigner struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

func GetEthereumSigner() (*EthereumSigner, error) {
	ethPK, err := ethCrypto.HexToECDSA(Config.Ethereum.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error initializing ethereum signer: %w", err)
	}

	ethAddress := ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex()

	return &EthereumSigner{
		PrivateKey: ethPK,
		Address:    ethAddress,
	}, nil
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
