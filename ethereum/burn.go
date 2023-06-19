package ethereum

// listen to events from the blockchain

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// burn monitor interface

type BurnMonitor interface {
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

type WPOKTBurnMonitor struct {
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	monitorInterval    time.Duration
	wpoktContract      WrappedPocketContract
}

func (b *WPOKTBurnMonitor) Stop() {
	log.Debug("Stopping burn monitor")
	b.stop <- true
}

func (b *WPOKTBurnMonitor) UpdateCurrentBlockNumber() {
	res, err := Client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Updated current ethereum blockNumber: ", res)
	b.currentBlockNumber = res
}

func (b *WPOKTBurnMonitor) HandleBurnEvent(event *WrappedPocketBurnAndBridge) bool {
	doc := models.Burn{
		BlockNumber:      event.Raw.BlockNumber,
		TransactionHash:  event.Raw.TxHash.String(),
		LogIndex:         uint64(event.Raw.Index),
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
	log.Debug("Storing burn event in db: ", event.Raw.TxHash, " ", event.Raw.Index)

	err := app.DB.InsertOne(models.CollectionBurns, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("Found duplicate burn event in db: ", event.Raw.TxHash, " ", event.Raw.Index)
			return true
		}
		log.Error("Error storing burn event in db: ", err)
		return false
	}

	log.Debug("Stored burn event in db: ", event.Raw.TxHash, " ", event.Raw.Index)
	return true
}

const MAX_QUERY_BLOCKS uint64 = 100000

func (b *WPOKTBurnMonitor) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := b.wpoktContract.FilterBurnAndBridge(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []*big.Int{}, []common.Address{}, []common.Address{})

	if err != nil {
		log.Error(err)
		return false
	}

	var success bool = true
	for filter.Next() {
		event := filter.Event()
		log.Debug("Found burn event: ", event.Raw.TxHash, " ", event.Raw.Index)
		success = success && b.HandleBurnEvent(event)
	}
	return success
}

func (b *WPOKTBurnMonitor) SyncTxs() bool {
	var success bool = true
	if (b.currentBlockNumber - b.startBlockNumber) > MAX_QUERY_BLOCKS {
		log.Debug("Syncing burn txs in chunks")
		for i := b.startBlockNumber; i < b.currentBlockNumber; i += MAX_QUERY_BLOCKS {
			endBlockNumber := i + MAX_QUERY_BLOCKS
			if endBlockNumber > b.currentBlockNumber {
				endBlockNumber = b.currentBlockNumber
			}
			log.Debug("Syncing burn txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && b.SyncBlocks(i, endBlockNumber)
		}
	} else {
		log.Debug("Syncing burn txs all at once")
		success = success && b.SyncBlocks(b.startBlockNumber, b.currentBlockNumber)
	}
	return success
}

func (b *WPOKTBurnMonitor) Start() {
	log.Debug("Starting burn monitor")
	stop := false
	for !stop {
		log.Debug("Starting burn Sync")

		b.UpdateCurrentBlockNumber()

		if (b.currentBlockNumber - b.startBlockNumber) > 0 {
			log.Debug("Syncing burn txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
			success := b.SyncTxs()
			if success {
				b.startBlockNumber = b.currentBlockNumber
			}
		} else {
			log.Debug("Already Synced up to blockNumber: ", b.currentBlockNumber)
		}

		log.Debug("Finished burn Sync")
		log.Debug("Sleeping for ", b.monitorInterval)
		log.Debug("Next burn Sync at: ", time.Now().Add(b.monitorInterval))

		select {
		case <-b.stop:
			stop = true
			log.Debug("Stopped burn monitor")
		case <-time.After(b.monitorInterval):
		}
	}
}

func NewBurnMonitor() BurnMonitor {
	log.Debug("Connecting to Wrapped Pocket contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), Client.GetClient())
	if err != nil {
		log.Error("Error connecting to Wrapped Pocket contract: ", err)
		panic(err)
	}
	log.Debug("Connected to Wrapped Pocket contract at: ", app.Config.Ethereum.WPOKTContractAddress)

	b := &WPOKTBurnMonitor{
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		monitorInterval:    time.Duration(app.Config.Ethereum.MonitorIntervalSecs) * time.Second,
		wpoktContract:      &WrappedPocketContractImpl{contract},
	}

	if app.Config.Ethereum.StartBlockNumber < 0 {
		b.UpdateCurrentBlockNumber()
		b.startBlockNumber = b.currentBlockNumber
	} else {
		b.startBlockNumber = uint64(app.Config.Ethereum.StartBlockNumber)
	}

	return b
}
