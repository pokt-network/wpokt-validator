package models

type MongoConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type PocketConfig struct {
	RPCURL string `yaml:"rpc_url"`
}

type EthereumConfig struct {
	RPCURL string `yaml:"rpc_url"`
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
