package client

import (
	"encoding/hex"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/libs/bytes"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	rpctypes "github.com/cometbft/cometbft/rpc/core/types"
	ctypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	"github.com/dan13ram/wpokt-validator/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"

	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	maxPageDepth = 100
)

type CosmosClient interface {
	Confirmations() uint64
	GetLatestBlockHeight() (int64, error)
	GetChainID() (string, error)
	GetTxsSentFromAddressAfterHeight(address string, height uint64) ([]*sdk.TxResponse, error)
	GetTxsSentToAddressAfterHeight(address string, height uint64) ([]*sdk.TxResponse, error)
	GetAccount(address string) (*auth.BaseAccount, error)
	Simulate(txBytes []byte) (*sdk.GasInfo, error)
	BroadcastTx(txBytes []byte) (string, error)
	GetTx(hash string) (*sdk.TxResponse, error)
	ValidateNetwork() error
}

type CosmosHTTPClient interface {
	Block(ctx context.Context, height *int64) (*rpctypes.ResultBlock, error)
	Status(ctx context.Context) (*rpctypes.ResultStatus, error)
	Tx(ctx context.Context, hash []byte, prove bool) (*rpctypes.ResultTx, error)
	TxSearch(ctx context.Context, query string, prove bool, page *int, limit *int, orderBy string) (*rpctypes.ResultTxSearch, error)
	ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*rpctypes.ResultABCIQuery, error)
	BroadcastTxSync(ctx context.Context, tx ctypes.Tx) (*rpctypes.ResultBroadcastTx, error)
	CheckTx(ctx context.Context, tx ctypes.Tx) (*rpctypes.ResultCheckTx, error)
}

type cosmosClient struct {
	grpcEnabled   bool
	confirmations uint64

	timeout      time.Duration
	config       models.PocketConfig
	bech32Prefix string
	coinDenom    string

	grpcConn  *grpc.ClientConn
	rpcClient CosmosHTTPClient

	logger *log.Entry
}

var cmtserviceNewServiceClient = cmtservice.NewServiceClient
var authNewQueryClient = auth.NewQueryClient
var txNewServiceClient = tx.NewServiceClient

func (c *cosmosClient) Confirmations() uint64 {
	return c.confirmations
}

