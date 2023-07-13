package models

type Config struct {
	Logger      LoggerConfig      `yaml:"logger"`
	MongoDB     MongoConfig       `yaml:"mongodb"`
	Ethereum    EthereumConfig    `yaml:"ethereum"`
	PoktMonitor PoktMonitorConfig `yaml:"pokt_monitor"`
	PoktSigner  PoktSignerConfig  `yaml:"pokt_signer"`
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
	RPCURL               string `yaml:"rpc_url"`
	RPCTimeOutSecs       uint64 `yaml:"rpc_timeout_secs"`
	ChainId              uint64 `yaml:"chain_id"`
	StartBlockNumber     int64  `yaml:"start_block_number"`
	WPOKTContractAddress string `yaml:"wpokt_contract_address"`
	MonitorIntervalSecs  uint64 `yaml:"monitor_interval_secs"`
	SignerIntervalSecs   uint64 `yaml:"signer_interval_secs"`
}

type PoktMonitorConfig struct {
	Enabled         bool   `yaml:"enabled"`
	RPCURL          string `yaml:"rpc_url"`
	RPCTimeOutSecs  uint64 `yaml:"rpc_timeout_secs"`
	ChainId         string `yaml:"chain_id"`
	StartHeight     int64  `yaml:"start_height"`
	IntervalSecs    uint64 `yaml:"monitor_interval_secs"`
	MultisigAddress string `yaml:"multisig_address"`
}

type PoktSignerConfig struct {
	Enabled            bool     `yaml:"enabled"`
	IntervalSecs       uint64   `yaml:"interval_secs"`
	PrivateKey         string   `yaml:"private_key"`
	MultisigAddress    string   `yaml:"multisig_address"`
	MultisigPublicKeys []string `yaml:"multisig_public_keys"`
}
