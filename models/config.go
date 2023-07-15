package models

type Config struct {
	Health        HealthConfig        `yaml:"health"`
	Logger        LoggerConfig        `yaml:"logger"`
	MongoDB       MongoConfig         `yaml:"mongodb"`
	Ethereum      EthereumConfig      `yaml:"ethereum"`
	WPOKTMonitor  WPOKTMonitorConfig  `yaml:"wpokt_monitor"`
	WPOKTSigner   WPOKTSignerConfig   `yaml:"wpokt_signer"`
	WPOKTExecutor WPOKTExecutorConfig `yaml:"wpokt_executor"`
	Pocket        PocketConfig        `yaml:"pocket"`
	PoktMonitor   PoktMonitorConfig   `yaml:"pokt_monitor"`
	PoktSigner    PoktSignerConfig    `yaml:"pokt_signer"`
	PoktExecutor  PoktExecutorConfig  `yaml:"pokt_executor"`
}

type HealthConfig struct {
	IntervalSecs uint64 `yaml:"interval_secs"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type MongoConfig struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	TimeOutSecs uint64 `yaml:"timeout_secs"`
}

type EthereumConfig struct {
	RPCURL                string   `yaml:"rpc_url"`
	RPCTimeOutSecs        uint64   `yaml:"rpc_timeout_secs"`
	ChainId               string   `yaml:"chain_id"`
	WPOKTContractAddress  string   `yaml:"wpokt_contract_address"`
	MintControllerAddress string   `yaml:"mint_controller_address"`
	ValidatorAddresses    []string `yaml:"validator_addresses"`
}

type WPOKTMonitorConfig struct {
	Enabled          bool   `yaml:"enabled"`
	StartBlockNumber int64  `yaml:"start_block_number"`
	IntervalSecs     uint64 `yaml:"interval_secs"`
}

type WPOKTSignerConfig struct {
	Enabled      bool   `yaml:"enabled"`
	IntervalSecs uint64 `yaml:"interval_secs"`
	PrivateKey   string `yaml:"private_key"`
}

type WPOKTExecutorConfig struct {
	Enabled          bool   `yaml:"enabled"`
	StartBlockNumber int64  `yaml:"start_block_number"`
	IntervalSecs     uint64 `yaml:"interval_secs"`
}

type PocketConfig struct {
	RPCURL             string   `yaml:"rpc_url"`
	RPCTimeOutSecs     uint64   `yaml:"rpc_timeout_secs"`
	ChainId            string   `yaml:"chain_id"`
	Fees               int64    `yaml:"fees"`
	MultisigPublicKeys []string `yaml:"multisig_public_keys"`
}

type PoktMonitorConfig struct {
	Enabled      bool   `yaml:"enabled"`
	StartHeight  int64  `yaml:"start_height"`
	IntervalSecs uint64 `yaml:"interval_secs"`
}

type PoktSignerConfig struct {
	Enabled      bool   `yaml:"enabled"`
	IntervalSecs uint64 `yaml:"interval_secs"`
	PrivateKey   string `yaml:"private_key"`
}

type PoktExecutorConfig struct {
	Enabled      bool   `yaml:"enabled"`
	IntervalSecs uint64 `yaml:"interval_secs"`
}
