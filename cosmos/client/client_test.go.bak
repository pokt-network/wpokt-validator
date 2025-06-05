package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	grpc "github.com/cosmos/gogoproto/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	rpctypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/cometbft/cometbft/p2p"
	ctypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dan13ram/wpokt-validator/common"
	mocks "github.com/dan13ram/wpokt-validator/cosmos/client/client_mocks"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	"github.com/dan13ram/wpokt-validator/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	goGRPC "google.golang.org/grpc"
)

func TestChain(t *testing.T) {
	chain := models.Chain{ChainID: "TestChainID", ChainName: "TestChain"}
	client := &cosmosClient{chain: chain}
	assert.Equal(t, chain, client.Chain())
}

func TestConfirmations(t *testing.T) {
	client := &cosmosClient{confirmations: 10}
	assert.Equal(t, uint64(10), client.Confirmations())
}

func TestGetLatestBlockHeight_GRPC(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	block := &cmtservice.Block{Header: cmtservice.Header{Height: 100}}
	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(&cmtservice.GetLatestBlockResponse{SdkBlock: block}, nil)

	height, err := client.GetLatestBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, int64(100), height)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetLatestBlockHeight_GRPC_Error(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	block := &cmtservice.Block{Header: cmtservice.Header{Height: 100}}
	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(&cmtservice.GetLatestBlockResponse{SdkBlock: block}, errors.New("error"))

	height, err := client.GetLatestBlockHeight()
	assert.Error(t, err)
	assert.Equal(t, int64(0), height)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetLatestBlockHeight_RPC(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	status := &rpctypes.ResultStatus{SyncInfo: rpctypes.SyncInfo{LatestBlockHeight: 100}}
	mockHTTPClient.On("Status", mock.Anything).Return(status, nil)

	height, err := client.GetLatestBlockHeight()
	assert.NoError(t, err)
	assert.Equal(t, int64(100), height)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetLatestBlockHeight_RPC_Error(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	mockHTTPClient.On("Status", mock.Anything).Return(nil, errors.New("error"))

	height, err := client.GetLatestBlockHeight()
	assert.Error(t, err)
	assert.Equal(t, int64(0), height)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentToAddressAfterHeight_AddressError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	recipientBech32 := "cosmos1test"

	txs, err := client.GetTxsSentToAddressAfterHeight(recipientBech32, 100)
	assert.Error(t, err)
	assert.Nil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentToAddressAfterHeight_GRPC(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	recipientAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	recipientBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, recipientAddress.Bytes())

	query := fmt.Sprintf("transfer.recipient='%s' AND tx.height>=100", recipientBech32)

	req := &tx.GetTxsEventRequest{
		Query:   query,
		OrderBy: tx.OrderBy_ORDER_BY_ASC,
		Page:    1,
		Limit:   50,
	}
	mockGRPCClient.On("GetTxsEvent", mock.Anything, req).Return(&tx.GetTxsEventResponse{Txs: []*tx.Tx{}}, nil)

	txs, err := client.GetTxsSentToAddressAfterHeight(recipientBech32, 100)
	assert.NoError(t, err)
	assert.NotNil(t, txs)
}

func TestGetTxsSentToAddressAfterHeight_GRPC_MultiPages(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	recipientAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	recipientBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, recipientAddress.Bytes())

	query := fmt.Sprintf("transfer.recipient='%s' AND tx.height>=100", recipientBech32)

	req1 := &tx.GetTxsEventRequest{
		Query:   query,
		OrderBy: tx.OrderBy_ORDER_BY_ASC,
		Page:    1,
		Limit:   50,
	}
	req2 := &tx.GetTxsEventRequest{
		Query:   query,
		OrderBy: tx.OrderBy_ORDER_BY_ASC,
		Page:    2,
		Limit:   50,
	}

	resTxs1 := []*sdk.TxResponse{
		{Height: 1},
		{Height: 2},
	}
	resTxs2 := []*sdk.TxResponse{
		{Height: 3},
		{Height: 4},
	}
	mockGRPCClient.On("GetTxsEvent", mock.Anything, req1).Return(&tx.GetTxsEventResponse{TxResponses: resTxs1, Total: 4}, nil).Once()
	mockGRPCClient.On("GetTxsEvent", mock.Anything, req2).Return(&tx.GetTxsEventResponse{TxResponses: resTxs2, Total: 4}, nil).Once()

	txs, err := client.GetTxsSentToAddressAfterHeight(recipientBech32, 100)
	assert.NoError(t, err)
	assert.NotNil(t, txs)
	assert.Len(t, txs, 4)
	assert.Equal(t, int64(1), txs[0].Height)
	assert.Equal(t, int64(2), txs[1].Height)
	assert.Equal(t, int64(3), txs[2].Height)
	assert.Equal(t, int64(4), txs[3].Height)
}

func TestGetTxsSentToAddressAfterHeight_GRPC_Error(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	recipientAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	recipientBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, recipientAddress.Bytes())

	query := fmt.Sprintf("transfer.recipient='%s' AND tx.height>=100", recipientBech32)

	req := &tx.GetTxsEventRequest{
		Query:   query,
		OrderBy: tx.OrderBy_ORDER_BY_ASC,
		Page:    1,
		Limit:   50,
	}
	mockGRPCClient.On("GetTxsEvent", mock.Anything, req).Return(nil, errors.New("error"))

	txs, err := client.GetTxsSentToAddressAfterHeight(recipientBech32, 100)
	assert.Error(t, err)
	assert.Nil(t, txs)
}

