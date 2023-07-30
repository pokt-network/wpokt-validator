package pokt

import (
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/dan13ram/wpokt-validator/pokt/util"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MintMonitorName = "mint monitor"
)

type MintMonitorService struct {
	wg            *sync.WaitGroup
	client        pokt.PocketClient
	stop          chan bool
	wpoktAddress  string
	vaultAddress  string
	interval      time.Duration
	startHeight   int64
	currentHeight int64

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *MintMonitorService) Start() {
	log.Info("[MINT MONITOR] Starting service")
	stop := false
	for !stop {
		log.Info("[MINT MONITOR] Starting sync")

		x.UpdateCurrentHeight()

		x.SyncTxs()

		x.UpdateHealth()

		log.Info("[MINT MONITOR] Finished sync, Sleeping for ", x.interval)

		select {
		case <-x.stop:
			stop = true
			log.Info("[MINT MONITOR] Stopped service")
		case <-time.After(x.interval):
		}
	}
	x.wg.Done()
}

func (x *MintMonitorService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *MintMonitorService) UpdateHealth() {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           MintMonitorName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     strconv.FormatInt(x.startHeight, 10),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (x *MintMonitorService) Stop() {
	log.Debug("[MINT MONITOR] Stopping service")
	x.stop <- true
}

func (x *MintMonitorService) UpdateCurrentHeight() {
	res, err := x.client.GetHeight()
	if err != nil {
		log.Error("[MINT MONITOR] Error getting current height: ", err)
		return
	}
	x.currentHeight = res.Height
	log.Info("[MINT MONITOR] Current height: ", x.currentHeight)
}

func (x *MintMonitorService) HandleInvalidMint(tx *pokt.TxResponse) bool {
	doc := util.CreateInvalidMint(tx, x.vaultAddress)

	log.Debug("[MINT MONITOR] Storing invalid mint tx")
	err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
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

func (x *MintMonitorService) HandleValidMint(tx *pokt.TxResponse, memo models.MintMemo) bool {
	doc := util.CreateMint(tx, memo, x.wpoktAddress, x.vaultAddress)

	log.Debug("[MINT MONITOR] Storing mint tx")
	err := app.DB.InsertOne(models.CollectionMints, doc)
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

func (x *MintMonitorService) SyncTxs() bool {

	if x.currentHeight <= x.startHeight {
		log.Info("[MINT MONITOR] No new blocks to sync")
		return true
	}

	txs, err := x.client.GetAccountTxsByHeight(x.vaultAddress, x.startHeight)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Info("[MINT MONITOR] Found ", len(txs), " txs to sync")
	var success bool = true
	for _, tx := range txs {
		memo, ok := util.ValidateMemo(tx.StdTx.Memo)

		if !ok {
			log.Info("[MINT MONITOR] Found invalid mint tx: ", tx.Hash, " with memo: ", "\""+tx.StdTx.Memo+"\"")
			success = x.HandleInvalidMint(tx) && success
			continue
		}

		log.Info("[MINT MONITOR] Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
		success = x.HandleValidMint(tx, memo) && success
	}

	if success {
		x.startHeight = x.currentHeight
	}

	return success
}

func NewMonitor(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if !app.Config.MintMonitor.Enabled {
		log.Debug("[MINT MONITOR] Pokt monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[MINT MONITOR] Initializing mint monitor")
	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[MINT MONITOR] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	multisigAddress := multisigPk.Address().String()
	log.Debug("[MINT EXECUTOR] Multisig address: ", multisigAddress)

	x := &MintMonitorService{
		wg:            wg,
		interval:      time.Duration(app.Config.MintMonitor.IntervalSecs) * time.Second,
		vaultAddress:  multisigAddress,
		wpoktAddress:  app.Config.Ethereum.WrappedPocketAddress,
		startHeight:   0,
		currentHeight: 0,
		stop:          make(chan bool),
		client:        pokt.NewClient(),
	}

	startHeight := (app.Config.Pocket.StartHeight)

	if (lastHealth.PoktHeight) != "" {
		if lastHeight, err := strconv.ParseInt(lastHealth.PoktHeight, 10, 64); err == nil {
			startHeight = lastHeight
		}
	}

	x.UpdateCurrentHeight()

	if startHeight > 0 {
		x.startHeight = startHeight
	} else {
		log.Info("[MINT MONITOR] Found invalid start height, using current height")
		x.startHeight = x.currentHeight
	}
	log.Info("[MINT MONITOR] Start height: ", x.startHeight)

	x.UpdateHealth()

	log.Info("[MINT MONITOR] Initialized mint monitor")

	return x
}
