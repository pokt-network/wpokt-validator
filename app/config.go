package app

import (
	"io/ioutil"
	"os"

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
	if Config.WPOKTSigner.PrivateKey == "" {
		Config.WPOKTSigner.PrivateKey = os.Getenv("ETH_PRIVATE_KEY")
	}
	if Config.Pocket.RPCURL == "" {
		Config.Pocket.RPCURL = os.Getenv("POKT_RPC_URL")
	}
	if Config.Pocket.ChainId == "" {
		Config.Pocket.ChainId = os.Getenv("POKT_CHAIN_ID")
	}
	if Config.PoktSigner.PrivateKey == "" {
		Config.PoktSigner.PrivateKey = os.Getenv("POKT_PRIVATE_KEY")
	}
}
