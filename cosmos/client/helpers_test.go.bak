package client

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	rpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cometbft/cometbft/types"

	clientMocks "github.com/dan13ram/wpokt-validator/cosmos/client/client_mocks"
	"github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
)

func TestGetBlocksForTxResults(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosHTTPClient(t)
	resTxs := []*rpctypes.ResultTx{
		{Height: 1},
		{Height: 2},
		{Height: 1},
	}

	resBlock1 := &rpctypes.ResultBlock{Block: &types.Block{Header: types.Header{Height: 1}}}
	resBlock2 := &rpctypes.ResultBlock{Block: &types.Block{Header: types.Header{Height: 2}}}

	mockClient.EXPECT().Block(mock.Anything, &resTxs[0].Height).Return(resBlock1, nil).Once()
	mockClient.EXPECT().Block(mock.Anything, &resTxs[1].Height).Return(resBlock2, nil).Once()

	resBlocks, err := getBlocksForTxResults(mockClient, resTxs)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(resBlocks))
	assert.Equal(t, resBlock1, resBlocks[1])
	assert.Equal(t, resBlock2, resBlocks[2])
}

func TestGetBlocksForTxResults_Error(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosHTTPClient(t)
	resTxs := []*rpctypes.ResultTx{
		{Height: 1},
		{Height: 2},
		{Height: 1},
	}

	resBlock1 := &rpctypes.ResultBlock{Block: &types.Block{Header: types.Header{Height: 1}}}

	mockClient.EXPECT().Block(mock.Anything, &resTxs[0].Height).Return(resBlock1, errors.New("error")).Once()

	resBlocks, err := getBlocksForTxResults(mockClient, resTxs)

	assert.Error(t, err)
	assert.Nil(t, resBlocks)
}

func TestFormatTxResults(t *testing.T) {
	resTxs := []*rpctypes.ResultTx{
		{Height: 1, Tx: []byte{1, 2, 3}},
	}

	resBlock := &rpctypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: 1,
				Time:   time.Now(),
			},
		},
	}

	resBlocks := map[int64]*rpctypes.ResultBlock{
		1: resBlock,
	}

	mockTx := clientMocks.NewMockAnyTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, nil
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	mockTx.EXPECT().AsAny().Return(&codectypes.Any{})

	res, err := formatTxResults("prefix", resTxs, resBlocks)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
}

func TestFormatTxResults_Error(t *testing.T) {
	resTxs := []*rpctypes.ResultTx{
		{Height: 1, Tx: []byte{1, 2, 3}},
	}

	resBlock := &rpctypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: 1,
				Time:   time.Now(),
			},
		},
	}

	resBlocks := map[int64]*rpctypes.ResultBlock{
		1: resBlock,
	}

	mockTx := clientMocks.NewMockAnyTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, errors.New("error")
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	res, err := formatTxResults("prefix", resTxs, resBlocks)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "decoding tx")
}

func TestMkTxResult(t *testing.T) {
	resTx := &rpctypes.ResultTx{Height: 1, Tx: []byte{1, 2, 3}}
	resBlock := &rpctypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: 1,
				Time:   time.Now(),
			},
		},
	}

	mockTx := clientMocks.NewMockAnyTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, nil
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	mockTx.EXPECT().AsAny().Return(&codectypes.Any{})

	txResult, err := mkTxResult("prefix", resTx, resBlock)
	assert.NoError(t, err)
	assert.NotNil(t, txResult)
}

func TestMkTxResult_ErrorDecoding(t *testing.T) {
	resTx := &rpctypes.ResultTx{Height: 1, Tx: []byte{1, 2, 3}}
	resBlock := &rpctypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: 1,
				Time:   time.Now(),
			},
		},
	}

	mockTx := clientMocks.NewMockAnyTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, errors.New("error")
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	txResult, err := mkTxResult("prefix", resTx, resBlock)
	assert.Error(t, err)
	assert.Nil(t, txResult)
	assert.Contains(t, err.Error(), "decoding tx")
}

func TestMkTxResult_InvalidTx(t *testing.T) {
	resTx := &rpctypes.ResultTx{Height: 1, Tx: []byte{1, 2, 3}}
	resBlock := &rpctypes.ResultBlock{
		Block: &types.Block{
			Header: types.Header{
				Height: 1,
				Time:   time.Now(),
			},
		},
	}

	mockTx := mocks.NewMockTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, nil
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	txResult, err := mkTxResult("prefix", resTx, resBlock)
	assert.Error(t, err)
	assert.Nil(t, txResult)
	assert.Contains(t, err.Error(), "expecting a type implementing")
}
