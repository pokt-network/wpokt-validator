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
	"github.com/dan13ram/wpokt-validator/eth/util"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	BurnMonitorName = "burn monitor"
)

type BurnMonitorService struct {
	wg                 *sync.WaitGroup
	stop               chan bool
	startBlockNumber   int64
	currentBlockNumber int64
	interval           time.Duration
	wpoktContract      *autogen.WrappedPocket
	client             eth.EthereumClient

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *BurnMonitorService) Start() {
	log.Info("[BURN MONITOR] Starting service")
	stop := false
	for !stop {
		log.Info("[BURN MONITOR] Starting sync")

		x.UpdateCurrentBlockNumber()

		x.SyncTxs()

		x.UpdateHealth()

		log.Info("[BURN MONITOR] Finished sync, Sleeping for ", x.interval)

		select {
		case <-x.stop:
			stop = true
			log.Info("[BURN MONITOR] Stopped service")
		case <-time.After(x.interval):
		}
	}
	x.wg.Done()
}

func (x *BurnMonitorService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *BurnMonitorService) UpdateHealth() {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           BurnMonitorName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     "",
		EthBlockNumber: strconv.FormatInt(x.startBlockNumber, 10),
		Healthy:        true,
	}
}

func (x *BurnMonitorService) Stop() {
	log.Debug("[BURN MONITOR] Stopping service")
	x.stop <- true
}

func (x *BurnMonitorService) UpdateCurrentBlockNumber() {
	res, err := x.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	x.currentBlockNumber = int64(res)
	log.Info("[BURN MONITOR] Current block number: ", x.currentBlockNumber)
}

func (x *BurnMonitorService) HandleBurnEvent(event *autogen.WrappedPocketBurnAndBridge) bool {
	doc := util.CreateBurn(event)

	// each event is a combination of transaction hash and log index
	log.Debug("[BURN MONITOR] Handling burn event: ", event.Raw.TxHash, " ", event.Raw.Index)

	err := app.DB.InsertOne(models.CollectionBurns, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Info("[BURN MONITOR] Found duplicate burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
			return true
		}
		log.Error("[BURN MONITOR] Error while storing burn event in db: ", err)
		return false
	}

	log.Info("[BURN MONITOR] Stored burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
	return true
}

func (x *BurnMonitorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := x.wpoktContract.FilterBurnAndBridge(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []*big.Int{}, []common.Address{}, []common.Address{})

	if err != nil {
		log.Error("[BURN MONITOR] Error while syncing burn events: ", err)
		return false
	}

	var success bool = true
	for filter.Next() {
		success = success && x.HandleBurnEvent(filter.Event)
	}
	return success
}

func (x *BurnMonitorService) SyncTxs() bool {
	if x.currentBlockNumber <= x.startBlockNumber {
		log.Info("[BURN MONITOR] [MINT EXECUTOR] No new blocks to sync")
		return true
	}

	var success bool = true
	if (x.currentBlockNumber - x.startBlockNumber) > eth.MAX_QUERY_BLOCKS {
		log.Debug("[BURN MONITOR] Syncing burn txs in chunks")
		for i := x.startBlockNumber; i < x.currentBlockNumber; i += eth.MAX_QUERY_BLOCKS {
			endBlockNumber := i + eth.MAX_QUERY_BLOCKS
			if endBlockNumber > x.currentBlockNumber {
				endBlockNumber = x.currentBlockNumber
			}
			log.Info("[BURN MONITOR] Syncing burn txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && x.SyncBlocks(uint64(i), uint64(endBlockNumber))
		}
	} else {
		log.Info("[BURN MONITOR] Syncing burn txs from blockNumber: ", x.startBlockNumber, " to blockNumber: ", x.currentBlockNumber)
		success = success && x.SyncBlocks(uint64(x.startBlockNumber), uint64(x.currentBlockNumber))
	}

	if success {
		x.startBlockNumber = x.currentBlockNumber
	}

	return success
}

func NewMonitor(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if app.Config.BurnMonitor.Enabled == false {
		log.Debug("[BURN MONITOR] BURN monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[BURN MONITOR] Initializing burn monitor")
	client, err := eth.NewClient()
	if err != nil {
		log.Fatal("[BURN MONITOR] Error initializing ethereum client: ", err)
	}
	log.Debug("[BURN MONITOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WrappedPocketAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WrappedPocketAddress), client.GetClient())
	if err != nil {
		log.Fatal("[BURN MONITOR] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[BURN MONITOR] Connected to wpokt contract")

	x := &BurnMonitorService{
		wg:                 wg,
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		interval:           time.Duration(app.Config.BurnMonitor.IntervalSecs) * time.Second,
		wpoktContract:      contract,
		client:             client,
	}

	if app.Config.BurnMonitor.Enabled == false {
		log.Debug("[BURN MONITOR] burn monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[BURN MONITOR] Initializing burn monitor")

	x.UpdateCurrentBlockNumber()

	startBlockNumber := int64(app.Config.Ethereum.StartBlockNumber)

	if lastBlockNumber, err := strconv.ParseInt(lastHealth.EthBlockNumber, 10, 64); err == nil {
		startBlockNumber = lastBlockNumber
	}

	if startBlockNumber > 0 {
		x.startBlockNumber = startBlockNumber
	} else {
		log.Warn("Found invalid start block number, updating to current block number")
		x.startBlockNumber = x.currentBlockNumber
	}

	log.Info("[BURN MONITOR] Start block number: ", x.startBlockNumber)

	x.UpdateHealth()

	log.Info("[BURN MONITOR] Initialized burn monitor")

	return x
}
