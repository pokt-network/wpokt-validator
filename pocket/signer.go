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

type PoktSignerService struct {
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

func (m *PoktSignerService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.Name(),
		LastSyncTime:   m.LastSyncTime(),
		NextSyncTime:   m.LastSyncTime().Add(m.Interval()),
		PoktHeight:     m.PoktHeight(),
		EthBlockNumber: m.EthBlockNumber(),
		Healthy:        true,
	}
}

func (m *PoktSignerService) PoktHeight() string {
	return strconv.FormatInt(m.poktHeight, 10)
}

func (m *PoktSignerService) EthBlockNumber() string {
	return strconv.FormatInt(m.ethBlockNumber, 10)
}

func (m *PoktSignerService) LastSyncTime() time.Time {
	return m.lastSyncTime
}

func (m *PoktSignerService) Interval() time.Duration {
	return m.interval
}

func (m *PoktSignerService) Name() string {
	return m.name
}

func (m *PoktSignerService) Start() {
	log.Debug("[POKT SIGNER] Starting pokt signer")
	stop := false
	for !stop {
		log.Debug("[POKT SIGNER] Starting pokt signer sync")
		m.lastSyncTime = time.Now()

		m.UpdateBlocks()

		m.SyncTxs()

		log.Debug("[POKT SIGNER] Finished pokt signer sync")
		log.Debug("[POKT SIGNER] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[POKT SIGNER] Stopped pokt signer")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *PoktSignerService) UpdateBlocks() {
	log.Debug("[POKT SIGNER] Updating blocks")

	poktHeight, err := m.poktClient.GetHeight()
	if err != nil {
		log.Error("[POKT SIGNER] Error fetching pokt block height: ", err)
		return
	}
	m.poktHeight = poktHeight.Height

	ethBlockNumber, err := m.ethClient.GetBlockNumber()
	if err != nil {
		log.Error("[POKT SIGNER] Error fetching eth block number: ", err)
		return
	}
	m.ethBlockNumber = int64(ethBlockNumber)

	log.Debug("[POKT SIGNER] Updated blocks")

}

func (m *PoktSignerService) Stop() {
	log.Debug("[POKT SIGNER] Stopping pokt signer")
	m.stop <- true
}

func (m *PoktSignerService) HandleInvalidMint(doc models.InvalidMint) bool {
	log.Debug("[POKT SIGNER] Handling invalid mint: ", doc.TransactionHash)

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
			log.Debug("[POKT SIGNER] Checking confirmations for invalid mint")
			mintHeight, err := strconv.ParseInt(doc.Height, 10, 64)
			if err != nil {
				log.Error("[POKT SIGNER] Error parsing height for invalid mint: ", err)
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
			log.Debug("[POKT SIGNER] Creating returnTx for invalid mint")
			amountWithFees, err := strconv.ParseInt(doc.Amount, 10, 64)
			if err != nil {
				log.Error("[POKT SIGNER] Error parsing amount for invalid mint: ", err)
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
				log.Error("[POKT SIGNER] Error creating tx for invalid mint: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[POKT SIGNER] Created tx for invalid mint")

		} else {
			log.Debug("[POKT SIGNER] Signing tx for invalid mint")

			returnTxBytes, err := SignMultisigTx(
				returnTx,
				app.Config.Pocket.ChainId,
				m.privateKey,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[POKT SIGNER] Error signing tx for invalid mint: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[POKT SIGNER] Signed tx for invalid mint")

		}
		signers = append(signers, m.privateKey.PublicKey().RawString())
		if len(signers) == m.numSigners && status == models.StatusConfirmed {
			status = models.StatusSigned
			log.Debug("[POKT SIGNER] Invalid mint fully signed")
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
		log.Error("[POKT SIGNER] Error updating invalid mint: ", err)
		return false
	}
	log.Debug("[POKT SIGNER] Updated invalid mint")
	return true
}

func (m *PoktSignerService) HandleBurn(doc models.Burn) bool {
	log.Debug("[POKT SIGNER] Handling burn: ", doc.TransactionHash)

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
			log.Debug("[POKT SIGNER] Checking confirmations for burn")
			burnBlockNumber, err := strconv.ParseInt(doc.BlockNumber, 10, 64)
			if err != nil {
				log.Error("[POKT SIGNER] Error parsing block number for burn: ", err)
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
			log.Debug("[POKT SIGNER] Creating returnTx for burn")
			amountWithFees, err := strconv.ParseInt(doc.Amount, 10, 64)
			if err != nil {
				log.Error("[POKT SIGNER] Error parsing amount for burn: ", err)
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
				log.Error("[POKT SIGNER] Error creating tx for burn: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[POKT SIGNER] Created tx for burn")

		} else {

			log.Debug("[POKT SIGNER] Signing tx for burn")

			returnTxBytes, err := SignMultisigTx(
				returnTx,
				app.Config.Pocket.ChainId,
				m.privateKey,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[POKT SIGNER] Error signing tx for burn: ", err)
				return false
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[POKT SIGNER] Signed tx for burn")

		}

		signers = append(signers, m.privateKey.PublicKey().RawString())

		if len(signers) == m.numSigners {
			status = models.StatusSigned
			log.Debug("[POKT SIGNER] Invalid burn fully signed")
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
		log.Error("[POKT SIGNER] Error updating burn: ", err)
		return false
	}
	log.Debug("[POKT SIGNER] Updated burn")

	return true
}

func (m *PoktSignerService) SyncTxs() bool {
	log.Debug("[POKT SIGNER] Syncing txs")
	filter := bson.M{
		"vault_address": m.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers":       bson.M{"$nin": []string{m.privateKey.PublicKey().RawString()}},
	}

	invalidMints := []models.InvalidMint{}
	err := app.DB.FindMany(models.CollectionInvalidMints, filter, &invalidMints)
	if err != nil {
		log.Error("[POKT SIGNER] Error fetching invalid mints: ", err)
		return false
	}
	log.Debug("[POKT SIGNER] Found invalid mints: ", len(invalidMints))

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
		log.Error("[POKT SIGNER] Error fetching burns: ", err)
		return false
	}
	log.Debug("[POKT SIGNER] Found burns: ", len(burns))
	for _, doc := range burns {
		success = m.HandleBurn(doc) && success
	}
	log.Debug("[POKT SIGNER] Synced txs")

	return success
}

func NewSigner(wg *sync.WaitGroup) models.Service {
	if !app.Config.PoktSigner.Enabled {
		log.Debug("[POKT SIGNER] Pokt signer disabled")
		return models.NewEmptyService(wg, "empty-pokt-signer")
	}

	log.Debug("[POKT SIGNER] Initializing pokt signer")

	pk, err := crypto.NewPrivateKey(app.Config.Pocket.PrivateKey)
	if err != nil {
		log.Fatal("[POKT SIGNER] Error initializing pokt signer: ", err)
	}
	log.Debug("[POKT SIGNER] Initialized pokt signer private key")
	log.Debug("[POKT SIGNER] Pokt signer public key: ", pk.PublicKey().RawString())
	log.Debug("[POKT SIGNER] Pokt signer address: ", pk.PublicKey().Address().String())

	var pks []crypto.PublicKey
	for _, pk := range app.Config.Pocket.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			log.Debug("[POKT SIGNER] Error parsing multisig public key: ", err)
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	log.Debug("[POKT SIGNER] Multisig address: ", multisigPk.Address().String())

	poktClient := pocket.NewClient()
	ethClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[POKT SIGNER] Error initializing ethereum client: ", err)
	}

	m := &PoktSignerService{
		wg:             wg,
		name:           "pokt-signer",
		interval:       time.Duration(app.Config.PoktSigner.IntervalSecs) * time.Second,
		stop:           make(chan bool),
		privateKey:     pk,
		multisigPubKey: multisigPk,
		numSigners:     len(pks),
		ethClient:      ethClient,
		poktClient:     poktClient,
		vaultAddress:   app.Config.Pocket.VaultAddress,
		wpoktAddress:   app.Config.Ethereum.WPOKTAddress,
	}

	log.Debug("[POKT SIGNER] Initialized pokt signer")

	return m
}
