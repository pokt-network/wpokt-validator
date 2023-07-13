package ethereum

// listen to events from the blockchain

import (
	"context"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type BurnAndBridgeIterator interface {
	Next() bool
	Event() *WrappedPocketBurnAndBridge
}

type BurnAndBridgeIteratorImpl struct {
	*WrappedPocketBurnAndBridgeIterator
}

func (i *BurnAndBridgeIteratorImpl) Event() *WrappedPocketBurnAndBridge {
	return i.WrappedPocketBurnAndBridgeIterator.Event
}

func (i *BurnAndBridgeIteratorImpl) Next() bool {
	return i.WrappedPocketBurnAndBridgeIterator.Next()
}

type WrappedPocketContract interface {
	FilterBurnAndBridge(opts *bind.FilterOpts, _amount []*big.Int, _from []common.Address, _poktAddress []common.Address) (BurnAndBridgeIterator, error)
}

type WrappedPocketContractImpl struct {
	*WrappedPocket
}

func (c *WrappedPocketContractImpl) FilterBurnAndBridge(opts *bind.FilterOpts, _amount []*big.Int, _from []common.Address, _poktAddress []common.Address) (BurnAndBridgeIterator, error) {
	iterator, err := c.WrappedPocket.FilterBurnAndBridge(opts, _amount, _from, _poktAddress)
	return &BurnAndBridgeIteratorImpl{iterator}, err
}

type WPoktMonitorService struct {
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	monitorInterval    time.Duration
	wpoktContract      WrappedPocketContract
}

func (b *WPoktMonitorService) Stop() {
	log.Debug("[WPOKT MONITOR] Stopping wpokt monitor")
	b.stop <- true
}

func (b *WPoktMonitorService) UpdateCurrentBlockNumber() {
	res, err := Client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = res
	log.Debug("[WPOKT MONITOR] Current block number: ", b.currentBlockNumber)
}

func (b *WPoktMonitorService) HandleBurnEvent(event *WrappedPocketBurnAndBridge) bool {
	doc := models.Burn{
		BlockNumber:      strconv.FormatInt(int64(event.Raw.BlockNumber), 10),
		TransactionHash:  event.Raw.TxHash.String(),
		LogIndex:         strconv.FormatInt(int64(event.Raw.Index), 10),
		SenderAddress:    event.From.String(),
		SenderChainId:    strconv.FormatInt(int64(app.Config.Ethereum.ChainId), 10),
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

const MAX_QUERY_BLOCKS uint64 = 100000

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
		log.Debug("[WPOKT MONITOR] Sleeping for ", b.monitorInterval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[WPOKT MONITOR] Stopped wpokt monitor")
		case <-time.After(b.monitorInterval):
		}
	}
}

func NewMonitor() models.Service {
	if app.Config.Ethereum.Enabled == false {
		log.Debug("[WPOKT MONITOR] Ethereum is disabled, not starting wpokt monitor")
		return models.NewEmptyService()
	}

	log.Debug("[WPOKT MONITOR] Initializing wpokt monitor")
	log.Debug("[WPOKT MONITOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), Client.GetClient())
	if err != nil {
		log.Fatal("[WPOKT MONITOR] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[WPOKT MONITOR] Connected to wpokt contract")

	b := &WPoktMonitorService{
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		monitorInterval:    time.Duration(app.Config.Ethereum.MonitorIntervalSecs) * time.Second,
		wpoktContract:      &WrappedPocketContractImpl{contract},
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
