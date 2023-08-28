package pokt

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
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestBurnSigner(t *testing.T, mockContract *eth.MockWrappedPocketContract, mockEthClient *eth.MockEthereumClient, mockPoktClient *pokt.MockPocketClient) *BurnSignerRunner {
	privateKey1, _ := crypto.NewPrivateKey("8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82")
	privateKey2, _ := crypto.NewPrivateKey("f2c227cd1299f62750e48d3e44c2d29cb3add4c8e9a171ae260b8fdeff49c761ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d055")
	privateKey3, _ := crypto.NewPrivateKey("05339b10520335644fe486e4d39ce33db4d079b5c1d3bceb725e75e4354f5ca7351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b")

	pks := []crypto.PublicKey{
		privateKey1.PublicKey(),
		privateKey2.PublicKey(),
		privateKey3.PublicKey(),
	}
	multisigPk := crypto.PublicKeyMultiSignature{PublicKeys: pks}

	app.Config.Pocket.TxFee = 10000

	x := &BurnSignerRunner{
		vaultAddress:   strings.ToLower(multisigPk.Address().String()),
		wpoktAddress:   "wpoktaddress",
		privateKey:     privateKey1,
		multisigPubKey: multisigPk,
		numSigners:     3,
		ethClient:      mockEthClient,
		poktClient:     mockPoktClient,
		poktHeight:     0,
		ethBlockNumber: 0,
		wpoktContract:  mockContract,
		minimumAmount:  big.NewInt(10000),
	}
	return x
}

func TestBurnSignerStatus(t *testing.T) {
	mockContract := eth.NewMockWrappedPocketContract(t)
	mockEthClient := eth.NewMockEthereumClient(t)
	mockPoktClient := pokt.NewMockPocketClient(t)
	x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "0")
	assert.Equal(t, status.PoktHeight, "0")
}

func TestBurnSignerUpdateBlocks(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, nil)
		mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(200))
		assert.Equal(t, x.ethBlockNumber, int64(200))
	})

	t.Run("With Error in Pokt Client", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(0))
		assert.Equal(t, x.ethBlockNumber, int64(0))
	})

	t.Run("With Error in Eth Client", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, nil)
		mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(200))
		assert.Equal(t, x.ethBlockNumber, int64(0))
	})

}

func TestBurnSignerValidateInvalidMint(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid transaction code", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		tx := &pokt.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg type", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		tx := &pokt.TxResponse{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg to address", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		ZERO_ADDRESS := "0x0000000000000000000000000000000000000000"

		mint := &models.InvalidMint{}

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

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg from address", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		ZERO_ADDRESS := "0x0000000000000000000000000000000000000000"

		mint := &models.InvalidMint{}

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

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg amount", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{
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

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Memo mismatch", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
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
						Amount:      "20000",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Memo is a valid mint memo", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
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
						Amount:      "100000",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Successful case", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          `{ "address": "0x0000000000000000000000000000000000000000", "chain_id": "31337" }`,
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
						Amount:      "20000",
					},
				},
				Memo: `{ "address": "0x0000000000000000000000000000000000000000", "chain_id": "31337" }`,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestBurnSignerValidateBurn(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(nil, errors.New("error"))

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid log index", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex: "index",
		}

		txReceipt := &types.Receipt{}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Log does not exist", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex: "0",
		}

		txReceipt := &types.Receipt{}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Error parsing log", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex: "0",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{
				{},
			},
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(nil, errors.New("error"))

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Amount mismatch", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex: "0",
			Amount:   "10",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(10),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Amount mismatch", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex: "0",
			Amount:   "100000",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(200000),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Sender mismatch", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex:      "0",
			Amount:        "20000",
			SenderAddress: "0xabcd",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
			From:   common.HexToAddress("0x1234"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Recipient mismatch", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: "abcd",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1234"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Successful case", func(t *testing.T) {

		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: "1c",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestBurnSignerHandleInvalidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		success := x.HandleInvalidMint(nil)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.Config.Pocket.Confirmations = 1

		invalidMint := &models.InvalidMint{
			Confirmations: "invalid",
			Height:        "invalid",
		}

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Error validating", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0

		invalidMint := &models.InvalidMint{
			Confirmations: "1",
			Height:        "99",
		}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation failure and update successful", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"

		address := common.HexToAddress("0x1234").Hex()

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			Confirmations: "1",
			Height:        "99",
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
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

	t.Run("Validation failure and update failed", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"

		address := common.HexToAddress("0x1234").Hex()

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			Confirmations: "1",
			Height:        "99",
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
						Amount:      "100",
					},
				},
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(errors.New("error")).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and signing failed", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
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
						Amount:      "20000",
					},
				},
				Memo: "invalid",
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and update successful", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		invalidMint := &models.InvalidMint{
			SenderAddress: x.privateKey.PublicKey().Address().String(),
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
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
						FromAddress: x.privateKey.PublicKey().Address().String(),
						Amount:      "20000",
					},
				},
				Memo: "invalid",
			},
		}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

	t.Run("Validation successful and mint pending and update successful", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 10
		app.Config.Ethereum.ChainId = "31337"

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
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
						Amount:      "20000",
					},
				},
				Memo: "invalid",
			},
		}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":        models.StatusPending,
				"confirmations": "1",
				"updated_at":    time.Now(),
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

}

