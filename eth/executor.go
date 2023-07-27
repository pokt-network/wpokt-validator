package eth

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MintExecutorName string = "mint executor"
)

type MintExecutorService struct {
	wg                 *sync.WaitGroup
	name               string
	stop               chan bool
	startBlockNumber   int64
	currentBlockNumber int64
	lastSyncTime       time.Time
	interval           time.Duration
	wpoktContract      *autogen.WrappedPocket
	mintControllerAbi  *abi.ABI
	client             eth.EthereumClient
	vaultAddress       string
	wpoktAddress       string
}

func (b *MintExecutorService) Start() {
	log.Info("[MINT EXECUTOR] Starting service")
	stop := false
	for !stop {
		log.Info("[MINT EXECUTOR] Starting sync")
		b.lastSyncTime = time.Now()

		b.UpdateCurrentBlockNumber()

		if (b.currentBlockNumber - b.startBlockNumber) > 0 {
			success := b.SyncTxs()
			if success {
				b.startBlockNumber = b.currentBlockNumber
			}
		} else {
			log.Info("[MINT EXECUTOR] No new blocks to sync")
		}

		log.Info("[MINT EXECUTOR] Finished sync, Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Info("[MINT EXECUTOR] Stopped service")
		case <-time.After(b.interval):
		}
	}
	b.wg.Done()
}

func (m *MintExecutorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.name,
		LastSyncTime:   m.lastSyncTime,
		NextSyncTime:   m.lastSyncTime.Add(m.interval),
		PoktHeight:     "",
		EthBlockNumber: strconv.FormatInt(m.startBlockNumber, 10),
		Healthy:        true,
	}
}

func (b *MintExecutorService) Stop() {
	log.Debug("[MINT EXECUTOR] Stopping service")
	b.stop <- true
}

func (b *MintExecutorService) InitStartBlockNumber(startBlockNumber int64) {
	if app.Config.Ethereum.StartBlockNumber > 0 {
		b.startBlockNumber = int64(app.Config.Ethereum.StartBlockNumber)
	} else {
		log.Warn("[MINT EXECUTOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Info("[MINT EXECUTOR] Start block number: ", b.startBlockNumber)
}

func (b *MintExecutorService) UpdateCurrentBlockNumber() {
	res, err := b.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = int64(res)
	log.Info("[MINT EXECUTOR] Current block number: ", b.currentBlockNumber)
}

func (b *MintExecutorService) HandleMintEvent(event *autogen.WrappedPocketMinted) bool {
	log.Debug("[MINT EXECUTOR] Handling mint event: ", event.Raw.TxHash, " ", event.Raw.Index)

	filter := bson.M{
		"wpokt_address":     b.wpoktAddress,
		"vault_address":     b.vaultAddress,
		"recipient_address": event.Recipient.Hex(),
		"amount":            event.Amount.String(),
		"nonce":             event.Nonce.String(),
		"status": bson.M{
			"$in": []string{models.StatusConfirmed, models.StatusSigned},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":       models.StatusSuccess,
			"mint_tx_hash": event.Raw.TxHash.String(),
			"updated_at":   time.Now(),
		},
	}

	err := app.DB.UpdateOne(models.CollectionMints, filter, update)

	if err != nil {
		log.Error("[MINT EXECUTOR] Error while updating mint: ", err)
		return false
	}

	log.Info("[MINT EXECUTOR] Mint event handled successfully")

	return true
}

func (b *MintExecutorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := b.wpoktContract.FilterMinted(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []common.Address{}, []*big.Int{}, []*big.Int{})

	if err != nil {
		log.Errorln("[MINT EXECUTOR] Error while syncing mint events: ", err)
		return false
	}

	var success bool = true
	for filter.Next() {
		success = success && b.HandleMintEvent(filter.Event)
	}

	return success
}

func (b *MintExecutorService) SyncTxs() bool {
	var success bool = true

	if (b.currentBlockNumber - b.startBlockNumber) > eth.MAX_QUERY_BLOCKS {
		log.Debug("[MINT EXECUTOR] Syncing mint txs in chunks")

		for i := b.startBlockNumber; i < b.currentBlockNumber; i += eth.MAX_QUERY_BLOCKS {
			endBlockNumber := i + eth.MAX_QUERY_BLOCKS
			if endBlockNumber > b.currentBlockNumber {
				endBlockNumber = b.currentBlockNumber
			}

			log.Info("[MINT EXECUTOR] Syncing mint txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && b.SyncBlocks(uint64(i), uint64(endBlockNumber))
		}

	} else {
		log.Info("[MINT EXECUTOR] Syncing mint txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
		success = success && b.SyncBlocks(uint64(b.startBlockNumber), uint64(b.currentBlockNumber))
	}
	return success
}

func newExecutor(wg *sync.WaitGroup) *MintExecutorService {
	client, err := eth.NewClient()
	if err != nil {
		log.Fatal("[MINT EXECUTOR] Error initializing ethereum client", err)
	}

	log.Debug("[MINT EXECUTOR] Connecting to mint contract at: ", app.Config.Ethereum.WrappedPocketAddress)

	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WrappedPocketAddress), client.GetClient())
	if err != nil {
		log.Fatal("[MINT EXECUTOR] Error initializing Wrapped Pocket contract", err)
	}

	log.Debug("[MINT EXECUTOR] Connected to mint contract")

	mintControllerAbi, err := autogen.MintControllerMetaData.GetAbi()
	if err != nil {
		log.Fatal("[MINT EXECUTOR] Error parsing MintController ABI", err)
	}

	log.Debug("[MINT EXECUTOR] Mint controller abi parsed")

	b := &MintExecutorService{
		wg:                 wg,
		name:               MintExecutorName,
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		interval:           time.Duration(app.Config.MintExecutor.IntervalSecs) * time.Second,
		wpoktContract:      contract,
		mintControllerAbi:  mintControllerAbi,
		client:             client,
		wpoktAddress:       app.Config.Ethereum.WrappedPocketAddress,
		vaultAddress:       app.Config.Pocket.VaultAddress,
	}

	return b
}

func NewExecutor(wg *sync.WaitGroup) models.Service {
	if app.Config.MintExecutor.Enabled == false {
		log.Debug("[MINT EXECUTOR] MINT executor disabled")
		return models.NewEmptyService(wg)
	}
	log.Debug("[MINT EXECUTOR] Initializing mint executor")

	m := newExecutor(wg)

	m.UpdateCurrentBlockNumber()

	m.InitStartBlockNumber(int64(app.Config.Ethereum.StartBlockNumber))

	log.Info("[MINT EXECUTOR] Initialized mint executor")

	return m
}

func NewExecutorWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if app.Config.MintExecutor.Enabled == false {
		log.Debug("[MINT EXECUTOR] MINT executor disabled")
		return models.NewEmptyService(wg)
	}
	log.Debug("[MINT EXECUTOR] Initializing mint executor with last health")

	m := newExecutor(wg)

	lastBlockNumber, err := strconv.ParseInt(lastHealth.EthBlockNumber, 10, 64)
	if err != nil {
		log.Error("[MINT EXECUTOR] Error parsing last block number from last health", err)
		lastBlockNumber = app.Config.Ethereum.StartBlockNumber
	}

	m.UpdateCurrentBlockNumber()

	m.InitStartBlockNumber(lastBlockNumber)

	log.Info("[MINT EXECUTOR] Initialized mint executor")

	return m
}
