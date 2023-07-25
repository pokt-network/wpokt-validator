package pocket

import (
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	ethereum "github.com/dan13ram/wpokt-backend/ethereum/client"
	"github.com/dan13ram/wpokt-backend/models"
	pocket "github.com/dan13ram/wpokt-backend/pocket/client"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	BurnSignerName = "burn-signer"
)

type BurnSignerService struct {
	wg             *sync.WaitGroup
	name           string
	stop           chan bool
	lastSyncTime   time.Time
	interval       time.Duration
	privateKey     crypto.PrivateKey
	multisigPubKey crypto.PublicKeyMultiSig
	numSigners     int
	ethClient      ethereum.EthereumClient
	poktClient     pocket.PocketClient
	poktHeight     int64
	ethBlockNumber int64
	vaultAddress   string
	wpoktAddress   string
}

func (m *BurnSignerService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.Name(),
		LastSyncTime:   m.LastSyncTime(),
		NextSyncTime:   m.LastSyncTime().Add(m.Interval()),
		PoktHeight:     m.PoktHeight(),
		EthBlockNumber: m.EthBlockNumber(),
		Healthy:        true,
	}
}

func (m *BurnSignerService) PoktHeight() string {
	return strconv.FormatInt(m.poktHeight, 10)
}

func (m *BurnSignerService) EthBlockNumber() string {
	return strconv.FormatInt(m.ethBlockNumber, 10)
}

func (m *BurnSignerService) LastSyncTime() time.Time {
	return m.lastSyncTime
}

func (m *BurnSignerService) Interval() time.Duration {
	return m.interval
}

func (m *BurnSignerService) Name() string {
	return m.name
}

func (m *BurnSignerService) Start() {
	log.Debug("[BURN SIGNER] Starting pokt signer")
	stop := false
	for !stop {
		log.Debug("[BURN SIGNER] Starting pokt signer sync")
		m.lastSyncTime = time.Now()

		m.UpdateBlocks()

		m.SyncTxs()

		log.Debug("[BURN SIGNER] Finished pokt signer sync")
		log.Debug("[BURN SIGNER] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[BURN SIGNER] Stopped pokt signer")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *BurnSignerService) UpdateBlocks() {
	log.Debug("[BURN SIGNER] Updating blocks")

	poktHeight, err := m.poktClient.GetHeight()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching pokt block height: ", err)
		return
	}
	m.poktHeight = poktHeight.Height

	ethBlockNumber, err := m.ethClient.GetBlockNumber()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching eth block number: ", err)
		return
	}
	m.ethBlockNumber = int64(ethBlockNumber)

	log.Debug("[BURN SIGNER] Updated blocks")

}

func (m *BurnSignerService) Stop() {
	log.Debug("[BURN SIGNER] Stopping pokt signer")
	m.stop <- true
}

func (m *BurnSignerService) HandleInvalidMint(doc models.InvalidMint) bool {
	log.Debug("[BURN SIGNER] Handling invalid mint: ", doc.TransactionHash)

	signers := doc.Signers
	returnTx := doc.ReturnTx

	if signers == nil {
		signers = []string{}
	}

	status := doc.Status
	confirmations, err := strconv.ParseInt(doc.Confirmations, 10, 64)
	if err != nil {
		confirmations = 0
	}

	if status == models.StatusPending {
		if app.Config.Pocket.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			log.Debug("[BURN SIGNER] Checking confirmations for invalid mint")
			mintHeight, err := strconv.ParseInt(doc.Height, 10, 64)
			if err != nil {
				log.Error("[BURN SIGNER] Error parsing height for invalid mint: ", err)
				return false
			}
			totalConfirmations := m.poktHeight - mintHeight
			if totalConfirmations >= app.Config.Pocket.Confirmations {
				status = models.StatusConfirmed
			}
			confirmations = totalConfirmations
		}
	}

	var update bson.M

	if status == models.StatusConfirmed {
		if returnTx == "" {
			log.Debug("[BURN SIGNER] Creating returnTx for invalid mint")
			amountWithFees, err := strconv.ParseInt(doc.Amount, 10, 64)
			if err != nil {
				log.Error("[BURN SIGNER] Error parsing amount for invalid mint: ", err)
				return false
			}
			amount := amountWithFees - app.Config.Pocket.Fees
			memo := doc.TransactionHash

			returnTxBytes, err := BuildMultiSigTxAndSign(
				doc.SenderAddress,
				memo,
				app.Config.Pocket.ChainId,
				amount,
				app.Config.Pocket.Fees,
				m.privateKey,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[BURN SIGNER] Error creating tx for invalid mint: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[BURN SIGNER] Created tx for invalid mint")

		} else {
			log.Debug("[BURN SIGNER] Signing tx for invalid mint")

			returnTxBytes, err := SignMultisigTx(
				returnTx,
				app.Config.Pocket.ChainId,
				m.privateKey,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[BURN SIGNER] Error signing tx for invalid mint: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[BURN SIGNER] Signed tx for invalid mint")

		}
		signers = append(signers, m.privateKey.PublicKey().RawString())
		if len(signers) == m.numSigners && status == models.StatusConfirmed {
			status = models.StatusSigned
			log.Debug("[BURN SIGNER] Invalid mint fully signed")
		}
		update = bson.M{
			"$set": bson.M{
				"return_tx":     returnTx,
				"signers":       signers,
				"status":        status,
				"confirmations": strconv.FormatInt(confirmations, 10),
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"status":        status,
				"confirmations": strconv.FormatInt(confirmations, 10),
			},
		}
	}

	filter := bson.M{
		"_id":           doc.Id,
		"vault_address": m.vaultAddress,
	}
	err = app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating invalid mint: ", err)
		return false
	}
	log.Debug("[BURN SIGNER] Updated invalid mint")
	return true
}