func TestGetTxsSentToAddressAfterHeight(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	recipientAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	recipientBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, recipientAddress.Bytes())

	query := fmt.Sprintf("transfer.recipient='%s' AND tx.height>=100", recipientBech32)
	mockHTTPClient.On("TxSearch", mock.Anything, query, false, mock.Anything, mock.Anything, "asc").Return(&rpctypes.ResultTxSearch{Txs: []*rpctypes.ResultTx{}}, nil)

	txs, err := client.GetTxsSentToAddressAfterHeight(recipientBech32, 100)
	assert.NoError(t, err)
	assert.NotNil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentFromAddressAfterHeight_AddressError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	senderBech32 := "cosmos1test"

	txs, err := client.GetTxsSentFromAddressAfterHeight(senderBech32, 100)
	assert.Error(t, err)
	assert.Nil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentFromAddressAfterHeight(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	senderAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	senderBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, senderAddress.Bytes())

	query := fmt.Sprintf("transfer.sender='%s' AND tx.height>=100", senderBech32)
	mockHTTPClient.On("TxSearch", mock.Anything, query, false, mock.Anything, mock.Anything, "asc").Return(&rpctypes.ResultTxSearch{Txs: []*rpctypes.ResultTx{}}, nil)

	txs, err := client.GetTxsSentFromAddressAfterHeight(senderBech32, 100)
	assert.NoError(t, err)
	assert.NotNil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentFromAddressAfterHeight_GetBlocksError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	senderAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	senderBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, senderAddress.Bytes())

	resTxs := []*rpctypes.ResultTx{
		{Height: 1},
		{Height: 2},
		{Height: 1},
	}

	query := fmt.Sprintf("transfer.sender='%s' AND tx.height>=100", senderBech32)
	mockHTTPClient.On("TxSearch", mock.Anything, query, false, mock.Anything, mock.Anything, "asc").Return(&rpctypes.ResultTxSearch{Txs: resTxs}, nil)
	mockHTTPClient.EXPECT().Block(mock.Anything, &resTxs[0].Height).Return(nil, errors.New("error")).Once()

	txs, err := client.GetTxsSentFromAddressAfterHeight(senderBech32, 100)
	assert.Error(t, err)
	assert.Nil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentFromAddressAfterHeight_FormatError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	senderAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	senderBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, senderAddress.Bytes())

	resTxs := []*rpctypes.ResultTx{
		{Height: 1},
		{Height: 2},
		{Height: 1},
	}

	resBlock1 := &rpctypes.ResultBlock{Block: &ctypes.Block{Header: ctypes.Header{Height: 1}}}
	resBlock2 := &rpctypes.ResultBlock{Block: &ctypes.Block{Header: ctypes.Header{Height: 2}}}

	mockHTTPClient.EXPECT().Block(mock.Anything, &resTxs[0].Height).Return(resBlock1, nil).Once()
	mockHTTPClient.EXPECT().Block(mock.Anything, &resTxs[1].Height).Return(resBlock2, nil).Once()

	query := fmt.Sprintf("transfer.sender='%s' AND tx.height>=100", senderBech32)
	mockHTTPClient.On("TxSearch", mock.Anything, query, false, mock.Anything, mock.Anything, "asc").Return(&rpctypes.ResultTxSearch{Txs: resTxs}, nil)

	mockTx := mocks.NewMockAnyTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, errors.New("error")
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	txs, err := client.GetTxsSentFromAddressAfterHeight(senderBech32, 100)
	assert.Error(t, err)
	assert.Nil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTxsSentFromAddressAfterHeight_SearchError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	senderAddress := ethcommon.BytesToAddress([]byte("cosmos1test"))
	senderBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, senderAddress.Bytes())

	query := fmt.Sprintf("transfer.sender='%s' AND tx.height>=100", senderBech32)
	mockHTTPClient.On("TxSearch", mock.Anything, query, false, mock.Anything, mock.Anything, "asc").Return(nil, errors.New("error"))

	txs, err := client.GetTxsSentFromAddressAfterHeight(senderBech32, 100)
	assert.Error(t, err)
	assert.Nil(t, txs)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetAccount_AddressError(t *testing.T) {
	originalAuthNewQueryClient := authNewQueryClient
	defer func() { authNewQueryClient = originalAuthNewQueryClient }()

	mockGRPCClient := mocks.NewMockAuthQueryClient(t)
	authNewQueryClient = func(conn grpc.ClientConn) auth.QueryClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	accountBech32 := "cosmos1account"

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid bech32 address", err.Error())

	mockGRPCClient.AssertExpectations(t)
}

