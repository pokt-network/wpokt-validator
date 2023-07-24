package models

type Config struct {
	Health        HealthConfig   `yaml:"health"`
	Logger        LoggerConfig   `yaml:"logger"`
	MongoDB       MongoConfig    `yaml:"mongodb"`
	Ethereum      EthereumConfig `yaml:"ethereum"`
	WPOKTMonitor  ServiceConfig  `yaml:"wpokt_monitor"`
	WPOKTSigner   ServiceConfig  `yaml:"wpokt_signer"`
	WPOKTExecutor ServiceConfig  `yaml:"wpokt_executor"`
	Pocket        PocketConfig   `yaml:"pocket"`
	PoktMonitor   ServiceConfig  `yaml:"pokt_monitor"`
	PoktSigner    ServiceConfig  `yaml:"pokt_signer"`
	PoktExecutor  ServiceConfig  `yaml:"pokt_executor"`
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
	StartBlockNumber      int64    `yaml:"start_block_number"`
	Confirmations         int64    `yaml:"confirmations"`
	PrivateKey            string   `yaml:"private_key"`
	RPCURL                string   `yaml:"rpc_url"`
	RPCTimeOutSecs        uint64   `yaml:"rpc_timeout_secs"`
	ChainId               string   `yaml:"chain_id"`
	WPOKTContractAddress  string   `yaml:"wpokt_contract_address"`
	MintControllerAddress string   `yaml:"mint_controller_address"`
	ValidatorAddresses    []string `yaml:"validator_addresses"`
}

type PocketConfig struct {
	StartHeight        int64    `yaml:"start_height"`
	Confirmations      int64    `yaml:"confirmations"`
	RPCURL             string   `yaml:"rpc_url"`
	PrivateKey         string   `yaml:"private_key"`
	RPCTimeOutSecs     uint64   `yaml:"rpc_timeout_secs"`
	ChainId            string   `yaml:"chain_id"`
	Fees               int64    `yaml:"fees"`
	MultisigPublicKeys []string `yaml:"multisig_public_keys"`
}

type ServiceConfig struct {
	Enabled      bool   `yaml:"enabled"`
	IntervalSecs uint64 `yaml:"interval_secs"`
}
