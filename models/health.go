package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionHealthChecks = "healthchecks"
)

type Health struct {
	Id               *primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PoktVaultAddress string              `bson:"pokt_vault_address" json:"pokt_vault_address"`
	PoktSigners      []string            `bson:"pokt_signers" json:"pokt_signers"`
	PoktPublicKey    string              `bson:"pokt_public_key" json:"pokt_public_key"`
	PoktAddress      string              `bson:"pokt_address" json:"pokt_address"`
	EthValidators    []string            `bson:"eth_validators" json:"eth_validators"`
	EthAddress       string              `bson:"eth_address" json:"eth_address"`
	WPoktAddress     string              `bson:"wpokt_address" json:"wpokt_address"`
	Hostname         string              `bson:"hostname" json:"hostname"`
	ValidatorId      string              `bson:"validator_id" json:"validator_id"`
	Healthy          bool                `bson:"healthy" json:"healthy"`
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`
	ServiceHealths   []ServiceHealth     `bson:"service_healths" json:"service_healths"`
}

type ServiceHealth struct {
	Name           string    `bson:"name" json:"name"`
	Healthy        bool      `bson:"healthy" json:"healthy"`
	EthBlockNumber string    `bson:"eth_block_number" json:"eth_block_number"` // not used for all services
	PoktHeight     string    `bson:"pokt_height" json:"pokt_height"`           // not used for all services
	LastSyncTime   time.Time `bson:"last_sync_time" json:"last_sync_time"`
	NextSyncTime   time.Time `bson:"next_sync_time" json:"next_sync_time"`
}
