package pocket

import (
	"encoding/json"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type MintMonitor interface {
	Stop()
	Start()
}

type WPOKTMintMonitor struct {
	stop            chan bool
	monitorInterval time.Duration
	startHeight     uint64
	currentHeight   uint64
}

func (m *WPOKTMintMonitor) Stop() {
	log.Debug("Stopping mint monitor")
	m.stop <- true
}

func (m *WPOKTMintMonitor) UpdateCurrentHeight() {
	res, err := Client.GetHeight()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Updated current pokt height: ", res.Height)
	m.currentHeight = uint64(res.Height)
}

func (m *WPOKTMintMonitor) HandleInvalidMint(tx *ResultTx) bool {
	doc := models.InvalidMint{
		Height:          uint64(tx.Height),
		TransactionHash: tx.Hash,
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.Pocket.ChainId,
		Amount:          tx.StdTx.Msg.Value.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
	}

	log.Debug("Storing invalid mint tx: ", tx.Hash, " in db")

	err := app.DB.InsertOne(models.CollectionInvalidMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("Found duplicate invalid mint tx: ", tx.Hash, " in db")
			return true
		}
		log.Error("Error storing invalid mint tx: ", tx.Hash, " in db: ", err)
		return false
	}

	log.Debug("Stored invalid mint tx: ", tx.Hash, " in db")
	return true
}

func (m *WPOKTMintMonitor) HandleValidMint(tx *ResultTx, memo models.MintMemo) bool {
	doc := models.Mint{
		Height:           uint64(tx.Height),
		TransactionHash:  tx.Hash,
		SenderAddress:    tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:    app.Config.Pocket.ChainId,
		RecipientAddress: memo.Address,
		RecipientChainId: memo.ChainId,
		Amount:           tx.StdTx.Msg.Value.Amount,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           models.StatusPending,
		Signers:          []string{},
	}

	log.Debug("Storing mint tx in db: ", tx.Hash)

	err := app.DB.InsertOne(models.CollectionMints, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("Found duplicate mint tx in db: ", tx.Hash)
			return true
		}
		log.Error("Error storing mint tx in db: ", err)
		return false
	}

	log.Debug("Stored mint tx in db: ", tx.Hash)
	return true
}

func (m *WPOKTMintMonitor) HandleTx(tx *ResultTx) bool {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId {
		log.Debug("Found invalid memo in mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
		return m.HandleInvalidMint(tx)
	}
	log.Debug("Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	return m.HandleValidMint(tx, memo)

}

func (m *WPOKTMintMonitor) SyncTxs() bool {
	txs, err := Client.GetAccountTxsByHeight(int64(m.startHeight))
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("Found ", len(txs), " mint txs")
	var success bool = true
	for _, tx := range txs {
		success = success && m.HandleTx(tx)
	}
	return success
}

func (m *WPOKTMintMonitor) Start() {
	log.Debug("Starting mint monitor")
	stop := false
	for !stop {
		log.Debug("Starting mint sync")

		m.UpdateCurrentHeight()

		if (m.currentHeight - m.startHeight) > 0 {
			log.Debug("Syncing mint txs from height: ", m.startHeight, " to height: ", m.currentHeight)
			success := m.SyncTxs()
			if success {
				m.startHeight = m.currentHeight
			}
		} else {
			log.Debug("Already synced up to height: ", m.currentHeight)
		}

		log.Debug("Finished mint sync")
		log.Debug("Sleeping mint monitor for: ", m.monitorInterval)
		log.Debug("Next mint sync will start at: ", time.Now().Add(m.monitorInterval))

		select {
		case <-m.stop:
			stop = true
			log.Debug("Stopped mint monitor")
		case <-time.After(m.monitorInterval):
		}
	}
}

func NewMintMonitor() MintMonitor {
	m := &WPOKTMintMonitor{
		monitorInterval: time.Duration(app.Config.Pocket.MonitorIntervalSecs) * time.Second,
		startHeight:     0,
		currentHeight:   0,
		stop:            make(chan bool),
	}
	if app.Config.Pocket.StartHeight < 0 {
		m.UpdateCurrentHeight()
		m.startHeight = m.currentHeight
	} else {
		m.startHeight = uint64(app.Config.Pocket.StartHeight)
	}
	return m
}
