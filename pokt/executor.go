package pokt

import (
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	BurnExecutorName = "burn executor"
)

type BurnExecutorService struct {
	wg              *sync.WaitGroup
	client          pokt.PocketClient
	wpoktAddress    string
	stop            chan bool
	interval        time.Duration
	multisigAddress string

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *BurnExecutorService) Start() {
	log.Info("[BURN EXECUTOR] Starting service")
	stop := false
	for !stop {
		log.Info("[BURN EXECUTOR] Starting sync")

		x.SyncTxs()

		x.UpdateHealth()

		log.Info("[BURN EXECUTOR] Finished sync, Sleeping for ", x.interval)

		select {
		case <-x.stop:
			stop = true
			log.Info("[BURN EXECUTOR] Stopped service")
		case <-time.After(x.interval):
		}
	}
	x.wg.Done()
}

func (x *BurnExecutorService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *BurnExecutorService) UpdateHealth() {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           BurnExecutorName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (x *BurnExecutorService) Stop() {
	log.Debug("[BURN EXECUTOR] Stopping service")
	x.stop <- true
}

func (x *BurnExecutorService) HandleInvalidMint(doc models.InvalidMint) bool {
	log.Debug("[BURN EXECUTOR] Handling invalid mint: ", doc.TransactionHash)

	var filter bson.M
	var update bson.M

	if doc.Status == models.StatusSigned {
		log.Debug("[BURN EXECUTOR] Submitting invalid mint")
		p := rpc.SendRawTxParams{
			Addr:        x.multisigAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res, err := x.client.SubmitRawTx(p)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error submitting transaction: ", err)
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		update = bson.M{
			"$set": bson.M{
				"status":         models.StatusSubmitted,
				"return_tx_hash": res.TransactionHash,
				"updated_at":     time.Now(),
			},
		}
	} else if doc.Status == models.StatusSubmitted {
		log.Debug("[BURN EXECUTOR] Checking invalid mint")
		_, err := x.client.GetTx(doc.ReturnTxHash)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error fetching transaction: ", err)
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update = bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}
	}

	err := app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error updating invalid mint: ", err)
		return false
	}
	log.Info("[BURN EXECUTOR] Handled invalid mint")

	return true
}

func (x *BurnExecutorService) HandleBurn(doc models.Burn) bool {
	log.Debug("[BURN EXECUTOR] Handling burn: ", doc.TransactionHash)

	var filter bson.M
	var update bson.M

	if doc.Status == models.StatusSigned {
		log.Debug("[BURN EXECUTOR] Submitting burn")
		p := rpc.SendRawTxParams{
			Addr:        x.multisigAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res, err := x.client.SubmitRawTx(p)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error submitting transaction: ", err)
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		update = bson.M{
			"$set": bson.M{
				"status":         models.StatusSubmitted,
				"return_tx_hash": res.TransactionHash,
				"updated_at":     time.Now(),
			},
		}
	} else if doc.Status == models.StatusSubmitted {
		log.Debug("[BURN EXECUTOR] Checking burn")
		_, err := x.client.GetTx(doc.ReturnTxHash)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error fetching transaction: ", err)
			return false
		}

		filter = bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update = bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}
	}

	err := app.DB.UpdateOne(models.CollectionBurns, filter, update)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error updating burn: ", err)
		return false
	}
	log.Info("[BURN EXECUTOR] Handled burn")
	return true
}

func (x *BurnExecutorService) SyncTxs() bool {
	filter := bson.M{
		"status": bson.M{
			"$in": []string{
				string(models.StatusSigned),
				string(models.StatusSubmitted),
			},
		},
		"vault_address": x.multisigAddress,
	}
	invalidMints := []models.InvalidMint{}

	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error fetching invalid mints: ", err)
		return false
	}

	log.Info("[BURN EXECUTOR] Found invalid mints: ", len(invalidMints))

	var success bool = true
	for _, doc := range invalidMints {
		success = x.HandleInvalidMint(doc) && success
	}

	filter = bson.M{
		"status": bson.M{
			"$in": []string{
				string(models.StatusSigned),
				string(models.StatusSubmitted),
			},
		},
		"wpokt_address": x.wpoktAddress,
	}
	burns := []models.Burn{}

	err = app.DB.FindMany(models.CollectionBurns, filter, &burns)
	if err != nil {
		log.Error("[BURN EXECUTOR] Error fetching burns: ", err)
		return false
	}

	log.Info("[BURN EXECUTOR] Found burns: ", len(burns))

	for _, doc := range burns {
		success = x.HandleBurn(doc) && success
	}

	return success
}

func NewExecutor(wg *sync.WaitGroup, health models.ServiceHealth) models.Service {
	if !app.Config.BurnExecutor.Enabled {
		log.Debug("[BURN EXECUTOR] Pokt executor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[BURN EXECUTOR] Initializing burn executor")

	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[BURN EXECUTOR] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	multisigAddress := multisigPk.Address().String()
	log.Debug("[BURN EXECUTOR] Multisig address: ", multisigAddress)
	if multisigAddress != app.Config.Pocket.VaultAddress {
		log.Fatal("[BURN EXECUTOR] Multisig address does not match vault address")
	}

	x := &BurnExecutorService{
		wg:              wg,
		interval:        time.Duration(app.Config.BurnExecutor.IntervalSecs) * time.Second,
		stop:            make(chan bool),
		multisigAddress: multisigAddress,
		wpoktAddress:    app.Config.Ethereum.WrappedPocketAddress,
		client:          pokt.NewClient(),
	}

	x.UpdateHealth()

	log.Info("[BURN EXECUTOR] Initialized burn executor")

	return x
}
