package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/eth/util"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MintSignerName = "mint signer"
)

type MintSignerService struct {
	wg                     *sync.WaitGroup
	name                   string
	stop                   chan bool
	address                string
	privateKey             *ecdsa.PrivateKey
	interval               time.Duration
	vaultAddress           string
	wpoktAddress           string
	wpoktContract          *autogen.WrappedPocket
	mintControllerContract *autogen.MintController
	numSigners             int
	domain                 util.DomainData
	poktClient             pokt.PocketClient
	ethClient              eth.EthereumClient
	poktHeight             int64

	healthMu sync.RWMutex
	health   models.ServiceHealth
}

func (x *MintSignerService) Start() {
	log.Info("[MINT SIGNER] Starting service")
	stop := false
	for !stop {
		log.Info("[MINT SIGNER] Starting sync")

		x.UpdateBlocks()

		x.SyncTxs()

		x.UpdateHealth()

		log.Info("[MINT SIGNER] Finished sync, Sleeping for ", x.interval)

		select {
		case <-x.stop:
			stop = true
			log.Info("[MINT SIGNER] Stopped service")
		case <-time.After(x.interval):
		}
	}
	x.wg.Done()
}

func (x *MintSignerService) Health() models.ServiceHealth {
	x.healthMu.RLock()
	defer x.healthMu.RUnlock()

	return x.health
}

func (x *MintSignerService) UpdateHealth() {
	x.healthMu.Lock()
	defer x.healthMu.Unlock()

	lastSyncTime := time.Now()

	x.health = models.ServiceHealth{
		Name:           MintSignerName,
		LastSyncTime:   lastSyncTime,
		NextSyncTime:   lastSyncTime.Add(x.interval),
		PoktHeight:     strconv.FormatInt(x.poktHeight, 10),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (x *MintSignerService) Stop() {
	log.Debug("[MINT SIGNER] Stopping service")
	x.stop <- true
}

func (x *MintSignerService) UpdateBlocks() {
	log.Debug("[MINT SIGNER] Updating blocks")
	poktHeight, err := x.poktClient.GetHeight()
	if err != nil {
		log.Error("[MINT SIGNER] Error fetching pokt block height: ", err)
		return
	}
	x.poktHeight = poktHeight.Height
}

func (x *MintSignerService) FindNonce(mint models.Mint) (*big.Int, error) {
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
		log.Info("[MINT SIGNER] Mint nonce not set, fetching from contract")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutSecs)*time.Second)
		defer cancel()
		opts := &bind.CallOpts{Context: ctx, Pending: false}
		currentNonce, err := x.wpoktContract.GetUserNonce(opts, common.HexToAddress(mint.RecipientAddress))
		if err != nil {
			log.Error("[MINT SIGNER] Error fetching nonce from contract: ", err)
			return nil, err
		}
		log.Debug("[MINT SIGNER] Current nonce: ", currentNonce, " for address: ", mint.RecipientAddress)

		var pendingMints []models.Mint
		filter := bson.M{
			"_id":               bson.M{"$ne": mint.Id},
			"vault_address":     x.vaultAddress,
			"wpokt_address":     x.wpoktAddress,
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

				pendingNonce := big.NewInt(nonces[len(nonces)-1])
				if currentNonce.Cmp(pendingNonce) < 0 {
					log.Debug("[MINT SIGNER] Pending nonce: ", pendingNonce)
					currentNonce = pendingNonce
				}
			}
		}

		nonce = currentNonce.Add(currentNonce, big.NewInt(1))
	}
	return nonce, nil
}

