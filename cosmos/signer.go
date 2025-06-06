package cosmos

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/common"
	cosmosClient "github.com/dan13ram/wpokt-validator/cosmos/client"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	BurnSignerName = "BURN SIGNER"
)

type BurnSignerRunner struct {
	signer         *app.PocketSigner
	ethClient      eth.EthereumClient
	poktClient     cosmosClient.CosmosClient
	poktHeight     int64
	ethBlockNumber int64
	vaultAddress   string
	wpoktAddress   string
	wpoktContract  eth.WrappedPocketContract
	minimumAmount  *big.Int
}

func (x *BurnSignerRunner) Run() {
	x.UpdateBlocks()
	x.SyncTxs()
}
func (x *BurnSignerRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{
		PoktHeight:     strconv.FormatInt(x.poktHeight, 10),
		EthBlockNumber: strconv.FormatInt(x.ethBlockNumber, 10),
	}
}

func (x *BurnSignerRunner) UpdateBlocks() {
	log.Debug("[BURN SIGNER] Updating blocks")

	poktHeight, err := x.poktClient.GetLatestBlockHeight()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching pokt block height: ", err)
		return
	}
	x.poktHeight = poktHeight

	ethBlockNumber, err := x.ethClient.GetBlockNumber()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching eth block number: ", err)
		return
	}
	x.ethBlockNumber = int64(ethBlockNumber)

	log.Info("[BURN SIGNER] Updated blocks")
}

func (x *BurnSignerRunner) ValidateInvalidMint(doc *models.InvalidMint) (bool, error) {
	log.Debug("[BURN SIGNER] Validating invalid mint: ", doc.TransactionHash)

	txResponse, err := x.poktClient.GetTx(doc.TransactionHash)
	if err != nil {
		return false, errors.New("Error fetching transaction: " + err.Error())
	}

	if txResponse == nil {
		return false, errors.New("Transaction not found")
	}
	result, err := util.ValidateTxToCosmosMultisig(txResponse, app.Config.Pocket, uint64(x.poktHeight))
	if err != nil {
		log.WithError(err).Errorf("Error validating tx")
		return false, nil
	}
	if result.TxStatus == models.TransactionStatusFailed {
		log.Debug("[BURN SIGNER] Invalid Mint Transaction is failed")
		return false, nil
	}

	// NOTE: If mint is disabled, we refund all txs
	if !app.Config.Pocket.MintDisabled && !result.NeedsRefund {
		log.Debug("[BURN SIGNER] Invalid Mint Transaction does not need refund")
		return false, nil
	}

	log.Debug("[BURN SIGNER] Validated invalid mint")
	return true, nil
}

func (x *BurnSignerRunner) FindMaxSequence() (uint64, error) {
	lockID, err := LockReadSequences()
	if err != nil {
		return 0, fmt.Errorf("could not lock sequences: %w", err)
	}
	//nolint:errcheck
	defer app.DB.Unlock(lockID)

	maxSequence, err := FindMaxSequence()
	if err != nil {
		return 0, err
	}
	account, err := x.poktClient.GetAccount(x.signer.MultisigAddress)
	if err != nil {
		return 0, err
	}
	if maxSequence == nil {
		return account.Sequence, nil
	}
	nextSequence := *maxSequence + 1
	if nextSequence > account.Sequence {
		return nextSequence, nil
	}

	return account.Sequence, nil
}

func (x *BurnSignerRunner) Sign(
	sequence *uint64,
	signatures []models.Signature,
	transactionBody string,
	toAddress []byte,
	amount sdk.Coin,
	memo string,
) (bson.M, error) {

	if sequence == nil {
		gotSequence, err := x.FindMaxSequence()
		if err != nil {
			return nil, fmt.Errorf("error getting sequence: %w", err)
		}
		sequence = &gotSequence
	}

	txBody, finalSignatures, err := SignTx(
		x.signer.Signer,
		app.Config.Pocket,
		x.poktClient,
		*sequence,
		signatures,
		transactionBody,
		toAddress,
		amount,
		memo,
	)

	if err != nil {
		return nil, err
	}

	update := bson.M{
		"status":                  models.StatusConfirmed,
		"return_transaction_body": string(txBody),
		"signatures":              finalSignatures,
		"sequence":                sequence,
		"updated_at":              time.Now(),
	}

	if len(finalSignatures) >= int(app.Config.Pocket.MultisigThreshold) {
		update["status"] = models.StatusSigned
	}

	return update, nil
}

