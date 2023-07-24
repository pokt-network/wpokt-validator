package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionHealthChecks = "healthchecks"
)

type Health struct {
	Id               *primitive.ObjectID `bson:"_id,omitempty"`
	PoktVaultAddress string              `bson:"pokt_vault_address"`
	PoktSigners      []string            `bson:"pokt_signers"`
	PoktPublicKey    string              `bson:"pokt_public_key"`
	PoktAddress      string              `bson:"pokt_address"`
	EthValidators    []string            `bson:"eth_validators"`
	EthAddress       string              `bson:"eth_address"`
	Hostname         string              `bson:"hostname"`
	Healthy          bool                `bson:"healthy"`
	CreatedAt        time.Time           `bson:"created_at"`
	ServiceHealths   []ServiceHealth     `bson:"service_healths"`
}

type ServiceHealth struct {
	Name           string    `bson:"name"`
	Healthy        bool      `bson:"healthy"`
	EthBlockNumber string    `bson:"eth_block_number"` // not used for all services
	PoktHeight     string    `bson:"pokt_height"`      // not used for all services
	LastSyncTime   time.Time `bson:"last_sync_time"`
	NextSyncTime   time.Time `bson:"next_sync_time"`
}
