package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionHealthChecks = "healthchecks"
)

type Health struct {
	Id            *primitive.ObjectID `bson:"_id,omitempty"`
	PoktPublicKey string              `bson:"pokt_public_key"`
	EthAddress    string              `bson:"eth_address"`
	Hostname      string              `bson:"hostname"`
	CreatedAt     time.Time           `bson:"created_at"`
}
