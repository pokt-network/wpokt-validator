package models

type Config struct {
	GoogleSecretManager GoogleSecretManagerConfig `yaml:"google_secret_manager" json:"google_secret_manager"`
	HealthCheck         HealthCheckConfig         `yaml:"health_check" json:"health_check"`
	Logger              LoggerConfig              `yaml:"logger" json:"logger"`
	MongoDB             MongoConfig               `yaml:"mongodb" json:"mongo_db"`
	Ethereum            EthereumConfig            `yaml:"ethereum" json:"ethereum"`
	Pocket              PocketConfig              `yaml:"pocket" json:"pocket"`
	MintMonitor         ServiceConfig             `yaml:"mint_monitor" json:"mint_monitor"`
	MintSigner          ServiceConfig             `yaml:"mint_signer" json:"mint_signer"`
	MintExecutor        ServiceConfig             `yaml:"mint_executor" json:"mint_executor"`
	BurnMonitor         ServiceConfig             `yaml:"burn_monitor" json:"burn_monitor"`
	BurnSigner          ServiceConfig             `yaml:"burn_signer" json:"burn_signer"`
	BurnExecutor        ServiceConfig             `yaml:"burn_executor" json:"burn_executor"`
}

type GoogleSecretManagerConfig struct {
	Enabled         bool   `yaml:"enabled" json:"enabled"`
	MongoSecretName string `yaml:"mongo_secret_name" json:"mongo_secret_name"`
	PoktSecretName  string `yaml:"pokt_secret_name" json:"pokt_secret_name"`
	EthSecretName   string `yaml:"eth_secret_name" json:"eth_secret_name"`
}

type HealthCheckConfig struct {
	IntervalSecs   int64 `yaml:"interval_secs" json:"interval_secs"`
	ReadLastHealth bool  `yaml:"read_last_health" json:"read_last_health"`
}

type LoggerConfig struct {
	Level string `yaml:"level" json:"level"`
}

type MongoConfig struct {
	URI         string `yaml:"uri" json:"uri"`
	Database    string `yaml:"database" json:"database"`
	TimeoutSecs int64  `yaml:"timeout_secs" json:"timeout_secs"`
}

type EthereumConfig struct {
	StartBlockNumber      int64    `yaml:"start_block_number" json:"start_block_number"`
	Confirmations         int64    `yaml:"confirmations" json:"confirmations"`
	PrivateKey            string   `yaml:"private_key" json:"private_key"`
	RPCURL                string   `yaml:"rpc_url" json:"rpcurl"`
	RPCTimeoutSecs        int64    `yaml:"rpc_timeout_secs" json:"rpc_time_out_secs"`
	ChainId               string   `yaml:"chain_id" json:"chain_id"`
	WrappedPocketAddress  string   `yaml:"wrapped_pocket_address" json:"wrapped_pocket_address"`
	MintControllerAddress string   `yaml:"mint_controller_address" json:"mint_controller_address"`
	ValidatorAddresses    []string `yaml:"validator_addresses" json:"validator_addresses"`
}

type PocketConfig struct {
	StartHeight        int64    `yaml:"start_height" json:"start_height"`
	Confirmations      int64    `yaml:"confirmations" json:"confirmations"`
	RPCURL             string   `yaml:"rpc_url" json:"rpcurl"`
	PrivateKey         string   `yaml:"private_key" json:"private_key"`
	RPCTimeoutSecs     int64    `yaml:"rpc_timeout_secs" json:"rpc_time_out_secs"`
	ChainId            string   `yaml:"chain_id" json:"chain_id"`
	TxFee              int64    `yaml:"tx_fee" json:"tx_fee"`
	VaultAddress       string   `yaml:"vault_address" json:"vault_address"`
	MultisigPublicKeys []string `yaml:"multisig_public_keys" json:"multisig_public_keys"`
}

type ServiceConfig struct {
	Enabled      bool  `yaml:"enabled" json:"enabled"`
	IntervalSecs int64 `yaml:"interval_secs" json:"interval_secs"`
}
