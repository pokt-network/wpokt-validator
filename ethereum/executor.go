package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"
	"reflect"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type WPoktExecutorService struct {
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	interval           time.Duration
	wpoktContract      WrappedPocketContract
	mintControllerAbi  *abi.ABI
}

func (b *WPoktExecutorService) Stop() {
	log.Debug("[WPOKT EXECUTOR] Stopping wpokt executor")
	b.stop <- true
}

func (b *WPoktExecutorService) UpdateCurrentBlockNumber() {
	res, err := Client.GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	b.currentBlockNumber = res
	log.Debug("[WPOKT EXECUTOR] Current block number: ", b.currentBlockNumber)
}

func decodeMintData(data interface{}) models.MintData {
	value := reflect.ValueOf(data)
	return models.MintData{
		Recipient: value.Field(0).Interface().(common.Address).String(),
		Amount:    value.Field(1).Interface().(*big.Int).String(),
		Nonce:     value.Field(2).Interface().(*big.Int).String(),
	}
}

func (b *WPoktExecutorService) HandleMintEvent(event *autogen.WrappedPocketTransfer) bool {
	log.Debug("[WPOKT EXECUTOR] Handling mint event: ", event.Raw.TxHash, " ", event.Raw.Index)

	if (event.From.String() != ZERO_ADDRESS) || (event.To.String() == ZERO_ADDRESS) {
		log.Debug("[WPOKT EXECUTOR] Skipping mint event: invalid from or to address")
		return true
	}

	recipient := event.To.String()
	amount := event.Value.String()

	tx, _, err := Client.GetTransactionByHash(event.Raw.TxHash.String())
	if err != nil {
		log.Error("[WPOKT EXECUTOR] Error while getting transaction by hash: ", err)
		return false
	}

	data := tx.Data()

	methodId := data[:4]

	mintMethod := b.mintControllerAbi.Methods["mintWrappedPocket"]

	if hex.EncodeToString(methodId) != hex.EncodeToString(mintMethod.ID) {
		log.Debug("[WPOKT EXECUTOR] Skipping mint event: invalid method id")
		return true
	}

	decoded, err := abi.Arguments(mintMethod.Inputs).UnpackValues(data[4:])
	if err != nil {
		log.Error("[WPOKT EXECUTOR] Error while decoding mint function input: ", err)
		return false
	}

	mintData := decodeMintData(decoded[0])

	if (mintData.Recipient != recipient) || (mintData.Amount != amount) {
		log.Debug("[WPOKT EXECUTOR] Skipping mint event: invalid recipient or amount")
		return true
	}

	filter := bson.M{
		"recipient_address": recipient,
		"amount":            amount,
		"data":              mintData,
		"status": bson.M{
			"$in": []string{models.StatusPending, models.StatusSigned},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":  models.StatusSuccess,
			"tx_hash": event.Raw.TxHash.String(),
		},
	}

	err = app.DB.UpdateOne(models.CollectionMints, filter, update)

	if err != nil {
		log.Error("[WPOKT EXECUTOR] Error while updating mint: ", err)
		return false
	}

	log.Debug("[WPOKT EXECUTOR] Mint event handled successfully")

	return true
}

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"

func (b *WPoktExecutorService) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := b.wpoktContract.FilterMinted(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []common.Address{}, []common.Address{})

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
}

func NewExecutor() models.Service {
	if app.Config.WPOKTExecutor.Enabled == false {
		log.Debug("[WPOKT EXECUTOR] WPOKT executor disabled")
		return models.NewEmptyService()
	}

	log.Debug("[WPOKT EXECUTOR] Initializing wpokt executor")
	log.Debug("[WPOKT EXECUTOR] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), Client.GetClient())
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
		stop:               make(chan bool),
		startBlockNumber:   0,
		currentBlockNumber: 0,
		interval:           time.Duration(app.Config.WPOKTExecutor.IntervalSecs) * time.Second,
		wpoktContract:      &WrappedPocketContractImpl{contract},
		mintControllerAbi:  mintControllerAbi,
	}

	b.UpdateCurrentBlockNumber()
	if app.Config.WPOKTExecutor.StartBlockNumber > 0 {
		b.startBlockNumber = uint64(app.Config.WPOKTExecutor.StartBlockNumber)
	} else {
		log.Debug("[WPOKT EXECUTOR] Found invalid start block number, updating to current block number")
		b.startBlockNumber = b.currentBlockNumber
	}

	log.Debug("[WPOKT EXECUTOR] Start block number: ", b.startBlockNumber)
	log.Debug("[WPOKT EXECUTOR] Initialized wpokt executor")

	return b
}
