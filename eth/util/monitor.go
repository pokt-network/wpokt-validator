package util

import (
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/dan13ram/wpokt-validator/models"

	"github.com/dan13ram/wpokt-validator/common"
)

func CreateBurn(event *autogen.WrappedPocketBurnAndBridge) models.Burn {
	recipientAddress, _ := common.Bech32FromBytes(app.Config.Pocket.Bech32Prefix, event.PoktAddress.Bytes())

	doc := models.Burn{
		BlockNumber:           strconv.FormatInt(int64(event.Raw.BlockNumber), 10),
		Confirmations:         "0",
		TransactionHash:       strings.ToLower(event.Raw.TxHash.String()),
		LogIndex:              strconv.FormatInt(int64(event.Raw.Index), 10),
		WPOKTAddress:          strings.ToLower(event.Raw.Address.String()),
		SenderAddress:         strings.ToLower(event.From.String()),
		SenderChainID:         app.Config.Ethereum.ChainID,
		RecipientAddress:      recipientAddress,
		RecipientChainID:      app.Config.Pocket.ChainID,
		Amount:                event.Amount.String(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		Status:                models.StatusPending,
		Signatures:            []models.Signature{},
		Sequence:              nil,
		ReturnTransactionHash: "",
		ReturnTransactionBody: "",
	}
	return doc
}
