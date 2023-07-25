package pocket

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	pocket "github.com/dan13ram/wpokt-backend/pocket/client"
	"github.com/ethereum/go-ethereum/common"
)

func createMint(tx *pocket.TxResponse, memo models.MintMemo, wpoktAddress string, vaultAddress string) models.Mint {
	return models.Mint{
		Height:           strconv.FormatInt(tx.Height, 10),
		Confirmations:    "0",
		TransactionHash:  tx.Hash,
		SenderAddress:    tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:    app.Config.Pocket.ChainId,
		RecipientAddress: memo.Address,
		RecipientChainId: memo.ChainId,
		WPOKTAddress:     wpoktAddress,
		VaultAddress:     vaultAddress,
		Amount:           tx.StdTx.Msg.Value.Amount,
		Memo:             &memo,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           models.StatusPending,
		Data:             nil,
	}
}

func validateMemo(txMemo string) (models.MintMemo, bool) {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(txMemo), &memo)

	address := common.HexToAddress(memo.Address).Hex()

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId || strings.ToLower(address) != strings.ToLower(memo.Address) {
		return memo, false
	}
	memo.Address = address
	return memo, true
}

func createInvalidMint(tx *pocket.TxResponse, vaultAddress string) models.InvalidMint {
	return models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: tx.Hash,
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.Pocket.ChainId,
		Memo:            tx.StdTx.Memo,
		Amount:          tx.StdTx.Msg.Value.Amount,
		VaultAddress:    vaultAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
		ReturnTx:        "",
		ReturnTxHash:    "",
	}
}
