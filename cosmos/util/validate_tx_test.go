package util

import (
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/dan13ram/wpokt-validator/app"
)

func TestValidateTxToCosmosMultisig(t *testing.T) {
	bech32Prefix := "pokt"
	app.Config.Ethereum.ChainID = "1"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	config := models.PocketConfig{
		Bech32Prefix:    bech32Prefix,
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
		TxFee:           100,
		Confirmations:   10,
	}
	currentCosmosBlockHeight := uint64(100)

	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}

	tx := &tx.Tx{
		Body: &tx.TxBody{
			Memo: `{"address": "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B", "chain_id": "1"}`,
		},
	}
	txValue, err := tx.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	txResponse.Tx = &codectypes.Any{Value: txValue}

	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusConfirmed, result.TxStatus)
	assert.Equal(t, uint64(10), result.Confirmations)
	assert.Equal(t, strings.ToLower("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"), result.Memo.Address)
	assert.Equal(t, "1", result.Memo.ChainID)
	assert.Equal(t, sdk.NewCoin("upokt", math.NewInt(1000)), result.Amount)
	assert.False(t, result.NeedsRefund)
}

func TestValidateTxToCosmosMultisig_TxWithNonZeroCode(t *testing.T) {
	bech32Prefix := "pokt"
	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   1,
		Events: []abci.Event{
			{
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
		},
	}
	config := models.PocketConfig{
		CoinDenom:    "upokt",
		Bech32Prefix: "pokt",
	}
	currentCosmosBlockHeight := uint64(100)

	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusFailed, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_ZeroCoins(t *testing.T) {
	bech32Prefix := "pokt"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "0upokt"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	currentCosmosBlockHeight := uint64(100)
	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusFailed, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_AmountTooLow(t *testing.T) {
	bech32Prefix := "pokt"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "50upokt"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
		TxFee:           100,
	}
	currentCosmosBlockHeight := uint64(100)

	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusFailed, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_InvalidAmount(t *testing.T) {
	bech32Prefix := "pokt"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "wrongpokt"},
				},
			},
		},
	}
	tx := &tx.Tx{}
	txValue, _ := tx.Marshal()
	txResponse.Tx = &codectypes.Any{Value: txValue}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	currentCosmosBlockHeight := uint64(100)
	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)
	assert.Equal(t, models.TransactionStatusFailed, result.TxStatus)
	assert.False(t, result.NeedsRefund)
}

func TestValidateTxToCosmosMultisig_FailedMemo(t *testing.T) {
	bech32Prefix := "pokt"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}

	tx := &tx.Tx{
		Body: &tx.TxBody{
			Memo: `{"address": "invalid", "chain_id": "1"}`,
		},
	}
	txValue, _ := tx.Marshal()
	txResponse.Tx = &codectypes.Any{Value: txValue}

	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	currentCosmosBlockHeight := uint64(100)

	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusConfirmed, result.TxStatus)
	assert.True(t, result.NeedsRefund)
}

func TestValidateTxToCosmosMultisig_ErrorUnmarshallingTx(t *testing.T) {
	bech32Prefix := "pokt"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	config := models.PocketConfig{
		Bech32Prefix:    bech32Prefix,
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
		TxFee:           100,
		Confirmations:   10,
	}
	currentCosmosBlockHeight := uint64(100)

	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}

	txResponse.Tx = &codectypes.Any{Value: []byte("invalid")}
	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusFailed, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_AmountTooHigh(t *testing.T) {
	bech32Prefix := "pokt"
	app.Config.Ethereum.ChainID = "1"
	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	config := models.PocketConfig{
		Bech32Prefix:    bech32Prefix,
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
		TxFee:           100,
		Confirmations:   10,
	}
	currentCosmosBlockHeight := uint64(100)

	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "1000000upokt"},
				},
			},
		},
	}

	tx := &tx.Tx{
		Body: &tx.TxBody{
			Memo: `{"address": "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B", "chain_id": "1"}`,
		},
	}
	txValue, err := tx.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	txResponse.Tx = &codectypes.Any{Value: txValue}

	minimum := math.NewInt(100)
	maximum := math.NewInt(10000)

	result := ValidateTxToCosmosMultisig(txResponse, config, currentCosmosBlockHeight, minimum, maximum)

	assert.Equal(t, models.TransactionStatusConfirmed, result.TxStatus)
	assert.Equal(t, uint64(10), result.Confirmations)
	assert.Equal(t, strings.ToLower("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"), result.Memo.Address)
	assert.Equal(t, "1", result.Memo.ChainID)
	assert.Equal(t, sdk.NewCoin("upokt", math.NewInt(1000000)), result.Amount)
	assert.True(t, result.NeedsRefund)
}
