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
)

func TestValidateTxToCosmosMultisig(t *testing.T) {
	bech32Prefix := "pokt"
	supportedChainIDsEthereum := map[uint32]bool{
		1: true,
	}
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
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
	txValue, _ := tx.Marshal()
	txResponse.Tx = &codectypes.Any{Value: txValue}

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusConfirmed, result.TxStatus)
	assert.Equal(t, uint64(10), result.Confirmations)
	assert.Equal(t, strings.ToLower("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"), result.Memo.Address)
	assert.Equal(t, "1", result.Memo.ChainID)
	assert.Equal(t, sdk.NewCoin("upokt", math.NewInt(1000)), result.Amount)
	assert.False(t, result.NeedsRefund)
}

func TestValidateTxToCosmosMultisig_ErrorParsingSender(t *testing.T) {
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{},
	}
	config := models.PocketConfig{}
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.Error(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_ErrorParsingSenderAddress(t *testing.T) {
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: "invalid_sender"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix: "pokt",
		CoinDenom:    "upokt",
	}
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.Error(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
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
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusFailed, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_ErrorParsingCoinsReceived(t *testing.T) {
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000invalid-anmount1"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_ErrorParsingCoinsSpent(t *testing.T) {
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
					{Key: "amount", Value: "1000invalid"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "0upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
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
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "50upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
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
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_InvalidSpender(t *testing.T) {
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: "invalid_sender"},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_SenderMismatch(t *testing.T) {
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

	spenderAddress := ethcommon.BytesToAddress([]byte("pokt1another"))
	spenderBech32, err := common.Bech32FromBytes(bech32Prefix, spenderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: spenderBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}
	config := models.PocketConfig{
		Bech32Prefix:    "pokt",
		CoinDenom:       "upokt",
		MultisigAddress: multisigBech32,
	}
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
}

func TestValidateTxToCosmosMultisig_InvalidCoins(t *testing.T) {
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
					{Key: "amount", Value: "500upokt"},
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
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusConfirmed, result.TxStatus)
	assert.True(t, result.NeedsRefund)
}

func TestValidateTxToCosmosMultisig_InvalidMemo(t *testing.T) {
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
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
	supportedChainIDsEthereum := map[uint32]bool{}
	currentCosmosBlockHeight := uint64(100)

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusConfirmed, result.TxStatus)
	assert.True(t, result.NeedsRefund)
}

func TestValidateTxToCosmosMultisig_ErrorUnmarshallingTx(t *testing.T) {
	bech32Prefix := "pokt"
	supportedChainIDsEthereum := map[uint32]bool{
		1: true,
	}
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
				Type: "message",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
				},
			},
			{
				Type: "coin_received",
				Attributes: []abci.EventAttribute{
					{Key: "receiver", Value: multisigBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
			{
				Type: "coin_spent",
				Attributes: []abci.EventAttribute{
					{Key: "spender", Value: senderBech32},
					{Key: "amount", Value: "1000upokt"},
				},
			},
		},
	}

	txResponse.Tx = &codectypes.Any{Value: []byte("invalid")}

	result, err := ValidateTxToCosmosMultisig(txResponse, config, supportedChainIDsEthereum, currentCosmosBlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, models.TransactionStatusInvalid, result.TxStatus)
}