func (m *BurnSignerService) HandleBurn(doc models.Burn) bool {
	log.Debug("[BURN SIGNER] Handling burn: ", doc.TransactionHash)

	signers := doc.Signers
	returnTx := doc.ReturnTx

	if signers == nil {
		signers = []string{}
	}

	status := doc.Status
	confirmations, err := strconv.ParseInt(doc.Confirmations, 10, 64)
	if err != nil {
		confirmations = 0
	}

	if status == models.StatusPending {
		if app.Config.Ethereum.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			log.Debug("[BURN SIGNER] Checking confirmations for burn")
			burnBlockNumber, err := strconv.ParseInt(doc.BlockNumber, 10, 64)
			if err != nil {
				log.Error("[BURN SIGNER] Error parsing block number for burn: ", err)
				return false
			}
			totalConfirmations := m.ethBlockNumber - burnBlockNumber
			if totalConfirmations >= app.Config.Pocket.Confirmations {
				status = models.StatusConfirmed
			}
			confirmations = totalConfirmations
		}
	}

	var update bson.M

	if status == models.StatusConfirmed {
		if returnTx == "" {
			log.Debug("[BURN SIGNER] Creating returnTx for burn")
			amountWithFees, err := strconv.ParseInt(doc.Amount, 10, 64)
			if err != nil {
				log.Error("[BURN SIGNER] Error parsing amount for burn: ", err)
				return false
			}
			amount := amountWithFees - app.Config.Pocket.Fees
			memo := doc.TransactionHash

			returnTxBytes, err := BuildMultiSigTxAndSign(
				doc.RecipientAddress,
				memo,
				app.Config.Pocket.ChainId,
				amount,
				app.Config.Pocket.Fees,
				m.privateKey,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[BURN SIGNER] Error creating tx for burn: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[BURN SIGNER] Created tx for burn")

		} else {

			log.Debug("[BURN SIGNER] Signing tx for burn")

			returnTxBytes, err := SignMultisigTx(
				returnTx,
				app.Config.Pocket.ChainId,
				m.privateKey,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[BURN SIGNER] Error signing tx for burn: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[BURN SIGNER] Signed tx for burn")

		}

		signers = append(signers, m.privateKey.PublicKey().RawString())

		if len(signers) == m.numSigners {
			status = models.StatusSigned
			log.Debug("[BURN SIGNER] Invalid burn fully signed")
		}

		update = bson.M{
			"$set": bson.M{
				"return_tx":     returnTx,
				"signers":       signers,
				"status":        status,
				"confirmations": strconv.FormatInt(confirmations, 10),
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"status":        status,
				"confirmations": strconv.FormatInt(confirmations, 10),
			},
		}
	}

	filter := bson.M{
		"_id":           doc.Id,
		"wpokt_address": m.wpoktAddress,
	}
	err = app.DB.UpdateOne(models.CollectionBurns, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating burn: ", err)
		return false
	}
	log.Debug("[BURN SIGNER] Updated burn")

	return true
}

func (m *BurnSignerService) SyncTxs() bool {
	log.Debug("[BURN SIGNER] Syncing txs")
	filter := bson.M{
		"vault_address": m.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers":       bson.M{"$nin": []string{m.privateKey.PublicKey().RawString()}},
	}

	invalidMints := []models.InvalidMint{}
	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching invalid mints: ", err)
		return false
	}
	log.Debug("[BURN SIGNER] Found invalid mints: ", len(invalidMints))

	var success bool = true
	for _, doc := range invalidMints {
		success = m.HandleInvalidMint(doc) && success
	}

	filter = bson.M{
		"wpokt_address": m.wpoktAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers":       bson.M{"$nin": []string{m.privateKey.PublicKey().RawString()}},
	}

	burns := []models.Burn{}
	err = app.DB.FindMany(models.CollectionBurns, filter, &burns)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching burns: ", err)
		return false
	}
	log.Debug("[BURN SIGNER] Found burns: ", len(burns))
	for _, doc := range burns {
		success = m.HandleBurn(doc) && success
	}
	log.Debug("[BURN SIGNER] Synced txs")

	return success
}

func newSigner(wg *sync.WaitGroup) models.Service {
	log.Debug("[BURN SIGNER] Initializing pokt signer")

	pk, err := crypto.NewPrivateKey(app.Config.Pocket.PrivateKey)
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing pokt signer: ", err)
	}
	log.Debug("[BURN SIGNER] Initialized pokt signer private key")
	log.Debug("[BURN SIGNER] Pokt signer public key: ", pk.PublicKey().RawString())
	log.Debug("[BURN SIGNER] Pokt signer address: ", pk.PublicKey().Address().String())

	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Debug("[BURN SIGNER] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	log.Debug("[BURN SIGNER] Multisig address: ", multisigPk.Address().String())

	poktClient := pocket.NewClient()
	ethClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing ethereum client: ", err)
	}

	m := &BurnSignerService{
		wg:             wg,
		name:           BurnSignerName,
		interval:       time.Duration(app.Config.BurnSigner.IntervalSecs) * time.Second,
		stop:           make(chan bool),
		privateKey:     pk,
		multisigPubKey: multisigPk,
		numSigners:     len(pks),
		ethClient:      ethClient,
		poktClient:     poktClient,
		vaultAddress:   app.Config.Pocket.VaultAddress,
		wpoktAddress:   app.Config.Ethereum.WPOKTAddress,
	}

	m.UpdateBlocks()

	log.Debug("[BURN SIGNER] Initialized pokt signer")

	return m
}

func NewSigner(wg *sync.WaitGroup) models.Service {
	return newSigner(wg)
}

func NewSignerWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	return newSigner(wg)
}
