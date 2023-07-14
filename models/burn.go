package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionBurns = "burns"
)

type Burn struct {
	Id               *primitive.ObjectID `bson:"_id,omitempty"`
	TransactionHash  string              `bson:"transaction_hash"`
	LogIndex         string              `bson:"log_index"`
	BlockNumber      string              `bson:"block_number"`
	SenderAddress    string              `bson:"sender_address"`
	SenderChainId    string              `bson:"sender_chain_id"`
	RecipientAddress string              `bson:"recipient_address"`
	RecipientChainId string              `bson:"recipient_chain_id"`
	Amount           string              `bson:"amount"`
	CreatedAt        time.Time           `bson:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at"`
	Status           string              `bson:"status"`
	ReturnTx         string              `bson:"return_tx"`
	Signers          []string            `bson:"signers"`
	ReturnTxHash     string              `bson:"return_tx_hash"`
}
