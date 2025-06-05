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
	MAX_QUERY_BLOCKS int64 = 499
)

type EthereumClient interface {
	ValidateNetwork()
	GetBlockNumber() (uint64, error)
	GetChainID() (*big.Int, error)
	GetClient() *ethclient.Client
	GetTransactionByHash(txHash string) (*types.Transaction, bool, error)
	GetTransactionReceipt(txHash string) (*types.Receipt, error)
}

type ethereumClient struct {
	client *ethclient.Client
}

var Client EthereumClient = &ethereumClient{}

func (c *ethereumClient) GetClient() *ethclient.Client {
	return c.client
}
func (c *ethereumClient) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
	defer cancel()

	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func (c *ethereumClient) GetChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
	defer cancel()

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return chainID, nil
}

func (c *ethereumClient) ValidateNetwork() {
	log.Debugln("[ETH]", "Validating network")
	log.Debugln("[ETH]", "uri", app.Config.Ethereum.RPCURL)
	client, err := ethclient.Dial(app.Config.Ethereum.RPCURL)
	if err != nil {
		log.Fatalln("[ETH]", "Failed to connect to Ethereum node:", err)
	}
	c.client = client

	chainID, err := c.GetChainID()
	if err != nil {
		log.Fatalln("[ETH]", "Failed to get chain ID:", err)
	}
	blockNumber, err := c.GetBlockNumber()
	if err != nil {
		log.Fatalln("[ETH]", "Failed to get block number:", err)
	}

	log.Debugln("[ETH]", "chainID", chainID.Uint64())

	if chainID.String() != app.Config.Ethereum.ChainID {
		log.Fatalln("[ETH]", "Chain ID Mismatch", "expected", app.Config.Ethereum.ChainID, "got", chainID.Uint64())
	}

	log.Debugln("[ETH]", "blockNumber", blockNumber)

	log.Infoln("[ETH]", "Validated network")
}

func (c *ethereumClient) GetTransactionByHash(txHash string) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
	defer cancel()

	tx, isPending, err := c.client.TransactionByHash(ctx, common.HexToHash(txHash))
	return tx, isPending, err
}

func (c *ethereumClient) GetTransactionReceipt(txHash string) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
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
