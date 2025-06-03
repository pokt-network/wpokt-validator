package util

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	// cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateMint(tx *sdk.TxResponse, memo models.MintMemo, wpoktAddress string, vaultAddress string) models.Mint {
	return models.Mint{
		Height:        strconv.FormatInt(tx.Height, 10),
		Confirmations: "0",
		// TransactionHash:     strings.ToLower(tx.Hash),
		// SenderAddress:       strings.ToLower(tx.StdTx.Msg.Value.FromAddress),
		SenderChainID:    app.Config.Pocket.ChainID,
		RecipientAddress: strings.ToLower(memo.Address),
		RecipientChainID: memo.ChainID,
		WPOKTAddress:     strings.ToLower(wpoktAddress),
		VaultAddress:     strings.ToLower(vaultAddress),
		// Amount:              tx.StdTx.Msg.Value.Amount,
		Memo:                &memo,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Status:              models.StatusPending,
		Data:                nil,
		MintTransactionHash: "",
		Signers:             []string{},
		Signatures:          []string{},
	}
}

func ValidateMemo(txMemo string) (models.MintMemo, bool) {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(txMemo), &memo)
	if err != nil {
		return memo, false
	}

	address := common.HexToAddress(memo.Address).Hex()
	if !strings.EqualFold(address, memo.Address) {
		return memo, false
	}

	if address == common.HexToAddress("").Hex() {
		return memo, false
	}
	memo.Address = address

	memoChainID, err := strconv.Atoi(memo.ChainID)
	if err != nil {
		return memo, false
	}

	appChainID, err := strconv.Atoi(app.Config.Ethereum.ChainID)
	if err != nil {
		return memo, false
	}

	if memoChainID != appChainID {
		return memo, false
	}
	memo.ChainID = app.Config.Ethereum.ChainID
	return memo, true
}

func CreateInvalidMint(tx *sdk.TxResponse, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:        strconv.FormatInt(tx.Height, 10),
		Confirmations: "0",
		// TransactionHash: strings.ToLower(tx.Hash),
		// SenderAddress:   strings.ToLower(tx.StdTx.Msg.Value.FromAddress),
		SenderChainID: app.Config.Pocket.ChainID,
		// Memo:            tx.StdTx.Memo,
		// Amount:          tx.StdTx.Msg.Value.Amount,
		VaultAddress: strings.ToLower(vaultAddress),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Status:       models.StatusPending,
		Signers:      []string{},
		ReturnTx:     "",
		ReturnTxHash: "",
	}
}

func CreateFailedMint(tx *sdk.TxResponse, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:        strconv.FormatInt(tx.Height, 10),
		Confirmations: "0",
		// TransactionHash: strings.ToLower(tx.Hash),
		// SenderAddress:   strings.ToLower(tx.StdTx.Msg.Value.FromAddress),
		SenderChainID: app.Config.Pocket.ChainID,
		// Memo:            tx.StdTx.Memo,
		// Amount:          tx.StdTx.Msg.Value.Amount,
		VaultAddress: strings.ToLower(vaultAddress),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Status:       models.StatusFailed,
		Signers:      []string{},
		ReturnTx:     "",
		ReturnTxHash: "",
	}
}