func TestGetAccount_GRPC(t *testing.T) {
	originalAuthNewQueryClient := authNewQueryClient
	defer func() { authNewQueryClient = originalAuthNewQueryClient }()

	mockGRPCClient := mocks.NewMockAuthQueryClient(t)
	authNewQueryClient = func(conn grpc.ClientConn) auth.QueryClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())
	account := &auth.BaseAccount{Address: accountBech32}
	accountBytes, err := account.Marshal()
	assert.NoError(t, err)

	mockGRPCClient.On("Account", mock.Anything, &auth.QueryAccountRequest{Address: accountBech32}).Return(&auth.QueryAccountResponse{Account: &codectypes.Any{Value: accountBytes}}, nil)

	result, err := client.GetAccount(accountBech32)
	assert.NoError(t, err)
	assert.Equal(t, account, result)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetAccount_GRPC_ClientError(t *testing.T) {
	originalAuthNewQueryClient := authNewQueryClient
	defer func() { authNewQueryClient = originalAuthNewQueryClient }()

	mockGRPCClient := mocks.NewMockAuthQueryClient(t)
	authNewQueryClient = func(conn grpc.ClientConn) auth.QueryClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())

	mockGRPCClient.On("Account", mock.Anything, &auth.QueryAccountRequest{Address: accountBech32}).Return(nil, errors.New("error"))

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetAccount_GRPC_UnmarshalError(t *testing.T) {
	originalAuthNewQueryClient := authNewQueryClient
	defer func() { authNewQueryClient = originalAuthNewQueryClient }()

	mockGRPCClient := mocks.NewMockAuthQueryClient(t)
	authNewQueryClient = func(conn grpc.ClientConn) auth.QueryClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  true,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())

	mockGRPCClient.On("Account", mock.Anything, &auth.QueryAccountRequest{Address: accountBech32}).Return(&auth.QueryAccountResponse{Account: &codectypes.Any{Value: []byte("invalid")}}, nil)

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetAccount_RPC(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())
	account := &auth.BaseAccount{Address: accountBech32}
	accountBytes, err := account.Marshal()
	assert.NoError(t, err)

	queryPath := "/cosmos.auth.v1beta1.Query/Account"
	queryData, err := util.NewProtoCodec(config.Bech32Prefix).Marshal(&auth.QueryAccountRequest{Address: accountBech32})
	assert.NoError(t, err)
	var queryDataHex bytes.HexBytes = queryData

	response := auth.QueryAccountResponse{Account: &codectypes.Any{Value: accountBytes}}
	responseBytes, err := response.Marshal()
	assert.NoError(t, err)

	mockHTTPClient.On("ABCIQuery", mock.Anything, queryPath, queryDataHex).Return(&rpctypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: responseBytes}}, nil)

	result, err := client.GetAccount(accountBech32)
	assert.NoError(t, err)
	assert.Equal(t, accountBech32, result.Address)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetAccount_RPC_RequestFailed(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())
	account := &auth.BaseAccount{Address: accountBech32}
	accountBytes, err := account.Marshal()
	assert.NoError(t, err)

	queryPath := "/cosmos.auth.v1beta1.Query/Account"
	queryData, err := util.NewProtoCodec(config.Bech32Prefix).Marshal(&auth.QueryAccountRequest{Address: accountBech32})
	assert.NoError(t, err)
	var queryDataHex bytes.HexBytes = queryData

	response := auth.QueryAccountResponse{Account: &codectypes.Any{Value: accountBytes}}
	responseBytes, err := response.Marshal()
	assert.NoError(t, err)

	mockHTTPClient.On("ABCIQuery", mock.Anything, queryPath, queryDataHex).Return(&rpctypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: responseBytes, Code: 1}}, nil)

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetAccount_RPC_ClientError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())

	queryPath := "/cosmos.auth.v1beta1.Query/Account"
	queryData, err := util.NewProtoCodec(config.Bech32Prefix).Marshal(&auth.QueryAccountRequest{Address: accountBech32})
	assert.NoError(t, err)
	var queryDataHex bytes.HexBytes = queryData

	mockHTTPClient.On("ABCIQuery", mock.Anything, queryPath, queryDataHex).Return(nil, errors.New("error"))

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetAccount_RPC_ResponseUnmarshalError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())

	queryPath := "/cosmos.auth.v1beta1.Query/Account"
	queryData, err := util.NewProtoCodec(config.Bech32Prefix).Marshal(&auth.QueryAccountRequest{Address: accountBech32})
	assert.NoError(t, err)
	var queryDataHex bytes.HexBytes = queryData

	mockHTTPClient.On("ABCIQuery", mock.Anything, queryPath, queryDataHex).Return(&rpctypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: []byte("invalid")}}, nil)

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetAccount_RPC_AccountUnmarshalError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled:  false,
		TimeoutMS:    5000,
		Bech32Prefix: "cosmos",
		ChainName:    "TestChain",
		ChainID:      "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		bech32Prefix:  config.Bech32Prefix,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	accountAddress := ethcommon.BytesToAddress([]byte("cosmos1account"))
	accountBech32, _ := common.Bech32FromBytes(config.Bech32Prefix, accountAddress.Bytes())

	queryPath := "/cosmos.auth.v1beta1.Query/Account"
	queryData, err := util.NewProtoCodec(config.Bech32Prefix).Marshal(&auth.QueryAccountRequest{Address: accountBech32})
	assert.NoError(t, err)
	var queryDataHex bytes.HexBytes = queryData

	response := auth.QueryAccountResponse{Account: &codectypes.Any{Value: []byte("invalid")}}
	responseBytes, err := response.Marshal()
	assert.NoError(t, err)

	mockHTTPClient.On("ABCIQuery", mock.Anything, queryPath, queryDataHex).Return(&rpctypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: responseBytes}}, nil)

	result, err := client.GetAccount(accountBech32)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockHTTPClient.AssertExpectations(t)
}

