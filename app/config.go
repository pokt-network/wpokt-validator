package app

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"
	"gopkg.in/yaml.v2"
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

func validateConfig() {
	log.Debug("[CONFIG] Validating config")
	{
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
	}

	{
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

		signer, err := GetEthereumSigner()
		if err != nil {
			log.Fatalf("[CONFIG] Error creating ethereum signer: %s", err.Error())
		}

		foundValidatorAddress := false
		for index, validatorAddress := range Config.Ethereum.ValidatorAddresses {
			if !common.IsValidEthereumAddress(validatorAddress) {
				log.Fatalf("[CONFIG] Ethereum.ValidatorAddresses[%d] is invalid: %s", index, validatorAddress)
			}
			if strings.EqualFold(validatorAddress, signer.Address) {
				foundValidatorAddress = true
				log.Debugf("[CONFIG] Found current ethereum validator address at index %d", index)
			}
		}
		if !foundValidatorAddress {
			log.Fatalf("[CONFIG] Ethereum.ValidatorAddresses does not contain validator address")
		}

	}

	// pocket

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

		if Config.Pocket.Mnemonic == "" {
			log.Fatal("[CONFIG] Pocket.Mnemonic is required")
		}

		_, err := GetPocketSignerAndMultisig()
		if err != nil {
			log.Fatalf("[CONFIG] Error creating pocket signer: %s", err.Error())
		}

	}

	{
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
	}

	{
		// health check
		if Config.HealthCheck.IntervalMillis == 0 {
			log.Fatal("[CONFIG] HealthCheck.Interval is required")
		}
	}

	log.Debug("[CONFIG] Config validated")
}
