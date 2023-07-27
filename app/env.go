package app

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func readConfigFromENV(envFile string) {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Warn("[ENV] Error loading .env file: ", err.Error())
		}
	}

	// mongodb
	if Config.MongoDB.URI == "" {
		Config.MongoDB.URI = os.Getenv("MONGODB_URI")
	}
	if Config.MongoDB.Database == "" {
		Config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	}
	if Config.MongoDB.TimeoutSecs == 0 {
		timeoutSecs, err := strconv.ParseInt(os.Getenv("MONGODB_TIMEOUT_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MONGODB_TIMEOUT_SECS: ", err.Error())
			log.Warn("[ENV] Setting MongoDB TimeoutSecs to 30")
			Config.MongoDB.TimeoutSecs = 30
		} else {
			Config.MongoDB.TimeoutSecs = timeoutSecs
		}
	}

	// ethereum
	if Config.Ethereum.RPCURL == "" {
		Config.Ethereum.RPCURL = os.Getenv("ETH_RPC_URL")
	}
	if Config.Ethereum.ChainId == "" {
		Config.Ethereum.ChainId = os.Getenv("ETH_CHAIN_ID")
	}
	if Config.Ethereum.PrivateKey == "" {
		Config.Ethereum.PrivateKey = os.Getenv("ETH_PRIVATE_KEY")
	}
	if Config.Ethereum.StartBlockNumber == 0 {
		blockNumber, err := strconv.ParseInt(os.Getenv("ETH_START_BLOCK_NUMBER"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_START_BLOCK_NUMBER: ", err.Error())
			log.Warn("[ENV] Setting Ethereum StartBlockNumber to 0")
			Config.Ethereum.StartBlockNumber = 0
		} else {
			Config.Ethereum.StartBlockNumber = blockNumber
		}
	}
	if Config.Ethereum.Confirmations == 0 {
		confirmations, err := strconv.ParseInt(os.Getenv("ETH_CONFIRMATIONS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_CONFIRMATIONS: ", err.Error())
			log.Warn("[ENV] Setting Ethereum Confirmations to 0")
			Config.Ethereum.Confirmations = 0
		} else {
			Config.Ethereum.Confirmations = confirmations
		}
	}
	if Config.Ethereum.RPCTimeoutSecs == 0 {
		timeoutSecs, err := strconv.ParseInt(os.Getenv("ETH_RPC_TIMEOUT_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_RPC_TIMEOUT_SECS: ", err.Error())
			log.Warn("[ENV] Setting Ethereum RPCTimeoutSecs to 30")
			Config.Ethereum.RPCTimeoutSecs = 30
		} else {
			Config.Ethereum.RPCTimeoutSecs = timeoutSecs
		}
	}
	if Config.Ethereum.WrappedPocketAddress == "" {
		Config.Ethereum.WrappedPocketAddress = os.Getenv("ETH_WRAPPED_POCKET_ADDRESS")
	}
	if Config.Ethereum.MintControllerAddress == "" {
		Config.Ethereum.MintControllerAddress = os.Getenv("ETH_MINT_CONTROLLER_ADDRESS")
	}
	if Config.Ethereum.ValidatorAddresses == nil || len(Config.Ethereum.ValidatorAddresses) == 0 {
		validatorAddresses := os.Getenv("ETH_VALIDATOR_ADDRESSES")
		if validatorAddresses != "" {
			Config.Ethereum.ValidatorAddresses = strings.Split(validatorAddresses, ",")
		}
	}

	// pocket
	if Config.Pocket.RPCURL == "" {
		Config.Pocket.RPCURL = os.Getenv("POKT_RPC_URL")
	}
	if Config.Pocket.ChainId == "" {
		Config.Pocket.ChainId = os.Getenv("POKT_CHAIN_ID")
	}
	if Config.Pocket.PrivateKey == "" {
		Config.Pocket.PrivateKey = os.Getenv("POKT_PRIVATE_KEY")
	}
	if Config.Pocket.StartHeight == 0 {
		startHeight, err := strconv.ParseInt(os.Getenv("POKT_START_HEIGHT"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_START_HEIGHT: ", err.Error())
			log.Warn("[ENV] Setting Pocket StartHeight to 0")
			Config.Pocket.StartHeight = 0
		} else {
			Config.Pocket.StartHeight = startHeight
		}
	}
	if Config.Pocket.Confirmations == 0 {
		confirmations, err := strconv.ParseInt(os.Getenv("POKT_CONFIRMATIONS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_CONFIRMATIONS: ", err.Error())
			log.Warn("[ENV] Setting Pocket Confirmations to 0")
		} else {
			Config.Pocket.Confirmations = confirmations
		}
	}
	if Config.Pocket.RPCTimeoutSecs == 0 {
		timeoutSecs, err := strconv.ParseInt(os.Getenv("POKT_RPC_TIMEOUT_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_RPC_TIMEOUT_SECS: ", err.Error())
			log.Warn("[ENV] Setting Pocket RPCTimeoutSecs to 30")
			Config.Pocket.RPCTimeoutSecs = 30
		} else {
			Config.Pocket.RPCTimeoutSecs = timeoutSecs
		}
	}
	if Config.Pocket.TxFee == 0 {
		txFee, err := strconv.ParseInt(os.Getenv("POKT_TX_FEE"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_TX_FEE: ", err.Error())
			log.Warn("[ENV] Setting Pocket TxFee to 10000")
			Config.Pocket.TxFee = 10000
		} else {
			Config.Pocket.TxFee = txFee
		}
	}
	if Config.Pocket.VaultAddress == "" {
		Config.Pocket.VaultAddress = os.Getenv("POKT_VAULT_ADDRESS")
	}
	if Config.Pocket.MultisigPublicKeys == nil || len(Config.Pocket.MultisigPublicKeys) == 0 {
		multisigPublicKeys := os.Getenv("POKT_MULTISIG_PUBLIC_KEYS")
		Config.Pocket.MultisigPublicKeys = strings.Split(multisigPublicKeys, ",")
	}

	// mint monitor
	if Config.MintMonitor.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("MINT_MONITOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_MONITOR_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting MintMonitor Enabled to false")
		} else {
			Config.MintMonitor.Enabled = enabled
		}
	}
	if Config.MintMonitor.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("MINT_MONITOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_MONITOR_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting MintMonitor IntervalSecs to 300")
			Config.MintMonitor.IntervalSecs = 300
		} else {
			Config.MintMonitor.IntervalSecs = intervalSecs
		}
	}

	// mint signer
	if Config.MintSigner.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("MINT_SIGNER_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_SIGNER_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting MintSigner Enabled to false")
		} else {
			Config.MintSigner.Enabled = enabled
		}
	}
	if Config.MintSigner.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("MINT_SIGNER_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_SIGNER_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting MintSigner IntervalSecs to 300")
			Config.MintSigner.IntervalSecs = 300
		} else {
			Config.MintSigner.IntervalSecs = intervalSecs
		}
	}

	// mint executor
	if Config.MintExecutor.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("MINT_EXECUTOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_EXECUTOR_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting MintExecutor Enabled to false")
		} else {
			Config.MintExecutor.Enabled = enabled
		}
	}
	if Config.MintExecutor.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("MINT_EXECUTOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_EXECUTOR_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting MintExecutor IntervalSecs to 300")
			Config.MintExecutor.IntervalSecs = 300
		} else {
			Config.MintExecutor.IntervalSecs = intervalSecs
		}
	}

	// burn monitor
	if Config.BurnMonitor.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("BURN_MONITOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_MONITOR_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting BurnMonitor Enabled to false")
		} else {
			Config.BurnMonitor.Enabled = enabled
		}
	}
	if Config.BurnMonitor.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("BURN_MONITOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_MONITOR_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting BurnMonitor IntervalSecs to 300")
			Config.BurnMonitor.IntervalSecs = 300
		} else {
			Config.BurnMonitor.IntervalSecs = intervalSecs
		}
	}

	// burn signer
	if Config.BurnSigner.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("BURN_SIGNER_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_SIGNER_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting BurnSigner Enabled to false")
		} else {
			Config.BurnSigner.Enabled = enabled
		}
	}
	if Config.BurnSigner.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("BURN_SIGNER_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_SIGNER_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting BurnSigner IntervalSecs to 300")
			Config.BurnSigner.IntervalSecs = 300
		} else {
			Config.BurnSigner.IntervalSecs = intervalSecs
		}
	}

	// burn executor
	if Config.BurnExecutor.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("BURN_EXECUTOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_EXECUTOR_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting BurnExecutor Enabled to false")
		} else {
			Config.BurnExecutor.Enabled = enabled
		}
	}
	if Config.BurnExecutor.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("BURN_EXECUTOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_EXECUTOR_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting BurnExecutor IntervalSecs to 300")
			Config.BurnExecutor.IntervalSecs = 300
		} else {
			Config.BurnExecutor.IntervalSecs = intervalSecs
		}
	}

	// health check
	if Config.HealthCheck.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("HEALTH_CHECK_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing HEALTH_CHECK_INTERVAL_SECS: ", err.Error())
			log.Warn("[ENV] Setting HealthCheck IntervalSecs to 300")
			Config.HealthCheck.IntervalSecs = 300
		} else {
			Config.HealthCheck.IntervalSecs = intervalSecs
		}
	}

	// logging
	if Config.Logger.Level == "" {
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			log.Warn("[ENV] Setting LogLevel to debug")
			Config.Logger.Level = "debug"
		} else {
			Config.Logger.Level = logLevel
		}
	}

	// google secret manager
	if Config.GoogleSecretManager.Enabled == false {
		enabled, err := strconv.ParseBool(os.Getenv("GOOGLE_SECRET_MANAGER_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing GOOGLE_SECRET_MANAGER_ENABLED: ", err.Error())
			log.Warn("[ENV] Setting GoogleSecretManager Enabled to false")
		} else {
			Config.GoogleSecretManager.Enabled = enabled
		}
	}
	if Config.GoogleSecretManager.ProjectId == "" {
		Config.GoogleSecretManager.ProjectId = os.Getenv("GOOGLE_PROJECT_ID")
	}
	if Config.GoogleSecretManager.PoktSecretName == "" {
		Config.GoogleSecretManager.PoktSecretName = os.Getenv("GOOGLE_POKT_SECRET_NAME")
	}
	if Config.GoogleSecretManager.EthSecretName == "" {
		Config.GoogleSecretManager.EthSecretName = os.Getenv("GOOGLE_ETH_SECRET_NAME")
	}

}
