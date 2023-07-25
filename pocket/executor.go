package pocket

import (
	"sync"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	pocket "github.com/dan13ram/wpokt-backend/pocket/client"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	PoktExecutorName = "pokt-executor"
)

type PoktExecutorService struct {
	wg              *sync.WaitGroup
	name            string
	client          pocket.PocketClient
	wpoktAddress    string
	stop            chan bool
	lastSyncTime    time.Time
	interval        time.Duration
	multisigAddress string
}

func (m *PoktExecutorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.Name(),
		LastSyncTime:   m.LastSyncTime(),
		NextSyncTime:   m.LastSyncTime().Add(m.Interval()),
		PoktHeight:     "",
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (m *PoktExecutorService) LastSyncTime() time.Time {
	return m.lastSyncTime
}

func (m *PoktExecutorService) Interval() time.Duration {
	return m.interval
}

func (m *PoktExecutorService) Name() string {
	return m.name
}

func (m *PoktExecutorService) Start() {
	log.Debug("[POKT EXECUTOR] Starting pokt executor")
	stop := false
	for !stop {
		log.Debug("[POKT EXECUTOR] Starting pokt executor sync")
		m.lastSyncTime = time.Now()

		m.SyncTxs()

		log.Debug("[POKT EXECUTOR] Finished pokt executor sync")
		log.Debug("[POKT EXECUTOR] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[POKT EXECUTOR] Stopped pokt executor")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *PoktExecutorService) Stop() {
	log.Debug("[POKT EXECUTOR] Stopping pokt executor")
	m.stop <- true
}

func (m *PoktExecutorService) HandleInvalidMint(doc models.InvalidMint) bool {
	log.Debug("[POKT EXECUTOR] Handling invalid mint: ", doc.TransactionHash)
	if doc.Status == models.StatusSigned {
		p := rpc.SendRawTxParams{
			Addr:        m.multisigAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res, err := m.client.SubmitRawTx(p)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error submitting transaction: ", err)
			return false
		}

		filter := bson.M{
			"_id":           doc.Id,
			"vault_address": m.multisigAddress,
		}
		update := bson.M{
			"$set": bson.M{
				"status":         models.StatusSubmitted,
				"return_tx_hash": res.TransactionHash,
			},
		}
		err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error updating invalid mint: ", err)
			return false
		}
		log.Debug("[POKT EXECUTOR] Submitted tx for invalid mint")

	} else if doc.Status == models.StatusSubmitted {
		_, err := m.client.GetTx(doc.ReturnTxHash)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error fetching transaction: ", err)
			return false
		}
		filter := bson.M{
			"_id":           doc.Id,
			"vault_address": m.multisigAddress,
		}
		update := bson.M{
			"$set": bson.M{
				"status": models.StatusSuccess,
			},
		}
		err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error updating invalid mint: ", err)
			return false
		}
		log.Debug("[POKT EXECUTOR] Executed return tx for invalid mint")
	}

	return true
}

func (m *PoktExecutorService) HandleBurn(doc models.Burn) bool {
	log.Debug("[POKT EXECUTOR] Handling burn: ", doc.TransactionHash)
	if doc.Status == models.StatusSigned {
		p := rpc.SendRawTxParams{
			Addr:        m.multisigAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res, err := m.client.SubmitRawTx(p)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error submitting transaction: ", err)
			return false
		}

		filter := bson.M{
			"_id":           doc.Id,
			"wpokt_address": m.wpoktAddress,
		}
		update := bson.M{
			"$set": bson.M{
				"status":         models.StatusSubmitted,
				"return_tx_hash": res.TransactionHash,
			},
		}
		err = app.DB.UpdateOne(models.CollectionBurns, filter, update)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error updating burn: ", err)
			return false
		}
		log.Debug("[POKT EXECUTOR] Submitted tx for burn")

	} else if doc.Status == models.StatusSubmitted {
		_, err := m.client.GetTx(doc.ReturnTxHash)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error fetching transaction: ", err)
			return false
		}
		filter := bson.M{
			"_id":           doc.Id,
			"wpokt_address": m.wpoktAddress,
		}
		update := bson.M{
			"$set": bson.M{
				"status": models.StatusSuccess,
			},
		}
		err = app.DB.UpdateOne(models.CollectionBurns, filter, update)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error updating burn: ", err)
			return false
		}
		log.Debug("[POKT EXECUTOR] Executed return tx for burn")
	}
	return true
}

func (m *PoktExecutorService) SyncTxs() bool {
	// filter for status signed or status submitted
	filter := bson.M{
		"status": bson.M{
			"$in": []string{
				string(models.StatusSigned),
				string(models.StatusSubmitted),
			},
		},
	}
	invalidMints := []models.InvalidMint{}
	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[POKT EXECUTOR] Error fetching invalid mints: ", err)
		return false
	}
	log.Debug("[POKT EXECUTOR] Found invalid mints: ", len(invalidMints))

	var success bool = true
	for _, doc := range invalidMints {
		success = m.HandleInvalidMint(doc) && success
	}

	burns := []models.Burn{}
	err = app.DB.FindMany(models.CollectionBurns, filter, &burns)
	if err != nil {
		log.Error("[POKT EXECUTOR] Error fetching burns: ", err)
		return false
	}

	log.Debug("[POKT EXECUTOR] Found burns: ", len(burns))

	for _, doc := range burns {
		success = m.HandleBurn(doc) && success
	}

	return success
}

func newExecutor(wg *sync.WaitGroup) models.Service {
	if !app.Config.PoktExecutor.Enabled {
		log.Debug("[POKT EXECUTOR] Pokt executor disabled")
		return models.NewEmptyService(wg, "empty-pokt-executor")
	}

	log.Debug("[POKT EXECUTOR] Initializing pokt executor")

	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	multisigAddress := multisigPk.Address().String()
	log.Debug("[POKT EXECUTOR] Multisig address: ", multisigAddress)
	if multisigAddress != app.Config.Pocket.VaultAddress {
		log.Fatal("[POKT EXECUTOR] Multisig address does not match vault address")
	}

	m := &PoktExecutorService{
		wg:              wg,
		name:            PoktExecutorName,
		interval:        time.Duration(app.Config.PoktExecutor.IntervalSecs) * time.Second,
		stop:            make(chan bool),
		multisigAddress: multisigAddress,
		wpoktAddress:    app.Config.Ethereum.WPOKTAddress,
		client:          pocket.NewClient(),
	}

	log.Debug("[POKT EXECUTOR] Initialized pokt executor")

	return m
}

func NewExecutor(wg *sync.WaitGroup) models.Service {
	return newExecutor(wg)
}

func NewExecutorWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	return newExecutor(wg)
}
