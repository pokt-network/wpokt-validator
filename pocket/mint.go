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

type Service interface {
	Stop()
	Start()
}

type MintMonitorService struct {
	stop            chan bool
	monitorInterval time.Duration
	startHeight     uint64
	currentHeight   uint64
}

func (m *MintMonitorService) Start() {
	log.Debug("[MINT MONITOR] Starting mint monitor")
	stop := false
	for !stop {
		log.Debug("[MINT MONITOR] Starting mint sync")

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

		log.Debug("[MINT MONITOR] Finished mint sync")
		log.Debug("[MINT MONITOR] Sleeping for ", m.monitorInterval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[MINT MONITOR] Stopped mint monitor")
		case <-time.After(m.monitorInterval):
		}
	}
}

func (m *MintMonitorService) Stop() {
	log.Debug("[MINT MONITOR] Stopping mint monitor")
	m.stop <- true
}

func (m *MintMonitorService) UpdateCurrentHeight() {
	res, err := Client.GetHeight()
	if err != nil {
		log.Error("[MINT MONITOR] Error getting current height: ", err)
		if log.GetLevel() == log.DebugLevel {
			log.Debug("[MINT MONITOR] Response: ", res)
		}
		return
	}
	log.Debug("[MINT MONITOR] Current height: ", m.currentHeight)
	m.currentHeight = uint64(res.Height)
}

func (m *MintMonitorService) HandleInvalidMint(tx *ResultTx) bool {
	doc := models.InvalidMint{
		Height:          strconv.FormatInt(tx.Height, 10),
		TransactionHash: tx.Hash,
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.Pocket.ChainId,
		Amount:          tx.StdTx.Msg.Value.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
	}

	log.Debug("[MINT MONITOR] Storing invalid mint tx: ", tx.Hash, " in db")

	err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[MINT MONITOR] Found duplicate invalid mint tx: ", tx.Hash, " in db")
			return true
		}
		log.Error("[MINT MONITOR] Error storing invalid mint tx: ", err)
		return false
	}

	log.Debug("[MINT MONITOR] Stored invalid mint tx: ", tx.Hash, " in db")
	return true
}

func (m *MintMonitorService) HandleValidMint(tx *ResultTx, memo models.MintMemo) bool {
	doc := models.Mint{
		Height:           strconv.FormatInt(tx.Height, 10),
		TransactionHash:  tx.Hash,
		SenderAddress:    tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:    app.Config.Pocket.ChainId,
		RecipientAddress: memo.Address,
		RecipientChainId: strconv.FormatInt(int64(memo.ChainId), 10),
		Amount:           tx.StdTx.Msg.Value.Amount,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           models.StatusPending,
		Signers:          []string{},
	}

	log.Debug("[MINT MONITOR] Storing mint tx: ", tx.Hash, " in db")

	err := app.DB.InsertOne(models.CollectionMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("[MINT MONITOR] Found duplicate mint tx: ", tx.Hash, " in db")
			return true
		}
		log.Error("[MINT MONITOR] Error storing mint tx: ", err)
		return false
	}

	log.Debug("[MINT MONITOR] Stored mint tx: ", tx.Hash, " in db")
	return true
}

func (m *MintMonitorService) HandleTx(tx *ResultTx) bool {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId {
		log.Debug("[MINT MONITOR] Found invalid mint tx: ", tx.Hash, " with memo: ", "\""+tx.StdTx.Memo+"\"")
		return m.HandleInvalidMint(tx)
	}
	log.Debug("[MINT MONITOR] Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	return m.HandleValidMint(tx, memo)

}

func (m *MintMonitorService) SyncTxs() bool {
	txs, err := Client.GetAccountTxsByHeight(int64(m.startHeight))
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("[MINT MONITOR] Found ", len(txs), " txs to sync")
	var success bool = true
	for _, tx := range txs {
		success = success && m.HandleTx(tx)
	}
	return success
}

func NewMintMonitor() Service {
	log.Debug("[MINT MONITOR] Initializing mint monitor")
	m := &MintMonitorService{
		monitorInterval: time.Duration(app.Config.Pocket.MonitorIntervalSecs) * time.Second,
		startHeight:     0,
		currentHeight:   0,
		stop:            make(chan bool),
	}

	m.UpdateCurrentHeight()
	if app.Config.Pocket.StartHeight > 0 {
		m.startHeight = uint64(app.Config.Pocket.StartHeight)
	} else {
		log.Debug("[MINT MONITOR] Found invalid start height, updating current height")
		m.startHeight = m.currentHeight
	}

	log.Debug("[MINT MONITOR] Start height: ", m.startHeight)
	log.Debug("[MINT MONITOR] Initialized mint monitor")

	return m
}
