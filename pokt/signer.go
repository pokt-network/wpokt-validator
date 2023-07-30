package pokt

import (
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/dan13ram/wpokt-validator/pokt/util"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	BurnSignerName = "burn signer"
)

type BurnSignerService struct {
	wg             *sync.WaitGroup
	stop           chan bool
	interval       time.Duration
	privateKey     crypto.PrivateKey
	multisigPubKey crypto.PublicKeyMultiSig
	numSigners     int
	ethClient      eth.EthereumClient
	poktClient     pokt.PocketClient
	poktHeight     int64
	ethBlockNumber int64
	vaultAddress   string
	wpoktAddress   string

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *BurnSignerService) Start() {
	log.Info("[BURN SIGNER] Starting service")
	stop := false
	for !stop {
		log.Info("[BURN SIGNER] Starting sync")

		x.UpdateBlocks()

		x.SyncTxs()

		x.UpdateHealth()

		log.Info("[BURN SIGNER] Finished sync, Sleeping for ", x.interval)

		select {
		case <-x.stop:
			stop = true
			log.Info("[BURN SIGNER] Stopped service")
		case <-time.After(x.interval):
		}
	}
	x.wg.Done()
}

func (x *BurnSignerService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *BurnSignerService) UpdateHealth() {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           BurnSignerName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     strconv.FormatInt(x.poktHeight, 10),
		EthBlockNumber: strconv.FormatInt(x.ethBlockNumber, 10),
		Healthy:        true,
	}
}

func (x *BurnSignerService) Stop() {
	log.Debug("[BURN SIGNER] Stopping service")
	x.stop <- true
}

func (x *BurnSignerService) UpdateBlocks() {
	log.Debug("[BURN SIGNER] Updating blocks")

	poktHeight, err := x.poktClient.GetHeight()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching pokt block height: ", err)
		return
	}
	x.poktHeight = poktHeight.Height

	ethBlockNumber, err := x.ethClient.GetBlockNumber()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching eth block number: ", err)
		return
	}
	x.ethBlockNumber = int64(ethBlockNumber)

	log.Info("[BURN SIGNER] Updated blocks")
}

func (x *BurnSignerService) HandleInvalidMint(doc models.InvalidMint) bool {
	log.Info("[BURN SIGNER] Handling invalid mint: ", doc.TransactionHash)

	doc, err := util.UpdateStatusAndConfirmationsForInvalidMint(doc, x.poktHeight)
	if err != nil {
		log.Error("[BURN SIGNER] Error getting invalid mint status: ", err)
		return false
	}

	var update bson.M

	if doc.Status == models.StatusConfirmed {
		log.Debug("[BURN SIGNER] Signing invalid mint")

		doc, err = util.SignInvalidMint(doc, x.privateKey, x.multisigPubKey, x.numSigners)
		if err != nil {
			log.Error("[BURN SIGNER] Error signing invalid mint: ", err)
			return false
		}

		update = bson.M{
			"$set": bson.M{
				"return_tx":     doc.ReturnTx,
				"signers":       doc.Signers,
				"status":        doc.Status,
				"confirmations": doc.Confirmations,
				"updated_at":    time.Now(),
			},
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

	filter := bson.M{"_id": doc.Id}
	err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating invalid mint: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Updated invalid mint")
	return true
}

func (x *BurnSignerService) HandleBurn(doc models.Burn) bool {
	log.Debug("[BURN SIGNER] Handling burn: ", doc.TransactionHash)

	doc, err := util.UpdateStatusAndConfirmationsForBurn(doc, x.ethBlockNumber)
	if err != nil {
		log.Error("[BURN SIGNER] Error getting burn status: ", err)
		return false
	}

	var update bson.M

	if doc.Status == models.StatusConfirmed {
		log.Debug("[BURN SIGNER] Signing burn")
		doc, err = util.SignBurn(doc, x.privateKey, x.multisigPubKey, x.numSigners)
		if err != nil {
			log.Error("[BURN SIGNER] Error signing burn: ", err)
			return false
		}

		update = bson.M{
			"$set": bson.M{
				"return_tx":     doc.ReturnTx,
				"signers":       doc.Signers,
				"status":        doc.Status,
				"confirmations": doc.Confirmations,
				"updated_at":    time.Now(),
			},
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

	filter := bson.M{"_id": doc.Id}
	err = app.DB.UpdateOne(models.CollectionBurns, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating burn: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Updated burn")

	return true
}

func (x *BurnSignerService) SyncTxs() bool {
	log.Info("[BURN SIGNER] Syncing txs")
	filter := bson.M{
		"vault_address": x.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers":       bson.M{"$nin": []string{x.privateKey.PublicKey().RawString()}},
	}

	invalidMints := []models.InvalidMint{}
	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching invalid mints: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Found invalid mints: ", len(invalidMints))

	var success bool = true
	for _, doc := range invalidMints {
		success = x.HandleInvalidMint(doc) && success
	}

	filter = bson.M{
		"wpokt_address": x.wpoktAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers":       bson.M{"$nin": []string{x.privateKey.PublicKey().RawString()}},
	}

	burns := []models.Burn{}
	err = app.DB.FindMany(models.CollectionBurns, filter, &burns)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching burns: ", err)
		return false
	}
	log.Info("[BURN SIGNER] Found burns: ", len(burns))
	for _, doc := range burns {
		success = x.HandleBurn(doc) && success
	}
	log.Info("[BURN SIGNER] Synced txs")

	return success
}

func NewSigner(wg *sync.WaitGroup, health models.ServiceHealth) models.Service {
	if app.Config.BurnSigner.Enabled == false {
		log.Debug("[BURN SIGNER] BURN signer disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[BURN SIGNER] Initializing burn signer")

	pk, err := crypto.NewPrivateKey(app.Config.Pocket.PrivateKey)
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing burn signer: ", err)
	}
	log.Info("[BURN SIGNER] POKT signer public key: ", pk.PublicKey().RawString())
	log.Debug("[BURN SIGNER] POKT signer address: ", pk.PublicKey().Address().String())
	log.Debug("[BURN SIGNER] Initialized burn signer private key")

	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[BURN SIGNER] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	log.Debug("[BURN SIGNER] Multisig address: ", multisigPk.Address().String())

	poktClient := pokt.NewClient()
	ethClient, err := eth.NewClient()
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing ethereum client: ", err)
	}

	x := &BurnSignerService{
		wg:             wg,
		interval:       time.Duration(app.Config.BurnSigner.IntervalSecs) * time.Second,
		stop:           make(chan bool),
		privateKey:     pk,
		multisigPubKey: multisigPk,
		numSigners:     len(pks),
		ethClient:      ethClient,
		poktClient:     poktClient,
		vaultAddress:   app.Config.Pocket.VaultAddress,
		wpoktAddress:   app.Config.Ethereum.WrappedPocketAddress,
	}

	x.UpdateBlocks()

	x.UpdateHealth()

	log.Info("[BURN SIGNER] Initialized burn signer")

	return x
}
