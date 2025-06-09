package util

import (
	"cosmossdk.io/math"
	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/tx"

	log "github.com/sirupsen/logrus"
)

type ValidateTxResult struct {
	TxValid       bool
	Tx            *tx.Tx
	TxHash        string
	Memo          models.MintMemo
	Amount        sdk.Coin
	SenderAddress string
	NeedsRefund   bool
}

func ValidateTxToCosmosMultisig(
	txResponse *sdk.TxResponse,
	config models.CosmosConfig,
	minAmount math.Int,
	maxAmount math.Int,
) *ValidateTxResult {
	logger := log.
		WithField("operation", "validateTxToCosmosMultisig").
		WithField("tx_hash", txResponse.TxHash)

	result := ValidateTxResult{
		TxValid:       false,
		Tx:            nil,
		TxHash:        common.Ensure0xPrefix(txResponse.TxHash),
		Amount:        sdk.Coin{},
		Memo:          models.MintMemo{},
		SenderAddress: "",
		NeedsRefund:   false,
	}

	if txResponse.Code != 0 {
		logger.Debugf("Found tx with non-zero code")
		return &result
	}

	transfers, err := ParseTransferEvents(
		txResponse.Events,
		config.MultisigAddress,
		config.CoinDenom,
	)

	if err != nil {
		logger.WithError(err).Debugf("Error parsing transfer events")
		return &result
	}

	if len(transfers) != 1 {
		logger.Debugf("Found tx with invalid transfers, expected 1, got %d", len(transfers))
		return &result
	}

	result.SenderAddress = transfers[0].Sender
	result.Amount = transfers[0].Amount

	if result.Amount.IsZero() {
		logger.Debugf("Found tx transfer with zero coins")
		return &result
	}

	if result.Amount.Amount.LTE(minAmount) {
		logger.Debugf("Found tx transfer with amount too low")
		return &result
	}

	tx := &tx.Tx{}
	err = tx.Unmarshal(txResponse.Tx.Value)
	if err != nil {
		logger.WithError(err).Errorf("Error unmarshalling tx")
		return &result
	}
	result.Tx = tx
	result.TxValid = true

	memo, err := ValidateMemo(tx.Body.Memo)
	if err != nil {
		logger.WithError(err).WithField("memo", tx.Body.Memo).Debugf("Found invalid memo")
		// refund
		result.NeedsRefund = true
		return &result
	}

	logger.WithField("memo", memo).Debugf("Found valid memo")
	result.Memo = memo

	if result.Amount.Amount.GT(maxAmount) {
		// refund any transactions that are too large since they can't be processed on ethereum due to the max mint limit
		logger.Debugf("Found tx transfer with amount too high")
		result.NeedsRefund = true
		return &result
	}

	return &result
}