func (x *BurnSignerRunner) HandleInvalidMint(doc *models.InvalidMint) bool {
	if doc == nil {
		log.Error("[BURN SIGNER] Invalid mint is nil")
		return false
	}
	log.Debug("[BURN SIGNER] Handling invalid mint: ", doc.TransactionHash)

	doc, err := util.UpdateStatusAndConfirmationsForInvalidMint(doc, x.poktHeight)
	if err != nil {
		log.Error("[BURN SIGNER] Error getting invalid mint status: ", err)
		return false
	}

	var update bson.M

	valid, err := x.ValidateInvalidMint(doc)
	if err != nil {
		log.Error("[BURN SIGNER] Error validating invalid mint: ", err)
		return false
	}

	if !valid {
		log.Error("[BURN SIGNER] Invalid mint failed validation")
		update = bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}
		if doc.Confirmations == "0" {
			log.Debug("[BURN SIGNER] Invalid mint has no confirmations, skipping")
			return false
		}
	} else {

		if doc.Status == models.StatusConfirmed {
			log.Debug("[BURN SIGNER] Signing invalid mint")

			amount, _ := math.NewIntFromString(doc.Amount)
			amountCoin := sdk.NewCoin(app.Config.Pocket.CoinDenom, amount)

			toAddress, err := common.AddressBytesFromBech32(app.Config.Pocket.Bech32Prefix, doc.SenderAddress)
			if err != nil {
				log.Error("[BURN SIGNER] Error parsing to address: ", err)
				return false
			}

			set, err := x.Sign(doc.Sequence, doc.Signatures, doc.ReturnTransactionBody, toAddress, amountCoin, "InvalidMint: "+doc.TransactionHash)

			if err != nil {
				log.Error("[BURN SIGNER] Error signing invalid mint: ", err)
				return false
			}

			update = bson.M{
				"$set": set,
			}

		} else {
			log.Debug("[BURN SIGNER] Not signing invalid mint")
			update = bson.M{
				"$set": bson.M{
					"status":        doc.Status,
					"confirmations": doc.Confirmations,
					"updated_at":    time.Now(),
				},
			}
		}
	}

	filter := bson.M{
		"_id":    doc.Id,
		"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
	}
	_, err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating invalid mint: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Handled invalid mint: ", doc.TransactionHash)
	return true
}

func (x *BurnSignerRunner) ValidateBurn(doc *models.Burn) (bool, error) {
	log.Debug("[BURN SIGNER] Validating burn: ", doc.TransactionHash)

	txReceipt, err := x.ethClient.GetTransactionReceipt(doc.TransactionHash)

	if err != nil {
		return false, errors.New("Error fetching transaction receipt: " + err.Error())
	}

	logIndex, err := strconv.Atoi(doc.LogIndex)
	if err != nil {
		log.Debug("[BURN SIGNER] Error converting log index to int: ", err)
		return false, nil
	}

	var burnLog *types.Log

	for _, log := range txReceipt.Logs {
		if log.Index == uint(logIndex) {
			burnLog = log
			break
		}
	}

	if burnLog == nil {
		log.Debug("[BURN SIGNER] Burn log not found")
		return false, nil
	}

	burnEvent, err := x.wpoktContract.ParseBurnAndBridge(*burnLog)
	if err != nil {
		log.Error("[BURN SIGNER] Error parsing burn event: ", err)
		return false, nil
	}

	amount, ok := new(big.Int).SetString(doc.Amount, 10)
	if !ok || amount.Cmp(x.minimumAmount) != 1 {
		log.Debug("[BURN SIGNER] Burn amount too low")
		return false, nil
	}

	if burnEvent.Amount.Cmp(amount) != 0 {
		log.Error("[BURN SIGNER] Invalid burn amount")
		return false, nil
	}
	if !strings.EqualFold(burnEvent.From.Hex(), doc.SenderAddress) {
		log.Error("[BURN SIGNER] Invalid burn sender")
		return false, nil
	}
	recipientBytes, _ := common.AddressBytesFromBech32(app.Config.Pocket.Bech32Prefix, doc.RecipientAddress)
	if !bytes.Equal(burnEvent.PoktAddress.Bytes(), recipientBytes) {
		log.Error("[BURN SIGNER] Invalid burn recipient")
		return false, nil
	}

	log.Debug("[BURN SIGNER] Validated burn")
	return true, nil
}

