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
	lastHeight      int64
	currentHeight   int64
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
	log.Info("Updated current pokt height: ", res.Height)
	m.currentHeight = res.Height
}

func (m *WPOKTMintMonitor) handleInvalidMint(tx *ResultTx) {
	doc := models.InvalidMint{
		Height:          tx.Height,
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
		Height:           tx.Height,
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
		log.Debug("Found invalid memo in tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
		m.handleInvalidMint(tx)
		return
	}
	log.Info("Found valid mint tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
	m.handleValidMint(tx, memo)

}

func (m *WPOKTMintMonitor) syncTxs() {
	log.Debug("Syncing txs from height: ", m.lastHeight, " to height: ", m.currentHeight)
	txs, err := GetAccountTransferTxs(m.lastHeight)
	if err != nil {
		log.Error(err)
	}
	log.Debug("Found ", len(txs), " txs")
	for _, tx := range txs {
		m.handleTx(tx)
	}
}

func (m *WPOKTMintMonitor) Start() {
	log.Debug("Starting mint monitor")
	stop := false
	for !stop {
		// Start
		m.updateCurrentHeight()
		if (m.currentHeight - m.lastHeight) > 0 {
			// Search for mint txs & store in db
			m.syncTxs()
			m.lastHeight = m.currentHeight
		}
		// End
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
		lastHeight:      app.Config.Pocket.StartHeight,
		currentHeight:   0,
		stop:            make(chan bool),
	}
}
