package util

import (
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/dan13ram/wpokt-validator/models"
)

func CreateBurn(event *autogen.WrappedPocketBurnAndBridge) models.Burn {
	doc := models.Burn{
		BlockNumber:      strconv.FormatInt(int64(event.Raw.BlockNumber), 10),
		Confirmations:    "0",
		TransactionHash:  strings.ToLower(event.Raw.TxHash.String()),
		LogIndex:         strconv.FormatInt(int64(event.Raw.Index), 10),
		WPOKTAddress:     strings.ToLower(event.Raw.Address.String()),
		SenderAddress:    strings.ToLower(event.From.String()),
		SenderChainId:    app.Config.Ethereum.ChainId,
		RecipientAddress: strings.ToLower(strings.Split(event.PoktAddress.String(), "0x")[1]),
		RecipientChainId: app.Config.Pocket.ChainId,
		Amount:           event.Amount.String(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           models.StatusPending,
		Signers:          []string{},
		ReturnTxHash:     "",
		ReturnTx:         "",
	}
	return doc
}