func (x *BurnSignerRunner) HandleBurn(doc *models.Burn) bool {
	if doc == nil {
		log.Error("[BURN SIGNER] Burn is nil")
		return false
	}
	log.Debug("[BURN SIGNER] Handling burn: ", doc.TransactionHash)

	doc, err := util.UpdateStatusAndConfirmationsForBurn(doc, x.ethBlockNumber)
	if err != nil {
		log.Error("[BURN SIGNER] Error getting burn status: ", err)
		return false
	}

	var update bson.M

	valid, err := x.ValidateBurn(doc)
	if err != nil {
		log.Error("[BURN SIGNER] Error validating burn: ", err)
		return false
	}
	if !valid {
		log.Error("[BURN SIGNER] Burn failed validation")
		update = bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}
		if doc.Confirmations == "0" {
			log.Debug("[BURN SIGNER] Burn has no confirmations, skipping")
			return false
		}
	} else {

		if doc.Status == models.StatusConfirmed {
			log.Debug("[BURN SIGNER] Signing burn")
			amount, _ := math.NewIntFromString(doc.Amount)
			amountCoin := sdk.NewCoin(app.Config.Pocket.CoinDenom, amount)

			toAddress, err := common.AddressBytesFromBech32(app.Config.Pocket.Bech32Prefix, doc.RecipientAddress)
			if err != nil {
				log.Error("[BURN SIGNER] Error parsing to address: ", err)
				return false
			}

			set, err := x.Sign(doc.Sequence, doc.Signatures, doc.ReturnTransactionBody, toAddress, amountCoin, "Burn: "+doc.TransactionHash)

			if err != nil {
				log.Error("[BURN SIGNER] Error signing invalid mint: ", err)
				return false
			}

			update = bson.M{
				"$set": set,
			}
		} else {
			log.Debug("[BURN SIGNER] Not signing burn")
			update = bson.M{
				"$set": bson.M{
					"status":        doc.Status,
					"confirmations": doc.Confirmations,
					"updated_at":    time.Now(),
				},
			}
		}
	}

	filter := bson.M{
		"_id":    doc.Id,
		"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
	}
	_, err = app.DB.UpdateOne(models.CollectionBurns, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating burn: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Handled burn: ", doc.TransactionHash)

	return true
}

func (x *BurnSignerRunner) SyncInvalidMints() bool {
	log.Debug("[BURN SIGNER] Syncing invalid mints")

	addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
	filter := bson.M{
		"$and": []bson.M{
			{
				"vault_address": x.vaultAddress,
			},
			{"$or": []bson.M{
				{"status": models.StatusPending},
				{"status": models.StatusConfirmed},
			}},
			{"$nor": []bson.M{
				{"signatures": bson.M{
					"$elemMatch": bson.M{"signer": addressHex},
				}},
			}},
		},
	}

	invalidMints := []models.InvalidMint{}
	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching invalid mints: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Found invalid mints: ", len(invalidMints))

	var success bool = true
	for i := range invalidMints {
		doc := invalidMints[i]

		resourceId := fmt.Sprintf("%s/%s", models.CollectionInvalidMints, doc.Id.Hex())
		lockId, err := app.DB.XLock(resourceId)
		if err != nil {
			log.Error("[BURN SIGNER] Error locking invalid mint: ", err)
			success = false
			continue
		}
		log.Debug("[BURN SIGNER] Locked invalid mint: ", doc.TransactionHash)

		success = x.HandleInvalidMint(&doc) && success

		if err = app.DB.Unlock(lockId); err != nil {
			log.Error("[BURN SIGNER] Error unlocking invalid mint: ", err)
			success = false
		} else {
			log.Debug("[BURN SIGNER] Unlocked invalid mint: ", doc.TransactionHash)
		}

	}

	log.Info("[BURN SIGNER] Synced invalid mints")
	return success
}

