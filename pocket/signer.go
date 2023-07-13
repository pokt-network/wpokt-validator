package pocket

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type PoktSignerService struct {
	stop           chan bool
	interval       time.Duration
	privateKey     crypto.PrivateKey
	signerAddress  string
	multisigPubKey crypto.PublicKeyMultiSig
	kb             keys.Keybase
	numSigners     int
}

func (m *PoktSignerService) Start() {
	log.Debug("[POKT SIGNER] Starting pokt signer")
	stop := false
	for !stop {
		log.Debug("[POKT SIGNER] Starting pokt signer sync")

		m.SyncTxs()

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
	// create a mongo filter where "status" is "pending" and "signers" array does not contain this signer
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

	for _, result := range invalidMints {
		// log the result
		log.Debug("[POKT SIGNER] Reading invalid mint: ", result.TransactionHash)

		signers := result.Signers
		returnTx := result.ReturnTx

		if signers == nil || len(signers) >= m.numSigners {
			signers = []string{}
		}

		if returnTx == "" {
			log.Debug("[POKT SIGNER] Creating returnTx for invalid mint: ", result.TransactionHash)
			amountWithFees, err := strconv.ParseInt(result.Amount, 10, 64)
			if err != nil {
				log.Error("[POKT SIGNER] Error parsing amount for invalid mint: ", err)
				continue
			}
			amount := amountWithFees - app.Config.Pocket.Fees

			memo := fmt.Sprintf(`{"txHash":"%s","msg":"returning pokt for invalid memo"}`, result.TransactionHash)
			returnTxBytes, err := BuildMultiSigTxAndSign(
				m.signerAddress,
				result.SenderAddress,
				memo,
				app.Config.Pocket.ChainId,
				amount,
				app.Config.Pocket.Fees,
				m.kb,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[POKT SIGNER] Error creating tx for invalid mint: ", err)
				continue
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[POKT SIGNER] Created tx for invalid mint: ", returnTx)

		} else {
			for _, element := range signers {
				if element == m.privateKey.PublicKey().RawString() {
					log.Debug("[POKT SIGNER] Invalid mint already signed: ", result.TransactionHash)
					continue
				}
			}

			log.Debug("[POKT SIGNER] Signing tx for invalid mint: ", result.TransactionHash)

			returnTxBytes, err := SignMultisigTx(
				m.signerAddress,
				returnTx,
				app.Config.Pocket.ChainId,
				m.kb,
				m.multisigPubKey,
			)
			if err != nil {
				log.Error("[POKT SIGNER] Error signing tx for invalid mint: ", err)
				continue
			}
			returnTx = hex.EncodeToString(returnTxBytes)
			log.Debug("[POKT SIGNER] Signed tx for invalid mint: ", returnTx)

		}

		signers = append(signers, m.privateKey.PublicKey().RawString())
		if err != nil {
			log.Error("[POKT SIGNER] Error creating tx for invalid mint: ", err)
			continue
		}

		status := models.StatusPending
		if len(signers) == m.numSigners {
			status = models.StatusSigned
			log.Debug("[POKT SIGNER] Invalid mint fully signed: ", result.TransactionHash)
		}

		// update the invalid mint
		update := bson.M{
			"$set": bson.M{
				"return_tx": returnTx,
				"signers":   signers,
				"status":    status,
			},
		}
		err = app.DB.UpdateOne(models.CollectionInvalidMints, bson.M{"_id": result.Id}, update)
		if err != nil {
			log.Error("[POKT SIGNER] Error updating invalid mint: ", err)
			continue
		}
		log.Debug("[POKT SIGNER] Updated invalid mint: ", result.TransactionHash)
	}

	return true
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
	for _, pk := range app.Config.PoktSigner.MultisigPublicKeys {
		p, err := crypto.NewPublicKey(pk)
		if err != nil {
			fmt.Println(fmt.Errorf("error creating the public key: %v", err))
			continue
		}
		pks = append(pks, p)
	}

	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}
	log.Debug("[POKT SIGNER] Initialized pokt signer multisig address: ", multisigPk.Address().String())

	keybase := keys.NewInMemory()
	passphrase := "PASSPHRASE"
	armorred, err := mintkey.EncryptArmorPrivKey(pk, passphrase, passphrase)
	if err != nil {
		log.Fatal("[POKT SIGNER] Error encrypting private key: ", err)
	}
	_, err = keybase.ImportPrivKey(armorred, passphrase, passphrase)
	if err != nil {
		log.Fatal("[POKT SIGNER] Error importing private key: ", err)
	}

	m := &PoktSignerService{
		interval:       time.Duration(app.Config.PoktSigner.IntervalSecs) * time.Second,
		stop:           make(chan bool),
		privateKey:     pk,
		signerAddress:  pk.PublicKey().Address().String(),
		kb:             keybase,
		multisigPubKey: multisigPk,
		numSigners:     len(pks),
	}

	log.Debug("[POKT SIGNER] Initialized pokt signer")

	return m
}
