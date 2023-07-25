package ethereum

import (
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/dan13ram/wpokt-backend/models"
)

func createBurn(event *autogen.WrappedPocketBurnAndBridge) models.Burn {
	doc := models.Burn{
		BlockNumber:      strconv.FormatInt(int64(event.Raw.BlockNumber), 10),
		Confirmations:    "0",
		TransactionHash:  event.Raw.TxHash.String(),
		LogIndex:         strconv.FormatInt(int64(event.Raw.Index), 10),
		WPOKTAddress:     event.Raw.Address.String(),
		SenderAddress:    event.From.String(),
		SenderChainId:    app.Config.Ethereum.ChainId,
		RecipientAddress: strings.Split(event.PoktAddress.String(), "0x")[1],
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
