package pocket

import (
	"fmt"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	log "github.com/sirupsen/logrus"
)

type MintMonitor interface {
	Cancel()
	Start()
}

type WPOKTMintMonitor struct {
	stop            chan bool
	monitorInterval time.Duration
	lastHeight      int64
	currentHeight   int64
}

func (m *WPOKTMintMonitor) Cancel() {
	m.stop <- true
}

func (m *WPOKTMintMonitor) updateCurrentHeight() {
	res, err := GetHeight()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Updated current pokt height: ", res.Height)
	m.currentHeight = res.Height
}

func (m *WPOKTMintMonitor) syncTxs() {
	log.Debug("Syncing txs from height: ", m.lastHeight, " to height: ", m.currentHeight)
	txs, err := GetAccountTransferTxs(m.lastHeight)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("TotalTxs: %d\n", len(txs))
	fmt.Println("Txs:")
	for _, tx := range txs {
		fmt.Printf("[%d]\tHash: %s\n", tx.Index, tx.Hash)
		fmt.Printf("\tHeight: %d\n", tx.Height)
		fmt.Printf("\tType: %s\n", tx.StdTx.Msg.Type)
		fmt.Printf("\tFrom: %s\n", tx.StdTx.Msg.Value.FromAddress)
		fmt.Printf("\tTo: %s\n", tx.StdTx.Msg.Value.ToAddress)
		fmt.Printf("\tAmount: %s\n", tx.StdTx.Msg.Value.Amount)
		fmt.Printf("\tMemo: %s\n", tx.StdTx.Memo)
		fmt.Printf("\tFee: %s %s\n", tx.StdTx.Fee[0].Amount, tx.StdTx.Fee[0].Denom)
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
			log.Debug("Stopping mint monitor")
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
