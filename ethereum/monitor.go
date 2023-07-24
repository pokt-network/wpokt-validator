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

type WPoktMonitorService struct {
	wg                 *sync.WaitGroup
	name               string
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	lastSyncTime       time.Time
	interval           time.Duration
	wpoktContract      WrappedPocketContract
	client             ethereum.EthereumClient
}

func (b *WPoktMonitorService) PoktHeight() string {
	return ""
}

func (b *WPoktMonitorService) EthBlockNumber() string {
	return strconv.FormatUint(b.startBlockNumber, 10)
}

func (b *WPoktMonitorService) Name() string {
	return b.name
}

func (b *WPoktMonitorService) Stop() {
	log.Debug("[WPOKT MONITOR] Stopping wpokt monitor")
	b.stop <- true
}

func (b *WPoktMonitorService) UpdateCurrentBlockNumber() {
	res, err := b.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = res
	log.Debug("[WPOKT MONITOR] Current block number: ", b.currentBlockNumber)
}

func (b *WPoktMonitorService) HandleBurnEvent(event *autogen.WrappedPocketBurnAndBridge) bool {
	doc := models.Burn{
		BlockNumber:      strconv.FormatInt(int64(event.Raw.BlockNumber), 10),
		TransactionHash:  event.Raw.TxHash.String(),
		LogIndex:         strconv.FormatInt(int64(event.Raw.Index), 10),
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
	log.Debug("[WPOKT MONITOR] Handling burn event: ", event.Raw.TxHash, " ", event.Raw.Index)

	err := app.DB.InsertOne(models.CollectionBurns, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[WPOKT MONITOR] Found duplicate burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
			return true
		}
		log.Error("[WPOKT MONITOR] Error while storing burn event in db: ", err)
		return false
	}

	log.Debug("[WPOKT MONITOR] Stored burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
	return true
}

func (b *WPoktMonitorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := b.wpoktContract.FilterBurnAndBridge(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []*big.Int{}, []common.Address{}, []common.Address{})

	if err != nil {
		log.Errorln("[WPOKT MONITOR] Error while syncing burn events: ", err)
		return false
	}

	var success bool = true
	for filter.Next() {
		event := filter.Event()
		success = success && b.HandleBurnEvent(event)
	}
	return success
}

func (b *WPoktMonitorService) SyncTxs() bool {
	var success bool = true
	if (b.currentBlockNumber - b.startBlockNumber) > MAX_QUERY_BLOCKS {
		log.Debug("[WPOKT MONITOR] Syncing burn txs in chunks")
		for i := b.startBlockNumber; i < b.currentBlockNumber; i += MAX_QUERY_BLOCKS {
			endBlockNumber := i + MAX_QUERY_BLOCKS
			if endBlockNumber > b.currentBlockNumber {
				endBlockNumber = b.currentBlockNumber
			}
			log.Debug("[WPOKT MONITOR] Syncing burn txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && b.SyncBlocks(i, endBlockNumber)
		}
	} else {
		log.Debug("[WPOKT MONITOR] Syncing burn txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
		success = success && b.SyncBlocks(b.startBlockNumber, b.currentBlockNumber)
	}
	return success
}

func (b *WPoktMonitorService) Start() {
	log.Debug("[WPOKT MONITOR] Starting wpokt monitor")
	stop := false
	for !stop {
		log.Debug("[WPOKT MONITOR] Starting burn sync")
		b.lastSyncTime = time.Now()

		b.UpdateCurrentBlockNumber()

		if (b.currentBlockNumber - b.startBlockNumber) > 0 {
			success := b.SyncTxs()
			if success {
				b.startBlockNumber = b.currentBlockNumber
			}
		} else {
			log.Debug("[WPOKT MONITOR] No new blocks to sync")
		}

		log.Debug("[WPOKT MONITOR] Finished burn sync")
		log.Debug("[WPOKT MONITOR] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[WPOKT MONITOR] Stopped wpokt monitor")
		case <-time.After(b.interval):
		}
	}
	b.wg.Done()
}

func (b *WPoktMonitorService) LastSyncTime() time.Time {
	return b.lastSyncTime
}

func (b *WPoktMonitorService) Interval() time.Duration {
	return b.interval
}

func NewMonitor(wg *sync.WaitGroup) models.Service {
	if app.Config.WPOKTMonitor.Enabled == false {
		log.Debug("[WPOKT MONITOR] WPOKT monitor disabled")
		return models.NewEmptyService(wg, "empty-wpokt-monitor")
	}

	log.Debug("[WPOKT MONITOR] Initializing wpokt monitor")
	client, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[WPOKT MONITOR] Error initializing ethereum client: ", err)
	}
	log.Debug("[WPOKT MONITOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), client.GetClient())
	if err != nil {
		log.Fatal("[WPOKT MONITOR] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[WPOKT MONITOR] Connected to wpokt contract")

	b := &WPoktMonitorService{
		wg:                 wg,
		name:               "wpokt-monitor",
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		interval:           time.Duration(app.Config.WPOKTMonitor.IntervalSecs) * time.Second,
		wpoktContract:      &WrappedPocketContractImpl{contract},
		client:             client,
	}

	b.UpdateCurrentBlockNumber()
	if app.Config.Ethereum.StartBlockNumber > 0 {
		b.startBlockNumber = uint64(app.Config.Ethereum.StartBlockNumber)
	} else {
		log.Debug("[WPOKT MONITOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Debug("[WPOKT MONITOR] Start block number: ", b.startBlockNumber)
	log.Debug("[WPOKT MONITOR] Initialized wpokt monitor")

	return b
}
