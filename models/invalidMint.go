package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionInvalidMints = "invalid_mints"
)

type InvalidMint struct {
	Id              *primitive.ObjectID `bson:"_id,omitempty"`
	TransactionHash string              `bson:"transaction_hash"`
	Height          string              `bson:"height"`
	Confirmations   string              `bson:"confirmations"`
	SenderAddress   string              `bson:"sender_address"`
	SenderChainId   string              `bson:"sender_chain_id"`
	Amount          string              `bson:"amount"`
	CreatedAt       time.Time           `bson:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at"`
	Status          string              `bson:"status"`
	ReturnTx        string              `bson:"return_tx"`
	Signers         []string            `bson:"signers"`
	ReturnTxHash    string              `bson:"return_tx_hash"`
}
