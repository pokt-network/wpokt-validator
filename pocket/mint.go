package pocket

import (
	"context"
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

func (m *WPOKTMintMonitor) updateCurrentHeight() {
	res, err := GetHeight()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Updated current pokt height: ", res.Height)
	m.currentHeight = uint64(res.Height)
}

func (m *WPOKTMintMonitor) handleInvalidMint(tx *ResultTx) {
	doc := models.InvalidMint{
		Height:          uint64(tx.Height),
		TransactionHash: tx.Hash.String(),
		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
		SenderChainId:   app.Config.Pocket.ChainId,
		Amount:          tx.StdTx.Msg.Value.Amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          models.StatusPending,
		Signers:         []string{},
	}

	log.Debug("Storing invalid mint tx: ", tx.Hash, " in db")

	col := app.DB.GetCollection(models.CollectionInvalidMints)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.Pocket.MonitorIntervalSecs))
	defer cancel()

	_, err := col.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Debug("Found duplicate invalid mint tx: ", tx.Hash, " in db")
			return
		}
		log.Error("Error storing invalid mint tx: ", tx.Hash, " in db: ", err)
		return
	}

	log.Debug("Stored invalid mint tx: ", tx.Hash, " in db")
}

func (m *WPOKTMintMonitor) handleValidMint(tx *ResultTx, memo models.MintMemo) {
	doc := models.Mint{
		Height:           uint64(tx.Height),
		TransactionHash:  tx.Hash.String(),
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

	log.Debug("Storing mint tx: ", tx.Hash, " in db")

	col := app.DB.GetCollection(models.CollectionMints)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.Pocket.MonitorIntervalSecs))
	defer cancel()

	_, err := col.InsertOne(ctx, doc)
	if err != nil {
		log.Error("Error storing mint tx: ", tx.Hash, " in db: ", err)
		return
	}

	log.Debug("Stored mint tx: ", tx.Hash, " in db")
}

func (m *WPOKTMintMonitor) handleTx(tx *ResultTx) {
	var memo models.MintMemo

	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId {
		log.Debug("Found invalid memo in mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
		m.handleInvalidMint(tx)
		return
	}
	log.Debug("Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	m.handleValidMint(tx, memo)

}

func (m *WPOKTMintMonitor) syncTxs() {
	txs, err := GetAccountTransferTxs(int64(m.startHeight))
	if err != nil {
		log.Error(err)
	}
	log.Debug("Found ", len(txs), " mint txs")
	for _, tx := range txs {
		m.handleTx(tx)
	}
}

func (m *WPOKTMintMonitor) Start() {
	log.Debug("Starting mint monitor")
	stop := false
	for !stop {
		log.Debug("Starting mint sync")

		m.updateCurrentHeight()

		if m.startHeight == 0 {
			m.startHeight = m.currentHeight
		}

		if (m.currentHeight - m.startHeight) > 0 {
			log.Debug("Syncing mint txs from height: ", m.startHeight, " to height: ", m.currentHeight)
			m.syncTxs()
			m.startHeight = m.currentHeight
		} else {
			log.Debug("Already synced up to height: ", m.currentHeight)
		}

		log.Debug("Finished mint sync")
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
	return &WPOKTMintMonitor{
		monitorInterval: time.Duration(app.Config.Pocket.MonitorIntervalSecs) * time.Second,
		startHeight:     app.Config.Pocket.StartHeight,
		currentHeight:   0,
		stop:            make(chan bool),
	}
}
