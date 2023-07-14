package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/dan13ram/wpokt-backend/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	beeCrypto "github.com/ethersphere/bee/pkg/crypto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var typesStandard = apitypes.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"MintControllerMintData": {
		{
			Name: "recipient",
			Type: "address",
		},
		{
			Name: "amount",
			Type: "uint256",
		},
		{
			Name: "nonce",
			Type: "uint256",
		},
	},
}

var domainStandard = apitypes.TypedDataDomain{
	Name:              "MintController",
	Version:           "1",
	ChainId:           math.NewHexOrDecimal256(5),
	VerifyingContract: "0x9742FC3ed4a14B556D11a64a2F9D4b72Df5f5a67",
}

const primaryType = "MintControllerMintData"

type WPoktSignerService struct {
	stop          chan bool
	address       string
	signer        beeCrypto.Signer
	interval      time.Duration
	wpoktContract *autogen.WrappedPocket
}

func (b *WPoktSignerService) Stop() {
	log.Debug("[WPOKT SIGNER] Stopping wpokt signer")
	b.stop <- true
}

func (b *WPoktSignerService) HandleMint(mint *models.Mint) bool {
	log.Debug("[WPOKT SIGNER] Handling mint: ", mint.TransactionHash)

	address := common.HexToAddress(mint.RecipientAddress)
	amount, ok := new(big.Int).SetString(mint.Amount, 10)
	if !ok {
		log.Error("[WPOKT SIGNER] Error converting decimal to big int")
		return false
	}
	nonce, err := b.wpoktContract.GetUserNonce(nil, address)
	if err != nil {
		log.Error("[WPOKT SIGNER] Error fetching nonce: ", err)
		return false
	}

	data := autogen.MintControllerMintData{
		Recipient: address,
		Amount:    amount,
		Nonce:     nonce,
	}

	message := apitypes.TypedDataMessage{
		"recipient": mint.RecipientAddress,
		"amount":    mint.Amount,
		"nonce":     nonce.String(),
	}

	typedData := apitypes.TypedData{
		Types:       typesStandard,
		PrimaryType: primaryType,
		Domain:      domainStandard,
		Message:     message,
	}

	signature, err := b.signer.SignTypedData(&typedData)

	if err != nil {
		log.Error("[WPOKT SIGNER] Error signing typed data: ", err)
		return false
	}

	signatureEncoded := hex.EncodeToString(signature)

	log.Debug("[WPOKT SIGNER] Data: ", data)
	log.Debug("[WPOKT SIGNER] Signature: ", signatureEncoded)

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

	signer := beeCrypto.NewDefaultSigner(privateKey)

	address := privateKeyToAddress(privateKey)
	log.Debug("[WPOKT SIGNER] Loaded private key for address: ", address)

	log.Debug("[WPOKT SIGNER] Connecting to wpokt contract at: ", app.Config.Ethereum.WPOKTContractAddress)
	contract, err := autogen.NewWrappedPocket(common.HexToAddress(app.Config.Ethereum.WPOKTContractAddress), Client.GetClient())
	if err != nil {
		log.Fatal("[WPOKT SIGNER] Error initializing Wrapped Pocket contract", err)
	}
	log.Debug("[WPOKT SIGNER] Connected to wpokt contract")

	b := &WPoktSignerService{
		stop:          make(chan bool),
		interval:      time.Duration(app.Config.WPOKTSigner.IntervalSecs) * time.Second,
		signer:        signer,
		address:       address,
		wpoktContract: contract,
	}

	log.Debug("[WPOKT SIGNER] Initialized wpokt signer")

	return b
}
