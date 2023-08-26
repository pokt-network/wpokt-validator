package eth

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/eth/util"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestMintSigner(t *testing.T, mockContract *eth.MockWrappedPocketContract, mockEthClient *eth.MockEthereumClient, mockPoktClient *pokt.MockPocketClient) *MintSignerRunner {
	pk, _ := crypto.HexToECDSA("1395eeb9c36ef43e9e05692c9ee34034c00a9bef301135a96d082b2a65fd1680")
	address := crypto.PubkeyToAddress(pk.PublicKey).Hex()

	x := &MintSignerRunner{
		address:      strings.ToLower(address),
		privateKey:   pk,
		vaultAddress: "vaultAddress",
		wpoktAddress: "wpoktAddress",
		numSigners:   3,
		domain: util.DomainData{
			Name:              "Test",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: common.HexToAddress(""),
		},
		wpoktContract: mockContract,
		ethClient:     mockEthClient,
		poktClient:    mockPoktClient,
		poktHeight:    100,
	}
	return x
}

func TestMintSignerStatus(t *testing.T) {
	mockContract := eth.NewMockWrappedPocketContract(t)
	mockEthClient := eth.NewMockEthereumClient(t)
	mockPoktClient := pokt.NewMockPocketClient(t)
	x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "100")
}

func TestMintSignerUpdateCurrentBlockNumber(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{
			Height: 200,
		}, nil)

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(200))
	})

	t.Run("With Error", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetHeight().Return(nil, errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(100))
	})

}

func TestMintSignerFindNonce(t *testing.T) {

	t.Run("Nonce already set", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{
			Nonce: "10",
		}

		gotNonce, err := x.FindNonce(mint)

		assert.Equal(t, gotNonce, big.NewInt(10))
		assert.Nil(t, err)

	})

	t.Run("No pending mints", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		nonce := big.NewInt(10)

		mockContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(nonce, nil)

		filter := bson.M{
			"_id":               bson.M{"$ne": mint.Id},
			"vault_address":     x.vaultAddress,
			"wpokt_address":     x.wpoktAddress,
			"recipient_address": strings.ToLower(mint.RecipientAddress),
			"status":            bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed, models.StatusSigned}},
		}

		mockDB.EXPECT().FindMany(models.CollectionMints, filter, mock.Anything).
			Run(func(collection string, _ interface{}, data interface{}) {
				d := data.(*[]models.Mint)
				*d = []models.Mint{}
			}).Return(nil)

		gotNonce, err := x.FindNonce(mint)

		assert.Equal(t, gotNonce, big.NewInt(11))
		assert.Nil(t, err)

	})

	t.Run("With pending mints but current nonce is higher", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		nonce := big.NewInt(5)

		mockContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(nonce, nil)

		filter := bson.M{
			"_id":               bson.M{"$ne": mint.Id},
			"vault_address":     x.vaultAddress,
			"wpokt_address":     x.wpoktAddress,
			"recipient_address": strings.ToLower(mint.RecipientAddress),
			"status":            bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed, models.StatusSigned}},
		}

		mockDB.EXPECT().FindMany(models.CollectionMints, filter, mock.Anything).
			Run(func(collection string, _ interface{}, data interface{}) {
				d := data.(*[]models.Mint)
				*d = []models.Mint{
					{
						Data: &models.MintData{
							Nonce: "invalid",
						},
					},
					{
						Data: &models.MintData{
							Nonce: "4",
						},
					},
					{
						Data: &models.MintData{
							Nonce: "5",
						},
					},
					{
						Data: &models.MintData{
							Nonce: "6",
						},
					},
				}
			}).Return(nil)

		gotNonce, err := x.FindNonce(mint)

		assert.Equal(t, gotNonce, big.NewInt(7))
		assert.Nil(t, err)

	})

	t.Run("Error with converting nonce", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{
			Nonce: "invalid",
		}

		gotNonce, err := x.FindNonce(mint)

		assert.NotNil(t, err)
		assert.Nil(t, gotNonce)

	})

	t.Run("Error finding current nonce", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		mockContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(big.NewInt(5), errors.New("error"))

		gotNonce, err := x.FindNonce(mint)

		assert.NotNil(t, err)
		assert.Nil(t, gotNonce)

	})

	t.Run("Error finding pending mints", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		mockContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(big.NewInt(5), nil)
		mockDB.EXPECT().FindMany(models.CollectionMints, mock.Anything, mock.Anything).Return(errors.New("error"))

		gotNonce, err := x.FindNonce(mint)

		assert.NotNil(t, err)
		assert.Nil(t, gotNonce)

	})

}