func TestBurnSignerHandleBurn(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		success := x.HandleBurn(nil)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		app.Config.Ethereum.Confirmations = 1

		burn := &models.Burn{
			Confirmations: "invalid",
			BlockNumber:   "invalid",
		}

		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Error validating", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0

		burn := &models.Burn{
			Confirmations: "1",
			BlockNumber:   "99",
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(nil, errors.New("error"))

		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Validation failure and update successful", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"

		burn := &models.Burn{
			SenderAddress: "abcd",
			Amount:        "100",
			Confirmations: "1",
			BlockNumber:   "99",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.True(t, success)
	})

	t.Run("Validation failure and update failed", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"

		burn := &models.Burn{
			SenderAddress: "abcd",
			Amount:        "100",
			Confirmations: "1",
			BlockNumber:   "99",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusFailed,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(errors.New("error")).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and signing failed", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"

		burn := &models.Burn{
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: "1c",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)
		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and update successful", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		burn := &models.Burn{
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: strings.ToLower(strings.Split(common.HexToAddress("0x1c").Hex(), "0x")[1]),
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.True(t, success)
	})

	t.Run("Validation successful and mint pending and update successful", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 10
		app.Config.Ethereum.ChainId = "31337"

		burn := &models.Burn{
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: strings.ToLower(strings.Split(common.HexToAddress("0x1c").Hex(), "0x")[1]),
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":        models.StatusPending,
				"confirmations": "1",
				"updated_at":    time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.True(t, success)
	})

}

func TestBurnSignerSyncInvalidMints(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filter := bson.M{
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filter, mock.Anything).Return(nil)

		success := x.SyncInvalidMints()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					{
						Id: &primitive.NilObjectID,
					},
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", errors.New("error"))
		success := x.SyncInvalidMints()

		assert.False(t, success)

	})

	t.Run("Error unlocking", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		invalidMint := &models.InvalidMint{
			Id:            &primitive.NilObjectID,
			SenderAddress: x.privateKey.PublicKey().Address().String(),
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
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
						FromAddress: x.privateKey.PublicKey().Address().String(),
						Amount:      "20000",
					},
				},
				Memo: "invalid",
			},
		}

		filterUpdate := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*invalidMint,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		invalidMint := &models.InvalidMint{
			Id:            &primitive.NilObjectID,
			SenderAddress: x.privateKey.PublicKey().Address().String(),
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
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
						FromAddress: x.privateKey.PublicKey().Address().String(),
						Amount:      "20000",
					},
				},
				Memo: "invalid",
			},
		}

		filterUpdate := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*invalidMint,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncInvalidMints()

		assert.True(t, success)
	})
}

func TestBurnSignerSyncBurns(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filter := bson.M{
			"wpokt_address": x.wpoktAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filter, mock.Anything).Return(nil)

		success := x.SyncBurns()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					{
						Id: &primitive.NilObjectID,
					},
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", errors.New("error"))
		success := x.SyncBurns()

		assert.False(t, success)

	})

	t.Run("Error unlocking", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		burn := &models.Burn{
			Id:               &primitive.NilObjectID,
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: strings.ToLower(strings.Split(common.HexToAddress("0x1c").Hex(), "0x")[1]),
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		filterUpdate := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*burn,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockContract := eth.NewMockWrappedPocketContract(t)
		mockEthClient := eth.NewMockEthereumClient(t)
		mockPoktClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		burn := &models.Burn{
			Id:               &primitive.NilObjectID,
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: strings.ToLower(strings.Split(common.HexToAddress("0x1c").Hex(), "0x")[1]),
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		filterUpdate := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*burn,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncBurns()

		assert.True(t, success)
	})

}

func TestBurnSignerRun(t *testing.T) {

	mockContract := eth.NewMockWrappedPocketContract(t)
	mockEthClient := eth.NewMockEthereumClient(t)
	mockPoktClient := pokt.NewMockPocketClient(t)
	mockDB := app.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

	mockPoktClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, nil)
	mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

	{
		filterFind := bson.M{
			"vault_address": x.vaultAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		invalidMint := &models.InvalidMint{
			Id:            &primitive.NilObjectID,
			SenderAddress: x.privateKey.PublicKey().Address().String(),
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
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
						FromAddress: x.privateKey.PublicKey().Address().String(),
						Amount:      "20000",
					},
				},
				Memo: "invalid",
			},
		}

		filterUpdate := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil).Once()

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*invalidMint,
				}
			}).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil).Once()

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil).Once()
	}

	{
		filterFind := bson.M{
			"wpokt_address": x.wpoktAddress,
			"status":        bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
			"signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainId = "31337"
		app.Config.Pocket.ChainId = "testnet"
		app.Config.Pocket.TxFee = 10000

		burn := &models.Burn{
			Id:               &primitive.NilObjectID,
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: strings.ToLower(strings.Split(common.HexToAddress("0x1c").Hex(), "0x")[1]),
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x1c"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil).Once()
		mockContract.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil).Once()

		filterUpdate := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				"signers":       []string{x.privateKey.PublicKey().RawString()},
				"status":        models.StatusConfirmed,
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*burn,
				}
			}).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil).Once()

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				returnTx := gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"]
				assert.NotEmpty(t, returnTx)
				gotUpdate.(bson.M)["$set"].(bson.M)["return_tx"] = ""
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil).Once()

	}

	x.Run()

}

func TestNewBurnSigner(t *testing.T) {

	t.Run("Disabled", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = false

		service := NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, app.EmptyServiceName)

	})

	t.Run("Invalid Private Key", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.PrivateKey = ""
		app.Config.Ethereum.RPCURL = ""

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid Multisig keys", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.PrivateKey = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = ""
		app.Config.Pocket.MultisigPublicKeys = []string{
			"invalid",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid Vault Address", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.PrivateKey = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = ""
		app.Config.Pocket.VaultAddress = ""
		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid ETH RPC", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.PrivateKey = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = ""
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Interval is 0", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.PrivateKey = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		service := NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})

		assert.Nil(t, service)
	})

	t.Run("Valid", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.PrivateKey = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.BurnSigner.IntervalMillis = 1
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		service := NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, BurnSignerName)

	})
}
