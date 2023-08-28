package util

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/ethereum/go-ethereum/common"
)

func CreateMint(tx *pokt.TxResponse, memo models.MintMemo, wpoktAddress string, vaultAddress string) models.Mint {
	return models.Mint{
		Height:              strconv.FormatInt(tx.Height, 10),
		Confirmations:       "0",
		TransactionHash:     strings.ToLower(tx.Hash),
		SenderAddress:       strings.ToLower(tx.StdTx.Msg.Value.FromAddress),
		SenderChainId:       app.Config.Pocket.ChainId,
		RecipientAddress:    strings.ToLower(memo.Address),
		RecipientChainId:    memo.ChainId,
		WPOKTAddress:        strings.ToLower(wpoktAddress),
		VaultAddress:        strings.ToLower(vaultAddress),
		Amount:              tx.StdTx.Msg.Value.Amount,
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
	if strings.ToLower(address) != strings.ToLower(memo.Address) {
		return memo, false
	}

	if address == common.HexToAddress("").Hex() {
		return memo, false
	}
	memo.Address = address

	memoChainId, err := strconv.Atoi(memo.ChainId)
	if err != nil {
		return memo, false
	}

	appChainId, err := strconv.Atoi(app.Config.Ethereum.ChainId)
	if err != nil {
		return memo, false
	}

	if memoChainId != appChainId {
		return memo, false
	}
	memo.ChainId = app.Config.Ethereum.ChainId
	return memo, true
}

func CreateInvalidMint(tx *pokt.TxResponse, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: strings.ToLower(tx.Hash),
		SenderAddress:   strings.ToLower(tx.StdTx.Msg.Value.FromAddress),
		SenderChainId:   app.Config.Pocket.ChainId,
		Memo:            tx.StdTx.Memo,
		Amount:          tx.StdTx.Msg.Value.Amount,
		VaultAddress:    strings.ToLower(vaultAddress),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
		ReturnTx:        "",
		ReturnTxHash:    "",
	}
}

func CreateFailedMint(tx *pokt.TxResponse, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: strings.ToLower(tx.Hash),
		SenderAddress:   strings.ToLower(tx.StdTx.Msg.Value.FromAddress),
		SenderChainId:   app.Config.Pocket.ChainId,
		Memo:            tx.StdTx.Memo,
		Amount:          tx.StdTx.Msg.Value.Amount,
		VaultAddress:    strings.ToLower(vaultAddress),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusFailed,
		Signers:         []string{},
		ReturnTx:        "",
		ReturnTxHash:    "",
	}
}
