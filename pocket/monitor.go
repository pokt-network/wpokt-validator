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

type PoktMonitorService struct {
	wg            *sync.WaitGroup
	name          string
	client        pocket.PocketClient
	stop          chan bool
	vaultAddress  string
	lastSyncTime  time.Time
	interval      time.Duration
	startHeight   uint64
	currentHeight uint64
}

func (m *PoktMonitorService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.Name(),
		LastSyncTime:   m.LastSyncTime(),
		NextSyncTime:   m.LastSyncTime().Add(m.Interval()),
		PoktHeight:     m.PoktHeight(),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (m *PoktMonitorService) PoktHeight() string {
	return strconv.FormatUint(m.startHeight, 10)
}

func (m *PoktMonitorService) LastSyncTime() time.Time {
	return m.lastSyncTime
}

func (m *PoktMonitorService) Interval() time.Duration {
	return m.interval
}

func (m *PoktMonitorService) Name() string {
	return m.name
}

func (m *PoktMonitorService) Start() {
	log.Debug("[POKT MONITOR] Starting pokt monitor")
	stop := false
	for !stop {
		log.Debug("[POKT MONITOR] Starting pokt monitor sync")
		m.lastSyncTime = time.Now()

		m.UpdateCurrentHeight()

		if (m.currentHeight - m.startHeight) > 0 {
			log.Debug("[POKT MONITOR] Syncing mint txs from height: ", m.startHeight, " to height: ", m.currentHeight)
			success := m.SyncTxs()
			if success {
				m.startHeight = m.currentHeight
			}
		} else {
			log.Debug("[POKT MONITOR] No new blocks to sync")
		}

		log.Debug("[POKT MONITOR] Finished pokt monitor sync")
		log.Debug("[POKT MONITOR] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[POKT MONITOR] Stopped pokt monitor")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *PoktMonitorService) Stop() {
	log.Debug("[POKT MONITOR] Stopping pokt monitor")
	m.stop <- true
}

func (m *PoktMonitorService) UpdateCurrentHeight() {
	res, err := m.client.GetHeight()
	if err != nil {
		log.Error("[POKT MONITOR] Error getting current height: ", err)
		if log.GetLevel() == log.DebugLevel {
			log.Debug("[POKT MONITOR] Response: ", res)
		}
		return
	}
	m.currentHeight = uint64(res.Height)
	log.Debug("[POKT MONITOR] Current height: ", m.currentHeight)
}

func (m *PoktMonitorService) HandleInvalidMint(tx *pocket.TxResponse) bool {
	doc := models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		Confirmations:   "0",
		TransactionHash: tx.Hash,
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.Pocket.ChainId,
		Amount:          tx.StdTx.Msg.Value.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
		ReturnTx:        "",
		ReturnTxHash:    "",
	}

	log.Debug("[POKT MONITOR] Storing invalid mint tx")

	err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[POKT MONITOR] Found duplicate invalid mint tx")
			return true
		}
		log.Error("[POKT MONITOR] Error storing invalid mint tx: ", err)
		return false
	}

	log.Debug("[POKT MONITOR] Stored invalid mint tx")
	return true
}

func (m *PoktMonitorService) HandleValidMint(tx *pocket.TxResponse, memo models.MintMemo) bool {
	doc := models.Mint{
		Height:              strconv.FormatInt(tx.Height, 10),
		TransactionHash:     tx.Hash,
		SenderAddress:       tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:       app.Config.Pocket.ChainId,
		RecipientAddress:    memo.Address,
		RecipientChainId:    memo.ChainId,
		Amount:              tx.StdTx.Msg.Value.Amount,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Status:              models.StatusPending,
		Data:                nil,
		Signers:             []string{},
		Signatures:          []string{},
		MintTransactionHash: "",
	}

	log.Debug("[POKT MONITOR] Storing mint tx")

	err := app.DB.InsertOne(models.CollectionMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[POKT MONITOR] Found duplicate mint tx")
			return true
		}
		log.Error("[POKT MONITOR] Error storing mint tx: ", err)
		return false
	}

	log.Debug("[POKT MONITOR] Stored mint tx")
	return true
}

func (m *PoktMonitorService) HandleTx(tx *pocket.TxResponse) bool {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

	address := common.HexToAddress(memo.Address).Hex()

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId || strings.ToLower(address) != strings.ToLower(memo.Address) {
		log.Debug("[POKT MONITOR] Found invalid mint tx: ", tx.Hash, " with memo: ", "\""+tx.StdTx.Memo+"\"")
		return m.HandleInvalidMint(tx)
	}
	memo.Address = address
	log.Debug("[POKT MONITOR] Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	return m.HandleValidMint(tx, memo)

}

func (m *PoktMonitorService) SyncTxs() bool {
	txs, err := m.client.GetAccountTxsByHeight(m.vaultAddress, int64(m.startHeight))
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("[POKT MONITOR] Found ", len(txs), " txs to sync")
	var success bool = true
	for _, tx := range txs {
		success = m.HandleTx(tx) && success
	}
	return success
}

func NewMonitor(wg *sync.WaitGroup) models.Service {
	if !app.Config.PoktMonitor.Enabled {
		log.Debug("[POKT MONITOR] Pokt monitor disabled")
		return models.NewEmptyService(wg, "empty-pokt-monitor")
	}

	log.Debug("[POKT MONITOR] Initializing pokt monitor")

	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Error("[POKT MONITOR] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	multisigAddress := multisigPk.Address().String()
	log.Debug("[POKT EXECUTOR] Multisig address: ", multisigAddress)

	m := &PoktMonitorService{
		wg:            wg,
		name:          "pokt-monitor",
		interval:      time.Duration(app.Config.PoktMonitor.IntervalSecs) * time.Second,
		vaultAddress:  multisigAddress,
		startHeight:   0,
		currentHeight: 0,
		stop:          make(chan bool),
		client:        pocket.NewClient(),
	}

	m.UpdateCurrentHeight()
	if app.Config.Pocket.StartHeight > 0 {
		m.startHeight = uint64(app.Config.Pocket.StartHeight)
	} else {
		log.Debug("[POKT MONITOR] Found invalid start height, updating current height")
		m.startHeight = m.currentHeight
	}

	log.Debug("[POKT MONITOR] Start height: ", m.startHeight)
	log.Debug("[POKT MONITOR] Initialized pokt monitor")

	return m
}
