package app

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/dan13ram/wpokt-backend/models"
	"gopkg.in/yaml.v2"
)

var (
	Config models.Config
)

func InitConfig(configFile string, envFile string) {
	var yamlFile, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("[CONFIG] Error reading config file %q: %s\n", configFile, err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("[CONFIG] Error unmarshalling config file %q: %s\n", configFile, err.Error())
	}
	readConfigFromEnv(envFile)
	validateConfig()
}

func readConfigFromEnv(envFile string) {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Warn("[CONFIG] Error loading .env file: ", err.Error())
		}
	}
	if Config.MongoDB.URI == "" {
		Config.MongoDB.URI = os.Getenv("MONGODB_URI")
	}
	if Config.MongoDB.Database == "" {
		Config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	}
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
			log.Warn("[CONFIG] Error parsing ETH_START_BLOCK_NUMBER: ", err.Error())
			log.Info("[CONFIG] Setting ETH_START_BLOCK_NUMBER to 0")
		} else {
			Config.Ethereum.StartBlockNumber = blockNumber
		}
	}
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
			log.Warn("[CONFIG] Error parsing POKT_START_HEIGHT: ", err.Error())
			log.Info("[CONFIG] Setting POKT_START_HEIGHT to 0")
		} else {
			Config.Pocket.StartHeight = startHeight
		}
	}
}

func validateConfig() {
	if Config.MongoDB.URI == "" {
		log.Fatal("[CONFIG] MongoDB.URI is required")
	}
	if Config.MongoDB.Database == "" {
		log.Fatal("[CONFIG] MongoDB.Database is required")
	}
	if Config.Ethereum.RPCURL == "" {
		log.Fatal("[CONFIG] Ethereum.RPCURL is required")
	}
	if Config.Ethereum.ChainId == "" {
		log.Fatal("[CONFIG] Ethereum.ChainId is required")
	}
	if Config.Ethereum.PrivateKey == "" {
		log.Fatal("[CONFIG] Ethereum.PrivateKey is required")
	}
	if Config.Ethereum.StartBlockNumber == 0 {
		log.Fatal("[CONFIG] Ethereum.StartBlockNumber is required")
	}
	if Config.Pocket.RPCURL == "" {
		log.Fatal("[CONFIG] Pocket.RPCURL is required")
	}
	if Config.Pocket.ChainId == "" {
		log.Fatal("[CONFIG] Pocket.ChainId is required")
	}
	if Config.Pocket.PrivateKey == "" {
		log.Fatal("[CONFIG] Pocket.PrivateKey is required")
	}
	if Config.Pocket.StartHeight == 0 {
		log.Fatal("[CONFIG] Pocket.StartHeight is required")
	}
}
