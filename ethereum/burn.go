package ethereum

// listen to events from the blockchain

import (
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	log "github.com/sirupsen/logrus"
)

// burn monitor interface

type BurnMonitor interface {
	Start()
	Stop()
}

type WPOKTBurnMonitor struct {
	stop               chan bool
	startBlockNumber   uint64
	currentBlockNumber uint64
	monitorInterval    time.Duration
}

func (m *WPOKTBurnMonitor) Stop() {
	log.Debug("Stopping burn monitor")
	m.stop <- true
}

func (m *WPOKTBurnMonitor) updateCurrentBlockNumber() {
	res, err := GetBlockNumber()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Updated current pokt blockNumber: ", res)
	m.currentBlockNumber = res
}

// func (m *WPOKTBurnMonitor) handleInvalidBurn(tx *ResultTx) {
// 	doc := models.InvalidBurn{
// 		BlockNumber:          tx.BlockNumber,
// 		TransactionHash: tx.Hash.String(),
// 		SenderAddress:   tx.StdTx.Msg.Value.FromAddress,
// 		SenderChainId:   app.Config.Pocket.ChainId,
// 		Amount:          tx.StdTx.Msg.Value.Amount,
// 		CreatedAt:       time.Now(),
// 		UpdatedAt:       time.Now(),
// 		Status:          models.StatusPending,
// 		Signers:         []string{},
// 	}

// 	log.Debug("Storing invalid burn tx: ", tx.Hash, " in db")

// 	col := app.DB.GetCollection(models.CollectionInvalidBurns)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.Pocket.MonitorIntervalSecs))
// 	defer cancel()

// 	_, err := col.InsertOne(ctx, doc)
// 	if err != nil {
// 		if mongo.IsDuplicateKeyError(err) {
// 			log.Debug("Found duplicate invalid burn tx: ", tx.Hash, " in db")
// 			return
// 		}
// 		log.Error("Error storing invalid burn tx: ", tx.Hash, " in db: ", err)
// 		return
// 	}

// 	log.Debug("Stored invalid burn tx: ", tx.Hash, " in db")
// }

// func (m *WPOKTBurnMonitor) handleValidBurn(tx *ResultTx, memo models.BurnMemo) {
// 	doc := models.Burn{
// 		BlockNumber:           tx.BlockNumber,
// 		TransactionHash:  tx.Hash.String(),
// 		SenderAddress:    tx.StdTx.Msg.Value.FromAddress,
// 		SenderChainId:    app.Config.Pocket.ChainId,
// 		RecipientAddress: memo.Address,
// 		RecipientChainId: memo.ChainId,
// 		Amount:           tx.StdTx.Msg.Value.Amount,
// 		CreatedAt:        time.Now(),
// 		UpdatedAt:        time.Now(),
// 		Status:           models.StatusPending,
// 		Signers:          []string{},
// 	}

// 	log.Debug("Storing burn tx: ", tx.Hash, " in db")

// 	col := app.DB.GetCollection(models.CollectionBurns)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.Config.Pocket.MonitorIntervalSecs))
// 	defer cancel()

// 	_, err := col.InsertOne(ctx, doc)
// 	if err != nil {
// 		log.Error("Error storing burn tx: ", tx.Hash, " in db: ", err)
// 		return
// 	}

// 	log.Debug("Stored burn tx: ", tx.Hash, " in db")
// }

// func (m *WPOKTBurnMonitor) handleTx(tx *ResultTx) {
// 	var memo models.BurnMemo

// 	err := json.Unmarshal([]byte(tx.StdTx.Memo), &memo)

// 	if err != nil || memo.ChainId != app.Config.Ethereum.ChainId {
// 		log.Debug("Found invalid memo in tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
// 		m.handleInvalidBurn(tx)
// 		return
// 	}
// 	log.Info("Found valid burn tx: ", tx.Hash, " with memo: ", tx.StdTx.Memo)
// 	m.handleValidBurn(tx, memo)

// }

func (m *WPOKTBurnMonitor) syncTxs() {
	log.Debug("Syncing burn txs from blockNumber: ", m.startBlockNumber, " to blockNumber: ", m.currentBlockNumber)
	// events, err := GetEvents(m.startBlockNumber, m.currentBlockNumber)
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Debug("Found ", len(txs), " txs")
	// for _, tx := range txs {
	// 	m.handleTx(tx)
	// }
}

func (m *WPOKTBurnMonitor) Start() {
	log.Debug("Starting burn monitor")
	stop := false
	for !stop {
		log.Debug("Starting burn sync")

		m.updateCurrentBlockNumber()

		if m.startBlockNumber == 0 {
			m.startBlockNumber = m.currentBlockNumber
		}

		if (m.currentBlockNumber - m.startBlockNumber) > 0 {
			log.Debug("Syncing burn txs from blockNumber: ", m.startBlockNumber, " to blockNumber: ", m.currentBlockNumber)
			m.syncTxs()
			m.startBlockNumber = m.currentBlockNumber
		}

		log.Debug("Finished burn sync")
		log.Debug("Sleeping for ", m.monitorInterval)
		log.Debug("Next burn sync at: ", time.Now().Add(m.monitorInterval))

		select {
		case <-m.stop:
			stop = true
			log.Debug("Stopped burn monitor")
		case <-time.After(m.monitorInterval):
		}
	}
}

func NewBurnMonitor() BurnMonitor {
	return &WPOKTBurnMonitor{
		stop:               make(chan bool),
		startBlockNumber:   app.Config.Ethereum.StartBlockNumber,
		currentBlockNumber: 0,
		monitorInterval:    time.Duration(app.Config.Pocket.MonitorIntervalSecs) * time.Second,
	}
}
