package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/eth/util"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	poktUtil "github.com/dan13ram/wpokt-validator/pokt/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MintSignerName = "MINT SIGNER"
)

type MintSignerRunner struct {
	address       string
	privateKey    *ecdsa.PrivateKey
	vaultAddress  string
	wpoktAddress  string
	wpoktContract eth.WrappedPocketContract
	numSigners    int
	domain        util.DomainData
	poktClient    pokt.PocketClient
	ethClient     eth.EthereumClient
	poktHeight    int64
	minimumAmount *big.Int
}

func (x *MintSignerRunner) Run() {
	x.UpdateBlocks()
	x.SyncTxs()
}

func (x *MintSignerRunner) Status() models.RunnerStatus {
	return models.RunnerStatus{
		PoktHeight: strconv.FormatInt(x.poktHeight, 10),
	}
}

func (x *MintSignerRunner) UpdateBlocks() {
	log.Debug("[MINT SIGNER] Updating blocks")
	poktHeight, err := x.poktClient.GetHeight()
	if err != nil {
		log.Error("[MINT SIGNER] Error fetching pokt block height: ", err)
		return
	}
	x.poktHeight = poktHeight.Height
}

func (x *MintSignerRunner) FindNonce(mint *models.Mint) (*big.Int, error) {
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
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
			"recipient_address": strings.ToLower(mint.RecipientAddress),
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

func (x *MintSignerRunner) ValidateMint(mint *models.Mint) (bool, error) {
	log.Debug("[MINT SIGNER] Validating mint: ", mint.TransactionHash)

	tx, err := x.poktClient.GetTx(mint.TransactionHash)
	if err != nil {
		return false, errors.New("Error fetching transaction: " + err.Error())
	}

	if tx.Tx == "" || tx.TxResult.Code != 0 {
		log.Debug("[MINT SIGNER] Transaction not found or failed")
		return false, nil
	}

	if tx.TxResult.MessageType != "send" || tx.StdTx.Msg.Type != "pos/Send" {
		log.Debug("[MINT SIGNER] Transaction message type is not send")
		return false, nil
	}

	if !strings.EqualFold(tx.StdTx.Msg.Value.ToAddress, x.vaultAddress) {
		log.Debug("[MINT SIGNER] Transaction recipient is not vault address")
		return false, nil
	}

	if !strings.EqualFold(tx.StdTx.Msg.Value.FromAddress, mint.SenderAddress) {
		log.Debug("[MINT SIGNER] Transaction signer is not sender address")
		return false, nil
	}

	amount, ok := new(big.Int).SetString(tx.StdTx.Msg.Value.Amount, 10)

	if !ok || amount.Cmp(x.minimumAmount) != 1 {
		log.Debug("[MINT SIGNER] Transaction amount too low")
		return false, nil
	}

	if tx.StdTx.Msg.Value.Amount != mint.Amount {
		log.Debug("[MINT SIGNER] Transaction amount does not match mint amount")
		return false, nil
	}

	memo, valid := poktUtil.ValidateMemo(tx.StdTx.Memo)
	if !valid {
		log.Debug("[MINT SIGNER] Memo failed validation")
		return false, nil
	}

	if !strings.EqualFold(memo.Address, mint.RecipientAddress) {
		log.Debug("[MINT SIGNER] Memo address does not match recipient address")
		return false, nil
	}

	if memo.ChainId != mint.RecipientChainId {
		log.Debug("[MINT SIGNER] Memo chain id does not match recipient chain id")
		return false, nil
	}

	log.Debug("[MINT SIGNER] Mint validated")
	return true, nil
}

func (x *MintSignerRunner) HandleMint(mint *models.Mint) bool {
	if mint == nil {
		log.Error("[MINT EXECUTOR] Invalid mint")
		return false
	}

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

	data := &autogen.MintControllerMintData{
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

	valid, err := x.ValidateMint(mint)
	if err != nil {
		log.Error("[MINT SIGNER] Error validating mint: ", err)
		return false
	}

	if !valid {
		log.Error("[MINT SIGNER] Mint failed validation")
		update = bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}
	} else {

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
						Recipient: strings.ToLower(data.Recipient.Hex()),
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
	log.Info("[MINT SIGNER] Handled mint: ", mint.TransactionHash)

	return true
}

func (x *MintSignerRunner) SyncTxs() bool {
	log.Debug("[MINT SIGNER] Syncing pending txs")

	filter := bson.M{
		"wpokt_address": x.wpoktAddress,
		"vault_address": x.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers": bson.M{
			"$nin": []string{x.address},
		},
	}

	var mints []models.Mint

	err := app.DB.FindMany(models.CollectionMints, filter, &mints)
	if err != nil {
		log.Error("[MINT SIGNER] Error fetching pending mints: ", err)
		return false
	}

	var success bool = true
	for i := range mints {
		mint := mints[i]

		resourceId := fmt.Sprintf("%s/%s", models.CollectionMints, strings.ToLower(mint.RecipientAddress))
		lockId, err := app.DB.XLock(resourceId)
		if err != nil {
			log.Error("[MINT SIGNER] Error locking mint: ", err)
			success = false
			continue
		}
		log.Debug("[MINT SIGNER] Locked mint: ", mint.TransactionHash)

		success = x.HandleMint(&mint) && success

		if err = app.DB.Unlock(lockId); err != nil {
			log.Error("[MINT SIGNER] Error unlocking mint: ", err)
			success = false
		} else {
			log.Debug("[MINT SIGNER] Unlocked mint: ", mint.TransactionHash)
		}

	}

	log.Debug("[MINT SIGNER] Finished syncing pending txs")
	return success
}

func NewMintSigner(wg *sync.WaitGroup, lastHealth models.ServiceHealth) app.Service {
	if !app.Config.MintSigner.Enabled {
		log.Debug("[MINT SIGNER] Disabled")
		return app.NewEmptyService(wg)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.Ethereum.RPCTimeoutMillis)*time.Millisecond)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx, Pending: false}
	domain, err := mintControllerContract.Eip712Domain(opts)

	if err != nil {
		log.Fatal("[MINT SIGNER] Error fetching mint controller domain data: ", err)
	}
	log.Debug("[MINT SIGNER] Fetched mint controller domain data")

	x := &MintSignerRunner{
		privateKey:    privateKey,
		address:       strings.ToLower(address),
		wpoktAddress:  strings.ToLower(app.Config.Ethereum.WrappedPocketAddress),
		vaultAddress:  strings.ToLower(app.Config.Pocket.VaultAddress),
		wpoktContract: eth.NewWrappedPocketContract(contract),
		numSigners:    len(app.Config.Ethereum.ValidatorAddresses),
		domain:        domain,
		ethClient:     ethClient,
		poktClient:    pokt.NewClient(),
		minimumAmount: big.NewInt(app.Config.Pocket.TxFee),
	}

	x.UpdateBlocks()

	log.Info("[MINT SIGNER] Initialized mint signer")

	return app.NewRunnerService(MintSignerName, x, wg, time.Duration(app.Config.MintSigner.IntervalMillis)*time.Millisecond)
}
