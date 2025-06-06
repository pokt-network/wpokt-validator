package eth

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MintExecutorName string = "MINT EXECUTOR"
)

type MintExecutorRunner struct {
	startBlockNumber   int64
	currentBlockNumber int64
	wpoktContract      eth.WrappedPocketContract
	mintControllerAbi  *abi.ABI
	client             eth.EthereumClient
	vaultAddress       string
	wpoktAddress       string
}

func (x *MintExecutorRunner) Run() {
	x.UpdateCurrentBlockNumber()
	x.SyncTxs()
}

func (x *MintExecutorRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{
		EthBlockNumber: strconv.FormatInt(x.startBlockNumber, 10),
	}
}

func (x *MintExecutorRunner) UpdateCurrentBlockNumber() {
	res, err := x.client.GetBlockNumber()
	if err != nil {
		log.Error("[MINT EXECUTOR] Error while getting current block number: ", err)
		return
	}

	x.currentBlockNumber = int64(res)
	log.Info("[MINT EXECUTOR] Current block number: ", x.currentBlockNumber)
}

func (x *MintExecutorRunner) HandleMintEvent(event *autogen.WrappedPocketMinted) bool {
	if event == nil {
		log.Error("[MINT EXECUTOR] Invalid mint event")
		return false
	}

	log.Debug("[MINT EXECUTOR] Handling mint event: ", event.Raw.TxHash, " ", event.Raw.Index)

	filter := bson.M{
		"wpokt_address":     x.wpoktAddress,
		"vault_address":     x.vaultAddress,
		"recipient_address": strings.ToLower(event.Recipient.Hex()),
		"amount":            event.Amount.String(),
		"nonce":             event.Nonce.String(),
		"status": bson.M{
			"$in": []string{models.StatusConfirmed, models.StatusSigned},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":       models.StatusSuccess,
			"mint_tx_hash": strings.ToLower(event.Raw.TxHash.String()),
			"updated_at":   time.Now(),
		},
	}

	_, err := app.DB.UpdateOne(models.CollectionMints, filter, update)

	if err != nil {
		log.Error("[MINT EXECUTOR] Error while updating mint: ", err)
		return false
	}

	log.Info("[MINT EXECUTOR] Mint event handled successfully")

	return true
}

func (x *MintExecutorRunner) SyncBlocks(startBlockNumber uint64, endBlockNumber uint64) bool {
	filter, err := x.wpoktContract.FilterMinted(&bind.FilterOpts{
		Start:   startBlockNumber,
		End:     &endBlockNumber,
		Context: context.Background(),
	}, []common.Address{}, []*big.Int{}, []*big.Int{})

	if filter != nil {
		defer filter.Close()
	}

	if err != nil {
		log.Errorln("[MINT EXECUTOR] Error while syncing mint events: ", err)
		return false
	}

	var success bool = true
	for filter.Next() {
		if err = filter.Error(); err != nil {
			success = false
			break
		}

		event := filter.Event()

		if event == nil {
			success = false
			continue
		}

		if event.Raw.Removed {
			continue
		}

		resourceId := fmt.Sprintf("%s/%s", models.CollectionMints, strings.ToLower(event.Recipient.Hex()))
		lockId, err := app.DB.XLock(resourceId)
		if err != nil {
			log.Error("[MINT EXECUTOR] Error locking mint: ", err)
			success = false
			continue
		}
		log.Debug("[MINT EXECUTOR] Locked mint: ", event.Raw.TxHash)

		success = x.HandleMintEvent(event) && success

		if err = app.DB.Unlock(lockId); err != nil {
			log.Error("[MINT EXECUTOR] Error unlocking mint: ", err)
			success = false
		} else {
			log.Debug("[MINT EXECUTOR] Unlocked mint: ", event.Raw.TxHash)
		}
	}

	if err = filter.Error(); err != nil {
		log.Errorln("[MINT EXECUTOR] Error while syncing mint events: ", err)
		return false
	}

	return success
}

func (x *MintExecutorRunner) SyncTxs() bool {

	if x.currentBlockNumber <= x.startBlockNumber {
		log.Info("[MINT EXECUTOR] No new blocks to sync")
		return true
	}

	var success bool = true

	if (x.currentBlockNumber - x.startBlockNumber) > eth.MAX_QUERY_BLOCKS {
		log.Debug("[MINT EXECUTOR] Syncing mint txs in chunks")

		for i := x.startBlockNumber; i < x.currentBlockNumber; i += eth.MAX_QUERY_BLOCKS {
			endBlockNumber := i + eth.MAX_QUERY_BLOCKS
			if endBlockNumber > x.currentBlockNumber {
				endBlockNumber = x.currentBlockNumber
			}

			log.Info("[MINT EXECUTOR] Syncing mint txs from blockNumber: ", i, " to blockNumber: ", endBlockNumber)
			success = success && x.SyncBlocks(uint64(i), uint64(endBlockNumber))
		}

	} else {
		log.Info("[MINT EXECUTOR] Syncing mint txs from blockNumber: ", x.startBlockNumber, " to blockNumber: ", x.currentBlockNumber)
		success = success && x.SyncBlocks(uint64(x.startBlockNumber), uint64(x.currentBlockNumber))
	}

	if success {
		x.startBlockNumber = x.currentBlockNumber
	}

	return success
}

func (x *MintExecutorRunner) InitStartBlockNumber(lastHealth models.ServiceHealth) {
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

	log.Info("[MINT EXECUTOR] Start block number: ", x.startBlockNumber)
}

func NewMintExecutor(wg *sync.WaitGroup, lastHealth models.ServiceHealth) app.Service {
	if !app.Config.MintExecutor.Enabled || app.Config.Pocket.MintDisabled {
		log.Debug("[MINT EXECUTOR] Disabled")
		return app.NewEmptyService(wg)
	}
	log.Debug("[MINT EXECUTOR] Initializing mint executor")

	client, err := eth.NewClient()
	if err != nil {
		log.Fatal("[MINT EXECUTOR] Error initializing ethereum client", err)
	}

	log.Debug("[MINT EXECUTOR] Connecting to mint contract at: ", app.Config.Ethereum.WrappedPocketAddress)

	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WrappedPocketAddress), client.GetClient())
	if err != nil {
		log.Fatal("[MINT EXECUTOR] Error initializing Wrapped Pocket contract", err)
	}

	log.Debug("[MINT EXECUTOR] Connected to mint contract")

	mintControllerAbi, err := autogen.MintControllerMetaData.GetAbi()
	if err != nil {
		log.Fatal("[MINT EXECUTOR] Error parsing MintController ABI", err)
	}

	log.Debug("[MINT EXECUTOR] Mint controller abi parsed")

	x := &MintExecutorRunner{
		startBlockNumber:   0,
		currentBlockNumber: 0,
		wpoktContract:      eth.NewWrappedPocketContract(contract),
		mintControllerAbi:  mintControllerAbi,
		client:             client,
		wpoktAddress:       strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		vaultAddress:       strings.ToLower(app.Config.Pocket.MultisigAddress),
	}

	x.UpdateCurrentBlockNumber()

	x.InitStartBlockNumber(lastHealth)

	log.Info("[MINT EXECUTOR] Initialized mint executor")

	return app.NewRunnerService(MintExecutorName, x, wg, time.Duration(app.Config.MintExecutor.IntervalMillis)*time.Millisecond)
}
