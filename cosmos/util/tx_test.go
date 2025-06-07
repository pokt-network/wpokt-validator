package util

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestNewSendTx(t *testing.T) {
	bech32Prefix := "pokt"
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.NoError(t, err)
	assert.NotEmpty(t, txBody)
}

func TestNewSendTx_ErrorFromAddress(t *testing.T) {
	bech32Prefix := "pokt"
	fromAddr := []byte{}
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	txBody, err := NewSendTx(bech32Prefix, fromAddr, toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error converting from address")
	assert.Empty(t, txBody)
}

func TestNewSendTx_ErrorToAddress(t *testing.T) {
	bech32Prefix := "pokt"
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := []byte{}
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr, amountIncludingFees, memo, feeAmount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error converting to address")
	assert.Empty(t, txBody)
}

func TestNewSendTx_ErrorSettingMsg(t *testing.T) {
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	bech32Prefix := "pokt"

	mockTxConfig := mocks.NewMockTxConfig(t)
	mockTxBuilder := mocks.NewMockTxBuilder(t)
	NewTxConfig = func(bech32Prefix string) client.TxConfig {
		return mockTxConfig
	}
	defer func() {
		NewTxConfig = newTxConfig
	}()
	mockTxConfig.EXPECT().NewTxBuilder().Return(mockTxBuilder)
	mockTxBuilder.EXPECT().SetMsgs(mock.Anything).Return(fmt.Errorf("error setting msg"))

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error setting msg")
	assert.Empty(t, txBody)
}

func TestNewSendTx_InvalidTxEncoder(t *testing.T) {
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	bech32Prefix := "pokt"

	mockTxConfig := mocks.NewMockTxConfig(t)
	mockTxBuilder := mocks.NewMockTxBuilder(t)
	NewTxConfig = func(bech32Prefix string) client.TxConfig {
		return mockTxConfig
	}
	defer func() {
		NewTxConfig = newTxConfig
	}()
	mockTxConfig.EXPECT().NewTxBuilder().Return(mockTxBuilder)
	mockTxBuilder.EXPECT().SetMsgs(mock.Anything).Return(nil)
	mockTxBuilder.EXPECT().SetMemo(memo)
	mockTxBuilder.EXPECT().SetFeeAmount(sdk.NewCoins(feeAmount))
	mockTxBuilder.EXPECT().SetGasLimit(SendGasLimit)
	mockTxConfig.EXPECT().TxJSONEncoder().Return(nil)

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting tx encoder")
	assert.Empty(t, txBody)
}

func TestNewSendTx_ErrorEncoding(t *testing.T) {
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	bech32Prefix := "pokt"

	mockTxConfig := mocks.NewMockTxConfig(t)
	mockTxBuilder := mocks.NewMockTxBuilder(t)
	NewTxConfig = func(bech32Prefix string) client.TxConfig {
		return mockTxConfig
	}
	defer func() {
		NewTxConfig = newTxConfig
	}()
	mockTxConfig.EXPECT().NewTxBuilder().Return(mockTxBuilder)
	mockTxBuilder.EXPECT().SetMsgs(mock.Anything).Return(nil)
	mockTxBuilder.EXPECT().SetMemo(memo)
	mockTxBuilder.EXPECT().SetFeeAmount(sdk.NewCoins(feeAmount))
	mockTxBuilder.EXPECT().SetGasLimit(SendGasLimit)
	mockTxBuilder.EXPECT().GetTx().Return(nil)

	var txJSONEncoder sdk.TxEncoder = func(tx sdk.Tx) ([]byte, error) {
		return nil, fmt.Errorf("error encoding tx")
	}
	mockTxConfig.EXPECT().TxJSONEncoder().Return(txJSONEncoder)

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error encoding tx")
	assert.Empty(t, txBody)
}

func TestParseTxBody(t *testing.T) {
	bech32Prefix := "pokt"
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.NoError(t, err)
	assert.NotEmpty(t, txBody)

	tx, err := ParseTxBody(bech32Prefix, txBody)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	msgs := tx.GetMsgs()
	assert.Equal(t, 1, len(msgs))
	assert.IsType(t, &banktypes.MsgSend{}, msgs[0])
}

func TestParseTxBody_Error(t *testing.T) {
	bech32Prefix := "pokt"
	invalidTxBody := "invalid_tx_body"

	tx, err := ParseTxBody(bech32Prefix, invalidTxBody)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding tx")
	assert.Nil(t, tx)
}

func TestParseTxBody_InvalidTxDecoder(t *testing.T) {
	bech32Prefix := "pokt"
	invalidTxBody := "valid_tx_body"

	mockTxConfig := mocks.NewMockTxConfig(t)
	NewTxConfig = func(bech32Prefix string) client.TxConfig {
		return mockTxConfig
	}
	defer func() {
		NewTxConfig = newTxConfig
	}()
	mockTxConfig.EXPECT().TxJSONDecoder().Return(nil)

	tx, err := ParseTxBody(bech32Prefix, invalidTxBody)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting tx decoder")
	assert.Nil(t, tx)
}

func TestWrapTxBuilder(t *testing.T) {
	bech32Prefix := "pokt"
	fromAddr := ethcommon.BytesToAddress([]byte{1, 2, 3})
	toAddr := ethcommon.BytesToAddress([]byte{4, 5, 6})
	amountIncludingFees := sdk.NewCoin("upokt", math.NewInt(1000))
	feeAmount := sdk.NewCoin("upokt", math.NewInt(100))
	memo := "Test Memo"

	txBody, err := NewSendTx(bech32Prefix, fromAddr[:], toAddr[:], amountIncludingFees, memo, feeAmount)
	assert.NoError(t, err)
	assert.NotEmpty(t, txBody)

	txBuilder, txConfig, err := WrapTxBuilder(bech32Prefix, txBody)
	assert.NoError(t, err)
	assert.NotNil(t, txBuilder)
	assert.NotNil(t, txConfig)
	assert.Equal(t, memo, txBuilder.GetTx().GetMemo())
	assert.Equal(t, feeAmount, txBuilder.GetTx().GetFee()[0])

	msgs := txBuilder.GetTx().GetMsgs()
	msg := msgs[0].(*banktypes.MsgSend)

	fromAddrBech32, _ := common.Bech32FromBytes(bech32Prefix, fromAddr[:])
	toAddrBech32, _ := common.Bech32FromBytes(bech32Prefix, toAddr[:])

	finalAmount := amountIncludingFees.Sub(feeAmount)

	assert.Equal(t, fromAddrBech32, msg.FromAddress)
	assert.Equal(t, toAddrBech32, msg.ToAddress)
	assert.Equal(t, finalAmount, msg.Amount[0])
}

func TestWrapTxBuilder_ErrorParsingTxBody(t *testing.T) {
	bech32Prefix := "pokt"
	invalidTxBody := "invalid_tx_body"

	txBuilder, txConfig, err := WrapTxBuilder(bech32Prefix, invalidTxBody)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing tx body")
	assert.Nil(t, txBuilder)
	assert.Nil(t, txConfig)
}
