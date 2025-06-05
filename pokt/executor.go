package pokt

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/common"
	cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	"github.com/dan13ram/wpokt-validator/models"

	"github.com/dan13ram/wpokt-validator/cosmos/util"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	BurnExecutorName = "BURN EXECUTOR"
)

type BurnExecutorRunner struct {
	signer       *app.PocketSigner
	client       cosmos.CosmosClient
	wpoktAddress string
	vaultAddress string
}

func (x *BurnExecutorRunner) Run() {
	x.SyncTxs()
}

func (x *BurnExecutorRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{}
}

func (x *BurnExecutorRunner) HandleInvalidMint(doc *models.InvalidMint) bool {

	if doc == nil || (doc.Status != models.StatusSigned && doc.Status != models.StatusSubmitted) {
		log.Error("[BURN EXECUTOR] Invalid mint is nil or has invalid status")
		return false
	}

	log.Debug("[BURN EXECUTOR] Handling invalid mint: ", doc.TransactionHash)

	var filter bson.M
	var update bson.M

	if doc.Status == models.StatusSigned {
		log.Debug("[BURN EXECUTOR] Submitting invalid mint")

		txBuilder, txCfg, err := util.WrapTxBuilder(app.Config.Pocket.Bech32Prefix, doc.ReturnTransactionBody)
		if err != nil {
			log.WithError(err).Errorf("Error wrapping tx builder")
			return false
		}

		txJSON, err := txCfg.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			log.WithError(err).Errorf("Error encoding tx")
			return false
		}

		txBytes, err := txCfg.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			log.WithError(err).Errorf("Error encoding tx")
			return false
		}

		txHash, err := x.client.BroadcastTx(txBytes)
		if err != nil {
			log.WithError(err).Errorf("Error broadcasting tx")
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		update = bson.M{
			"$set": bson.M{
				"status":                  models.StatusSubmitted,
				"return_transaction_body": string(txJSON),
				"return_transaction_hash": common.Ensure0xPrefix(txHash),
				"updated_at":              time.Now(),
			},
		}
	} else if doc.Status == models.StatusSubmitted {
		log.Debug("[BURN EXECUTOR] Checking invalid mint")
		tx, err := x.client.GetTx(doc.ReturnTransactionHash)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error fetching transaction: ", err)
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		if tx.Code != 0 {
			log.Error("[BURN EXECUTOR] Invalid mint return tx failed: ", tx.TxHash)
			update = bson.M{
				"$set": bson.M{
					"status":                  models.StatusConfirmed,
					"updated_at":              time.Now(),
					"return_transaction_hash": "",
					"return_transaction_body": "",
					"signatures":              []models.Signature{},
					"sequence":                nil,
				},
			}
		} else {
			log.Debug("[BURN EXECUTOR] Invalid mint return tx succeeded: ", tx.TxHash)
			update = bson.M{
				"$set": bson.M{
					"status":     models.StatusSuccess,
					"updated_at": time.Now(),
				},
			}
		}
	}

	if _, err := app.DB.UpdateOne(models.CollectionInvalidMints, filter, update); err != nil {
		log.Error("[BURN EXECUTOR] Error updating invalid mint: ", err)
		return false
	}

	log.Info("[BURN EXECUTOR] Handled invalid mint")
	return true
}

func (x *BurnExecutorRunner) HandleBurn(doc *models.Burn) bool {

	if doc == nil || (doc.Status != models.StatusSigned && doc.Status != models.StatusSubmitted) {
		log.Error("[BURN EXECUTOR] Burn is nil or has invalid status")
		return false
	}

	log.Debug("[BURN EXECUTOR] Handling burn: ", doc.TransactionHash)

	var filter bson.M
	var update bson.M

	if doc.Status == models.StatusSigned {
		log.Debug("[BURN EXECUTOR] Submitting burn")

		txBuilder, txCfg, err := util.WrapTxBuilder(app.Config.Pocket.Bech32Prefix, doc.ReturnTransactionBody)
		if err != nil {
			log.WithError(err).Errorf("Error wrapping tx builder")
			return false
		}

		txJSON, err := txCfg.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			log.WithError(err).Errorf("Error encoding tx")
			return false
		}

		txBytes, err := txCfg.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			log.WithError(err).Errorf("Error encoding tx")
			return false
		}

		txHash, err := x.client.BroadcastTx(txBytes)
		if err != nil {
			log.WithError(err).Errorf("Error broadcasting tx")
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		update = bson.M{
			"$set": bson.M{
				"status":                  models.StatusSubmitted,
				"return_transaction_body": string(txJSON),
				"return_transaction_hash": common.Ensure0xPrefix(txHash),
				"updated_at":              time.Now(),
			},
		}
	} else if doc.Status == models.StatusSubmitted {
		log.Debug("[BURN EXECUTOR] Checking burn")
		tx, err := x.client.GetTx(doc.ReturnTransactionHash)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error fetching transaction: ", err)
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		if tx.Code != 0 {
			log.Error("[BURN EXECUTOR] Burn return tx failed: ", tx.TxHash)
			update = bson.M{
				"$set": bson.M{
					"status":                  models.StatusConfirmed,
					"updated_at":              time.Now(),
					"return_transaction_hash": "",
					"return_transaction_body": "",
					"signatures":              []models.Signature{},
					"sequence":                nil,
				},
			}
		} else {
			log.Debug("[BURN EXECUTOR] Burn return tx succeeded: ", tx.TxHash)
			update = bson.M{
				"$set": bson.M{
					"status":     models.StatusSuccess,
					"updated_at": time.Now(),
				},
			}
		}
	}

	if _, err := app.DB.UpdateOne(models.CollectionBurns, filter, update); err != nil {
		log.Error("[BURN EXECUTOR] Error updating burn: ", err)
		return false
	}

	log.Info("[BURN EXECUTOR] Handled burn")
	return true
}

