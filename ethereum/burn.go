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

type Service interface {
	Start()
	Stop()
}

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

type BurnMonitorService struct {
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	monitorInterval    time.Duration
	wpoktContract      WrappedPocketContract
}

func (b *BurnMonitorService) Stop() {
	log.Debug("Stopping burn monitor")
	b.stop <- true
}

func (b *BurnMonitorService) UpdateCurrentBlockNumber() {
	res, err := Client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("[BURN MONITOR] Current block number: ", res)
	b.currentBlockNumber = res
}

func (b *BurnMonitorService) HandleBurnEvent(event *WrappedPocketBurnAndBridge) bool {
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

const MAX_QUERY_BLOCKS uint64 = 100000

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
			success = success && b.SyncBlocks(i, endBlockNumber)
		}
	} else {
		log.Debug("[BURN MONITOR] Syncing burn txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
		success = success && b.SyncBlocks(b.startBlockNumber, b.currentBlockNumber)
	}
	return success
}

func (b *BurnMonitorService) Start() {
	log.Debug("[BURN MONITOR] Starting burn monitor")
	stop := false
	for !stop {
		log.Debug("[BURN MONITOR] Starting burn sync")

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
		log.Debug("[BURN MONITOR] Sleeping for ", b.monitorInterval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[BURN MONITOR] Stopped burn monitor")
		case <-time.After(b.monitorInterval):
		}
	}
}

func NewBurnMonitor() Service {
	log.Debug("[BURN MONITOR] Initializing burn monitor")
	log.Debug("[BURN MONITOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), Client.GetClient())
	if err != nil {
		log.Error("[BURN MONITOR] Error initializing Wrapped Pocket contract", err)
		panic(err)
	}
	log.Debug("[BURN MONITOR] Connected to wpokt contract")

	b := &BurnMonitorService{
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
		log.Debug("[BURN MONITOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Debug("[BURN MONITOR] Start block number: ", b.startBlockNumber)
	log.Debug("[BURN MONITOR] Initialized burn monitor")

	return b
}
