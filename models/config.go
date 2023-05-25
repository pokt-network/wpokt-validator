package models

type MongoConfig struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	TimeOutSecs int64  `yaml:"timeout_secs"`
}

type PocketConfig struct {
	RPCURL              string `yaml:"rpc_url"`
	RPCTimeOutSecs      int64  `yaml:"rpc_timeout_secs"`
	ChainId             string `yaml:"chain_id"`
	StartHeight         int64  `yaml:"start_height"`
	MonitorIntervalSecs int64  `yaml:"monitor_interval_secs"`
}

type EthereumConfig struct {
	RPCURL         string `yaml:"rpc_url"`
	RPCTimeOutSecs int64  `yaml:"rpc_timeout_secs"`
	ChainId        int64  `yaml:"chain_id"`
}

type CopperConfig struct {
	VaultAddress string `yaml:"vault_address"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

//Config model
type Config struct {
	MongoDB  MongoConfig    `yaml:"mongodb"`
	Pocket   PocketConfig   `yaml:"pocket"`
	Ethereum EthereumConfig `yaml:"ethereum"`
	Copper   CopperConfig   `yaml:"copper"`
	Logger   LoggerConfig   `yaml:"logger"`
}
