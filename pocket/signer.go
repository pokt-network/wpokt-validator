package pocket

import (
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	log "github.com/sirupsen/logrus"
)

type PoktSignerService struct {
	stop     chan bool
	interval time.Duration
}

func (m *PoktSignerService) Start() {
	log.Debug("[POKT SIGNER] Starting pokt signer")
	stop := false
	for !stop {
		log.Debug("[POKT SIGNER] Starting pokt signer sync")

		// find all invalid mints and burns from the db

		log.Debug("[POKT SIGNER] Finished mint sync")
		log.Debug("[POKT SIGNER] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[POKT SIGNER] Stopped pokt signer")
		case <-time.After(m.interval):
		}
	}
}

func (m *PoktSignerService) Stop() {
	log.Debug("[POKT SIGNER] Stopping pokt signer")
	m.stop <- true
}

func (m *PoktSignerService) HandleInvalidMint(doc models.InvalidMint) bool {
	return true
}

func (m *PoktSignerService) HandleBurn(mint models.Burn) bool {
	return true
}

func (m *PoktSignerService) SyncTxs() bool {

	return true
}

func NewPoktSigner() Service {
	log.Debug("[POKT SIGNER] Initializing pokt signer")
	m := &PoktSignerService{
		interval: time.Duration(app.Config.PoktSigner.IntervalSecs) * time.Second,
		stop:     make(chan bool),
	}

	log.Debug("[POKT SIGNER] Initialized pokt signer")

	return m
}
