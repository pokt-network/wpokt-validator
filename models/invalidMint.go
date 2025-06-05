package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionInvalidMints = "shannonInvalidMints"
)

type InvalidMint struct {
	Id              *primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TransactionHash string              `bson:"transaction_hash" json:"transaction_hash"`
	Height          string              `bson:"height" json:"height"`
	Confirmations   string              `bson:"confirmations" json:"confirmations"`
	SenderAddress   string              `bson:"sender_address" json:"sender_address"`
	SenderChainID   string              `bson:"sender_chain_id" json:"sender_chain_id"`
	VaultAddress    string              `bson:"vault_address" json:"vault_address"`
	Amount          string              `bson:"amount" json:"amount"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at" json:"updated_at"`
	Status          string              `bson:"status" json:"status"`
	Memo            string              `bson:"memo" json:"memo"`

	ReturnTransactionBody string      `json:"return_transaction_body" bson:"transaction_body"`
	Signatures            []Signature `json:"signatures" bson:"signatures"`
	Sequence              *uint64     `json:"sequence" bson:"sequence"` // account sequence for submitting the transaction
	ReturnTransactionHash string      `json:"return_transaction_hash" bson:"transaction_hash"`
}
