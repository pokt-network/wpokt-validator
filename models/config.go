package models

type Config struct {
	Health       HealthConfig   `yaml:"health" json:"health"`
	Logger       LoggerConfig   `yaml:"logger" json:"logger"`
	MongoDB      MongoConfig    `yaml:"mongodb" json:"mongo_db"`
	Ethereum     EthereumConfig `yaml:"ethereum" json:"ethereum"`
	Pocket       PocketConfig   `yaml:"pocket" json:"pocket"`
	MintMonitor  ServiceConfig  `yaml:"mint_monitor" json:"mint_monitor"`
	MintSigner   ServiceConfig  `yaml:"mint_signer" json:"mint_signer"`
	MintExecutor ServiceConfig  `yaml:"mint_executor" json:"mint_executor"`
	BurnMonitor  ServiceConfig  `yaml:"burn_monitor" json:"burn_monitor"`
	BurnSigner   ServiceConfig  `yaml:"burn_signer" json:"burn_signer"`
	BurnExecutor ServiceConfig  `yaml:"burn_executor" json:"burn_executor"`
}

type HealthConfig struct {
	IntervalSecs uint64 `yaml:"interval_secs" json:"interval_secs"`
}

type LoggerConfig struct {
	Level string `yaml:"level" json:"level"`
}

type MongoConfig struct {
	URI         string `yaml:"uri" json:"uri"`
	Database    string `yaml:"database" json:"database"`
	TimeOutSecs uint64 `yaml:"timeout_secs" json:"time_out_secs"`
}

type EthereumConfig struct {
	StartBlockNumber      int64    `yaml:"start_block_number" json:"start_block_number"`
	Confirmations         int64    `yaml:"confirmations" json:"confirmations"`
	PrivateKey            string   `yaml:"private_key" json:"private_key"`
	RPCURL                string   `yaml:"rpc_url" json:"rpcurl"`
	RPCTimeOutSecs        uint64   `yaml:"rpc_timeout_secs" json:"rpc_time_out_secs"`
	ChainId               string   `yaml:"chain_id" json:"chain_id"`
	WPOKTAddress          string   `yaml:"wpokt_address" json:"wpokt_address"`
	MintControllerAddress string   `yaml:"mint_controller_address" json:"mint_controller_address"`
	ValidatorAddresses    []string `yaml:"validator_addresses" json:"validator_addresses"`
}

type PocketConfig struct {
	StartHeight        int64    `yaml:"start_height" json:"start_height"`
	Confirmations      int64    `yaml:"confirmations" json:"confirmations"`
	RPCURL             string   `yaml:"rpc_url" json:"rpcurl"`
	PrivateKey         string   `yaml:"private_key" json:"private_key"`
	RPCTimeOutSecs     uint64   `yaml:"rpc_timeout_secs" json:"rpc_time_out_secs"`
	ChainId            string   `yaml:"chain_id" json:"chain_id"`
	Fees               int64    `yaml:"fees" json:"fees"`
	VaultAddress       string   `yaml:"vault_address" json:"vault_address"`
	MultisigPublicKeys []string `yaml:"multisig_public_keys" json:"multisig_public_keys"`
}

type ServiceConfig struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	IntervalSecs uint64 `yaml:"interval_secs" json:"interval_secs"`
}
