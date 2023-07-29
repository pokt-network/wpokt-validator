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

	if os.Getenv("MONGODB_URI") != "" {
		Config.MongoDB.URI = os.Getenv("MONGODB_URI")
	}
	if os.Getenv("MONGODB_DATABASE") != "" {
		Config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	}
	if os.Getenv("MONGODB_TIMEOUT_SECS") != "" {
		timeoutSecs, err := strconv.ParseInt(os.Getenv("MONGODB_TIMEOUT_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MONGODB_TIMEOUT_SECS: ", err.Error())
		} else {
			Config.MongoDB.TimeoutSecs = timeoutSecs
		}
	}

	// ethereum
	if os.Getenv("ETH_RPC_URL") != "" {
		Config.Ethereum.RPCURL = os.Getenv("ETH_RPC_URL")
	}
	if os.Getenv("ETH_CHAIN_ID") != "" {
		Config.Ethereum.ChainId = os.Getenv("ETH_CHAIN_ID")
	}
	if os.Getenv("ETH_PRIVATE_KEY") != "" {
		Config.Ethereum.PrivateKey = os.Getenv("ETH_PRIVATE_KEY")
	}
	if os.Getenv("ETH_START_BLOCK_NUMBER") != "" {
		blockNumber, err := strconv.ParseInt(os.Getenv("ETH_START_BLOCK_NUMBER"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_START_BLOCK_NUMBER: ", err.Error())
		} else {
			Config.Ethereum.StartBlockNumber = blockNumber
		}
	}
	if os.Getenv("ETH_CONFIRMATIONS") != "" {
		confirmations, err := strconv.ParseInt(os.Getenv("ETH_CONFIRMATIONS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_CONFIRMATIONS: ", err.Error())
		} else {
			Config.Ethereum.Confirmations = confirmations
		}
	}
	if os.Getenv("ETH_RPC_TIMEOUT_SECS") != "" {
		timeoutSecs, err := strconv.ParseInt(os.Getenv("ETH_RPC_TIMEOUT_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_RPC_TIMEOUT_SECS: ", err.Error())
		} else {
			Config.Ethereum.RPCTimeoutSecs = timeoutSecs
		}
	}
	if os.Getenv("ETH_WRAPPED_POCKET_ADDRESS") != "" {
		Config.Ethereum.WrappedPocketAddress = os.Getenv("ETH_WRAPPED_POCKET_ADDRESS")
	}
	if os.Getenv("ETH_MINT_CONTROLLER_ADDRESS") != "" {
		Config.Ethereum.MintControllerAddress = os.Getenv("ETH_MINT_CONTROLLER_ADDRESS")
	}
	if os.Getenv("ETH_VALIDATOR_ADDRESSES") != "" {
		Config.Ethereum.ValidatorAddresses = strings.Split(os.Getenv("ETH_VALIDATOR_ADDRESSES"), ",")
	}

	// pocket
	if os.Getenv("POKT_RPC_URL") != "" {
		Config.Pocket.RPCURL = os.Getenv("POKT_RPC_URL")
	}
	if os.Getenv("POKT_CHAIN_ID") != "" {
		Config.Pocket.ChainId = os.Getenv("POKT_CHAIN_ID")
	}
	if os.Getenv("POKT_PRIVATE_KEY") != "" {
		Config.Pocket.PrivateKey = os.Getenv("POKT_PRIVATE_KEY")
	}
	if os.Getenv("POKT_START_HEIGHT") != "" {
		startHeight, err := strconv.ParseInt(os.Getenv("POKT_START_HEIGHT"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_START_HEIGHT: ", err.Error())
		} else {
			Config.Pocket.StartHeight = startHeight
		}
	}
	if os.Getenv("POKT_CONFIRMATIONS") != "" {
		confirmations, err := strconv.ParseInt(os.Getenv("POKT_CONFIRMATIONS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_CONFIRMATIONS: ", err.Error())
		} else {
			Config.Pocket.Confirmations = confirmations
		}
	}
	if os.Getenv("POKT_RPC_TIMEOUT_SECS") != "" {
		timeoutSecs, err := strconv.ParseInt(os.Getenv("POKT_RPC_TIMEOUT_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_RPC_TIMEOUT_SECS: ", err.Error())
		} else {
			Config.Pocket.RPCTimeoutSecs = timeoutSecs
		}
	}
	if os.Getenv("POKT_TX_FEE") != "" {
		txFee, err := strconv.ParseInt(os.Getenv("POKT_TX_FEE"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_TX_FEE: ", err.Error())
		} else {
			Config.Pocket.TxFee = txFee
		}
	}
	if os.Getenv("POKT_VAULT_ADDRESS") != "" {
		Config.Pocket.VaultAddress = os.Getenv("POKT_VAULT_ADDRESS")
	}
	if Config.Pocket.MultisigPublicKeys == nil || len(Config.Pocket.MultisigPublicKeys) == 0 {
		multisigPublicKeys := os.Getenv("POKT_MULTISIG_PUBLIC_KEYS")
		Config.Pocket.MultisigPublicKeys = strings.Split(multisigPublicKeys, ",")
	}

	// mint monitor
	if os.Getenv("MINT_MONITOR_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("MINT_MONITOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_MONITOR_ENABLED: ", err.Error())
		} else {
			Config.MintMonitor.Enabled = enabled
		}
	}
	if os.Getenv("MINT_MONITOR_INTERVAL_SECS") != "" {
		intervalSecs, err := strconv.ParseInt(os.Getenv("MINT_MONITOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_MONITOR_INTERVAL_SECS: ", err.Error())
		} else {
			Config.MintMonitor.IntervalSecs = intervalSecs
		}
	}

	// mint signer
	if os.Getenv("MINT_SIGNER_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("MINT_SIGNER_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_SIGNER_ENABLED: ", err.Error())
		} else {
			Config.MintSigner.Enabled = enabled
		}
	}
	if os.Getenv("MINT_SIGNER_INTERVAL_SECS") != "" {
		intervalSecs, err := strconv.ParseInt(os.Getenv("MINT_SIGNER_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_SIGNER_INTERVAL_SECS: ", err.Error())
		} else {
			Config.MintSigner.IntervalSecs = intervalSecs
		}
	}

	// mint executor
	if os.Getenv("MINT_EXECUTOR_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("MINT_EXECUTOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_EXECUTOR_ENABLED: ", err.Error())
		} else {
			Config.MintExecutor.Enabled = enabled
		}
	}
	if os.Getenv("MINT_EXECUTOR_INTERVAL_SECS") != "" {
		intervalSecs, err := strconv.ParseInt(os.Getenv("MINT_EXECUTOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_EXECUTOR_INTERVAL_SECS: ", err.Error())
		} else {
			Config.MintExecutor.IntervalSecs = intervalSecs
		}
	}

	// burn monitor
	if os.Getenv("BURN_MONITOR_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("BURN_MONITOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_MONITOR_ENABLED: ", err.Error())
		} else {
			Config.BurnMonitor.Enabled = enabled
		}
	}
	if os.Getenv("BURN_MONITOR_INTERVAL_SECS") != "" {
		intervalSecs, err := strconv.ParseInt(os.Getenv("BURN_MONITOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_MONITOR_INTERVAL_SECS: ", err.Error())
		} else {
			Config.BurnMonitor.IntervalSecs = intervalSecs
		}
	}

	// burn signer
	if os.Getenv("BURN_SIGNER_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("BURN_SIGNER_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_SIGNER_ENABLED: ", err.Error())
		} else {
			Config.BurnSigner.Enabled = enabled
		}
	}
	if os.Getenv("BURN_SIGNER_INTERVAL_SECS") != "" {
		intervalSecs, err := strconv.ParseInt(os.Getenv("BURN_SIGNER_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_SIGNER_INTERVAL_SECS: ", err.Error())
		} else {
			Config.BurnSigner.IntervalSecs = intervalSecs
		}
	}

	// burn executor
	if os.Getenv("BURN_EXECUTOR_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("BURN_EXECUTOR_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_EXECUTOR_ENABLED: ", err.Error())
		} else {
			Config.BurnExecutor.Enabled = enabled
		}
	}
	if os.Getenv("BURN_EXECUTOR_INTERVAL_SECS") != "" {
		intervalSecs, err := strconv.ParseInt(os.Getenv("BURN_EXECUTOR_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_EXECUTOR_INTERVAL_SECS: ", err.Error())
		} else {
			Config.BurnExecutor.IntervalSecs = intervalSecs
		}
	}

	// health check
	if Config.HealthCheck.IntervalSecs == 0 {
		intervalSecs, err := strconv.ParseInt(os.Getenv("HEALTH_CHECK_INTERVAL_SECS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing HEALTH_CHECK_INTERVAL_SECS: ", err.Error())
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
