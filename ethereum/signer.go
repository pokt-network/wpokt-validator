package ethereum

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	ethereum "github.com/dan13ram/wpokt-backend/ethereum/client"
	"github.com/dan13ram/wpokt-backend/models"
	pocket "github.com/dan13ram/wpokt-backend/pocket/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MintSignerName = "mint-signer"
)

type MintSignerService struct {
	wg                     *sync.WaitGroup
	name                   string
	stop                   chan bool
	address                string
	privateKey             *ecdsa.PrivateKey
	lastSyncTime           time.Time
	interval               time.Duration
	vaultAddress           string
	wpoktAddress           string
	wpoktContract          *autogen.WrappedPocket
	mintControllerContract *autogen.MintController
	numSigners             int
	domain                 DomainData
	poktClient             pocket.PocketClient
	ethClient              ethereum.EthereumClient
	poktHeight             int64
}

func (m *MintSignerService) Start() {
	log.Debug("[MINT SIGNER] Starting wpokt signer")
	stop := false
	for !stop {
		log.Debug("[MINT SIGNER] Starting wpokt signer sync")
		m.lastSyncTime = time.Now()

		m.UpdateBlocks()
		m.SyncTxs()

		log.Debug("[MINT SIGNER] Finished wpokt signer sync")
		log.Debug("[MINT SIGNER] Sleeping for ", m.interval)

		select {
		case <-m.stop:
			stop = true
			log.Debug("[MINT SIGNER] Stopped wpokt signer")
		case <-time.After(m.interval):
		}
	}
	m.wg.Done()
}

func (m *MintSignerService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           m.name,
		LastSyncTime:   m.lastSyncTime,
		NextSyncTime:   m.lastSyncTime.Add(m.interval),
		PoktHeight:     strconv.FormatInt(m.poktHeight, 10),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (m *MintSignerService) Stop() {
	log.Debug("[MINT SIGNER] Stopping wpokt signer")
	m.stop <- true
}

func (m *MintSignerService) UpdateBlocks() {
	log.Debug("[MINT SIGNER] Updating blocks")
	poktHeight, err := m.poktClient.GetHeight()
	if err != nil {
		log.Error("[MINT SIGNER] Error fetching pokt block height: ", err)
		return
	}
	m.poktHeight = poktHeight.Height
}

func (m *MintSignerService) FindNonce(mint models.Mint) (*big.Int, error) {
	log.Debug("[MINT SIGNER] Finding nonce for mint: ", mint.TransactionHash)
	var nonce *big.Int

	if mint.Nonce != "" {
		mintNonce, ok := new(big.Int).SetString(mint.Nonce, 10)
		if !ok {
			log.Error("[MINT SIGNER] Error converting decimal to big int")
			return nil, errors.New("error converting decimal to big int")
		}
		nonce = mintNonce
	}

	if nonce == nil || nonce.Cmp(big.NewInt(0)) == 0 {
		log.Debug("[MINT SIGNER] Mint nonce not set, fetching from contract")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
		defer cancel()
		opts := &bind.CallOpts{Context: ctx, Pending: false}
		currentNonce, err := m.wpoktContract.GetUserNonce(opts, common.HexToAddress(mint.RecipientAddress))

		if err != nil {
			log.Error("[MINT SIGNER] Error fetching nonce from contract: ", err)
			return nil, err
		}

		var pendingMints []models.Mint
		filter := bson.M{
			"_id":               bson.M{"$ne": mint.Id},
			"vault_address":     m.vaultAddress,
			"wpokt_address":     m.wpoktAddress,
			"recipient_address": mint.RecipientAddress,
			"status":            bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed, models.StatusSigned}},
		}
		err = app.DB.FindMany(models.CollectionMints, filter, &pendingMints)
		if err != nil {
			log.Error("[MINT SIGNER] Error fetching pending mints: ", err)
			return nil, err
		}

		if len(pendingMints) > 0 {
			var nonces []int64

			for _, pendingMint := range pendingMints {
				if pendingMint.Data != nil {
					nonce, err := strconv.ParseInt(pendingMint.Data.Nonce, 10, 64)
					if err != nil {
						log.Error("[MINT SIGNER] Error converting nonce to int: ", err)
						continue
					}
					nonces = append(nonces, nonce)
				}
			}

			if len(nonces) > 0 {
				sort.Slice(nonces, func(i, j int) bool {
					return nonces[i] < nonces[j]
				})

				newNonce := big.NewInt(nonces[len(nonces)-1])
				if currentNonce.Cmp(newNonce) < 0 {
					currentNonce = newNonce
				}
			}
		}

		nonce = currentNonce.Add(currentNonce, big.NewInt(1))
	}
	return nonce, nil
}

