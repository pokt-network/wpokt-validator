package models

type ChainType string

const (
	ChainTypeEthereum ChainType = "ethereum"
	ChainTypeCosmos   ChainType = "cosmos"
)

type Chain struct {
	ChainID   string    `bson:"chain_id" json:"chain_id"`
	ChainType ChainType `bson:"chain_type" json:"chain_type"`
}
