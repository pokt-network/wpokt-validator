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

type WPOKTBurnMonitor struct {
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	monitorInterval    time.Duration
	wpoktContract      *WrappedPocket
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
	log.Info("Updated current ethereum blockNumber: ", res)
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

	col := app.DB.GetCollection(models.CollectionBurns)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.Pocket.MonitorIntervalSecs))
	defer cancel()

	_, err := col.InsertOne(ctx, doc)
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

func (b *WPOKTBurnMonitor) SyncTxs() bool {
	filter, err := b.wpoktContract.WrappedPocketFilterer.FilterBurnAndBridge(&bind.FilterOpts{
		Start:   b.startBlockNumber,
		End:     &b.currentBlockNumber,
		Context: context.Background(),
	}, []*big.Int{}, []common.Address{}, []common.Address{})

	if err != nil {
		log.Error(err)
		return false
	}

	var success bool = true
	for filter.Next() {
		log.Debug("Found burn event: ", filter.Event.Raw.TxHash, " ", filter.Event.Raw.Index)
		success = success && b.HandleBurnEvent(filter.Event)
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
		monitorInterval:    time.Duration(app.Config.Pocket.MonitorIntervalSecs) * time.Second,
		wpoktContract:      contract,
	}

	if app.Config.Ethereum.StartBlockNumber < 0 {
		b.UpdateCurrentBlockNumber()
		b.startBlockNumber = b.currentBlockNumber
	} else {
		b.startBlockNumber = uint64(app.Config.Ethereum.StartBlockNumber)
	}

	return b
}