func (c *cosmosClient) getLatestBlockGRPC() (*cmtservice.Block, error) {
	client := cmtserviceNewServiceClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &cmtservice.GetLatestBlockRequest{}

	resp, err := client.GetLatestBlock(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.SdkBlock, nil
}

func (c *cosmosClient) getStatusRPC() (*rpctypes.ResultStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.rpcClient.Status(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *cosmosClient) GetLatestBlockHeight() (int64, error) {
	if c.grpcEnabled {
		block, err := c.getLatestBlockGRPC()
		if err != nil {
			return 0, err
		}
		return block.Header.Height, nil
	}

	status, err := c.getStatusRPC()

	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil

}

func (c *cosmosClient) GetTxsSentToAddressAfterHeight(address string, height uint64) ([]*sdk.TxResponse, error) {
	if !common.IsValidBech32Address(c.bech32Prefix, address) {
		return nil, fmt.Errorf("invalid bech32 address")
	}

	query := fmt.Sprintf("transfer.recipient='%s' AND tx.height>=%d", address, height)

	return c.getTxsByEvents(query)
}

func (c *cosmosClient) GetTxsSentFromAddressAfterHeight(address string, height uint64) ([]*sdk.TxResponse, error) {
	if !common.IsValidBech32Address(c.bech32Prefix, address) {
		return nil, fmt.Errorf("invalid bech32 address")
	}

	query := fmt.Sprintf("transfer.sender='%s' AND tx.height>=%d", address, height)

	return c.getTxsByEvents(query)
}

func (c *cosmosClient) getTxsByEventsPerPageGRPC(query string, page uint64) ([]*sdk.TxResponse, uint64, error) {
	client := txNewServiceClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &tx.GetTxsEventRequest{
		Query:   query,
		OrderBy: tx.OrderBy_ORDER_BY_ASC,
		Page:    page,
		Limit:   50,
	}

	resp, err := client.GetTxsEvent(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return resp.TxResponses, resp.Total, nil
}

func (c *cosmosClient) getTxsByEventsPerPageRPC(query string, page uint64) ([]*sdk.TxResponse, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	limit := 50
	pageint := int(page)

	resTxs, err := c.rpcClient.TxSearch(ctx, query, false, &pageint, &limit, "asc")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get txs: %s", err)
	}

	resBlocks, err := getBlocksForTxResults(c.rpcClient, resTxs.Txs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get blocks for txs: %s", err)
	}

	txs, err := formatTxResults(c.bech32Prefix, resTxs.Txs, resBlocks)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to format tx results: %s", err)
	}

	return txs, uint64(resTxs.TotalCount), err
}

func (c *cosmosClient) getTxsByEvents(query string) ([]*sdk.TxResponse, error) {
	var page uint64 = 1
	var txs = make([]*sdk.TxResponse, 0)
	for {

		var respTxs []*sdk.TxResponse
		var err error
		var total uint64

		if c.grpcEnabled {
			respTxs, total, err = c.getTxsByEventsPerPageGRPC(query, page)
		} else {
			respTxs, total, err = c.getTxsByEventsPerPageRPC(query, page)
		}

		if err != nil {
			return nil, err
		}

		txs = append(txs, respTxs...)

		if len(respTxs) == 0 || len(txs) >= int(total) || page >= maxPageDepth {
			break
		}
		page++
	}

	return txs, nil
}

func (c *cosmosClient) getTxGRPC(hash string) (*sdk.TxResponse, error) {
	client := txNewServiceClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &tx.GetTxRequest{
		Hash: hash,
	}

	resp, err := client.GetTx(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %s", err)
	}

	return resp.TxResponse, nil
}

func (c *cosmosClient) getTxRPC(hash string) (*sdk.TxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash: %s", err)
	}

	resTx, err := c.rpcClient.Tx(ctx, hashBytes, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx: %s", err)
	}

	resBlocks, err := getBlocksForTxResults(c.rpcClient, []*rpctypes.ResultTx{resTx})
	if err != nil {
		return nil, fmt.Errorf("failed to get blocks for tx: %s", err)
	}

	out, err := mkTxResult(c.bech32Prefix, resTx, resBlocks[resTx.Height])
	if err != nil {
		return nil, fmt.Errorf("failed to format tx result: %s", err)
	}

	return out, nil
}

func (c *cosmosClient) GetTx(hash string) (*sdk.TxResponse, error) {
	hash = strings.TrimPrefix(hash, "0x")
	if c.grpcEnabled {
		return c.getTxGRPC(hash)
	}
	return c.getTxRPC(hash)
}

func (c *cosmosClient) getAccountGRPC(address string) (*auth.BaseAccount, error) {
	client := authNewQueryClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := auth.QueryAccountRequest{
		Address: address,
	}

	resp, err := client.Account(ctx, &req)
	if err != nil {
		return nil, err
	}

	var account auth.BaseAccount
	if err := account.Unmarshal(resp.Account.Value); err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *cosmosClient) getAccountRPC(address string) (*auth.BaseAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	reqBz, _ := util.NewProtoCodec(c.bech32Prefix).Marshal(&auth.QueryAccountRequest{Address: address}) // no reason to fail since account address is validated

	res, err := c.rpcClient.ABCIQuery(ctx, "/cosmos.auth.v1beta1.Query/Account", reqBz)

	if err != nil {
		return nil, fmt.Errorf("failed to get account: %s", err)
	}

	if res.Response.Code != 0 {
		return nil, fmt.Errorf("failed to get account, got code %d: %s", res.Response.Code, res.Response.Log)
	}

	var account auth.QueryAccountResponse
	if err := account.Unmarshal(res.Response.Value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %s", err)
	}

	var baseAccount auth.BaseAccount
	if err := baseAccount.Unmarshal(account.Account.Value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base account: %s", err)
	}

	return &baseAccount, nil
}

func (c *cosmosClient) GetAccount(address string) (*auth.BaseAccount, error) {
	if !common.IsValidBech32Address(c.bech32Prefix, address) {
		return nil, fmt.Errorf("invalid bech32 address")
	}
	if c.grpcEnabled {
		return c.getAccountGRPC(address)
	}
	return c.getAccountRPC(address)
}

func (c *cosmosClient) broadcastTxGRPC(txBytes []byte) (string, error) {
	client := txNewServiceClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	}

	resp, err := client.BroadcastTx(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast tx: %s", err)
	}

	if resp.TxResponse.Code != 0 {
		return "", fmt.Errorf("failed to broadcast tx, got code %d: %s", resp.TxResponse.Code, resp.TxResponse.RawLog)
	}

	return resp.TxResponse.TxHash, nil
}