func (x *MintSignerService) HandleMint(mint models.Mint) bool {
	log.Debug("[MINT SIGNER] Handling mint: ", mint.TransactionHash)

	address := common.HexToAddress(mint.RecipientAddress)
	amount, ok := new(big.Int).SetString(mint.Amount, 10)
	if !ok {
		log.Error("[MINT SIGNER] Error converting decimal to big int")
		return false
	}

	nonce, err := x.FindNonce(mint)

	if err != nil {
		log.Error("[MINT SIGNER] Error fetching nonce: ", err)
		return false
	}

	if nonce == nil || nonce.Cmp(big.NewInt(0)) == 0 {
		log.Error("[MINT SIGNER] Error fetching nonce")
		return false
	}
	log.Debug("[MINT SIGNER] Found Nonce: ", nonce)

	data := autogen.MintControllerMintData{
		Recipient: address,
		Amount:    amount,
		Nonce:     nonce,
	}

	mint, err = util.UpdateStatusAndConfirmationsForMint(mint, x.poktHeight)
	if err != nil {
		log.Error("[MINT SIGNER] Error updating status and confirmations for mint: ", err)
		return false
	}

	var update bson.M
	if mint.Status == models.StatusConfirmed {
		log.Debug("[MINT SIGNER] Mint confirmed, signing")

		mint, err := util.SignMint(mint, data, x.domain, x.privateKey, x.numSigners)
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
				"updated_at":    time.Now(),
			},
		}

	} else {
		log.Debug("[MINT SIGNER] Mint pending confirmation, not signing")
		update = bson.M{
			"$set": bson.M{
				"status":        mint.Status,
				"confirmations": mint.Confirmations,
				"updated_at":    time.Now(),
			},
		}
	}

	filter := bson.M{
		"_id":    mint.Id,
		"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
	}

	err = app.DB.UpdateOne(models.CollectionMints, filter, update)
	if err != nil {
		log.Error("[MINT SIGNER] Error updating mint: ", err)
		return false
	}
	log.Info("[MINT SIGNER] Mint updated with signature")

	return true
}

func (x *MintSignerService) SyncTxs() bool {
	log.Debug("[MINT SIGNER] Syncing pending txs")

	filter := bson.M{
		"wpokt_address": x.wpoktAddress,
		"vault_address": x.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers": bson.M{
			"$nin": []string{x.address},
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
		success = x.HandleMint(mint) && success
	}

	log.Debug("[MINT SIGNER] Finished syncing pending txs")
	return success
}

func NewSigner(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	if app.Config.MintSigner.Enabled == false {
		log.Debug("[MINT SIGNER] BURN signer disabled")
		return models.NewEmptyService(wg)
	}

	log.Debug("[MINT SIGNER] Initializing mint signer")

	privateKey, err := crypto.HexToECDSA(app.Config.Ethereum.PrivateKey)
	if err != nil {
		log.Fatal("[MINT SIGNER] Error loading private key: ", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	log.Info("[MINT SIGNER] ETH signer address: ", address)

	ethClient, err := eth.NewClient()
	if err != nil {
		log.Fatal("[MINT SIGNER] Error initializing ethereum client: ", err)
	}

	log.Debug("[MINT SIGNER] Connecting to wpokt contract at: ", app.Config.Ethereum.WrappedPocketAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WrappedPocketAddress), ethClient.GetClient())
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutSecs)*time.Second)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx, Pending: false}
	domain, err := mintControllerContract.Eip712Domain(opts)

	if err != nil {
		log.Fatal("[MINT SIGNER] Error fetching mint controller domain data: ", err)
	}
	log.Debug("[MINT SIGNER] Fetched mint controller domain data")

	x := &MintSignerService{
		wg:                     wg,
		name:                   MintSignerName,
		stop:                   make(chan bool),
		interval:               time.Duration(app.Config.MintSigner.IntervalSecs) * time.Second,
		privateKey:             privateKey,
		address:                address,
		wpoktAddress:           app.Config.Ethereum.WrappedPocketAddress,
		vaultAddress:           app.Config.Pocket.VaultAddress,
		wpoktContract:          contract,
		mintControllerContract: mintControllerContract,
		numSigners:             len(app.Config.Ethereum.ValidatorAddresses),
		domain:                 domain,
		ethClient:              ethClient,
		poktClient:             pokt.NewClient(),
	}

	x.UpdateBlocks()

	x.UpdateHealth()

	log.Info("[MINT SIGNER] Initialized mint signer")

	return x
}
