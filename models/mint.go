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
	Height           string              `bson:"height"`
	SenderAddress    string              `bson:"sender_address"`
	SenderChainId    string              `bson:"sender_chain_id"`
	RecipientAddress string              `bson:"recipient_address"`
	RecipientChainId string              `bson:"recipient_chain_id"`
	Amount           string              `bson:"amount"`
	CreatedAt        time.Time           `bson:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at"`
	Status           string              `bson:"status"`
	Data             *MintData           `bson:"data"`
	Signers          []string            `bson:"signers"`
	Signatures       []string            `bson:"signatures"`
}

type MintMemo struct {
	Address string `json:"address"`
	ChainId string `json:"chain_id"`
}

type MintData struct {
	Recipient string `json:"recipient" bson:"recipient"`
	Amount    string `json:"amount" bson:"amount"`
	Nonce     string `json:"nonce" bson:"nonce"`
}
