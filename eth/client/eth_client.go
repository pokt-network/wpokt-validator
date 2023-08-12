package client

import (
	"context"
	"time"

	"math/big"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

const (
	MAX_QUERY_BLOCKS int64 = 100000
)

type EthereumClient interface {
	ValidateNetwork()
	GetBlockNumber() (uint64, error)
	GetChainId() (*big.Int, error)
	GetClient() *ethclient.Client
	GetTransactionByHash(txHash string) (*types.Transaction, bool, error)
	GetTransactionReceipt(txHash string) (*types.Receipt, error)
}

type ethereumClient struct {
	client  *ethclient.Client
	timeout time.Duration
}

var Client EthereumClient = &ethereumClient{}

func (c *ethereumClient) GetClient() *ethclient.Client {
	return c.client
}
func (c *ethereumClient) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func (c *ethereumClient) GetChainId() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	chainId, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return chainId, nil
}

func (c *ethereumClient) ValidateNetwork() {
	log.Debugln("[ETH]", "Validating network")
	log.Debugln("[ETH]", "uri", app.Config.Ethereum.RPCURL)
	client, err := ethclient.Dial(app.Config.Ethereum.RPCURL)
	if err != nil {
		log.Fatalln("[ETH]", "Failed to connect to Ethereum node:", err)
	}
	c.client = client
	c.timeout = time.Duration(app.Config.Ethereum.RPCTimeoutMillis) * time.Millisecond

	chainId, err := c.GetChainId()
	if err != nil {
		log.Fatalln("[ETH]", "Failed to get chain ID:", err)
	}
	blockNumber, err := c.GetBlockNumber()
	if err != nil {
		log.Fatalln("[ETH]", "Failed to get block number:", err)
	}

	log.Debugln("[ETH]", "chainId", chainId.Uint64())

	if chainId.String() != app.Config.Ethereum.ChainId {
		log.Fatalln("[ETH]", "Chain ID Mismatch", "expected", app.Config.Ethereum.ChainId, "got", chainId.Uint64())
	}

	log.Debugln("[ETH]", "blockNumber", blockNumber)

	log.Infoln("[ETH]", "Validated network")
}

func (c *ethereumClient) GetTransactionByHash(txHash string) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tx, isPending, err := c.client.TransactionByHash(ctx, common.HexToHash(txHash))
	return tx, isPending, err
}

func (c *ethereumClient) GetTransactionReceipt(txHash string) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	receipt, err := c.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	return receipt, err
}

func NewClient() (EthereumClient, error) {
	client, err := ethclient.Dial(app.Config.Ethereum.RPCURL)
	return &ethereumClient{
		client: client,
	}, err
}
