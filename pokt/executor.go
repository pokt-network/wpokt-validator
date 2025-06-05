package pokt

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/common"
	cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	"github.com/dan13ram/wpokt-validator/models"

	// "github.com/pokt-network/pocket-core/app/cmd/rpc"
	// "github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	BurnExecutorName = "BURN EXECUTOR"
)

type BurnExecutorRunner struct {
	client               cosmos.CosmosClient
	wpoktAddress         string
	vaultAddress         string
	multisigAddressBytes []byte
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

	// if doc.Status == models.StatusSigned {
	// 	log.Debug("[BURN EXECUTOR] Submitting invalid mint")

	// var err error

	// p := rpc.SendRawTxParams{
	// 	Addr:        x.vaultAddress,
	// 	RawHexBytes: doc.ReturnTx,
	// }
	//
	// res, err := x.client.SubmitRawTx(p)
	// if err != nil {
	// 	log.Error("[BURN EXECUTOR] Error submitting transaction: ", err)
	// 	return false
	// }
	//
	// if res == nil || strings.TrimSpace(res.TransactionHash) == "" {
	// 	log.Error("[BURN EXECUTOR] Invalid mint return tx hash not found")
	// 	return false
	// }
	//
	// filter = bson.M{
	// 	"_id":    doc.Id,
	// 	"status": models.StatusSigned,
	// }
	//
	// update = bson.M{
	// 	"$set": bson.M{
	// 		"status":         models.StatusSubmitted,
	// 		"return_tx_hash": res.TransactionHash,
	// 		"updated_at":     time.Now(),
	// 	},
	// }
	// } else if doc.Status == models.StatusSubmitted {
	// log.Debug("[BURN EXECUTOR] Checking invalid mint")
	// tx, err := x.client.GetTx(doc.ReturnTxHash)
	// if err != nil {
	// 	log.Error("[BURN EXECUTOR] Error fetching transaction: ", err)
	// 	return false
	// }
	// if tx == nil || tx.Tx == "" {
	// 	log.Error("[BURN EXECUTOR] Invalid mint return tx not found: ", doc.ReturnTxHash)
	// 	return false
	// }
	//
	// filter = bson.M{
	// 	"_id":    doc.Id,
	// 	"status": models.StatusSubmitted,
	// }
	//
	// if tx.TxResult.Code != 0 {
	// 	log.Error("[BURN EXECUTOR] Invalid mint return tx failed: ", tx.Hash)
	// 	update = bson.M{
	// 		"$set": bson.M{
	// 			"status":         models.StatusConfirmed,
	// 			"updated_at":     time.Now(),
	// 			"return_tx_hash": "",
	// 			"return_tx":      "",
	// 			"signers":        []string{},
	// 		},
	// 	}
	// } else {
	// 	log.Debug("[BURN EXECUTOR] Invalid mint return tx succeeded: ", tx.Hash)
	// 	update = bson.M{
	// 		"$set": bson.M{
	// 			"status":     models.StatusSuccess,
	// 			"updated_at": time.Now(),
	// 		},
	// 	}
	// }
	// }

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

		// p := rpc.SendRawTxParams{
		// 	Addr:        x.vaultAddress,
		// 	RawHexBytes: doc.ReturnTx,
		// }
		//
		// res, err := x.client.SubmitRawTx(p)
		// if err != nil {
		// 	log.Error("[BURN EXECUTOR] Error submitting transaction: ", err)
		// 	return false
		// }
		//
		// if res == nil || strings.TrimSpace(res.TransactionHash) == "" {
		// 	log.Error("[BURN EXECUTOR] Burn return tx hash not found")
		// 	return false
		// }
		//
		// filter = bson.M{
		// 	"_id":    doc.Id,
		// 	"status": models.StatusSigned,
		// }
		//
		// update = bson.M{
		// 	"$set": bson.M{
		// 		"status":         models.StatusSubmitted,
		// 		"return_tx_hash": res.TransactionHash,
		// 		"updated_at":     time.Now(),
		// 	},
		// }
	} else if doc.Status == models.StatusSubmitted {
		log.Debug("[BURN EXECUTOR] Checking burn")
		// tx, err := x.client.GetTx(doc.ReturnTxHash)
		// if err != nil {
		// 	log.Error("[BURN EXECUTOR] Error fetching transaction: ", err)
		// 	return false
		// }
		//
		// if tx == nil || tx.Tx == "" {
		// 	log.Error("[BURN EXECUTOR] Burn return tx not found: ", doc.ReturnTxHash)
		// 	return false
		// }
		//
		// filter = bson.M{
		// 	"_id":    doc.Id,
		// 	"status": models.StatusSubmitted,
		// }
		//
		// if tx.TxResult.Code != 0 {
		// 	log.Error("[BURN EXECUTOR] Burn return tx failed: ", tx.Hash)
		// 	update = bson.M{
		// 		"$set": bson.M{
		// 			"status":         models.StatusConfirmed,
		// 			"updated_at":     time.Now(),
		// 			"return_tx_hash": "",
		// 			"return_tx":      "",
		// 			"signers":        []string{},
		// 		},
		// 	}
		// } else {
		// 	log.Debug("[BURN EXECUTOR] Burn return tx succeeded: ", tx.Hash)
		// 	update = bson.M{
		// 		"$set": bson.M{
		// 			"status":     models.StatusSuccess,
		// 			"updated_at": time.Now(),
		// 		},
		// 	}
		// }
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

	config := app.Config.Pocket

	log.Debug("[BURN EXECTOR] Initializing")
	var pks []crypto.PubKey
	for _, pk := range config.MultisigPublicKeys {
		pKey, err := common.CosmosPublicKeyFromHex(pk)
		if err != nil {
			log.Fatalf("Error parsing public key: %s", err)
		}
		pks = append(pks, pKey)
	}

	multisigPk := multisig.NewLegacyAminoPubKey(int(config.MultisigThreshold), pks)
	multisigAddressBytes := multisigPk.Address().Bytes()
	multisigAddress, _ := common.Bech32FromBytes(config.Bech32Prefix, multisigAddressBytes)

	if !strings.EqualFold(multisigAddress, config.MultisigAddress) {
		log.Fatalf("Multisig address does not match config")
	}

	client, err := cosmosNewClient(config)
	if err != nil {
		log.Fatalf("Error creating pokt client: %s", err)
	}

	x := &BurnExecutorRunner{
		multisigAddressBytes: multisigAddressBytes,
		vaultAddress:         multisigAddress,
		wpoktAddress:         strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		client:               client,
	}

	log.Info("[BURN EXECUTOR] Initialized")

	return app.NewRunnerService(BurnExecutorName, x, wg, time.Duration(app.Config.BurnExecutor.IntervalMillis)*time.Millisecond)
}
