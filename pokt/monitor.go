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
	name          string
	client        pokt.PocketClient
	stop          chan bool
	wpoktAddress  string
	vaultAddress  string
	lastSyncTime  time.Time
	interval      time.Duration
	startHeight   int64
	currentHeight int64
}

func (m *MintMonitorService) Start() {
	log.Info("[MINT MONITOR] Starting service")
	stop := false
	for !stop {
		log.Info("[MINT MONITOR] Starting sync")
		m.lastSyncTime = time.Now()

		m.UpdateCurrentHeight()

		if (m.currentHeight - m.startHeight) > 0 {
			log.Info("[MINT MONITOR] Syncing mint txs from height: ", m.startHeight, " to height: ", m.currentHeight)
			success := m.SyncTxs()
			if success {
				m.startHeight = m.currentHeight
			}
		} else {
			log.Info("[MINT MONITOR] No new blocks to sync")
		}

		log.Info("[MINT MONITOR] Finished sync, Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Info("[MINT MONITOR] Stopped service")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *MintMonitorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.name,
		LastSyncTime:   m.lastSyncTime,
		NextSyncTime:   m.lastSyncTime.Add(m.interval),
		PoktHeight:     strconv.FormatInt(m.startHeight, 10),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (m *MintMonitorService) Stop() {
	log.Debug("[MINT MONITOR] Stopping service")
	m.stop <- true
}

func (m *MintMonitorService) InitStartHeight(height int64) {
	if height > 0 {
		m.startHeight = height
	} else {
		log.Info("[MINT MONITOR] Found invalid start height, using current height")
		m.startHeight = m.currentHeight
	}
	log.Info("[MINT MONITOR] Start height: ", m.startHeight)
}

func (m *MintMonitorService) UpdateCurrentHeight() {
	res, err := m.client.GetHeight()
	if err != nil {
		log.Error("[MINT MONITOR] Error getting current height: ", err)
		return
	}
	m.currentHeight = res.Height
	log.Info("[MINT MONITOR] Current height: ", m.currentHeight)
}

func (m *MintMonitorService) HandleInvalidMint(tx *pokt.TxResponse) bool {
	doc := util.CreateInvalidMint(tx, m.vaultAddress)

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

func (m *MintMonitorService) HandleValidMint(tx *pokt.TxResponse, memo models.MintMemo) bool {
	doc := util.CreateMint(tx, memo, m.wpoktAddress, m.vaultAddress)

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

func (m *MintMonitorService) SyncTxs() bool {
	txs, err := m.client.GetAccountTxsByHeight(m.vaultAddress, int64(m.startHeight))
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
			success = m.HandleInvalidMint(tx) && success
			continue
		}

		log.Info("[MINT MONITOR] Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
		success = m.HandleValidMint(tx, memo) && success
	}
	return success
}

func newMonitor(wg *sync.WaitGroup) *MintMonitorService {
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

	m := &MintMonitorService{
		wg:            wg,
		name:          MintMonitorName,
		interval:      time.Duration(app.Config.MintMonitor.IntervalSecs) * time.Second,
		vaultAddress:  multisigAddress,
		wpoktAddress:  app.Config.Ethereum.WrappedPocketAddress,
		startHeight:   0,
		currentHeight: 0,
		stop:          make(chan bool),
		client:        pokt.NewClient(),
	}

	return m
}

func NewMonitor(wg *sync.WaitGroup) models.Service {
	if !app.Config.MintMonitor.Enabled {
		log.Debug("[MINT MONITOR] Pokt monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[MINT MONITOR] Initializing mint monitor")

	m := newMonitor(wg)

	m.UpdateCurrentHeight()

	m.InitStartHeight(app.Config.Pocket.StartHeight)

	log.Info("[MINT MONITOR] Initialized mint monitor")

	return m
}

func NewMonitorWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if !app.Config.MintMonitor.Enabled {
		log.Debug("[MINT MONITOR] Pokt monitor disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[MINT MONITOR] Initializing mint monitor with last health")

	m := newMonitor(wg)

	lastHeight, err := strconv.ParseInt(lastHealth.PoktHeight, 10, 64)
	if err != nil {
		log.Error("[MINT MONITOR] Error parsing last height: ", err)
		lastHeight = app.Config.Pocket.StartHeight
	}

	m.UpdateCurrentHeight()

	m.InitStartHeight(lastHeight)

	log.Info("[MINT MONITOR] Initialized mint monitor")

	return m
}