func TestValidateMint(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid transaction code", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		tx := &pokt.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg type", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg to address", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		ZERO_ADDRESS := "0x0000000000000000000000000000000000000000"

		mint := &models.Mint{}

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress: ZERO_ADDRESS,
					},
				},
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg from address", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		ZERO_ADDRESS := "0x0000000000000000000000000000000000000000"

		mint := &models.Mint{}

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: ZERO_ADDRESS,
					},
				},
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg amount", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{
			SenderAddress: "abcd",
			Amount:        "100",
		}

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
					},
				},
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction memo", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{
			SenderAddress: "abcd",
			Amount:        "100",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: `{ "address": "0x0000000000000000000000000000000000000000", "chain_id": "31337" }`,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction memo address", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress: "abcd",
			Amount:        "100",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction memo chainId", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Valid transaction and mint", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			RecipientChainId: "31337",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateMint(mint)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestMintSignerHandleMint(t *testing.T) {

	t.Run("Nil mint", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		success := x.HandleMint(nil)

		assert.False(t, success)
	})

	t.Run("Invalid mint amount", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "invalid",
			RecipientChainId: "31337",
		}

		success := x.HandleMint(mint)

		assert.False(t, success)
	})

	t.Run("Error finding nonce", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "invalid",
			RecipientChainId: "31337",
		}

		success := x.HandleMint(mint)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 1

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "invalid",
			Confirmations:    "invalid",
		}

		success := x.HandleMint(mint)

		assert.False(t, success)
	})

	t.Run("Error validating mint", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.False(t, success)
	})

	t.Run("Validating mint returned false", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        10,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    mint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}
		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filter, mock.Anything).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Return(nil)

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.True(t, success)
	})

	t.Run("Validating mint returned true, mint pending", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 10

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    mint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}
		update := bson.M{
			"$set": bson.M{
				"status":        models.StatusPending,
				"confirmations": "1",
				"updated_at":    time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filter, mock.Anything).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Return(nil)

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusPending, mint.Status)
		assert.Equal(t, mint.Confirmations, "1")

		assert.True(t, success)
	})

	t.Run("Validating mint returned true, mint confirmed, signing failed", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.domain = util.DomainData{
			ChainId: big.NewInt(1),
		}

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.False(t, success)
	})

	t.Run("Error updating mint", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    mint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}
		update := bson.M{
			"$set": bson.M{
				"data": models.MintData{
					Recipient: strings.ToLower(mint.RecipientAddress),
					Amount:    mint.Amount,
					Nonce:     mint.Nonce,
				},
				"nonce":         mint.Nonce,
				"signatures":    mint.Signatures,
				"signers":       []string{x.address},
				"status":        models.StatusConfirmed,
				"confirmations": "0",
				"updated_at":    time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filter, mock.Anything).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				gotUpdate.(bson.M)["$set"].(bson.M)["signatures"] = update["$set"].(bson.M)["signatures"]
				assert.Equal(t, update, gotUpdate)
			}).Return(errors.New("error"))

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    mint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}
		update := bson.M{
			"$set": bson.M{
				"data": models.MintData{
					Recipient: strings.ToLower(mint.RecipientAddress),
					Amount:    mint.Amount,
					Nonce:     mint.Nonce,
				},
				"nonce":         mint.Nonce,
				"signatures":    mint.Signatures,
				"signers":       []string{x.address},
				"status":        models.StatusConfirmed,
				"confirmations": "0",
				"updated_at":    time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filter, mock.Anything).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				gotUpdate.(bson.M)["$set"].(bson.M)["signatures"] = update["$set"].(bson.M)["signatures"]
				assert.Equal(t, update, gotUpdate)
			}).Return(nil)

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.True(t, success)
	})

}

func TestMintSignerSyncTxs(t *testing.T) {

	t.Run("Error finding mints", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncTxs()

		assert.False(t, success)

	})

	t.Run("No mints to handle", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		filter := bson.M{
			"wpokt_address": x.wpoktAddress,
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers": bson.M{
				"$nin": []string{x.address},
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionMints, filter, mock.Anything).Return(nil)

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers": bson.M{
				"$nin": []string{x.address},
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Mint)
				*v = []models.Mint{
					{},
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", errors.New("error"))
		success := x.SyncTxs()

		assert.False(t, success)

	})

	t.Run("Error unlocking", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers": bson.M{
				"$nin": []string{x.address},
			},
		}

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		filterUpdate := bson.M{
			"_id":    mint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}
		update := bson.M{
			"$set": bson.M{
				"data": models.MintData{
					Recipient: strings.ToLower(mint.RecipientAddress),
					Amount:    mint.Amount,
					Nonce:     mint.Nonce,
				},
				"nonce":         mint.Nonce,
				"signatures":    mint.Signatures,
				"signers":       []string{x.address},
				"status":        models.StatusConfirmed,
				"confirmations": "0",
				"updated_at":    time.Now(),
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Mint)
				*v = []models.Mint{
					*mint,
				}
			})

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filterUpdate, mock.Anything).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				gotUpdate.(bson.M)["$set"].(bson.M)["signatures"] = update["$set"].(bson.M)["signatures"]
				assert.Equal(t, update, gotUpdate)
			}).Return(nil)

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers": bson.M{
				"$nin": []string{x.address},
			},
		}

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "100",
			Nonce:            "1",
			RecipientChainId: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainId = "31337"

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Msg: pokt.Msg{
					Type: "pos/Send",
					Value: pokt.Value{
						ToAddress:   x.vaultAddress,
						FromAddress: "abcd",
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		filterUpdate := bson.M{
			"_id":    mint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}
		update := bson.M{
			"$set": bson.M{
				"data": models.MintData{
					Recipient: strings.ToLower(mint.RecipientAddress),
					Amount:    mint.Amount,
					Nonce:     mint.Nonce,
				},
				"nonce":         mint.Nonce,
				"signatures":    mint.Signatures,
				"signers":       []string{x.address},
				"status":        models.StatusConfirmed,
				"confirmations": "0",
				"updated_at":    time.Now(),
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Mint)
				*v = []models.Mint{
					*mint,
				}
			})

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filterUpdate, mock.Anything).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				gotUpdate.(bson.M)["$set"].(bson.M)["signatures"] = update["$set"].(bson.M)["signatures"]
				assert.Equal(t, update, gotUpdate)
			}).Return(nil)

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncTxs()

		assert.True(t, success)
	})

}

