package pocket

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	pocket "github.com/dan13ram/wpokt-backend/pocket/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MintMonitorName = "mint-monitor"
)

type MintMonitorService struct {
	wg            *sync.WaitGroup
	name          string
	client        pocket.PocketClient
	stop          chan bool
	wpoktAddress  string
	vaultAddress  string
	lastSyncTime  time.Time
	interval      time.Duration
	startHeight   int64
	currentHeight int64
}

func (m *MintMonitorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.Name(),
		LastSyncTime:   m.LastSyncTime(),
		NextSyncTime:   m.LastSyncTime().Add(m.Interval()),
		PoktHeight:     m.PoktHeight(),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (m *MintMonitorService) InitStartHeight(height int64) {
	m.UpdateCurrentHeight()
	if height > 0 {
		m.startHeight = height
	} else {
		log.Debug("[MINT MONITOR] Found invalid start height, updating current height")
		m.startHeight = m.currentHeight
	}
	log.Debug("[MINT MONITOR] Start height: ", m.startHeight)
}

func (m *MintMonitorService) PoktHeight() string {
	return strconv.FormatInt(m.startHeight, 10)
}

func (m *MintMonitorService) LastSyncTime() time.Time {
	return m.lastSyncTime
}

func (m *MintMonitorService) Interval() time.Duration {
	return m.interval
}

func (m *MintMonitorService) Name() string {
	return m.name
}

func (m *MintMonitorService) Start() {
	log.Debug("[MINT MONITOR] Starting pokt monitor")
	stop := false
	for !stop {
		log.Debug("[MINT MONITOR] Starting pokt monitor sync")
		m.lastSyncTime = time.Now()

		m.UpdateCurrentHeight()

		if (m.currentHeight - m.startHeight) > 0 {
			log.Debug("[MINT MONITOR] Syncing mint txs from height: ", m.startHeight, " to height: ", m.currentHeight)
			success := m.SyncTxs()
			if success {
				m.startHeight = m.currentHeight
			}
		} else {
			log.Debug("[MINT MONITOR] No new blocks to sync")
		}

		log.Debug("[MINT MONITOR] Finished pokt monitor sync")
		log.Debug("[MINT MONITOR] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[MINT MONITOR] Stopped pokt monitor")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *MintMonitorService) Stop() {
	log.Debug("[MINT MONITOR] Stopping pokt monitor")
	m.stop <- true
}

func (m *MintMonitorService) UpdateCurrentHeight() {
	res, err := m.client.GetHeight()
	if err != nil {
		log.Error("[MINT MONITOR] Error getting current height: ", err)
		if log.GetLevel() == log.DebugLevel {
			log.Debug("[MINT MONITOR] Response: ", res)
		}
		return
	}
	m.currentHeight = res.Height
	log.Debug("[MINT MONITOR] Current height: ", m.currentHeight)
}

func (m *MintMonitorService) HandleInvalidMint(tx *pocket.TxResponse) bool {
	doc := models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: tx.Hash,
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.Pocket.ChainId,
		Memo:            tx.StdTx.Memo,
		Amount:          tx.StdTx.Msg.Value.Amount,
		VaultAddress:    m.vaultAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
		ReturnTx:        "",
		ReturnTxHash:    "",
	}

	log.Debug("[MINT MONITOR] Storing invalid mint tx")

	err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[MINT MONITOR] Found duplicate invalid mint tx")
			return true
		}
		log.Error("[MINT MONITOR] Error storing invalid mint tx: ", err)
		return false
	}

	log.Debug("[MINT MONITOR] Stored invalid mint tx")
	return true
}

func (m *MintMonitorService) HandleValidMint(tx *pocket.TxResponse, memo models.MintMemo) bool {
	doc := models.Mint{
		Height:              strconv.FormatInt(tx.Height, 10),
		Confirmations:       "0",
		TransactionHash:     tx.Hash,
		SenderAddress:       tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:       app.Config.Pocket.ChainId,
		RecipientAddress:    memo.Address,
		RecipientChainId:    memo.ChainId,
		WPOKTAddress:        m.wpoktAddress,
		VaultAddress:        m.vaultAddress,
		Amount:              tx.StdTx.Msg.Value.Amount,
		Memo:                &memo,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Status:              models.StatusPending,
		Data:                nil,
		Signers:             []string{},
		Signatures:          []string{},
		MintTransactionHash: "",
	}

	log.Debug("[MINT MONITOR] Storing mint tx")

	err := app.DB.InsertOne(models.CollectionMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[MINT MONITOR] Found duplicate mint tx")
			return true
		}
		log.Error("[MINT MONITOR] Error storing mint tx: ", err)
		return false
	}

	log.Debug("[MINT MONITOR] Stored mint tx")
	return true
}

func (m *MintMonitorService) HandleTx(tx *pocket.TxResponse) bool {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

	address := common.HexToAddress(memo.Address).Hex()

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId || strings.ToLower(address) != strings.ToLower(memo.Address) {
		log.Debug("[MINT MONITOR] Found invalid mint tx: ", tx.Hash, " with memo: ", "\""+tx.StdTx.Memo+"\"")
		return m.HandleInvalidMint(tx)
	}
	memo.Address = address
	log.Debug("[MINT MONITOR] Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	return m.HandleValidMint(tx, memo)

}

func (m *MintMonitorService) SyncTxs() bool {
	txs, err := m.client.GetAccountTxsByHeight(m.vaultAddress, int64(m.startHeight))
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("[MINT MONITOR] Found ", len(txs), " txs to sync")
	var success bool = true
	for _, tx := range txs {
		success = m.HandleTx(tx) && success
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
		wpoktAddress:  app.Config.Ethereum.WPOKTAddress,
		startHeight:   0,
		currentHeight: 0,
		stop:          make(chan bool),
		client:        pocket.NewClient(),
	}

	return m
}

func NewMonitor(wg *sync.WaitGroup) models.Service {
	if !app.Config.MintMonitor.Enabled {
		log.Debug("[MINT MONITOR] Pokt monitor disabled")
		return models.NewEmptyService(wg, "empty-pokt-monitor")
	}

	log.Debug("[MINT MONITOR] Initializing pokt monitor")

	m := newMonitor(wg)

	m.InitStartHeight(app.Config.Pocket.StartHeight)

	log.Debug("[MINT MONITOR] Initialized pokt monitor")

	return m
}

func NewMonitorWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if !app.Config.MintMonitor.Enabled {
		log.Debug("[MINT MONITOR] Pokt monitor disabled")
		return models.NewEmptyService(wg, "empty-pokt-monitor")
	}

	log.Debug("[MINT MONITOR] Initializing pokt monitor with last health")

	m := newMonitor(wg)

	lastHeight, err := strconv.ParseInt(lastHealth.PoktHeight, 10, 64)
	if err != nil {
		log.Error("[MINT MONITOR] Error parsing last height: ", err)
		lastHeight = app.Config.Pocket.StartHeight
	}

	m.InitStartHeight(lastHeight)

	log.Debug("[MINT MONITOR] Initialized pokt monitor")

	return m
}