func TestBroadcastTx_GRPC(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	txBytes := []byte("txBytes")
	mockGRPCClient.On("BroadcastTx", mock.Anything, &tx.BroadcastTxRequest{TxBytes: txBytes, Mode: tx.BroadcastMode_BROADCAST_MODE_SYNC}).Return(&tx.BroadcastTxResponse{TxResponse: &sdk.TxResponse{TxHash: "txHash", Code: 0}}, nil)

	txHash, err := client.BroadcastTx(txBytes)
	assert.NoError(t, err)
	assert.Equal(t, "txHash", txHash)

	mockGRPCClient.AssertExpectations(t)
}

func TestBroadcastTx_GRPC_ClientError(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	txBytes := []byte("txBytes")
	mockGRPCClient.On("BroadcastTx", mock.Anything, &tx.BroadcastTxRequest{TxBytes: txBytes, Mode: tx.BroadcastMode_BROADCAST_MODE_SYNC}).Return(nil, errors.New("error"))

	txHash, err := client.BroadcastTx(txBytes)
	assert.Error(t, err)
	assert.Empty(t, txHash)

	mockGRPCClient.AssertExpectations(t)
}

func TestBroadcastTx_GRPC_TxFailed(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	txBytes := []byte("txBytes")
	mockGRPCClient.On("BroadcastTx", mock.Anything, &tx.BroadcastTxRequest{TxBytes: txBytes, Mode: tx.BroadcastMode_BROADCAST_MODE_SYNC}).Return(&tx.BroadcastTxResponse{TxResponse: &sdk.TxResponse{TxHash: "txHash", Code: 1}}, nil)

	txHash, err := client.BroadcastTx(txBytes)
	assert.Error(t, err)
	assert.Empty(t, txHash)

	mockGRPCClient.AssertExpectations(t)
}

