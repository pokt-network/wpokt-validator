package util

import (
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateMint(tx *sdk.TxResponse, result *ValidateTxResult, wpoktAddress string, vaultAddress string) models.Mint {
	return models.Mint{
		Height:              strconv.FormatInt(tx.Height, 10),
		Confirmations:       "0",
		TransactionHash:     common.Ensure0xPrefix(tx.TxHash),
		SenderAddress:       result.SenderAddress,
		SenderChainID:       app.Config.Pocket.ChainID,
		RecipientAddress:    strings.ToLower(result.Memo.Address),
		RecipientChainID:    result.Memo.ChainID,
		WPOKTAddress:        strings.ToLower(wpoktAddress),
		VaultAddress:        strings.ToLower(vaultAddress),
		Amount:              result.Amount.Amount.String(),
		Memo:                &result.Memo,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Status:              models.StatusPending,
		Data:                nil,
		MintTransactionHash: "",
		Signers:             []string{},
		Signatures:          []string{},
	}
}

func CreateInvalidMint(tx *sdk.TxResponse, result *ValidateTxResult, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: common.Ensure0xPrefix(tx.TxHash),
		SenderAddress:   result.SenderAddress,
		SenderChainID:   app.Config.Pocket.ChainID,
		Memo:            result.Tx.Body.Memo,
		Amount:          result.Amount.Amount.String(),
		VaultAddress:    strings.ToLower(vaultAddress),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,

		Signatures:            []models.Signature{},
		Sequence:              nil,
		ReturnTransactionHash: "",
		ReturnTransactionBody: "",
	}
}

func CreateFailedMint(tx *sdk.TxResponse, result *ValidateTxResult, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: common.Ensure0xPrefix(tx.TxHash),
		SenderAddress:   result.SenderAddress,
		SenderChainID:   app.Config.Pocket.ChainID,
		Memo:            result.Tx.Body.Memo,
		Amount:          result.Amount.Amount.String(),
		VaultAddress:    strings.ToLower(vaultAddress),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusFailed,

		Signatures:            []models.Signature{},
		Sequence:              nil,
		ReturnTransactionHash: "",
		ReturnTransactionBody: "",
	}
}