func (m *MintSignerService) HandleMint(mint models.Mint) bool {
	log.Debug("[MINT SIGNER] Handling mint: ", mint.TransactionHash)

	address := common.HexToAddress(mint.RecipientAddress)
	amount, ok := new(big.Int).SetString(mint.Amount, 10)
	if !ok {
		log.Error("[MINT SIGNER] Error converting decimal to big int")
		return false
	}

	nonce, err := m.FindNonce(mint)

	if err != nil {
		log.Error("[MINT SIGNER] Error fetching nonce: ", err)
		return false
	}

	if nonce == nil || nonce.Cmp(big.NewInt(0)) == 0 {
		log.Error("[MINT SIGNER] Error fetching nonce")
		return false
	}

	data := autogen.MintControllerMintData{
		Recipient: address,
		Amount:    amount,
		Nonce:     nonce,
	}

	mint, err = updateStatusAndConfirmationsForMint(mint, m.poktHeight)
	if err != nil {
		log.Error("[MINT SIGNER] Error updating status and confirmations for mint: ", err)
		return false
	}

	var update bson.M
	if mint.Status == models.StatusConfirmed {
		log.Debug("[MINT SIGNER] Mint confirmed, signing")

		mint, err := signMint(mint, data, m.domain, m.privateKey, m.numSigners)
		if err != nil {
			log.Error("[MINT SIGNER] Error signing mint: ", err)
			return false
		}

		update = bson.M{
			"$set": bson.M{
				"data": models.MintData{
					Recipient: data.Recipient.Hex(),
					Amount:    data.Amount.String(),
					Nonce:     data.Nonce.String(),
				},
				"nonce":         data.Nonce.String(),
				"signatures":    mint.Signatures,
				"signers":       mint.Signers,
				"status":        mint.Status,
				"confirmations": mint.Confirmations,
			},
		}

	} else {
		log.Debug("[MINT SIGNER] Mint pending confirmation, not signing")
		update = bson.M{
			"$set": bson.M{
				"status":        mint.Status,
				"confirmations": mint.Confirmations,
			},
		}
	}

	filter := bson.M{
		"_id":           mint.Id,
		"wpokt_address": m.wpoktAddress,
		"vault_address": m.vaultAddress,
	}

	err = app.DB.UpdateOne(models.CollectionMints, filter, update)
	if err != nil {
		log.Error("[MINT SIGNER] Error updating mint: ", err)
		return false
	}
	log.Debug("[MINT SIGNER] Mint updated with signature")

	return true
}

func (m *MintSignerService) SyncTxs() bool {
	log.Debug("[MINT SIGNER] Syncing pending txs")

	filter := bson.M{
		"wpokt_address": m.wpoktAddress,
		"vault_address": m.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers": bson.M{
			"$nin": []string{m.address},
		},
	}

	var results []models.Mint

	err := app.DB.FindMany(models.CollectionMints, filter, &results)
	if err != nil {
		log.Error("[MINT SIGNER] Error fetching pending mints: ", err)
		return false
	}

	var success bool = true
	for _, mint := range results {
		success = m.HandleMint(mint) && success

	}

	log.Debug("[MINT SIGNER] Finished syncing pending txs")
	return success
}

func newSigner(wg *sync.WaitGroup) models.Service {
	if app.Config.MintSigner.Enabled == false {
		log.Debug("[MINT SIGNER] BURN signer disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[MINT SIGNER] Initializing wpokt signer")

	privateKey, err := crypto.HexToECDSA(app.Config.Ethereum.PrivateKey)
	if err != nil {
		log.Fatal("[MINT SIGNER] Error loading private key: ", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	log.Debug("[MINT SIGNER] Loaded private key for address: ", address)

	ethClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[MINT SIGNER] Error initializing ethereum client: ", err)
	}

	log.Debug("[MINT SIGNER] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTAddress), ethClient.GetClient())
	if err != nil {
		log.Fatal("[MINT SIGNER] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[MINT SIGNER] Connected to wpokt contract")

	log.Debug("[MINT SIGNER] Connecting to mint controller contract at: ", app.Config.Ethereum.MintControllerAddress)
	mintControllerContract, err := autogen.NewMintController(common.HexToAddress(app.Config.Ethereum.MintControllerAddress), ethClient.GetClient())
	if err != nil {
		log.Fatal("[MINT SIGNER] Error initializing Mint Controller contract", err)
	}
	log.Debug("[MINT SIGNER] Connected to mint controller contract")

	log.Debug("[MINT SIGNER] Fetching mint controller domain data")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx, Pending: false}
	domain, err := mintControllerContract.Eip712Domain(opts)

	if err != nil {
		log.Fatal("[MINT SIGNER] Error fetching mint controller domain data: ", err)
	}
	log.Debug("[MINT SIGNER] Fetched mint controller domain data")

	b := &MintSignerService{
		wg:                     wg,
		name:                   MintSignerName,
		stop:                   make(chan bool),
		interval:               time.Duration(app.Config.MintSigner.IntervalSecs) * time.Second,
		privateKey:             privateKey,
		address:                address,
		wpoktAddress:           app.Config.Ethereum.WPOKTAddress,
		vaultAddress:           app.Config.Pocket.VaultAddress,
		wpoktContract:          contract,
		mintControllerContract: mintControllerContract,
		numSigners:             len(app.Config.Ethereum.ValidatorAddresses),
		domain:                 domain,
		ethClient:              ethClient,
		poktClient:             pocket.NewClient(),
	}

	log.Debug("[MINT SIGNER] Initialized wpokt signer")

	return b
}
func NewSigner(wg *sync.WaitGroup) models.Service {
	return newSigner(wg)
}

func NewSignerWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	return newSigner(wg)
}