func TestBroadcastTx_RPC(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	var txBytes ctypes.Tx = []byte("txBytes")
	mockHTTPClient.On("BroadcastTxSync", mock.Anything, txBytes).Return(&rpctypes.ResultBroadcastTx{Hash: []byte("txHash"), Code: 0}, nil)

	txHash, err := client.BroadcastTx(txBytes)
	assert.NoError(t, err)
	assert.Equal(t, hex.EncodeToString([]byte("txHash")), txHash)

	mockHTTPClient.AssertExpectations(t)
}

func TestBroadcastTx_RPC_ClientError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	var txBytes ctypes.Tx = []byte("txBytes")
	mockHTTPClient.On("BroadcastTxSync", mock.Anything, txBytes).Return(nil, errors.New("error"))

	txHash, err := client.BroadcastTx(txBytes)
	assert.Error(t, err)
	assert.Empty(t, txHash)

	mockHTTPClient.AssertExpectations(t)
}

func TestBroadcastTx_RPC_TxFailed(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	var txBytes ctypes.Tx = []byte("txBytes")
	mockHTTPClient.On("BroadcastTxSync", mock.Anything, txBytes).Return(&rpctypes.ResultBroadcastTx{Hash: []byte("txHash"), Code: 1}, nil)

	txHash, err := client.BroadcastTx(txBytes)
	assert.Error(t, err)
	assert.Empty(t, txHash)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTx_GRPC(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	txHash := "txHash"
	txResponse := &sdk.TxResponse{TxHash: txHash, Code: 0}
	mockGRPCClient.On("GetTx", mock.Anything, &tx.GetTxRequest{Hash: txHash}).Return(&tx.GetTxResponse{TxResponse: txResponse}, nil)

	result, err := client.GetTx(txHash)
	assert.NoError(t, err)
	assert.Equal(t, txResponse, result)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetTx_GRPC_Error(t *testing.T) {
	originalTxNewServiceClient := txNewServiceClient
	defer func() { txNewServiceClient = originalTxNewServiceClient }()

	mockGRPCClient := mocks.NewMockTxServiceClient(t)
	txNewServiceClient = func(conn grpc.ClientConn) tx.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		logger:        log.NewEntry(log.New()),
	}

	txHash := "txHash"
	mockGRPCClient.On("GetTx", mock.Anything, &tx.GetTxRequest{Hash: txHash}).Return(nil, errors.New("error"))

	result, err := client.GetTx(txHash)
	assert.Error(t, err)
	assert.Nil(t, result)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetTx_RPC(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	txHash := "0102"
	hashBytes, err := hex.DecodeString(txHash)
	assert.NoError(t, err)

	txResponse := &rpctypes.ResultTx{Hash: hashBytes, Height: 100}
	mockHTTPClient.On("Tx", mock.Anything, hashBytes, true).Return(txResponse, nil)
	mockHTTPClient.On("Block", mock.Anything, &txResponse.Height).Return(&rpctypes.ResultBlock{Block: &ctypes.Block{Header: ctypes.Header{Time: time.Now()}}}, nil)

	result, err := client.GetTx(txHash)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTx_RPC_DecodeError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	txHash := "hash"
	result, err := client.GetTx(txHash)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to decode hash")

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTx_RPC_ClientError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	txHash := "0102"
	hashBytes, err := hex.DecodeString(txHash)
	assert.NoError(t, err)

	mockHTTPClient.On("Tx", mock.Anything, hashBytes, true).Return(nil, errors.New("error"))

	result, err := client.GetTx(txHash)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get tx")

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTx_RPC_BlockError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	txHash := "0102"
	hashBytes, err := hex.DecodeString(txHash)
	assert.NoError(t, err)

	txResponse := &rpctypes.ResultTx{Hash: hashBytes, Height: 100}
	mockHTTPClient.On("Tx", mock.Anything, hashBytes, true).Return(txResponse, nil)
	mockHTTPClient.On("Block", mock.Anything, &txResponse.Height).Return(nil, errors.New("error"))

	result, err := client.GetTx(txHash)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get blocks for tx")

	mockHTTPClient.AssertExpectations(t)
}

func TestGetTx_RPC_FormatError(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	txHash := "0102"
	hashBytes, err := hex.DecodeString(txHash)
	assert.NoError(t, err)

	txResponse := &rpctypes.ResultTx{Hash: hashBytes, Height: 100}
	mockHTTPClient.On("Tx", mock.Anything, hashBytes, true).Return(txResponse, nil)
	mockHTTPClient.On("Block", mock.Anything, &txResponse.Height).Return(&rpctypes.ResultBlock{Block: &ctypes.Block{Header: ctypes.Header{Time: time.Now()}}}, nil)

	mockTx := mocks.NewMockAnyTx(t)
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return mockTx, errors.New("error")
	}
	utilNewTxDecoder = func(bech32Prefix string) sdk.TxDecoder {
		return txDecoder
	}
	defer func() {
		utilNewTxDecoder = util.NewTxDecoder
	}()

	result, err := client.GetTx(txHash)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to format tx result")

	mockHTTPClient.AssertExpectations(t)
}

func TestGetChainID_GRPC(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	block := &cmtservice.Block{Header: cmtservice.Header{Height: 100, ChainID: config.ChainID}}
	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(&cmtservice.GetLatestBlockResponse{SdkBlock: block}, nil)

	chainID, err := client.GetChainID()
	assert.NoError(t, err)
	assert.Equal(t, config.ChainID, chainID)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetChainID_GRPC_Error(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	chainID, err := client.GetChainID()
	assert.Error(t, err)
	assert.Equal(t, "", chainID)

	mockGRPCClient.AssertExpectations(t)
}

func TestGetChainID_RPC(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	status := &rpctypes.ResultStatus{SyncInfo: rpctypes.SyncInfo{LatestBlockHeight: 100}, NodeInfo: p2p.DefaultNodeInfo{Network: config.ChainID}}
	mockHTTPClient.On("Status", mock.Anything).Return(status, nil)

	chainID, err := client.GetChainID()
	assert.NoError(t, err)
	assert.Equal(t, config.ChainID, chainID)

	mockHTTPClient.AssertExpectations(t)
}

func TestGetChainID_RPC_Error(t *testing.T) {
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	mockHTTPClient.On("Status", mock.Anything).Return(nil, errors.New("error"))

	chainID, err := client.GetChainID()
	assert.Error(t, err)
	assert.Equal(t, "", chainID)

	mockHTTPClient.AssertExpectations(t)
}

func TestValidateNetwork(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		logger:        log.NewEntry(log.New()),
	}

	block := &cmtservice.Block{Header: cmtservice.Header{Height: 100, ChainID: config.ChainID}}
	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(&cmtservice.GetLatestBlockResponse{SdkBlock: block}, nil)

	err := client.ValidateNetwork()
	assert.NoError(t, err)

	mockGRPCClient.AssertExpectations(t)
}

func TestValidateNetwork_ErrorGetChainID(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(nil, errors.New("error getting chain id"))

	err := client.ValidateNetwork()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get latest block: error getting chain id")

	mockGRPCClient.AssertExpectations(t)
}

func TestValidateNetwork_InvalidChainID(t *testing.T) {
	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)
	mockHTTPClient := mocks.NewMockCosmosHTTPClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client := &cosmosClient{
		grpcEnabled:   config.GRPCEnabled,
		confirmations: config.Confirmations,
		timeout:       time.Duration(config.TimeoutMS) * time.Millisecond,
		chain:         models.Chain{ChainID: config.ChainID, ChainName: config.ChainName},
		grpcConn:      nil,
		rpcClient:     mockHTTPClient,
		logger:        log.NewEntry(log.New()),
	}

	block := &cmtservice.Block{Header: cmtservice.Header{Height: 100, ChainID: "InvalidChainID"}}
	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(&cmtservice.GetLatestBlockResponse{SdkBlock: block}, nil)

	err := client.ValidateNetwork()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected chain id TestChainID, got InvalidChainID")

	mockGRPCClient.AssertExpectations(t)
}

