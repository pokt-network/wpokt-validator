package pocket

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type PoktMonitorService struct {
	stop          chan bool
	interval      time.Duration
	startHeight   uint64
	currentHeight uint64
}

func (m *PoktMonitorService) Start() {
	log.Debug("[POKT MONITOR] Starting pokt monitor")
	stop := false
	for !stop {
		log.Debug("[POKT MONITOR] Starting pokt monitor sync")

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
}

func (m *PoktMonitorService) Stop() {
	log.Debug("[POKT MONITOR] Stopping pokt monitor")
	m.stop <- true
}

func (m *PoktMonitorService) UpdateCurrentHeight() {
	res, err := Client.GetHeight()
	if err != nil {
		log.Error("[POKT MONITOR] Error getting current height: ", err)
		if log.GetLevel() == log.DebugLevel {
			log.Debug("[POKT MONITOR] Response: ", res)
		}
		return
	}
	log.Debug("[POKT MONITOR] Current height: ", m.currentHeight)
	m.currentHeight = uint64(res.Height)
}

func (m *PoktMonitorService) HandleInvalidMint(tx *ResultTx) bool {
	doc := models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		TransactionHash: tx.Hash,
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.PoktMonitor.ChainId,
		Amount:          tx.StdTx.Msg.Value.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
	}

	log.Debug("[POKT MONITOR] Storing invalid mint tx: ", tx.Hash, " in db")

	err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[POKT MONITOR] Found duplicate invalid mint tx: ", tx.Hash, " in db")
			return true
		}
		log.Error("[POKT MONITOR] Error storing invalid mint tx: ", err)
		return false
	}

	log.Debug("[POKT MONITOR] Stored invalid mint tx: ", tx.Hash, " in db")
	return true
}

func (m *PoktMonitorService) HandleValidMint(tx *ResultTx, memo models.MintMemo) bool {
	doc := models.Mint{
		Height:           strconv.FormatInt(tx.Height, 10),
		TransactionHash:  tx.Hash,
		SenderAddress:    tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:    app.Config.PoktMonitor.ChainId,
		RecipientAddress: memo.Address,
		RecipientChainId: strconv.FormatInt(int64(memo.ChainId), 10),
		Amount:           tx.StdTx.Msg.Value.Amount,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           models.StatusPending,
		Signers:          []string{},
	}

	log.Debug("[POKT MONITOR] Storing mint tx: ", tx.Hash, " in db")

	err := app.DB.InsertOne(models.CollectionMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[POKT MONITOR] Found duplicate mint tx: ", tx.Hash, " in db")
			return true
		}
		log.Error("[POKT MONITOR] Error storing mint tx: ", err)
		return false
	}

	log.Debug("[POKT MONITOR] Stored mint tx: ", tx.Hash, " in db")
	return true
}

func (m *PoktMonitorService) HandleTx(tx *ResultTx) bool {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId {
		log.Debug("[POKT MONITOR] Found invalid mint tx: ", tx.Hash, " with memo: ", "\""+tx.StdTx.Memo+"\"")
		return m.HandleInvalidMint(tx)
	}
	log.Debug("[POKT MONITOR] Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	return m.HandleValidMint(tx, memo)

}

func (m *PoktMonitorService) SyncTxs() bool {
	txs, err := Client.GetAccountTxsByHeight(int64(m.startHeight))
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("[POKT MONITOR] Found ", len(txs), " txs to sync")
	var success bool = true
	for _, tx := range txs {
		success = success && m.HandleTx(tx)
	}
	return success
}

func NewMonitor() Service {
	if !app.Config.PoktMonitor.Enabled {
		log.Debug("[POKT MONITOR] Pokt monitor disabled")
		return nil
	}

	log.Debug("[POKT MONITOR] Initializing pokt monitor")
	m := &PoktMonitorService{
		interval:      time.Duration(app.Config.PoktMonitor.IntervalSecs) * time.Second,
		startHeight:   0,
		currentHeight: 0,
		stop:          make(chan bool),
	}

	m.UpdateCurrentHeight()
	if app.Config.PoktMonitor.StartHeight > 0 {
		m.startHeight = uint64(app.Config.PoktMonitor.StartHeight)
	} else {
		log.Debug("[POKT MONITOR] Found invalid start height, updating current height")
		m.startHeight = m.currentHeight
	}

	log.Debug("[POKT MONITOR] Start height: ", m.startHeight)
	log.Debug("[POKT MONITOR] Initialized pokt monitor")

	return m
}
