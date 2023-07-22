package ethereum

import (
	"context"
	"time"

	"math/big"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

type EthereumClient interface {
	ValidateNetwork()
	GetBlockNumber() (uint64, error)
	GetChainId() (*big.Int, error)
	GetClient() *ethclient.Client
	GetTransactionByHash(txHash string) (*types.Transaction, bool, error)
}

type ethereumClient struct {
	client *ethclient.Client
}

var Client EthereumClient = &ethereumClient{}

func (c *ethereumClient) GetClient() *ethclient.Client {
	return c.client
}
func (c *ethereumClient) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
	defer cancel()

	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func (c *ethereumClient) GetChainId() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
	defer cancel()

	chainId, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return chainId, nil
}

func (c *ethereumClient) ValidateNetwork() {
	log.Debugln("[ETH]", "Validating network")
	log.Debugln("[ETH]", "URL", app.Config.Ethereum.RPCURL)
	client, err := ethclient.Dial(app.Config.Ethereum.RPCURL)
	if err != nil {
		panic(err)
	}
	c.client = client

	blockNumber, err := c.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	log.Debugln("[ETH]", "Validating network", "blockNumber", blockNumber)

	chainId, err := c.GetChainId()
	if err != nil {
		panic(err)
	}
	log.Debugln("[ETH]", "Validating network", "chainId", chainId.Uint64())

	if chainId.String() != app.Config.Ethereum.ChainId {
		log.Debugln("[ETH]", "Chain ID Mismatch", "expected", app.Config.Ethereum.ChainId, "got", chainId.Uint64())
		panic("[ETH] Chain ID Mismatch")
	}
	log.Debugln("[ETH]", "Validated network")
}

func (c *ethereumClient) GetTransactionByHash(txHash string) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
	defer cancel()

	tx, isPending, err := c.client.TransactionByHash(ctx, common.HexToHash(txHash))
	return tx, isPending, err
}