func (c *cosmosClient) broadcastTxRPC(txBytes []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.rpcClient.BroadcastTxSync(ctx, txBytes)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast tx: %s", err)
	}

	if res.Code != 0 {
		return "", fmt.Errorf("failed to broadcast tx, got code %d: %s", res.Code, res.Log)
	}

	return res.Hash.String(), nil
}

func (c *cosmosClient) BroadcastTx(txBytes []byte) (string, error) {
	if c.grpcEnabled {
		return c.broadcastTxGRPC(txBytes)
	}
	return c.broadcastTxRPC(txBytes)
}

func (c *cosmosClient) simulateGRPC(txBytes []byte) (*sdk.GasInfo, error) {
	client := txNewServiceClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &tx.SimulateRequest{
		TxBytes: txBytes,
	}

	resp, err := client.Simulate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to simulate tx: %s", err)
	}

	return resp.GasInfo, nil
}

func (c *cosmosClient) simulateRPC(txBytes []byte) (*sdk.GasInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.rpcClient.CheckTx(ctx, txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to simulate tx: %s", err)
	}

	if res.Code != 0 {
		return nil, fmt.Errorf("failed to simulate tx, got code %d: %s", res.Code, res.Log)
	}

	return &sdk.GasInfo{
		GasWanted: uint64(res.GasWanted),
		GasUsed:   uint64(res.GasUsed),
	}, nil
}

func (c *cosmosClient) Simulate(txBytes []byte) (*sdk.GasInfo, error) {
	if c.grpcEnabled {
		return c.simulateGRPC(txBytes)
	}
	return c.simulateRPC(txBytes)
}

func (c *cosmosClient) GetChainID() (string, error) {
	var chainID string
	if c.grpcEnabled {
		res, err := c.getLatestBlockGRPC()
		if err != nil {
			return "", fmt.Errorf("failed to get latest block: %s", err)
		}
		chainID = res.Header.ChainID
	} else {
		status, err := c.getStatusRPC()
		if err != nil {
			return "", fmt.Errorf("failed to get status: %s", err)
		}
		chainID = status.NodeInfo.Network
	}

	return chainID, nil
}

func (c *cosmosClient) ValidateNetwork() error {
	c.logger.Debugf("Validating network")
	chainID, err := c.GetChainID()
	if err != nil {
		return err
	}
	if chainID != c.config.ChainID {
		return fmt.Errorf("expected chain id %s, got %s", c.config.ChainID, chainID)
	}
	c.logger.Debugf("Validated network")
	return nil
}

var grpcNewClient = grpc.NewClient
var rpchttpNew = func(url, endpoint string) (CosmosHTTPClient, error) {
	return rpchttp.New(url, endpoint)
}

func NewClient(config models.PocketConfig) (CosmosClient, error) {
	var connection *grpc.ClientConn
	var client CosmosHTTPClient

	logger := log.
		WithField("module", "cosmos").
		WithField("package", "client").
		WithField("chain_id", strings.ToLower(config.ChainID))

	if config.GRPCEnabled {
		grpcURL := fmt.Sprintf("%s:%d", config.GRPCHost, config.GRPCPort)
		conn, err := grpcNewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.WithError(err).Error("failed to connect to grpc")
			return nil, fmt.Errorf("failed to connect to grpc")
		}
		connection = conn
		client = nil
	} else {
		c, err := rpchttpNew(config.RPCURL, "/websocket")
		if err != nil {
			logger.WithError(err).Error("failed to connect to rpc")
			return nil, fmt.Errorf("failed to connect to rpc")
		}
		client = c
		connection = nil
	}

	c := &cosmosClient{
		grpcEnabled: config.GRPCEnabled,
		config:      config,

		timeout:      time.Duration(config.RPCTimeoutMillis) * time.Millisecond,
		bech32Prefix: config.Bech32Prefix,
		coinDenom:    config.CoinDenom,

		confirmations: uint64(config.Confirmations),

		grpcConn:  connection,
		rpcClient: client,

		logger: logger,
	}

	err := c.ValidateNetwork()
	if err != nil {
		logger.WithError(err).Error("failed to validate network")
		return nil, fmt.Errorf("failed to validate network")
	}

	return c, nil
}
