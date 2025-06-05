package pokt

import (
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dan13ram/wpokt-validator/app"
	cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	"github.com/dan13ram/wpokt-validator/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MintMonitorName = "MINT MONITOR"
)

type MintMonitorRunner struct {
	client        cosmos.CosmosClient
	wpoktAddress  string
	vaultAddress  string
	startHeight   int64
	currentHeight int64
	minimumAmount *big.Int
}

func (x *MintMonitorRunner) Run() {
	x.UpdateCurrentHeight()
	x.SyncTxs()
}

func (x *MintMonitorRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{
		PoktHeight: strconv.FormatInt(x.startHeight, 10),
	}
}

func (x *MintMonitorRunner) UpdateCurrentHeight() {
	res, err := x.client.GetLatestBlockHeight()
	if err != nil {
		log.Error("[MINT MONITOR] Error getting current height: ", err)
		return
	}
	x.currentHeight = res
	log.Info("[MINT MONITOR] Current height: ", x.currentHeight)
}

func (x *MintMonitorRunner) HandleFailedMint(tx *sdk.TxResponse, result *util.ValidateTxResult) bool {
	if tx == nil || result == nil {
		log.Debug("[MINT MONITOR] Invalid tx response")
		return false
	}

	doc := util.CreateFailedMint(tx, result, x.vaultAddress)

	log.Debug("[MINT MONITOR] Storing failed mint tx")
	_, err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Info("[MINT MONITOR] Found duplicate failed mint tx")
			return true
		}
		log.Error("[MINT MONITOR] Error storing failed mint tx: ", err)
		return false
	}

	log.Info("[MINT MONITOR] Stored failed mint tx")
	return true
}

func (x *MintMonitorRunner) HandleInvalidMint(tx *sdk.TxResponse, result *util.ValidateTxResult) bool {
	if tx == nil {
		log.Debug("[MINT MONITOR] Invalid tx response")
		return false
	}

	doc := util.CreateInvalidMint(tx, result, x.vaultAddress)

	log.Debug("[MINT MONITOR] Storing invalid mint tx")
	_, err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Info("[MINT MONITOR] Found duplicate invalid mint tx")
			return true
		}
		log.Error("[MINT MONITOR] Error storing invalid mint tx: ", err)
		return false
	}

	log.Info("[MINT MONITOR] Stored invalid mint tx")
	return true
}

func (x *MintMonitorRunner) HandleValidMint(tx *sdk.TxResponse, result *util.ValidateTxResult) bool {
	if tx == nil {
		log.Debug("[MINT MONITOR] Invalid tx response")
		return false
	}

	doc := util.CreateMint(tx, result, x.wpoktAddress, x.vaultAddress)

	log.Debug("[MINT MONITOR] Storing mint tx")
	_, err := app.DB.InsertOne(models.CollectionMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Info("[MINT MONITOR] Found duplicate mint tx")
			return true
		}
		log.Error("[MINT MONITOR] Error storing mint tx: ", err)
		return false
	}

	log.Info("[MINT MONITOR] Stored mint tx")
	return true
}

func (x *MintMonitorRunner) SyncTxs() bool {

	if x.currentHeight <= x.startHeight {
		log.Info("[MINT MONITOR] No new blocks to sync")
		return true
	}

	// txs, err := x.client.GetAccountTxsByHeight(x.vaultAddress, x.startHeight)
	txResponses, err := x.client.GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight))
	if err != nil {
		log.Error("[MINT MONITOR] Error getting txs: ", err)
		return false
	}
	log.Info("[MINT MONITOR] Found ", len(txResponses), " txs to sync")
	var success bool = true
	for _, txResponse := range txResponses {

		result, err := util.ValidateTxToCosmosMultisig(txResponse, app.Config.Pocket, uint64(x.currentHeight))
		if err != nil {
			log.WithError(err).Errorf("Error validating tx")
			success = false
			continue
		}

		if result.TxStatus == models.TransactionStatusFailed {
			log.Info("[MINT MONITOR] Found failed mint tx: ", result.TxHash)
			success = x.HandleFailedMint(txResponse, result) && success
			continue
		}

		if result.TxStatus == models.TransactionStatusInvalid {
			log.Info("[MINT MONITOR] Found invalid mint tx: ", result.TxHash)
			success = x.HandleInvalidMint(txResponse, result) && success
			continue
		}

		log.Info("[MINT MONITOR] Found valid mint tx: ", result.TxHash)
		success = x.HandleValidMint(txResponse, result) && success

	}

	if success {
		x.startHeight = x.currentHeight
	}

	return success
}

func (x *MintMonitorRunner) InitStartHeight(lastHealth models.ServiceHealth) {
	startHeight := (app.Config.Pocket.StartHeight)

	if (lastHealth.PoktHeight) != "" {
		if lastHeight, err := strconv.ParseInt(lastHealth.PoktHeight, 10, 64); err == nil {
			startHeight = lastHeight
		}
	}
	if startHeight > 0 {
		x.startHeight = startHeight
	} else {
		log.Info("[MINT MONITOR] Found invalid start height, using current height")
		x.startHeight = x.currentHeight
	}
	log.Info("[MINT MONITOR] Start height: ", x.startHeight)
}

func NewMintMonitor(wg *sync.WaitGroup, lastHealth models.ServiceHealth) app.Service {
	if !app.Config.MintMonitor.Enabled {
		log.Debug("[MINT MONITOR] Disabled")
		return app.NewEmptyService(wg)
	}

	log.Debug("[MINT MONITOR] Initializing")

	signer, err := app.GetPocketSignerAndMultisig()
	if err != nil {
		log.Fatal("[MINT MONITOR] Error getting signer and multisig: ", err)
	}

	client, err := cosmosNewClient(app.Config.Pocket)
	if err != nil {
		log.Fatalf("Error creating pokt client: %s", err)
	}

	x := &MintMonitorRunner{
		vaultAddress:  signer.MultisigAddress,
		wpoktAddress:  strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		startHeight:   0,
		currentHeight: 0,
		client:        client,
		minimumAmount: big.NewInt(app.Config.Pocket.TxFee),
	}

	x.UpdateCurrentHeight()

	x.InitStartHeight(lastHealth)

	log.Info("[MINT MONITOR] Initialized")

	return app.NewRunnerService(MintMonitorName, x, wg, time.Duration(app.Config.MintMonitor.IntervalMillis)*time.Millisecond)
}

var cosmosNewClient = cosmos.NewClient
