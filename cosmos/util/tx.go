package util

import (
	"fmt"

	"github.com/dan13ram/wpokt-validator/common"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	SendGasLimit uint64 = 200000
)

func NewSendTx(
	bech32Prefix string,
	fromAddr []byte,
	toAddr []byte,
	amountIncludingFees sdk.Coin,
	memo string,
	feeAmount sdk.Coin,
) (string, error) {

	finalAmount := amountIncludingFees.Sub(feeAmount)

	fromAddress, err := common.Bech32FromBytes(bech32Prefix, fromAddr)
	if err != nil {
		return "", fmt.Errorf("error converting from address: %w", err)
	}
	toAddress, err := common.Bech32FromBytes(bech32Prefix, toAddr)
	if err != nil {
		return "", fmt.Errorf("error converting to address: %w", err)
	}

	msg := &banktypes.MsgSend{FromAddress: fromAddress, ToAddress: toAddress, Amount: sdk.NewCoins(finalAmount)}

	txConfig := NewTxConfig(bech32Prefix)

	refundTx := txConfig.NewTxBuilder()

	err = refundTx.SetMsgs(msg)
	if err != nil {
		return "", fmt.Errorf("error setting msg: %w", err)
	}

	refundTx.SetMemo(memo)
	refundTx.SetFeeAmount(sdk.NewCoins(feeAmount))
	refundTx.SetGasLimit(SendGasLimit)

	txEncoder := txConfig.TxJSONEncoder()

	if txEncoder == nil {
		return "", fmt.Errorf("error getting tx encoder")
	}

	txBody, err := txEncoder(refundTx.GetTx())
	if err != nil {
		return "", fmt.Errorf("error encoding tx: %w", err)
	}
	return string(txBody), nil
}

func ParseTxBody(
	bech32Prefix string,
	txBody string,
) (sdk.Tx, error) {
	txConfig := NewTxConfig(bech32Prefix)
	txDecoder := txConfig.TxJSONDecoder()

	if txDecoder == nil {
		return nil, fmt.Errorf("error getting tx decoder")
	}

	tx, err := txDecoder([]byte(txBody))
	if err != nil {
		return nil, fmt.Errorf("error decoding tx: %w", err)
	}

	return tx, nil
}

func WrapTxBuilder(
	bech32Prefix string,
	txBody string,
) (client.TxBuilder, client.TxConfig, error) {
	tx, err := ParseTxBody(bech32Prefix, txBody)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing tx body: %w", err)
	}
	txConfig := NewTxConfig(bech32Prefix)
	txBuilder, err := txConfig.WrapTxBuilder(tx)
	return txBuilder, txConfig, err
}
