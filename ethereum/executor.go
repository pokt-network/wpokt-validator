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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type WPoktExecutorService struct {
	wg                 *sync.WaitGroup
	name               string
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	lastSyncTime       time.Time
	interval           time.Duration
	wpoktContract      WrappedPocketContract
	mintControllerAbi  *abi.ABI
	client             ethereum.EthereumClient
}

func (b *WPoktExecutorService) PoktHeight() string {
	return ""
}

func (b *WPoktExecutorService) EthBlockNumber() string {
	return strconv.FormatUint(b.startBlockNumber, 10)
}

func (b *WPoktExecutorService) LastSyncTime() time.Time {
	return b.lastSyncTime
}

func (b *WPoktExecutorService) Interval() time.Duration {
	return b.interval
}

func (b *WPoktExecutorService) Name() string {
	return b.name
}

func (b *WPoktExecutorService) Stop() {
	log.Debug("[WPOKT EXECUTOR] Stopping wpokt executor")
	b.stop <- true
}

func (b *WPoktExecutorService) UpdateCurrentBlockNumber() {
	res, err := b.client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = res
	log.Debug("[WPOKT EXECUTOR] Current block number: ", b.currentBlockNumber)
}

func (b *WPoktExecutorService) HandleMintEvent(event *autogen.WrappedPocketMinted) bool {
	log.Debug("[WPOKT EXECUTOR] Handling mint event: ", event.Raw.TxHash, " ", event.Raw.Index)

	recipient := event.Recipient.Hex()
	amount := event.Amount.String()
	nonce := event.Nonce.String()

	filter := bson.M{
		"recipient_address": recipient,
		"amount":            amount,
		"nonce":             nonce,
	}

	update := bson.M{
		"$set": bson.M{
			"status":       models.StatusSuccess,
			"mint_tx_hash": event.Raw.TxHash.String(),
		},
	}

	err := app.DB.UpdateOne(models.CollectionMints, filter, update)

	if err != nil {
		log.Error("[WPOKT EXECUTOR] Error while updating mint: ", err)
		return false
	}

	log.Debug("[WPOKT EXECUTOR] Mint event handled successfully")

	return true
}

func (b *WPoktExecutorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := b.wpoktContract.FilterMinted(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []common.Address{}, []*big.Int{}, []*big.Int{})

	if err != nil {
		log.Errorln("[WPOKT EXECUTOR] Error while syncing mint events: ", err)
		return false
	}

	var success bool = true
	for filter.Next() {
		event := filter.Event()
		success = success && b.HandleMintEvent(event)
	}
	return success
}

func (b *WPoktExecutorService) SyncTxs() bool {
	var success bool = true
	if (b.currentBlockNumber - b.startBlockNumber) > MAX_QUERY_BLOCKS {
		log.Debug("[WPOKT EXECUTOR] Syncing mint txs in chunks")
		for i := b.startBlockNumber; i < b.currentBlockNumber; i += MAX_QUERY_BLOCKS {
			endBlockNumber := i + MAX_QUERY_BLOCKS
			if endBlockNumber > b.currentBlockNumber {
				endBlockNumber = b.currentBlockNumber
			}
			log.Debug("[WPOKT EXECUTOR] Syncing mint txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && b.SyncBlocks(i, endBlockNumber)
		}
	} else {
		log.Debug("[WPOKT EXECUTOR] Syncing mint txs from blockNumber: ", b.startBlockNumber, " to blockNumber: ", b.currentBlockNumber)
		success = success && b.SyncBlocks(b.startBlockNumber, b.currentBlockNumber)
	}
	return success
}

func (b *WPoktExecutorService) Start() {
	log.Debug("[WPOKT EXECUTOR] Starting wpokt executor")
	stop := false
	for !stop {
		log.Debug("[WPOKT EXECUTOR] Starting mint sync")
		b.lastSyncTime = time.Now()

		b.UpdateCurrentBlockNumber()

		if (b.currentBlockNumber - b.startBlockNumber) > 0 {
			success := b.SyncTxs()
			if success {
				b.startBlockNumber = b.currentBlockNumber
			}
		} else {
			log.Debug("[WPOKT EXECUTOR] No new blocks to sync")
		}

		log.Debug("[WPOKT EXECUTOR] Finished mint sync")
		log.Debug("[WPOKT EXECUTOR] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[WPOKT EXECUTOR] Stopped wpokt executor")
		case <-time.After(b.interval):
		}
	}
	b.wg.Done()
}

func NewExecutor(wg *sync.WaitGroup) models.Service {
	if app.Config.WPOKTExecutor.Enabled == false {
		log.Debug("[WPOKT EXECUTOR] WPOKT executor disabled")
		return models.NewEmptyService(wg, "empty-wpokt-executor")
	}

	log.Debug("[WPOKT EXECUTOR] Initializing wpokt executor")
	client, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[WPOKT EXECUTOR] Error initializing ethereum client", err)
	}
	log.Debug("[WPOKT EXECUTOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), client.GetClient())
	if err != nil {
		log.Fatal("[WPOKT EXECUTOR] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[WPOKT EXECUTOR] Connected to wpokt contract")

	mintControllerAbi, err := autogen.MintControllerMetaData.GetAbi()
	if err != nil {
		log.Fatal("[WPOKT EXECUTOR] Error parsing MintController ABI", err)
	}

	log.Debug("[WPOKT EXECUTOR] Mint controller abi parsed")

	b := &WPoktExecutorService{
		wg:                 wg,
		name:               "wpokt-executor",
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		interval:           time.Duration(app.Config.WPOKTExecutor.IntervalSecs) * time.Second,
		wpoktContract:      &WrappedPocketContractImpl{contract},
		mintControllerAbi:  mintControllerAbi,
		client:             client,
	}

	b.UpdateCurrentBlockNumber()
	if app.Config.Ethereum.StartBlockNumber > 0 {
		b.startBlockNumber = uint64(app.Config.Ethereum.StartBlockNumber)
	} else {
		log.Debug("[WPOKT EXECUTOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Debug("[WPOKT EXECUTOR] Start block number: ", b.startBlockNumber)
	log.Debug("[WPOKT EXECUTOR] Initialized wpokt executor")

	return b
}