func (x *BurnExecutorRunner) SyncInvalidMints() bool {
	log.Debug("[BURN EXECUTOR] Syncing invalid mints")

	filter := bson.M{
		"status": bson.M{
			"$in": []string{
				string(models.StatusSigned),
				string(models.StatusSubmitted),
			},
		},
		"vault_address": x.vaultAddress,
	}
	invalidMints := []models.InvalidMint{}

	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error fetching invalid mints: ", err)
		return false
	}

	log.Info("[BURN EXECUTOR] Found invalid mints: ", len(invalidMints))

	var success bool = true
	for i := range invalidMints {
		doc := invalidMints[i]

		resourceId := fmt.Sprintf("%s/%s", models.CollectionInvalidMints, doc.Id.Hex())
		lockId, err := app.DB.XLock(resourceId)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error locking invalid mint: ", err)
			success = false
			continue
		}
		log.Debug("[BURN EXECUTOR] Locked invalid mint: ", doc.TransactionHash)

		success = x.HandleInvalidMint(&doc) && success

		if err := app.DB.Unlock(lockId); err != nil {
			log.Error("[BURN EXECUTOR] Error unlocking invalid mint: ", err)
			success = false
		} else {
			log.Debug("[BURN EXECUTOR] Unlocked invalid mint: ", doc.TransactionHash)
		}

	}

	log.Debug("[BURN EXECUTOR] Synced invalid mints")
	return success
}

func (x *BurnExecutorRunner) SyncBurns() bool {
	log.Debug("[BURN EXECUTOR] Syncing burns")

	filter := bson.M{
		"status": bson.M{
			"$in": []string{
				string(models.StatusSigned),
				string(models.StatusSubmitted),
			},
		},
		"wpokt_address": x.wpoktAddress,
	}
	burns := []models.Burn{}

	err := app.DB.FindMany(models.CollectionBurns, filter, &burns)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error fetching burns: ", err)
		return false
	}

	log.Info("[BURN EXECUTOR] Found burns: ", len(burns))

	var success bool = true

	for i := range burns {
		doc := burns[i]

		resourceId := fmt.Sprintf("%s/%s", models.CollectionBurns, doc.Id.Hex())
		lockId, err := app.DB.XLock(resourceId)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error locking burn: ", err)
			success = false
			continue
		}
		log.Debugln("[BURN EXECUTOR] Locked burn:", doc.TransactionHash, doc.LogIndex)

		success = x.HandleBurn(&doc) && success

		if err := app.DB.Unlock(lockId); err != nil {
			log.Error("[BURN EXECUTOR] Error unlocking burn: ", err)
			success = false
		} else {
			log.Debugln("[BURN EXECUTOR] Unlocked burn:", doc.TransactionHash, doc.LogIndex)
		}

	}

	log.Debug("[BURN EXECUTOR] Synced burns")
	return success
}

func (x *BurnExecutorRunner) SyncTxs() bool {
	log.Debug("[BURN EXECUTOR] Syncing")

	success := x.SyncInvalidMints()
	success = x.SyncBurns() && success

	log.Info("[BURN EXECUTOR] Synced txs")
	return success
}

func NewBurnExecutor(wg *sync.WaitGroup, health models.ServiceHealth) app.Service {
	if !app.Config.BurnExecutor.Enabled {
		log.Debug("[BURN EXECUTOR] Disabled")
		return app.NewEmptyService(wg)
	}

	log.Debug("[BURN EXECTOR] Initializing")
	signer, err := app.GetPocketSignerAndMultisig()
	if err != nil {
		log.Fatal("[BURN SIGNER] Error getting signer and multisig: ", err)
	}

	client, err := cosmosNewClient(app.Config.Pocket)
	if err != nil {
		log.Fatalf("Error creating pokt client: %s", err)
	}

	x := &BurnExecutorRunner{
		signer:       signer,
		vaultAddress: signer.MultisigAddress,
		wpoktAddress: strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		client:       client,
	}

	log.Info("[BURN EXECUTOR] Initialized")

	return app.NewRunnerService(BurnExecutorName, x, wg, time.Duration(app.Config.BurnExecutor.IntervalMillis)*time.Millisecond)
}
