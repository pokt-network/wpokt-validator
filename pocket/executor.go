package pocket

import (
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type PoktExecutorService struct {
	stop            chan bool
	interval        time.Duration
	multisigAddress string
}

func (m *PoktExecutorService) Start() {
	log.Debug("[POKT EXECUTOR] Starting pokt executor")
	stop := false
	for !stop {
		log.Debug("[POKT EXECUTOR] Starting pokt executor sync")

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

		res, err := Client.SubmitRawTx(p)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error submitting transaction: ", err)
			return false
		}

		filter := bson.M{"_id": doc.Id}
		update := bson.M{"$set": bson.M{"status": models.StatusSubmitted, "return_tx_hash": res.TransactionHash}}
		err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error updating invalid mint: ", err)
			return false
		}
		log.Debug("[POKT EXECUTOR] Submitted tx for invalid mint")

	} else if doc.Status == models.StatusSubmitted {
		_, err := Client.GetTx(doc.ReturnTxHash)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error fetching transaction: ", err)
			return false
		}
		filter := bson.M{"_id": doc.Id}
		update := bson.M{"$set": bson.M{"status": models.StatusSuccess}}
		err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
		if err != nil {
			log.Error("[POKT EXECUTOR] Error updating invalid mint: ", err)
			return false
		}
		log.Debug("[POKT EXECUTOR] Executed return tx for invalid mint")
	}

	return true
}

func (m *PoktExecutorService) HandleBurn(mint models.Burn) bool {
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

	return success
}

func NewExecutor() models.Service {
	if !app.Config.PoktExecutor.Enabled {
		log.Debug("[POKT EXECUTOR] Pokt executor disabled")
		return models.NewEmptyService()
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

	m := &PoktExecutorService{
		interval:        time.Duration(app.Config.PoktExecutor.IntervalSecs) * time.Second,
		stop:            make(chan bool),
		multisigAddress: multisigAddress,
	}

	log.Debug("[POKT EXECUTOR] Initialized pokt executor")

	return m
}