func (x *BurnSignerRunner) SyncBurns() bool {
	log.Debug("[BURN SIGNER] Syncing burns")

	addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
	filter := bson.M{
		"$and": []bson.M{
			{
				"wpokt_address": x.wpoktAddress,
			},
			{"$or": []bson.M{
				{"status": models.StatusPending},
				{"status": models.StatusConfirmed},
			}},
			{"$nor": []bson.M{
				{"signatures": bson.M{
					"$elemMatch": bson.M{"signer": addressHex},
				}},
			}},
		},
	}

	burns := []models.Burn{}
	err := app.DB.FindMany(models.CollectionBurns, filter, &burns)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching burns: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Found burns: ", len(burns))

	var success bool = true

	for i := range burns {
		doc := burns[i]

		resourceId := fmt.Sprintf("%s/%s", models.CollectionBurns, doc.Id.Hex())
		lockId, err := app.DB.XLock(resourceId)
		if err != nil {
			log.Error("[BURN SIGNER] Error locking burn: ", err)
			success = false
			continue
		}
		log.Debug("[BURN SIGNER] Locked burn: ", doc.TransactionHash)

		success = x.HandleBurn(&doc) && success

		if err = app.DB.Unlock(lockId); err != nil {
			log.Error("[BURN SIGNER] Error unlocking burn: ", err)
			success = false
		} else {
			log.Debug("[BURN SIGNER] Unlocked burn: ", doc.TransactionHash)
		}

	}

	log.Info("[BURN SIGNER] Synced burns")
	return success
}

func (x *BurnSignerRunner) SyncTxs() bool {
	log.Debug("[BURN SIGNER] Syncing")

	success := x.SyncInvalidMints()
	success = x.SyncBurns() && success

	log.Info("[BURN SIGNER] Synced txs")
	return success
}

func NewBurnSigner(wg *sync.WaitGroup, health models.ServiceHealth) app.Service {
	if !app.Config.BurnSigner.Enabled {
		log.Debug("[BURN SIGNER] Disabled")
		return app.NewEmptyService(wg)
	}

	log.Debug("[BURN SIGNER] Initializing")

	signer, err := app.GetPocketSignerAndMultisig()
	if err != nil {
		log.Fatal("[BURN SIGNER] Error getting signer and multisig: ", err)
	}

	ethClient, err := eth.NewClient()
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing ethereum client: ", err)
	}

	log.Debug("[BURN SIGNER] Connecting to wpokt contract at: ", app.Config.Ethereum.WrappedPocketAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WrappedPocketAddress), ethClient.GetClient())
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[BURN SIGNER] Connected to wpokt contract")

	poktClient, err := cosmosNewClient(app.Config.Pocket)
	if err != nil {
		log.Fatalf("Error creating pokt client: %s", err)
	}

	x := &BurnSignerRunner{
		signer:        signer,
		ethClient:     ethClient,
		poktClient:    poktClient,
		vaultAddress:  signer.MultisigAddress,
		wpoktAddress:  strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		wpoktContract: eth.NewWrappedPocketContract(contract),
		minimumAmount: big.NewInt(app.Config.Pocket.TxFee),
	}

	x.UpdateBlocks()

	log.Info("[BURN SIGNER] Initialized")

	return app.NewRunnerService(BurnSignerName, x, wg, time.Duration(app.Config.BurnSigner.IntervalMillis)*time.Millisecond)
}
