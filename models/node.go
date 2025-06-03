package models

type ChainType string

const (
	ChainTypeEthereum ChainType = "ethereum"
	ChainTypeCosmos   ChainType = "cosmos"
)

type Chain struct {
	ChainID     string    `bson:"chain_id" json:"chain_id"`
	ChainName   string    `bson:"chain_name" json:"chain_name"`
	ChainDomain uint32    `bson:"chain_domain" json:"chain_domain"`
	ChainType   ChainType `bson:"chain_type" json:"chain_type"`
}
