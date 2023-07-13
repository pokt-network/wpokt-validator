package models

type MongoConfig struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	TimeOutSecs uint64 `yaml:"timeout_secs"`
}

type PocketConfig struct {
	RPCURL              string `yaml:"rpc_url"`
	RPCTimeOutSecs      uint64 `yaml:"rpc_timeout_secs"`
	ChainId             string `yaml:"chain_id"`
	StartHeight         int64  `yaml:"start_height"`
	MonitorIntervalSecs uint64 `yaml:"monitor_interval_secs"`
	VaultAddress        string `yaml:"vault_address"`
}

type EthereumConfig struct {
	RPCURL               string `yaml:"rpc_url"`
	RPCTimeOutSecs       uint64 `yaml:"rpc_timeout_secs"`
	ChainId              uint64 `yaml:"chain_id"`
	StartBlockNumber     int64  `yaml:"start_block_number"`
	WPOKTContractAddress string `yaml:"wpokt_contract_address"`
	MonitorIntervalSecs  uint64 `yaml:"monitor_interval_secs"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	MongoDB  MongoConfig    `yaml:"mongodb"`
	Pocket   PocketConfig   `yaml:"pocket"`
	Ethereum EthereumConfig `yaml:"ethereum"`
	Logger   LoggerConfig   `yaml:"logger"`
}
