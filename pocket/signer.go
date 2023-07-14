package pocket

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/pokt-network/pocket-core/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type PoktSignerService struct {
	stop           chan bool
	interval       time.Duration
	privateKey     crypto.PrivateKey
	multisigPubKey crypto.PublicKeyMultiSig
	numSigners     int
}

func (m *PoktSignerService) Start() {
	log.Debug("[POKT SIGNER] Starting pokt signer")
	stop := false
	for !stop {
		log.Debug("[POKT SIGNER] Starting pokt signer sync")

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
}

func (m *PoktSignerService) Stop() {
	log.Debug("[POKT SIGNER] Stopping pokt signer")
	m.stop <- true
}

func (m *PoktSignerService) HandleInvalidMint(doc models.InvalidMint) bool {
	log.Debug("[POKT SIGNER] Handling invalid mint: ", doc.TransactionHash)

	signers := doc.Signers
	returnTx := doc.ReturnTx

	if signers == nil || len(signers) >= m.numSigners {
		signers = []string{}
	}

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

	status := models.StatusPending
	if len(signers) == m.numSigners {
		status = models.StatusSigned
		log.Debug("[POKT SIGNER] Invalid mint fully signed")
	}

	filter := bson.M{"_id": doc.Id}
	update := bson.M{
		"$set": bson.M{
			"return_tx": returnTx,
			"signers":   signers,
			"status":    status,
		},
	}
	err := app.DB.UpdateOne(models.CollectionInvalidMints, filter, update)
	if err != nil {
		log.Error("[POKT SIGNER] Error updating invalid mint: ", err)
		return false
	}
	log.Debug("[POKT SIGNER] Updated invalid mint")
	return true
}

func (m *PoktSignerService) HandleBurn(mint models.Burn) bool {
	return true
}

func (m *PoktSignerService) SyncTxs() bool {
	filter := bson.M{
		"status":  models.StatusPending,
		"signers": bson.M{"$nin": []string{m.privateKey.PublicKey().RawString()}},
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

	return success
}

func NewSigner() models.Service {
	if !app.Config.PoktSigner.Enabled {
		log.Debug("[POKT SIGNER] Pokt signer disabled")
		return models.NewEmptyService()
	}

	log.Debug("[POKT SIGNER] Initializing pokt signer")

	pk, err := crypto.NewPrivateKey(app.Config.PoktSigner.PrivateKey)
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

	m := &PoktSignerService{
		interval:       time.Duration(app.Config.PoktSigner.IntervalSecs) * time.Second,
		stop:           make(chan bool),
		privateKey:     pk,
		multisigPubKey: multisigPk,
		numSigners:     len(pks),
	}

	log.Debug("[POKT SIGNER] Initialized pokt signer")

	return m
}
