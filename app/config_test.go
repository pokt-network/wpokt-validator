package app

import (
	"fmt"
	"io"
	"testing"

	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestReadConfigFromConfigFile(t *testing.T) {
	t.Run("Config File Provided", func(t *testing.T) {
		configFile := "../config/config.sample.yml"

		read := readConfigFromConfigFile(configFile)

		assert.Equal(t, read, true)
		assert.Equal(t, Config.MongoDB.Database, "mongodb-database")
		assert.Equal(t, Config.MongoDB.TimeoutMillis, int64(2000))
	})

	t.Run("No Config File Provided", func(t *testing.T) {
		configFile := ""

		read := readConfigFromConfigFile(configFile)
		assert.Equal(t, read, false)
	})

	t.Run("Invalid Config File Path", func(t *testing.T) {
		configFile := "../config.sample.invalid.yml"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { readConfigFromConfigFile(configFile) })
	})

	t.Run("Invalid Config File Contents", func(t *testing.T) {
		configFile := "../README.md"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { readConfigFromConfigFile(configFile) })
	})

}

func TestInitConfig(t *testing.T) {
	t.Run("Config Initialization Success", func(t *testing.T) {
		configFile := "../config/config.sample.yml"
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

	})

	t.Run("Config Initialization No Config File", func(t *testing.T) {
		configFile := ""
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {

		configFile := "../config/config.sample.yml"
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

		validateConfig()

	})

	t.Run("Without MongoDB URI", func(t *testing.T) {
		Config = models.Config{}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })

	})

	t.Run("Without MongoDB Database", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })

	})

	t.Run("Without MongoDB Timeout", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 0

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth RPC URL", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth ChainID", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth RPC Timeout", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth Private Key", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth wPOKT Address", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth Mint Controller Address", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Eth Validator Addresses", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt RPC URL", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt ChainID", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt RPC Timeout", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt Private Key", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt Tx Fee", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt Vault Address", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without Pokt Multisig PublicKeys", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without MintMonitor Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without MintSigner Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true
		Config.MintMonitor.IntervalMillis = 1000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without MintExecutor Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true
		Config.MintMonitor.IntervalMillis = 1000
		Config.MintSigner.IntervalMillis = 1000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without BurnMonitor Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true
		Config.MintMonitor.IntervalMillis = 1000
		Config.MintSigner.IntervalMillis = 1000
		Config.MintExecutor.IntervalMillis = 1000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without BurnSigner Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true
		Config.MintMonitor.IntervalMillis = 1000
		Config.MintSigner.IntervalMillis = 1000
		Config.MintExecutor.IntervalMillis = 1000
		Config.BurnMonitor.IntervalMillis = 1000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without BurnExecutor Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true
		Config.MintMonitor.IntervalMillis = 1000
		Config.MintSigner.IntervalMillis = 1000
		Config.MintExecutor.IntervalMillis = 1000
		Config.BurnMonitor.IntervalMillis = 1000
		Config.BurnSigner.IntervalMillis = 1000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

	t.Run("Without HealthCheck Interval", func(t *testing.T) {
		Config = models.Config{}
		Config.MongoDB.URI = "mongodb://localhost:27017"
		Config.MongoDB.Database = "mongodb-database"
		Config.MongoDB.TimeoutMillis = 2000
		Config.Ethereum.RPCURL = "http://localhost:8545"
		Config.Ethereum.ChainID = "31337"
		Config.Ethereum.RPCTimeoutMillis = 2000
		Config.Ethereum.PrivateKey = "abcd"
		Config.Ethereum.WrappedPocketAddress = "0x1234"
		Config.Ethereum.MintControllerAddress = "0x1234"
		Config.Ethereum.ValidatorAddresses = []string{"0x1234"}
		Config.Pocket.RPCURL = "http://localhost:8081"
		Config.Pocket.ChainID = "localnet"
		Config.Pocket.RPCTimeoutMillis = 2000
		Config.Pocket.Mnemonic = "abcd"
		Config.Pocket.TxFee = 10000
		Config.Pocket.MultisigAddress = "0x1234"
		Config.Pocket.MultisigPublicKeys = []string{"1234"}
		Config.MintMonitor.Enabled = true
		Config.MintSigner.Enabled = true
		Config.MintExecutor.Enabled = true
		Config.BurnMonitor.Enabled = true
		Config.BurnSigner.Enabled = true
		Config.BurnExecutor.Enabled = true
		Config.MintMonitor.IntervalMillis = 1000
		Config.MintSigner.IntervalMillis = 1000
		Config.MintExecutor.IntervalMillis = 1000
		Config.BurnMonitor.IntervalMillis = 1000
		Config.BurnSigner.IntervalMillis = 1000
		Config.BurnExecutor.IntervalMillis = 1000

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() })
	})

}
