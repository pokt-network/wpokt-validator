package util

import (
	"bytes"

	"cosmossdk.io/math"
	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/tx"

	log "github.com/sirupsen/logrus"
)

type ValidateTxResult struct {
	Memo          models.MintMemo
	Confirmations uint64
	TxStatus      models.TransactionStatus
	Tx            *tx.Tx
	Amount        sdk.Coin
	SenderAddress []byte
	NeedsRefund   bool
}

func ValidateTxToCosmosMultisig(
	txResponse *sdk.TxResponse,
	config models.PocketConfig,
	currentCosmosBlockHeight uint64,
) (*ValidateTxResult, error) {
	logger := log.
		WithField("operation", "validateTxToCosmosMultisig").
		WithField("tx_hash", txResponse.TxHash)

	result := ValidateTxResult{
		Memo:          models.MintMemo{},
		TxStatus:      models.TransactionStatusInvalid,
		Tx:            nil,
		Amount:        sdk.Coin{},
		SenderAddress: nil,
		NeedsRefund:   false,
	}

	sender, err := ParseMessageSenderEvent(txResponse.Events)
	if err != nil {
		logger.WithError(err).Errorf("Error parsing message sender")
		return &result, err
	}

	senderAddress, err := common.AddressBytesFromBech32(config.Bech32Prefix, sender)
	if err != nil {
		logger.WithError(err).Errorf("Error parsing sender address")
		return &result, err
	}

	result.SenderAddress = senderAddress

	if txResponse.Code != 0 {
		logger.Debugf("Found tx with non-zero code")
		result.TxStatus = models.TransactionStatusFailed
		return &result, nil
	}

	coinsReceived, err := ParseCoinsReceivedEvents(config.CoinDenom, config.MultisigAddress, txResponse.Events)
	if err != nil {
		logger.WithError(err).Debugf("Error parsing coins received events")
		return &result, nil
	}

	coinsSpentSender, coinsSpent, err := ParseCoinsSpentEvents(config.CoinDenom, txResponse.Events)
	if err != nil {
		logger.WithError(err).Debugf("Error parsing coins spent events")
		return &result, nil
	}

	if coinsReceived.IsZero() || coinsSpent.IsZero() {
		logger.Debugf("Found tx with zero coins")
		return &result, nil
	}

	if coinsReceived.Amount.LTE(math.NewIntFromUint64(uint64(config.TxFee))) {
		logger.Debugf("Found tx with amount too low")
		return &result, nil
	}

	spenderAddress, err := common.AddressBytesFromBech32(config.Bech32Prefix, coinsSpentSender)
	if err != nil {
		logger.WithError(err).Errorf("Error parsing spender address")
		return &result, nil
	}
	if !bytes.Equal(senderAddress, spenderAddress) {
		logger.Errorf("Sender address does not match spender address")
		return &result, nil
	}

	tx := &tx.Tx{}
	err = tx.Unmarshal(txResponse.Tx.Value)
	if err != nil {
		logger.WithError(err).Errorf("Error unmarshalling tx")
		return &result, nil
	}
	result.Tx = tx

	result.TxStatus = models.TransactionStatusPending

	result.Confirmations = currentCosmosBlockHeight - uint64(txResponse.Height)
	if result.Confirmations >= uint64(config.Confirmations) {
		result.TxStatus = models.TransactionStatusConfirmed
	}

	result.Amount = coinsSpent

	if !coinsSpent.Amount.Equal(coinsReceived.Amount) {
		logger.Debugf("Found tx with invalid coins")
		// refund
		result.NeedsRefund = true
		return &result, nil
	}

	memo, err := ValidateMemo(tx.Body.Memo)
	if err != nil {
		logger.WithError(err).WithField("memo", tx.Body.Memo).Debugf("Found invalid memo")
		// refund
		result.NeedsRefund = true
		return &result, nil
	}

	logger.WithField("memo", memo).Debugf("Found valid memo")
	result.Memo = memo

	return &result, nil
}