func TestMintSignerRun(t *testing.T) {

	mockContract := eth.NewMockWrappedPocketContract(t)
	mockEthClient := eth.NewMockEthereumClient(t)
	mockPoktClient := pokt.NewMockPocketClient(t)
	mockDB := app.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestMintSigner(t, mockContract, mockEthClient, mockPoktClient)

	filterFind := bson.M{
		"wpokt_address": x.wpoktAddress,
		"vault_address": x.vaultAddress,
		"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		"signers": bson.M{
			"$nin": []string{x.address},
		},
	}

	address := common.HexToAddress("0x1234").Hex()

	app.Config.Pocket.Confirmations = 0

	mint := &models.Mint{
		SenderAddress:    "abcd",
		RecipientAddress: address,
		Amount:           "100",
		Nonce:            "1",
		RecipientChainId: "31337",
		Height:           "99",
		Confirmations:    "invalid",
	}

	app.Config.Ethereum.ChainId = "31337"

	tx := &pokt.TxResponse{
		Tx: "abcd",
		TxResult: pokt.TxResult{
			Code:        0,
			MessageType: "send",
		},
		StdTx: pokt.StdTx{
			Msg: pokt.Msg{
				Type: "pos/Send",
				Value: pokt.Value{
					ToAddress:   x.vaultAddress,
					FromAddress: "abcd",
					Amount:      "100",
				},
			},
			Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
		},
	}

	mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

	filterUpdate := bson.M{
		"_id":    mint.Id,
		"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
	}
	update := bson.M{
		"$set": bson.M{
			"data": models.MintData{
				Recipient: strings.ToLower(mint.RecipientAddress),
				Amount:    mint.Amount,
				Nonce:     mint.Nonce,
			},
			"nonce":         mint.Nonce,
			"signatures":    mint.Signatures,
			"signers":       []string{x.address},
			"status":        models.StatusConfirmed,
			"confirmations": "0",
			"updated_at":    time.Now(),
		},
	}

	mockDB.EXPECT().FindMany(models.CollectionMints, filterFind, mock.Anything).Return(nil).
		Run(func(_ string, _ interface{}, result interface{}) {
			v := result.(*[]models.Mint)
			*v = []models.Mint{
				*mint,
			}
		})

	mockDB.EXPECT().UpdateOne(models.CollectionMints, filterUpdate, mock.Anything).
		Run(func(_ string, _ interface{}, gotUpdate interface{}) {
			gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
			gotUpdate.(bson.M)["$set"].(bson.M)["signatures"] = update["$set"].(bson.M)["signatures"]
			assert.Equal(t, update, gotUpdate)
		}).Return(nil)

	mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
	mockDB.EXPECT().Unlock("lockId").Return(nil)

	mockPoktClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{
		Height: 200,
	}, nil)

	x.Run()

}

func TestNewMintSigner(t *testing.T) {

	t.Run("Disabled", func(t *testing.T) {

		app.Config.MintSigner.Enabled = false

		service := NewMintSigner(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, app.EmptyServiceName)

	})

	t.Run("Invalid private key", func(t *testing.T) {

		app.Config.MintSigner.Enabled = true
		app.Config.Ethereum.PrivateKey = "0xinvalid"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})

	})

	t.Run("Invalid ETH RPC", func(t *testing.T) {

		app.Config.MintSigner.Enabled = true
		app.Config.Ethereum.PrivateKey = "1395eeb9c36ef43e9e05692c9ee34034c00a9bef301135a96d082b2a65fd1680"
		app.Config.Ethereum.RPCURL = ""

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})

	})

	t.Run("Error fetching domain data", func(t *testing.T) {

		app.Config.MintSigner.Enabled = true
		app.Config.Ethereum.PrivateKey = "1395eeb9c36ef43e9e05692c9ee34034c00a9bef301135a96d082b2a65fd1680"
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

}