func TestNewClient(t *testing.T) {
	originalGRPCDial := grpcDial
	grpcDial = func(target string, opts ...goGRPC.DialOption) (*goGRPC.ClientConn, error) {
		return nil, nil
	}
	defer func() { grpcDial = originalGRPCDial }()

	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	block := &cmtservice.Block{Header: cmtservice.Header{Height: 100, ChainID: config.ChainID}}
	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(&cmtservice.GetLatestBlockResponse{SdkBlock: block}, nil)

	client, err := NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	mockGRPCClient.AssertExpectations(t)
}

func TestNewClient_ValidateError(t *testing.T) {
	originalGRPCDial := grpcDial
	grpcDial = func(target string, opts ...goGRPC.DialOption) (*goGRPC.ClientConn, error) {
		return nil, nil
	}
	defer func() { grpcDial = originalGRPCDial }()

	originalCmtserviceNewServiceClient := cmtserviceNewServiceClient
	defer func() { cmtserviceNewServiceClient = originalCmtserviceNewServiceClient }()

	mockGRPCClient := mocks.NewMockCMTServiceClient(t)

	cmtserviceNewServiceClient = func(conn grpc.ClientConn) cmtservice.ServiceClient {
		return mockGRPCClient
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	mockGRPCClient.On("GetLatestBlock", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	client, err := NewClient(config)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to validate network")

	mockGRPCClient.AssertExpectations(t)
}

func TestNewClient_GRPCERror(t *testing.T) {
	originalGRPCDial := grpcDial
	grpcDial = func(target string, opts ...goGRPC.DialOption) (*goGRPC.ClientConn, error) {
		return nil, errors.New("error")
	}
	defer func() { grpcDial = originalGRPCDial }()
	config := models.CosmosNetworkConfig{
		GRPCEnabled: true,
		GRPCHost:    "invalid",
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client, err := NewClient(config)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_RPCError(t *testing.T) {
	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		RPCURL:      "invalid",
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client, err := NewClient(config)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_RPCDialError(t *testing.T) {
	originalrpchttpNew := rpchttpNew
	defer func() { rpchttpNew = originalrpchttpNew }()
	rpchttpNew = func(url string, endpoint string) (CosmosHTTPClient, error) {
		return nil, errors.New("error")
	}

	config := models.CosmosNetworkConfig{
		GRPCEnabled: false,
		RPCURL:      "invalid",
		TimeoutMS:   5000,
		ChainName:   "TestChain",
		ChainID:     "TestChainID",
	}

	client, err := NewClient(config)
	assert.Error(t, err)
	assert.Nil(t, client)
}
