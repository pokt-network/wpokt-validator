package cosmos

import (
	"context"
	"cosmossdk.io/math"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dan13ram/wpokt-validator/app"
	cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/ethereum/go-ethereum/common"
)

const (
	MintMonitorName = "MINT MONITOR"
)

type MintMonitorRunner struct {
	client                 cosmos.CosmosClient
	ethClient              eth.EthereumClient
	mintControllerContract eth.MintControllerContract
	wpoktAddress           string
	vaultAddress           string
	startHeight            int64
	currentHeight          int64
	minimumAmount          math.Int
	maximumAmount          math.Int
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
	if tx == nil || result == nil {
		log.Debug("[MINT MONITOR] Invalid tx response")
		return false
	}

	doc := util.CreateInvalidMint(tx, result, x.vaultAddress)

	if app.Config.Pocket.MintDisabled {
		// ensure that existing mints are not counted as invalid mints after mint is disabled
		err := app.DB.FindOne(models.CollectionMints, bson.M{"transaction_hash": doc.TransactionHash}, &models.Mint{})
		if err == nil {
			log.Warn("[MINT MONITOR] Ignoring invalid mint since it exists as a valid mint")
			return true
		}
	}

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
	if tx == nil || result == nil {
		log.Debug("[MINT MONITOR] Invalid tx response")
		return false
	}

	if app.Config.Pocket.MintDisabled {
		log.Error("[MINT MONITOR] HandleValidMint called when mint is disabled")
		return true
	}

	doc := util.CreateMint(tx, result, x.wpoktAddress, x.vaultAddress)

	// ensure that existing invalid mints are not counted as valid mints after mint is enabled
	if err := app.DB.FindOne(models.CollectionInvalidMints, bson.M{"transaction_hash": doc.TransactionHash}, &models.InvalidMint{}); err == nil {
		log.Warn("[MINT MONITOR] Ignoring valid mint since it exists as an invalid mint")
		return true
	}

	log.Debug("[MINT MONITOR] Storing mint tx")
	if _, err := app.DB.InsertOne(models.CollectionMints, doc); err != nil {
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

var utilValidateTxToCosmosMultisig = util.ValidateTxToCosmosMultisig

func (x *MintMonitorRunner) SyncTxs() bool {

	if x.currentHeight <= x.startHeight {
		log.Info("[MINT MONITOR] No new blocks to sync")
		return true
	}

	txResponses, err := x.client.GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight))
	if err != nil {
		log.Error("[MINT MONITOR] Error getting txs: ", err)
		return false
	}
	log.Info("[MINT MONITOR] Found ", len(txResponses), " txs to sync")
	var success = true
	for _, txResponse := range txResponses {

		result := utilValidateTxToCosmosMultisig(txResponse, app.Config.Pocket, uint64(x.currentHeight), x.minimumAmount, x.maximumAmount)

		if result.TxStatus == models.TransactionStatusFailed {
			log.Info("[MINT MONITOR] Found failed mint tx: ", result.TxHash)
			success = x.HandleFailedMint(txResponse, result) && success
			continue
		}

		if result.NeedsRefund || app.Config.Pocket.MintDisabled {
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

func (x *MintMonitorRunner) UpdateMaxMintLimit() {
	log.Debug("[MINT MONITOR] Fetching mint controller max mint limit")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx, Pending: false}
	mintLimit, err := x.mintControllerContract.MaxMintLimit(opts)

	if err != nil {
		log.Error("[MINT MONITOR] Error fetching mint controller max mint limit: ", err)
		return
	}
	log.Debug("[MINT MONITOR] Fetched mint controller max mint limit")
	x.maximumAmount = math.NewIntFromBigInt(mintLimit)
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
		log.Fatalf("[MINT MONITOR] Error creating pokt client: %s", err)
	}

	ethClient, err := eth.NewClient()
	if err != nil {
		log.Fatal("[MINT MONITOR] Error initializing ethereum client: ", err)
	}

	log.Debug("[MINT MONITOR] Connecting to mint controller contract at: ", app.Config.Ethereum.MintControllerAddress)
	mintControllerContract, err := autogen.NewMintController(common.HexToAddress(app.Config.Ethereum.MintControllerAddress), ethClient.GetClient())
	if err != nil {
		log.Fatal("[MINT MONITOR] Error initializing Mint Controller contract", err)
	}
	log.Debug("[MINT MONITOR] Connected to mint controller contract")

	x := &MintMonitorRunner{
		vaultAddress:           signer.MultisigAddress,
		wpoktAddress:           strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		startHeight:            0,
		currentHeight:          0,
		client:                 client,
		ethClient:              ethClient,
		minimumAmount:          math.NewIntFromUint64(uint64(app.Config.Pocket.TxFee)),
		mintControllerContract: eth.NewMintControllerContract(mintControllerContract),
	}

	x.UpdateCurrentHeight()

	x.InitStartHeight(lastHealth)

	x.UpdateMaxMintLimit()

	if x.maximumAmount.LT(x.minimumAmount) {
		log.Fatal("[MINT MONITOR] Invalid max mint limit")
	}

	log.Info("[MINT MONITOR] Initialized")

	return app.NewRunnerService(MintMonitorName, x, wg, time.Duration(app.Config.MintMonitor.IntervalMillis)*time.Millisecond)
}

var cosmosNewClient = cosmos.NewClient
