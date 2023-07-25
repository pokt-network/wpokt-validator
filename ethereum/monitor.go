package ethereum

import (
	"context"
	"math/big"
	"strconv"
	"strings"
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
	BurnMonitorName = "burn-monitor"
)

type BurnMonitorService struct {
	wg                 *sync.WaitGroup
	name               string
	stop               chan bool
	startBlockNumber   int64
	currentBlockNumber int64
	lastSyncTime       time.Time
	interval           time.Duration
	wpoktContract      WrappedPocketContract
	client             ethereum.EthereumClient
}

func (b *BurnMonitorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           b.Name(),
		LastSyncTime:   b.LastSyncTime(),
		NextSyncTime:   b.LastSyncTime().Add(b.Interval()),
		PoktHeight:     b.PoktHeight(),
		EthBlockNumber: b.EthBlockNumber(),
		Healthy:        true,
	}
}

func (b *BurnMonitorService) PoktHeight() string {
	return ""
}

func (b *BurnMonitorService) EthBlockNumber() string {
	return strconv.FormatInt(b.startBlockNumber, 10)
}

func (b *BurnMonitorService) Name() string {
	return b.name
}

func (b *BurnMonitorService) Stop() {
	log.Debug("[BURN MONITOR] Stopping wpokt monitor")
	b.stop <- true
}

func (b *BurnMonitorService) UpdateCurrentBlockNumber() {
	res, err := b.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = int64(res)
	log.Debug("[BURN MONITOR] Current block number: ", b.currentBlockNumber)
}

func (b *BurnMonitorService) HandleBurnEvent(event *autogen.WrappedPocketBurnAndBridge) bool {
	doc := models.Burn{
		BlockNumber:      strconv.FormatInt(int64(event.Raw.BlockNumber), 10),
		Confirmations:    "0",
		TransactionHash:  event.Raw.TxHash.String(),
		LogIndex:         strconv.FormatInt(int64(event.Raw.Index), 10),
		WPOKTAddress:     event.Raw.Address.String(),
		SenderAddress:    event.From.String(),
		SenderChainId:    app.Config.Ethereum.ChainId,
		RecipientAddress: strings.Split(event.PoktAddress.String(), "0x")[1],
		RecipientChainId: app.Config.Pocket.ChainId,
		Amount:           event.Amount.String(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           models.StatusPending,
		Signers:          []string{},
	}

	// each event is a combination of transaction hash and log index
	log.Debug("[BURN MONITOR] Handling burn event: ", event.Raw.TxHash, " ", event.Raw.Index)

	err := app.DB.InsertOne(models.CollectionBurns, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[BURN MONITOR] Found duplicate burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
			return true
		}
		log.Error("[BURN MONITOR] Error while storing burn event in db: ", err)
		return false
	}

	log.Debug("[BURN MONITOR] Stored burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
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
		event := filter.Event()
		success = success && b.HandleBurnEvent(event)
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
			log.Debug("[BURN MONITOR] Syncing burn txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && b.SyncBlocks(uint64(i), uint64(endBlockNumber))
		}
	} else {
		log.Debug("[BURN MONITOR] Syncing burn txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
		success = success && b.SyncBlocks(uint64(b.startBlockNumber), uint64(b.currentBlockNumber))
	}
	return success
}

func (b *BurnMonitorService) Start() {
	log.Debug("[BURN MONITOR] Starting wpokt monitor")
	stop := false
	for !stop {
		log.Debug("[BURN MONITOR] Starting burn sync")
		b.lastSyncTime = time.Now()

		b.UpdateCurrentBlockNumber()

		if (b.currentBlockNumber - b.startBlockNumber) > 0 {
			success := b.SyncTxs()
			if success {
				b.startBlockNumber = b.currentBlockNumber
			}
		} else {
			log.Debug("[BURN MONITOR] No new blocks to sync")
		}

		log.Debug("[BURN MONITOR] Finished burn sync")
		log.Debug("[BURN MONITOR] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[BURN MONITOR] Stopped wpokt monitor")
		case <-time.After(b.interval):
		}
	}
	b.wg.Done()
}

func (b *BurnMonitorService) LastSyncTime() time.Time {
	return b.lastSyncTime
}

func (b *BurnMonitorService) Interval() time.Duration {
	return b.interval
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
		wpoktContract:      &WrappedPocketContractImpl{contract},
		client:             client,
	}

	return b
}

func (b *BurnMonitorService) InitStartBlockNumber(startBlockNumber int64) {
	b.UpdateCurrentBlockNumber()
	if app.Config.Ethereum.StartBlockNumber > 0 {
		b.startBlockNumber = int64(app.Config.Ethereum.StartBlockNumber)
	} else {
		log.Debug("[BURN MONITOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Debug("[BURN EXECUTOR] Start block number: ", b.startBlockNumber)
}

func NewMonitor(wg *sync.WaitGroup) models.Service {
	if app.Config.BurnMonitor.Enabled == false {
		log.Debug("[BURN MONITOR] BURN monitor disabled")
		return models.NewEmptyService(wg, "empty-wpokt-monitor")
	}

	log.Debug("[BURN MONITOR] Initializing wpokt monitor")

	m := newMonitor(wg)

	m.InitStartBlockNumber(int64(app.Config.Ethereum.StartBlockNumber))
	log.Debug("[BURN MONITOR] Initialized wpokt monitor")

	return m
}

func NewMonitorWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if app.Config.BurnMonitor.Enabled == false {
		log.Debug("[BURN MONITOR] BURN monitor disabled")
		return models.NewEmptyService(wg, "empty-wpokt-monitor")
	}

	log.Debug("[BURN MONITOR] Initializing wpokt monitor")

	m := newMonitor(wg)

	lastBlockNumber, err := strconv.ParseInt(lastHealth.EthBlockNumber, 10, 64)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error parsing last block number from last health", err)
		lastBlockNumber = app.Config.Ethereum.StartBlockNumber
	}

	m.InitStartBlockNumber(lastBlockNumber)
	log.Debug("[BURN MONITOR] Initialized wpokt monitor")

	return m
}
