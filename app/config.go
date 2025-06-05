package app

import (
	"encoding/hex"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"
	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
)

var (
	Config models.Config
)

func InitConfig(configFile string, envFile string) {
	log.Debug("[CONFIG] Initializing config")
	readConfigFromConfigFile(configFile)
	readConfigFromENV(envFile)
	readKeysFromGSM()
	validateConfig()
	log.Info("[CONFIG] Config initialized")
}

func readConfigFromConfigFile(configFile string) bool {
	if configFile == "" {
		log.Debug("[CONFIG] No config file provided")
		return false
	}
	log.Debugf("[CONFIG] Reading config file %s", configFile)
	var yamlFile, err = os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("[CONFIG] Error reading config file %q: %s\n", configFile, err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("[CONFIG] Error unmarshalling config file %q: %s\n", configFile, err.Error())
	}
	log.Debugf("[CONFIG] Config loaded from %s", configFile)
	return true
}

func validateAndCreateSigner(mnemonic string) (common.Signer, error) {
	log.Debug("Initializing signer")

	return common.NewMnemonicSigner(mnemonic)

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

func validateConfig() {
	log.Debug("[CONFIG] Validating config")
	// mongodb
	if Config.MongoDB.URI == "" {
		log.Fatal("[CONFIG] MongoDB.URI is required")
	}
	if Config.MongoDB.Database == "" {
		log.Fatal("[CONFIG] MongoDB.Database is required")
	}
	if Config.MongoDB.TimeoutMillis == 0 {
		log.Fatal("[CONFIG] MongoDB.TimeoutMillis is required")
	}

	// ethereum
	if Config.Ethereum.RPCURL == "" {
		log.Fatal("[CONFIG] Ethereum.RPCURL is required")
	}
	if Config.Ethereum.ChainID == "" {
		log.Fatal("[CONFIG] Ethereum.ChainID is required")
	}
	if Config.Ethereum.RPCTimeoutMillis == 0 {
		log.Fatal("[CONFIG] Ethereum.RPCTimeoutMillis is required")
	}
	if Config.Ethereum.PrivateKey == "" {
		log.Fatal("[CONFIG] Ethereum.PrivateKey is required")
	} else {
		Config.Ethereum.PrivateKey = strings.TrimPrefix(Config.Ethereum.PrivateKey, "0x")
	}

	if Config.Ethereum.WrappedPocketAddress == "" {
		log.Fatal("[CONFIG] Ethereum.WrappedPocketAddress is required")
	}
	if Config.Ethereum.MintControllerAddress == "" {
		log.Fatal("[CONFIG] Ethereum.MintControllerAddress is required")
	}
	if len(Config.Ethereum.ValidatorAddresses) == 0 {
		log.Fatal("[CONFIG] Ethereum.ValidatorAddresses is required")
	}

	// pocket
	if Config.Pocket.Mnemonic == "" {
		log.Fatal("[CONFIG] Pocket.Mnemonic is required")
	}

	signer, err := validateAndCreateSigner(Config.Pocket.Mnemonic)
	if err != nil {
		log.Fatalf("[CONFIG] Error creating signer: %s", err.Error())
	}

	cosmosPubKeyHex := hex.EncodeToString(signer.CosmosPublicKey().Bytes())
	{
		// cosmos
		if Config.Pocket.StartHeight == 0 {
			log.Warn("Pocket.StartBlockHeight is 0")
		}
		if Config.Pocket.Confirmations == 0 {
			log.Warn("Pocket.Confirmations is 0")
		}
		if Config.Pocket.GRPCEnabled {
			if Config.Pocket.GRPCHost == "" {
				log.Fatal("Pocket.GRPCHost is required when GRPCEnabled is true")
			}
			if Config.Pocket.GRPCPort == 0 {
				log.Fatal("Pocket.GRPCPort is required when GRPCEnabled is true")
			}
		} else {
			if Config.Pocket.RPCURL == "" {
				log.Fatal("Pocket.RPCURL is required when GRPCEnabled is false")
			}
		}
		if Config.Pocket.RPCTimeoutMillis == 0 {
			log.Fatal("Pocket.TimeoutMS is required")
		}
		if Config.Pocket.ChainID == "" {
			log.Fatal("Pocket.ChainID is required")
		}
		if Config.Pocket.TxFee == 0 {
			log.Warn("Pocket.TxFee is 0")
		}
		if Config.Pocket.Bech32Prefix == "" {
			log.Fatal("Pocket.Bech32Prefix is required")
		}
		if Config.Pocket.CoinDenom == "" {
			log.Fatal("Pocket.CoinDenom is required")
		}
		if !common.IsValidBech32Address(Config.Pocket.Bech32Prefix, Config.Pocket.MultisigAddress) {
			log.Fatal("Pocket.MultisigAddress is invalid")
		}
		if len(Config.Pocket.MultisigPublicKeys) <= 1 {
			log.Fatal("Pocket.MultisigPublicKeys is required and must have at least 2 public keys")
		}
		foundPublicKey := false
		seen := make(map[string]bool)
		var pKeys []crypto.PubKey
		for j, publicKey := range Config.Pocket.MultisigPublicKeys {
			publicKey = strings.ToLower(publicKey)
			if !common.IsValidCosmosPublicKey(publicKey) {
				log.Fatalf("Pocket.MultisigPublicKeys[%d] is invalid", j)
			}
			if strings.EqualFold(publicKey, cosmosPubKeyHex) {
				foundPublicKey = true
			}
			pKey, _ := common.CosmosPublicKeyFromHex(publicKey) // cannot fail because public key is valid
			pKeys = append(pKeys, pKey)
			if seen[publicKey] {
				log.Fatalf("Pocket.MultisigPublicKeys[%d] is duplicated", j)
			}
			seen[publicKey] = true
		}
		if !foundPublicKey {
			log.Fatal("Pocket.MultisigPublicKeys must contain the public key of this oracle")
		}
		if Config.Pocket.MultisigThreshold == 0 || Config.Pocket.MultisigThreshold > uint64(len(Config.Pocket.MultisigPublicKeys)) {
			log.Fatal("Pocket.MultisigThreshold is invalid")
		}
		multisigPk := multisig.NewLegacyAminoPubKey(int(Config.Pocket.MultisigThreshold), pKeys)
		multisigBech32, _ := bech32.ConvertAndEncode(Config.Pocket.Bech32Prefix, multisigPk.Address().Bytes()) // no reason it should fail
		if !strings.EqualFold(Config.Pocket.MultisigAddress, multisigBech32) {
			log.Fatal("Pocket.MultisigAddress is not valid for the given public keys and threshold")
		}

	}

	// services
	if Config.MintMonitor.Enabled && Config.MintMonitor.IntervalMillis == 0 {
		log.Fatal("[CONFIG] MintMonitor.Interval is required")
	}
	if Config.MintSigner.Enabled && Config.MintSigner.IntervalMillis == 0 {
		log.Fatal("[CONFIG] MintSigner.Interval is required")
	}
	if Config.MintExecutor.Enabled && Config.MintExecutor.IntervalMillis == 0 {
		log.Fatal("[CONFIG] MintExecutor.Interval is required")
	}
	if Config.BurnMonitor.Enabled && Config.BurnMonitor.IntervalMillis == 0 {
		log.Fatal("[CONFIG] BurnMonitor.Interval is required")
	}
	if Config.BurnSigner.Enabled && Config.BurnSigner.IntervalMillis == 0 {
		log.Fatal("[CONFIG] BurnSigner.Interval is required")
	}
	if Config.BurnExecutor.Enabled && Config.BurnExecutor.IntervalMillis == 0 {
		log.Fatal("[CONFIG] BurnExecutor.Interval is required")
	}

	if Config.HealthCheck.IntervalMillis == 0 {
		log.Fatal("[CONFIG] HealthCheck.Interval is required")
	}

	log.Debug("[CONFIG] Config validated")
}
