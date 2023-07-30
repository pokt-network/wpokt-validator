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
	stop               chan bool
	startBlockNumber   int64
	currentBlockNumber int64
	interval           time.Duration
	wpoktContract      *autogen.WrappedPocket
	mintControllerAbi  *abi.ABI
	client             eth.EthereumClient
	vaultAddress       string
	wpoktAddress       string

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *MintExecutorService) Start() {
	log.Info("[MINT EXECUTOR] Starting service")
	stop := false
	for !stop {
		log.Info("[MINT EXECUTOR] Starting sync")

		x.UpdateCurrentBlockNumber()

		x.SyncTxs()

		x.UpdateHealth()

		log.Info("[MINT EXECUTOR] Finished sync, Sleeping for ", x.interval)

		select {
		case <-x.stop:
			stop = true
			log.Info("[MINT EXECUTOR] Stopped service")
		case <-time.After(x.interval):
		}
	}
	x.wg.Done()
}

func (x *MintExecutorService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *MintExecutorService) UpdateHealth() {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           MintExecutorName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     "",
		EthBlockNumber: strconv.FormatInt(x.startBlockNumber, 10),
		Healthy:        true,
	}
}

func (x *MintExecutorService) Stop() {
	log.Debug("[MINT EXECUTOR] Stopping service")
	x.stop <- true
}

func (x *MintExecutorService) UpdateCurrentBlockNumber() {
	res, err := x.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}

	x.currentBlockNumber = int64(res)
	log.Info("[MINT EXECUTOR] Current block number: ", x.currentBlockNumber)
}

func (x *MintExecutorService) HandleMintEvent(event *autogen.WrappedPocketMinted) bool {
	log.Debug("[MINT EXECUTOR] Handling mint event: ", event.Raw.TxHash, " ", event.Raw.Index)

	filter := bson.M{
		"wpokt_address":     x.wpoktAddress,
		"vault_address":     x.vaultAddress,
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

func (x *MintExecutorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := x.wpoktContract.FilterMinted(&bind.FilterOpts{
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
		success = success && x.HandleMintEvent(filter.Event)
	}

	return success
}

func (x *MintExecutorService) SyncTxs() bool {

	if x.currentBlockNumber <= x.startBlockNumber {
		log.Info("[MINT EXECUTOR] No new blocks to sync")
		return true
	}

	var success bool = true

	if (x.currentBlockNumber - x.startBlockNumber) > eth.MAX_QUERY_BLOCKS {
		log.Debug("[MINT EXECUTOR] Syncing mint txs in chunks")

		for i := x.startBlockNumber; i < x.currentBlockNumber; i += eth.MAX_QUERY_BLOCKS {
			endBlockNumber := i + eth.MAX_QUERY_BLOCKS
			if endBlockNumber > x.currentBlockNumber {
				endBlockNumber = x.currentBlockNumber
			}

			log.Info("[MINT EXECUTOR] Syncing mint txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && x.SyncBlocks(uint64(i), uint64(endBlockNumber))
		}

	} else {
		log.Info("[MINT EXECUTOR] Syncing mint txs from blockNumber: ", x.startBlockNumber, " to blockNumber: ", x.currentBlockNumber)
		success = success && x.SyncBlocks(uint64(x.startBlockNumber), uint64(x.currentBlockNumber))
	}

	if success {
		x.startBlockNumber = x.currentBlockNumber
	}

	return success
}

func NewExecutor(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if app.Config.MintExecutor.Enabled == false {
		log.Debug("[MINT EXECUTOR] MINT executor disabled")
		return models.NewEmptyService(wg)
	}
	log.Debug("[MINT EXECUTOR] Initializing mint executor with last health")

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

	x := &MintExecutorService{
		wg:                 wg,
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

	x.UpdateCurrentBlockNumber()

	startBlockNumber := int64(app.Config.Ethereum.StartBlockNumber)

	if lastBlockNumber, err := strconv.ParseInt(lastHealth.EthBlockNumber, 10, 64); err == nil {
		startBlockNumber = lastBlockNumber
	}

	if startBlockNumber > 0 {
		x.startBlockNumber = startBlockNumber
	} else {
		log.Warn("[MINT EXECUTOR] Found invalid start block number, updating to current block number")
		x.startBlockNumber = x.currentBlockNumber
	}

	log.Info("[MINT EXECUTOR] Start block number: ", x.startBlockNumber)

	x.UpdateHealth()

	log.Info("[MINT EXECUTOR] Initialized mint executor")

	return x
}
