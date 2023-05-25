package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionMints = "mints"
)

type Mint struct {
	Id               *primitive.ObjectID `bson:"_id,omitempty"`
	TransactionHash  string              `bson:"transaction_hash"`
	Height           int64               `bson:"height"`
	SenderAddress    string              `bson:"sender_address"`
	SenderChainId    string              `bson:"sender_chain_id"`
	RecipientAddress string              `bson:"recipient_address"`
	RecipientChainId int64               `bson:"recipient_chain_id"`
	Amount           string              `bson:"amount"`
	CreatedAt        time.Time           `bson:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at"`
	Status           string              `bson:"status"`
	Signers          []string            `bson:"signers"`
	Order            *Order              `bson:"order"`
}

type MintMemo struct {
	Address string `json:"address"`
	ChainId int64  `json:"chain_id"`
}
