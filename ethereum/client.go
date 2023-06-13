package ethereum

import (
	"context"
	"time"

	"math/big"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

type EthereumClient interface {
	ValidateNetwork()
	GetBlockNumber() (uint64, error)
	GetChainId() (*big.Int, error)
	GetClient() *ethclient.Client
}

type ethereumClient struct {
	client *ethclient.Client
}

var Client EthereumClient = &ethereumClient{}

func (c *ethereumClient) GetClient() *ethclient.Client {
	return c.client
}

func (c *ethereumClient) ValidateNetwork() {
	log.Debugln("Connecting to Ethereum node", "url", app.Config.Ethereum.RPCURL)
	client, err := ethclient.Dial(app.Config.Ethereum.RPCURL)
	if err != nil {
		panic(err)
	}
	c.client = client

	blockNumber, err := c.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	log.Debugln("Connected to Ethereum node", "blockNumber", blockNumber)

	chainId, err := c.GetChainId()
	if err != nil {
		panic(err)
	}

	if chainId.Uint64() != app.Config.Ethereum.ChainId {
		log.Debugln("ethereum chainId mismatch", "config", app.Config.Ethereum.ChainId, "node", chainId.Int64())
		panic("ethereum chain id mismatch")
	}
	log.Debugln("Connected to Ethereum node", "chainId", chainId.String())
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
