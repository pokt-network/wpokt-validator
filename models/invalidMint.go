package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionInvalidMints = "invalidMints"
)

type InvalidMint struct {
	Id              *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionHash string              `bson:"transaction_hash" json:"transaction_hash"`
	Height          string              `bson:"height" json:"height"`
	Confirmations   string              `bson:"confirmations" json:"confirmations"`
	SenderAddress   string              `bson:"sender_address" json:"sender_address"`
	SenderChainId   string              `bson:"sender_chain_id" json:"sender_chain_id"`
	VaultAddress    string              `bson:"vault_address" json:"vault_address"`
	Amount          string              `bson:"amount" json:"amount"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at" json:"updated_at"`
	Status          string              `bson:"status" json:"status"`
	ReturnTx        string              `bson:"return_tx" json:"return_tx"`
	Signers         []string            `bson:"signers" json:"signers"`
	ReturnTxHash    string              `bson:"return_tx_hash" json:"return_tx_hash"`
	Memo            string              `bson:"memo" json:"memo"`
}
