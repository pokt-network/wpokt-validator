package ethereum

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	ethereum "github.com/dan13ram/wpokt-backend/ethereum/client"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	BurnMonitorName         = "burn monitor"
	MAX_QUERY_BLOCKS int64  = 100000
	ZERO_ADDRESS     string = "0x0000000000000000000000000000000000000000"
)

type BurnMonitorService struct {
	wg                 *sync.WaitGroup
	name               string
	stop               chan bool
	startBlockNumber   int64
	currentBlockNumber int64
	lastSyncTime       time.Time
	interval           time.Duration
	wpoktContract      *autogen.WrappedPocket
	client             ethereum.EthereumClient
}

func (b *BurnMonitorService) Start() {
	log.Info("[BURN MONITOR] Starting service")
	stop := false
	for !stop {
		log.Info("[BURN MONITOR] Starting sync")
		b.lastSyncTime = time.Now()

		b.UpdateCurrentBlockNumber()

		if (b.currentBlockNumber - b.startBlockNumber) > 0 {
			success := b.SyncTxs()
			if success {
				b.startBlockNumber = b.currentBlockNumber
			}
		} else {
			log.Info("[BURN MONITOR] No new blocks to sync")
		}

		log.Info("[BURN MONITOR] Finished sync, Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Info("[BURN MONITOR] Stopped service")
		case <-time.After(b.interval):
		}
	}
	b.wg.Done()
}

func (b *BurnMonitorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           b.name,
		LastSyncTime:   b.lastSyncTime,
		NextSyncTime:   b.lastSyncTime.Add(b.interval),
		PoktHeight:     "",
		EthBlockNumber: strconv.FormatInt(b.startBlockNumber, 10),
		Healthy:        true,
	}
}

func (b *BurnMonitorService) Stop() {
	log.Debug("[BURN MONITOR] Stopping service")
	b.stop <- true
}

func (b *BurnMonitorService) InitStartBlockNumber(startBlockNumber int64) {
	if app.Config.Ethereum.StartBlockNumber > 0 {
		b.startBlockNumber = int64(app.Config.Ethereum.StartBlockNumber)
	} else {
		log.Warn("[BURN MONITOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Info("[BURN EXECUTOR] Start block number: ", b.startBlockNumber)
}

func (b *BurnMonitorService) UpdateCurrentBlockNumber() {
	res, err := b.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = int64(res)
	log.Info("[BURN MONITOR] Current block number: ", b.currentBlockNumber)
}

func (b *BurnMonitorService) HandleBurnEvent(event *autogen.WrappedPocketBurnAndBridge) bool {
	doc := createBurn(event)

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

func (b *BurnMonitorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := b.wpoktContract.FilterBurnAndBridge(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []*big.Int{}, []common.Address{}, []common.Address{})

	if err != nil {
		log.Errorln("[BURN MONITOR] Error while syncing burn events: ", err)
		return false
	}

	var success bool = true
	for filter.Next() {
		success = success && b.HandleBurnEvent(filter.Event)
	}
	return success
}

func (b *BurnMonitorService) SyncTxs() bool {
	var success bool = true
	if (b.currentBlockNumber - b.startBlockNumber) > MAX_QUERY_BLOCKS {
		log.Debug("[BURN MONITOR] Syncing burn txs in chunks")
		for i := b.startBlockNumber; i < b.currentBlockNumber; i += MAX_QUERY_BLOCKS {
			endBlockNumber := i + MAX_QUERY_BLOCKS
			if endBlockNumber > b.currentBlockNumber {
				endBlockNumber = b.currentBlockNumber
			}
			log.Info("[BURN MONITOR] Syncing burn txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && b.SyncBlocks(uint64(i), uint64(endBlockNumber))
		}
	} else {
		log.Info("[BURN MONITOR] Syncing burn txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
		success = success && b.SyncBlocks(uint64(b.startBlockNumber), uint64(b.currentBlockNumber))
	}
	return success
}

func newMonitor(wg *sync.WaitGroup) *BurnMonitorService {
	client, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[BURN MONITOR] Error initializing ethereum client: ", err)
	}
	log.Debug("[BURN MONITOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTAddress), client.GetClient())
	if err != nil {
		log.Fatal("[BURN MONITOR] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[BURN MONITOR] Connected to wpokt contract")

	b := &BurnMonitorService{
		wg:                 wg,
		name:               BurnMonitorName,
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		interval:           time.Duration(app.Config.BurnMonitor.IntervalSecs) * time.Second,
		wpoktContract:      contract,
		client:             client,
	}

	return b
}

func NewMonitor(wg *sync.WaitGroup) models.Service {
	if app.Config.BurnMonitor.Enabled == false {
		log.Debug("[BURN MONITOR] BURN monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[BURN MONITOR] Initializing burn monitor")

	m := newMonitor(wg)

	m.UpdateCurrentBlockNumber()

	m.InitStartBlockNumber(int64(app.Config.Ethereum.StartBlockNumber))

	log.Info("[BURN MONITOR] Initialized burn monitor")

	return m
}

func NewMonitorWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if app.Config.BurnMonitor.Enabled == false {
		log.Debug("[BURN MONITOR] BURN monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[BURN MONITOR] Initializing burn monitor")

	m := newMonitor(wg)

	lastBlockNumber, err := strconv.ParseInt(lastHealth.EthBlockNumber, 10, 64)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error parsing last block number from last health", err)
		lastBlockNumber = app.Config.Ethereum.StartBlockNumber
	}

	m.UpdateCurrentBlockNumber()

	m.InitStartBlockNumber(lastBlockNumber)

	log.Info("[BURN MONITOR] Initialized burn monitor")

	return m
}
