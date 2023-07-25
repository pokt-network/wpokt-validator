package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
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
	validators             []string
	domain                 DomainData
	poktClient             pocket.PocketClient
	ethClient              ethereum.EthereumClient
	poktHeight             int64
}

func (b *MintSignerService) Health() models.ServiceHealth {
	return models.ServiceHealth{
		Name:           b.Name(),
		LastSyncTime:   b.LastSyncTime(),
		NextSyncTime:   b.LastSyncTime().Add(b.Interval()),
		PoktHeight:     b.PoktHeight(),
		EthBlockNumber: "",
		Healthy:        true,
	}
}

func (b *MintSignerService) PoktHeight() string {
	return strconv.FormatInt(b.poktHeight, 10)
}

func (b *MintSignerService) LastSyncTime() time.Time {
	return b.lastSyncTime
}

func (b *MintSignerService) Interval() time.Duration {
	return b.interval
}

func (b *MintSignerService) Name() string {
	return b.name
}

func (b *MintSignerService) Stop() {
	log.Debug("[BURN SIGNER] Stopping wpokt signer")
	b.stop <- true
}

func (b *MintSignerService) UpdateBlocks() {
	log.Debug("[BURN SIGNER] Updating blocks")
	poktHeight, err := b.poktClient.GetHeight()
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching pokt block height: ", err)
		return
	}
	b.poktHeight = poktHeight.Height
}

// finds nonce for mint transaction
func (b *MintSignerService) FindNonce(mint *models.Mint) (*big.Int, error) {
	log.Debug("[BURN SIGNER] Finding nonce for mint: ", mint.TransactionHash)
	var nonce *big.Int

	if mint.Nonce != "" {
		mintNonce, ok := new(big.Int).SetString(mint.Nonce, 10)
		if !ok {
			log.Error("[BURN SIGNER] Error converting decimal to big int")
			return nil, errors.New("error converting decimal to big int")
		}
		nonce = mintNonce
	}

	if nonce == nil || nonce.Cmp(big.NewInt(0)) == 0 {
		log.Debug("[BURN SIGNER] Mint nonce not set, fetching from contract")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
		defer cancel()
		opts := &bind.CallOpts{Context: ctx, Pending: false}
		currentNonce, err := b.wpoktContract.GetUserNonce(opts, common.HexToAddress(mint.RecipientAddress))

		if err != nil {
			log.Error("[BURN SIGNER] Error fetching nonce from contract: ", err)
			return nil, err
		}

		var pendingMints []models.Mint
		filter := bson.M{
			"_id":               bson.M{"$ne": mint.Id},
			"vault_address":     b.vaultAddress,
			"wpokt_address":     b.wpoktAddress,
			"recipient_address": mint.RecipientAddress,
			"status":            bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed, models.StatusSigned}},
		}
		err = app.DB.FindMany(models.CollectionMints, filter, &pendingMints)
		if err != nil {
			log.Error("[BURN SIGNER] Error fetching pending mints: ", err)
			return nil, err
		}

		if len(pendingMints) > 0 {
			var nonces []int64

			for _, pendingMint := range pendingMints {
				if pendingMint.Data != nil {
					nonce, err := strconv.ParseInt(pendingMint.Data.Nonce, 10, 64)
					if err != nil {
						log.Error("[BURN SIGNER] Error converting nonce to int: ", err)
						continue
					}
					nonces = append(nonces, nonce)
				}
			}

			if len(nonces) > 0 {
				sort.Slice(nonces, func(i, j int) bool {
					return nonces[i] < nonces[j]
				})

				currentNonce = big.NewInt(nonces[len(nonces)-1])
			}
		}

		nonce = currentNonce.Add(currentNonce, big.NewInt(1))
	}
	return nonce, nil
}

func (b *MintSignerService) HandleMint(mint *models.Mint) bool {
	log.Debug("[BURN SIGNER] Handling mint: ", mint.TransactionHash)

	address := common.HexToAddress(mint.RecipientAddress)
	amount, ok := new(big.Int).SetString(mint.Amount, 10)
	if !ok {
		log.Error("[BURN SIGNER] Error converting decimal to big int")
		return false
	}

	nonce, err := b.FindNonce(mint)

	if err != nil {
		log.Error("[BURN SIGNER] Error fetching nonce: ", err)
		return false
	}

	if nonce == nil || nonce.Cmp(big.NewInt(0)) == 0 {
		log.Error("[BURN SIGNER] Error fetching nonce")
		return false
	}

	data := autogen.MintControllerMintData{
		Recipient: address,
		Amount:    amount,
		Nonce:     nonce,
	}

	status := mint.Status
	confirmations, err := strconv.ParseInt(mint.Confirmations, 10, 64)
	if err != nil {
		confirmations = 0
	}

	if status == models.StatusPending {
		if app.Config.Pocket.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			log.Debug("[BURN SIGNER] Mint pending confirmation")
			mintHeight, err := strconv.ParseInt(mint.Height, 10, 64)
			if err != nil {
				log.Error("[BURN SIGNER] Error converting mint height to int: ", err)
				return false
			}
			totalConfirmations := b.poktHeight - mintHeight
			if totalConfirmations >= app.Config.Pocket.Confirmations {
				status = models.StatusConfirmed
			}
			confirmations = totalConfirmations
		}
	}

	var update bson.M
	if status == models.StatusConfirmed {
		log.Debug("[BURN SIGNER] Mint confirmed")
		signature, err := SignTypedData(b.domain, data, b.privateKey)
		if err != nil {
			log.Error("[BURN SIGNER] Error signing typed data: ", err)
			return false
		}

		log.Debug("[BURN SIGNER] Mint signed")

		signatureEncoded := "0x" + hex.EncodeToString(signature)
		signatures := append(mint.Signatures, signatureEncoded)
		signers := append(mint.Signers, b.address)
		sortedSigners := sortAddresses(signers)

		sortedSignatures := make([]string, len(signatures))

		for i, signature := range signatures {
			signer := signers[i]
			index := -1
			for j, validator := range sortedSigners {
				if validator == signer {
					index = j
					break
				}
			}
			sortedSignatures[index] = signature
		}

		if len(sortedSignatures) == len(b.validators) {
			log.Debug("[BURN SIGNER] Mint fully signed")
			status = models.StatusSigned
		}

		update = bson.M{
			"$set": bson.M{
				"data": models.MintData{
					Recipient: data.Recipient.Hex(),
					Amount:    data.Amount.String(),
					Nonce:     data.Nonce.String(),
				},
				"nonce":         data.Nonce.String(),
				"signatures":    sortedSignatures,
				"signers":       sortedSigners,
				"status":        status,
				"confirmations": strconv.FormatInt(confirmations, 10),
			},
		}

	} else {
		log.Debug("[BURN SIGNER] Mint pending confirmation")
		update = bson.M{
			"$set": bson.M{
				"status":        status,
				"confirmations": strconv.FormatInt(confirmations, 10),
			},
		}
	}

	filter := bson.M{
		"_id":           mint.Id,
		"wpokt_address": b.wpoktAddress,
		"vault_address": b.vaultAddress,
	}

	err = app.DB.UpdateOne(models.CollectionMints, filter, update)
	if err != nil {
		log.Error("[BURN SIGNER] Error updating mint: ", err)
		return false
	}
	log.Debug("[BURN SIGNER] Mint updated with signature")

	return true
}

