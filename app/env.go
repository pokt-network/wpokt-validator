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
		} else {
			log.Debug("[ENV] .env file loaded from: ", envFile)
		}
	} else {
		log.Debug("[ENV] No .env file provided")
	}

	log.Debug("[ENV] Reading config from ENV variables")

	if os.Getenv("MONGODB_URI") != "" {
		Config.MongoDB.URI = os.Getenv("MONGODB_URI")
	}
	if os.Getenv("MONGODB_DATABASE") != "" {
		Config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	}
	if os.Getenv("MONGODB_TIMEOUT_MS") != "" {
		timeoutMillis, err := strconv.ParseInt(os.Getenv("MONGODB_TIMEOUT_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MONGODB_TIMEOUT_MS: ", err.Error())
		} else {
			Config.MongoDB.TimeoutMillis = timeoutMillis
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
	if os.Getenv("ETH_RPC_TIMEOUT_MS") != "" {
		timeoutMillis, err := strconv.ParseInt(os.Getenv("ETH_RPC_TIMEOUT_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing ETH_RPC_TIMEOUT_MS: ", err.Error())
		} else {
			Config.Ethereum.RPCTimeoutMillis = timeoutMillis
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
	// TODO: fix env support for latest config
	// if os.Getenv("POKT_PRIVATE_KEY") != "" {
	// 	Config.Pocket.PrivateKey = os.Getenv("POKT_PRIVATE_KEY")
	// }
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
	if os.Getenv("POKT_RPC_TIMEOUT_MS") != "" {
		timeoutMillis, err := strconv.ParseInt(os.Getenv("POKT_RPC_TIMEOUT_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing POKT_RPC_TIMEOUT_MS: ", err.Error())
		} else {
			Config.Pocket.RPCTimeoutMillis = timeoutMillis
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
	// if os.Getenv("POKT_VAULT_ADDRESS") != "" {
	// 	Config.Pocket.VaultAddress = os.Getenv("POKT_VAULT_ADDRESS")
	// }
	if os.Getenv("POKT_MULTISIG_PUBLIC_KEYS") != "" {
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
	if os.Getenv("MINT_MONITOR_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("MINT_MONITOR_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_MONITOR_INTERVAL_MS: ", err.Error())
		} else {
			Config.MintMonitor.IntervalMillis = intervalMillis
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
	if os.Getenv("MINT_SIGNER_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("MINT_SIGNER_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_SIGNER_INTERVAL_MS: ", err.Error())
		} else {
			Config.MintSigner.IntervalMillis = intervalMillis
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
	if os.Getenv("MINT_EXECUTOR_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("MINT_EXECUTOR_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing MINT_EXECUTOR_INTERVAL_MS: ", err.Error())
		} else {
			Config.MintExecutor.IntervalMillis = intervalMillis
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
	if os.Getenv("BURN_MONITOR_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("BURN_MONITOR_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_MONITOR_INTERVAL_MS: ", err.Error())
		} else {
			Config.BurnMonitor.IntervalMillis = intervalMillis
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
	if os.Getenv("BURN_SIGNER_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("BURN_SIGNER_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_SIGNER_INTERVAL_MS: ", err.Error())
		} else {
			Config.BurnSigner.IntervalMillis = intervalMillis
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
	if os.Getenv("BURN_EXECUTOR_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("BURN_EXECUTOR_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing BURN_EXECUTOR_INTERVAL_MS: ", err.Error())
		} else {
			Config.BurnExecutor.IntervalMillis = intervalMillis
		}
	}

	// health check
	if os.Getenv("HEALTH_CHECK_INTERVAL_MS") != "" {
		intervalMillis, err := strconv.ParseInt(os.Getenv("HEALTH_CHECK_INTERVAL_MS"), 10, 64)
		if err != nil {
			log.Warn("[ENV] Error parsing HEALTH_CHECK_INTERVAL_MS: ", err.Error())
		} else {
			Config.HealthCheck.IntervalMillis = intervalMillis
		}
	}
	if os.Getenv("HEALTH_CHECK_READ_LAST_HEALTH") != "" {
		readLastHealth, err := strconv.ParseBool(os.Getenv("HEALTH_CHECK_READ_LAST_HEALTH"))
		if err != nil {
			log.Warn("[ENV] Error parsing HEALTH_CHECK_READ_LAST_HEALTH: ", err.Error())
		} else {
			Config.HealthCheck.ReadLastHealth = readLastHealth
		}
	}

	// logging
	if os.Getenv("LOG_LEVEL") != "" {
		Config.Logger.Level = os.Getenv("LOG_LEVEL")
	}

	// google secret manager
	if os.Getenv("GOOGLE_SECRET_MANAGER_ENABLED") != "" {
		enabled, err := strconv.ParseBool(os.Getenv("GOOGLE_SECRET_MANAGER_ENABLED"))
		if err != nil {
			log.Warn("[ENV] Error parsing GOOGLE_SECRET_MANAGER_ENABLED: ", err.Error())
		} else {
			Config.GoogleSecretManager.Enabled = enabled
		}
	}
	if os.Getenv("GOOGLE_MONGO_SECRET_NAME") != "" {
		Config.GoogleSecretManager.MongoSecretName = os.Getenv("GOOGLE_MONGO_SECRET_NAME")
	}
	if os.Getenv("GOOGLE_POKT_SECRET_NAME") != "" {
		Config.GoogleSecretManager.PoktSecretName = os.Getenv("GOOGLE_POKT_SECRET_NAME")
	}
	if os.Getenv("GOOGLE_ETH_SECRET_NAME") != "" {
		Config.GoogleSecretManager.EthSecretName = os.Getenv("GOOGLE_ETH_SECRET_NAME")
	}

	log.Debug("[ENV] Config read from env variables")
}
