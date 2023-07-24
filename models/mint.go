package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionMints = "mints"
)

type Mint struct {
	Id                  *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionHash     string              `bson:"transaction_hash" json:"transaction_hash"`
	Height              string              `bson:"height" json:"height"`
	Confirmations       string              `bson:"confirmations" json:"confirmations"`
	SenderAddress       string              `bson:"sender_address" json:"sender_address"`
	SenderChainId       string              `bson:"sender_chain_id" json:"sender_chain_id"`
	RecipientAddress    string              `bson:"recipient_address" json:"recipient_address"`
	RecipientChainId    string              `bson:"recipient_chain_id" json:"recipient_chain_id"`
	WPOKTAddress        string              `bson:"wpokt_address" json:"wpokt_address"`
	VaultAddress        string              `bson:"vault_address" json:"vault_address"`
	Amount              string              `bson:"amount" json:"amount"`
	Nonce               string              `bson:"nonce" json:"nonce"`
	Memo                *MintMemo           `bson:"memo" json:"memo"`
	CreatedAt           time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time           `bson:"updated_at" json:"updated_at"`
	Status              string              `bson:"status" json:"status"`
	Data                *MintData           `bson:"data" json:"data"`
	Signers             []string            `bson:"signers" json:"signers"`
	Signatures          []string            `bson:"signatures" json:"signatures"`
	MintTransactionHash string              `bson:"mint_tx_hash" json:"mint_transaction_hash"`
}

type MintMemo struct {
	Address string `bson:"address" json:"address"`
	ChainId string `bson:"chain_id" json:"chain_id"`
}

type MintData struct {
	Recipient string `bson:"recipient" json:"recipient"`
	Amount    string `bson:"amount" json:"amount"`
	Nonce     string `bson:"nonce" json:"nonce"`
}