func (b *MintSignerService) SyncTxs() bool {
	log.Debug("[BURN SIGNER] Syncing pending txs")

	filter := bson.M{
		"wpokt_address": b.wpoktAddress,
		"vault_address": b.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers": bson.M{
			"$nin": []string{b.address},
		},
	}

	var results []models.Mint

	err := app.DB.FindMany(models.CollectionMints, filter, &results)
	if err != nil {
		log.Error("[BURN SIGNER] Error fetching pending mints: ", err)
		return false
	}

	var success bool = true
	for _, mint := range results {
		success = b.HandleMint(&mint) && success

	}

	log.Debug("[BURN SIGNER] Finished syncing pending txs")
	return success
}

func (b *MintSignerService) Start() {
	log.Debug("[BURN SIGNER] Starting wpokt signer")
	stop := false
	for !stop {
		log.Debug("[BURN SIGNER] Starting wpokt signer sync")
		b.lastSyncTime = time.Now()

		b.UpdateBlocks()
		b.SyncTxs()

		log.Debug("[BURN SIGNER] Finished wpokt signer sync")
		log.Debug("[BURN SIGNER] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[BURN SIGNER] Stopped wpokt signer")
		case <-time.After(b.interval):
		}
	}
	b.wg.Done()
}

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("[BURN SIGNER] Error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address
}

func newSigner(wg *sync.WaitGroup) models.Service {
	if app.Config.MintSigner.Enabled == false {
		log.Debug("[BURN SIGNER] BURN signer disabled")
		return models.NewEmptyService(wg, "empty-wpokt-signer")
	}

	log.Debug("[BURN SIGNER] Initializing wpokt signer")

	privateKey, err := crypto.HexToECDSA(app.Config.Ethereum.PrivateKey)
	if err != nil {
		log.Fatal("[BURN SIGNER] Error loading private key: ", err)
	}

	address := privateKeyToAddress(privateKey)
	log.Debug("[BURN SIGNER] Loaded private key for address: ", address)
	ethClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing ethereum client: ", err)
	}

	log.Debug("[BURN SIGNER] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTAddress), ethClient.GetClient())
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[BURN SIGNER] Connected to wpokt contract")

	log.Debug("[BURN SIGNER] Connecting to mint controller contract at: ", app.Config.Ethereum.MintControllerAddress)
	mintControllerContract, err := autogen.NewMintController(common.HexToAddress(app.Config.Ethereum.MintControllerAddress), ethClient.GetClient())
	if err != nil {
		log.Fatal("[BURN SIGNER] Error initializing Mint Controller contract", err)
	}
	log.Debug("[BURN SIGNER] Connected to mint controller contract")

	log.Debug("[BURN SIGNER] Fetching mint controller domain data")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeOutSecs)*time.Second)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx, Pending: false}
	domain, err := mintControllerContract.Eip712Domain(opts)

	if err != nil {
		log.Fatal("[BURN SIGNER] Error fetching mint controller domain data: ", err)
	}
	log.Debug("[BURN SIGNER] Fetched mint controller domain data")

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
		validators:             sortAddresses(app.Config.Ethereum.ValidatorAddresses),
		domain:                 domain,
		ethClient:              ethClient,
		poktClient:             pocket.NewClient(),
	}

	log.Debug("[BURN SIGNER] Initialized wpokt signer")

	return b
}

func sortAddresses(addresses []string) []string {
	for i, address := range addresses {
		addresses[i] = common.HexToAddress(address).Hex()
	}
	sort.Slice(addresses, func(i, j int) bool {
		return common.HexToAddress(addresses[i]).Big().Cmp(common.HexToAddress(addresses[j]).Big()) == -1
	})
	return addresses
}

func NewSigner(wg *sync.WaitGroup) models.Service {
	return newSigner(wg)
}

func NewSignerWithLastHealth(wg *sync.WaitGroup, lastHealth models.ServiceHealth) models.Service {
	return newSigner(wg)
}
