package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"sort"
	"strconv"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type WPoktSignerService struct {
	stop                   chan bool
	address                string
	privateKey             *ecdsa.PrivateKey
	interval               time.Duration
	wpoktContract          *autogen.WrappedPocket
	mintControllerContract *autogen.MintController
	validators             []string
	domain                 DomainData
}

func (b *WPoktSignerService) Stop() {
	log.Debug("[WPOKT SIGNER] Stopping wpokt signer")
	b.stop <- true
}

func (b *WPoktSignerService) DetermineNonce(mint *models.Mint) {

}

func (b *WPoktSignerService) HandleMint(mint *models.Mint) bool {
	log.Debug("[WPOKT SIGNER] Handling mint: ", mint.TransactionHash)

	address := common.HexToAddress(mint.RecipientAddress)
	amount, ok := new(big.Int).SetString(mint.Amount, 10)
	if !ok {
		log.Error("[WPOKT SIGNER] Error converting decimal to big int")
		return false
	}

	var nonce *big.Int

	if mint.Data == nil {
		log.Debug("[WPOKT SIGNER] Mint data not set, fetching from contract")
		currentNonce, err := b.wpoktContract.GetUserNonce(nil, address)
		if err != nil {
			log.Error("[WPOKT SIGNER] Error fetching nonce: ", err)
			return false
		}
		var pendingMints []models.Mint
		filter := bson.M{
			"_id":               bson.M{"$ne": mint.Id},
			"recipient_address": mint.RecipientAddress,
			"status":            bson.M{"$in": []string{models.StatusPending, models.StatusSigned}},
		}
		err = app.DB.FindMany(models.CollectionMints, filter, &pendingMints)
		if err != nil {
			log.Error("[WPOKT SIGNER] Error fetching pending mints: ", err)
			return false
		}

		if len(pendingMints) > 0 {
			var nonces []int64

			for _, pendingMint := range pendingMints {
				if pendingMint.Data != nil {
					nonce, err := strconv.ParseInt(pendingMint.Data.Nonce, 10, 64)
					if err != nil {
						log.Error("[WPOKT SIGNER] Error converting nonce to int: ", err)
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
	} else {
		nonce, ok = new(big.Int).SetString(mint.Data.Nonce, 10)
		if !ok {
			log.Error("[WPOKT SIGNER] Error converting decimal to big int")
			return false
		}
	}

	data := autogen.MintControllerMintData{
		Recipient: address,
		Amount:    amount,
		Nonce:     nonce,
	}

	signature, err := SignTypedData(b.domain, data, b.privateKey)
	if err != nil {
		log.Error("[WPOKT SIGNER] Error signing typed data: ", err)
		return false
	}

	log.Debug("[WPOKT SIGNER] Mint signed")

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

	status := models.StatusPending // TODO confirmations

	if len(sortedSignatures) == len(b.validators) {
		log.Debug("[WPOKT SIGNER] Mint fully signed")
		status = models.StatusSigned
	}

	update := bson.M{
		"$set": bson.M{
			"data": models.MintData{
				Recipient: data.Recipient.Hex(),
				Amount:    data.Amount.String(),
				Nonce:     data.Nonce.String(),
			},
			"signatures": sortedSignatures,
			"signers":    sortedSigners,
			"status":     status,
		},
	}

	filter := bson.M{
		"_id": mint.Id,
	}

	err = app.DB.UpdateOne(models.CollectionMints, filter, update)
	if err != nil {
		log.Error("[WPOKT SIGNER] Error updating mint: ", err)
		return false
	}
	log.Debug("[WPOKT SIGNER] Mint updated with signature")

	return true
}

func (b *WPoktSignerService) SyncTxs() bool {
	// TODO confirmations

	filter := bson.M{
		"status": models.StatusPending,
		"signers": bson.M{
			"$nin": []string{b.address},
		},
	}

	var results []models.Mint

	err := app.DB.FindMany(models.CollectionMints, filter, &results)
	if err != nil {
		log.Error("[WPOKT SIGNER] Error fetching pending mints: ", err)
		return false
	}

	var success bool = true
	for _, mint := range results {
		success = b.HandleMint(&mint) && success

	}

	return success
}

func (b *WPoktSignerService) Start() {
	log.Debug("[WPOKT SIGNER] Starting wpokt signer")
	stop := false
	for !stop {
		log.Debug("[WPOKT SIGNER] Starting wpokt signer sync")

		b.SyncTxs()

		log.Debug("[WPOKT SIGNER] Finished wpokt signer sync")
		log.Debug("[WPOKT SIGNER] Sleeping for ", b.interval)

		select {
		case <-b.stop:
			stop = true
			log.Debug("[WPOKT SIGNER] Stopped wpokt signer")
		case <-time.After(b.interval):
		}
	}
}

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("[WPOKT SIGNER] Error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address
}

func NewSigner() models.Service {
	if app.Config.WPOKTSigner.Enabled == false {
		log.Debug("[WPOKT SIGNER] WPOKT signer disabled")
		return models.NewEmptyService()
	}

	log.Debug("[WPOKT SIGNER] Initializing wpokt signer")

	privateKey, err := crypto.HexToECDSA(app.Config.WPOKTSigner.PrivateKey)
	if err != nil {
		log.Fatal("[WPOKT SIGNER] Error loading private key: ", err)
	}

	address := privateKeyToAddress(privateKey)
	log.Debug("[WPOKT SIGNER] Loaded private key for address: ", address)

	log.Debug("[WPOKT SIGNER] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), Client.GetClient())
	if err != nil {
		log.Fatal("[WPOKT SIGNER] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[WPOKT SIGNER] Connected to wpokt contract")

	log.Debug("[WPOKT SIGNER] Connecting to mint controller contract at: ", app.Config.Ethereum.MintControllerAddress)
	mintControllerContract, err := autogen.NewMintController(common.HexToAddress(app.Config.Ethereum.MintControllerAddress), Client.GetClient())
	if err != nil {
		log.Fatal("[WPOKT SIGNER] Error initializing Mint Controller contract", err)
	}
	log.Debug("[WPOKT SIGNER] Connected to mint controller contract")

	log.Debug("[WPOKT SIGNER] Fetching mint controller domain data")
	domain, err := mintControllerContract.Eip712Domain(nil)
	if err != nil {
		log.Fatal("[WPOKT SIGNER] Error fetching mint controller domain data: ", err)
	}
	log.Debug("[WPOKT SIGNER] Fetched mint controller domain data")

	b := &WPoktSignerService{
		stop:                   make(chan bool),
		interval:               time.Duration(app.Config.WPOKTSigner.IntervalSecs) * time.Second,
		privateKey:             privateKey,
		address:                address,
		wpoktContract:          contract,
		mintControllerContract: mintControllerContract,
		validators:             sortAddresses(app.Config.Ethereum.ValidatorAddresses),
		domain:                 domain,
	}

	log.Debug("[WPOKT SIGNER] Initialized wpokt signer")

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
