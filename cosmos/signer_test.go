package cosmos

import (
	"errors"
	// "fmt"
	"io"
	// "math/big"
	// "strings"
	// "sync"
	"testing"
	// "time"

	"github.com/dan13ram/wpokt-validator/app"
	// appMocks "github.com/dan13ram/wpokt-validator/app/mocks"
	// cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	cosmosMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	// "github.com/dan13ram/wpokt-validator/eth/autogen"
	// eth "github.com/dan13ram/wpokt-validator/eth/client"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	ethMocks "github.com/dan13ram/wpokt-validator/eth/client/mocks"
	// "github.com/dan13ram/wpokt-validator/models"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	"cosmossdk.io/math"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestBurnSigner(t *testing.T,
	mockWPOKT *ethMocks.MockWrappedPocketContract,
	mockMintController *ethMocks.MockMintControllerContract,
	mockEthClient *ethMocks.MockEthereumClient,
	mockPoktClient *cosmosMocks.MockCosmosClient,
) *BurnSignerRunner {
	app.Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	app.Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
	app.Config.Pocket.MultisigPublicKeys = []string{
		"0223aa679d6d5344e201e0df9f02ab15a84726eee0dfb4e953c46a9e2cb52349dc",
		"02faaaf0f385bb17381f36dcd86ab2486e8ff8d93440436496665ac007953076c2",
		"02cae233806460db75a941a269490ca5165a620b43241edb8bc72e169f4143a6df",
	}
	app.Config.Pocket.MultisigAddress = "pokt10r5n6x28p9qntchsmhxd4ftq9lk6vzcx3dv4gx"
	app.Config.Pocket.MultisigThreshold = 2
	app.Config.Pocket.Bech32Prefix = "pokt"
	app.Config.Pocket.TxFee = 10000

	signer, err := app.GetPocketSignerAndMultisig()
	assert.Nil(t, err)

	x := &BurnSignerRunner{
		signer:                 signer,
		ethClient:              mockEthClient,
		poktClient:             mockPoktClient,
		poktHeight:             0,
		ethBlockNumber:         0,
		vaultAddress:           signer.MultisigAddress,
		wpoktAddress:           "wpoktaddress",
		wpoktContract:          mockWPOKT,
		mintControllerContract: mockMintController,
		minimumAmount:          math.NewInt(10000),
		maximumAmount:          math.NewInt(20000),
	}
	return x
}

func TestBurnSignerStatus(t *testing.T) {
	mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
	mockMintController := ethMocks.NewMockMintControllerContract(t)
	mockEthClient := ethMocks.NewMockEthereumClient(t)
	mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
	x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockPoktClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "0")
	assert.Equal(t, status.PoktHeight, "0")
}

func TestBurnSignerUpdateBlocks(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetLatestBlockHeight().Return(200, nil)
		mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(200))
		assert.Equal(t, x.ethBlockNumber, int64(200))
	})

	t.Run("With Error in Pokt Client", func(t *testing.T) {
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetLatestBlockHeight().Return(200, errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(0))
		assert.Equal(t, x.ethBlockNumber, int64(0))
	})

	t.Run("With Error in Eth Client", func(t *testing.T) {
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetLatestBlockHeight().Return(200, nil)
		mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.poktHeight, int64(200))
		assert.Equal(t, x.ethBlockNumber, int64(0))
	})

}

/*
func TestBurnSignerValidateInvalidMint(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid transaction code", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		tx := &sdk.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg type", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		tx := &sdk.TxResponse{
			Tx: "abcd",
			TxResult: sdk.TxResult{
				Code: 0,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg to address", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		ZERO_ADDRESS := "0000000000000000000000000000000000000000"

		mint := &models.InvalidMint{}

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress: ZERO_ADDRESS,
			// 		},
			// 	},
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})
	t.Run("Incorrect transaction msg to address", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{}

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress: "abcd",
			// 		},
			// 	},
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg from address", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		ZERO_ADDRESS := "0x0000000000000000000000000000000000000000"

		mint := &models.InvalidMint{}

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: ZERO_ADDRESS,
			// 		},
			// 	},
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg amount", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
		}

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: "abcd",
			// 		},
			// 	},
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Memo mismatch", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: "abcd",
			// 			Amount:      "20000",
			// 		},
			// 	},
			// 	Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Memo is a valid mint memo", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: "abcd",
			// 			Amount:      "100000",
			// 		},
			// 	},
			// 	Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Successful case", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          `{ "address": "0x0000000000000000000000000000000000000000", "chain_id": "31337" }`,
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: "abcd",
			// 			Amount:      "20000",
			// 		},
			// 	},
			// 	Memo: `{ "address": "0x0000000000000000000000000000000000000000", "chain_id": "31337" }`,
			// },
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestBurnSignerValidateBurn(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		burn := &models.Burn{}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(nil, errors.New("error"))

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid log index", func(t *testing.T) {

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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

		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		success := x.HandleInvalidMint(nil)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

		address := common.HexToAddress("0x1234").Hex()

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			Confirmations: "1",
			Height:        "99",
		}

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: "abcd",
			// 			Amount:      "100",
			// 		},
			// 	},
			// 	Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			// },
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
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

	t.Run("Validation failure and update failed", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

		address := common.HexToAddress("0x1234").Hex()

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			Confirmations: "1",
			Height:        "99",
		}

		tx := &sdk.TxResponse{
			// Tx: "abcd",
			// TxResult: sdk.TxResult{
			// 	Code:        0,
			// 	MessageType: "send",
			// },
			// StdTx: sdk.StdTx{
			// 	Msg: sdk.Msg{
			// 		Type: "pos/Send",
			// 		Value: sdk.Value{
			// 			ToAddress:   x.vaultAddress,
			// 			FromAddress: "abcd",
			// 			Amount:      "100",
			// 		},
			// 	},
			// 	Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			// },
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
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), errors.New("error")).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and signing failed", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		tx := &sdk.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and update successful", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000

		invalidMint := &models.InvalidMint{
			// SenderAddress: x.privateKey.PublicKey().Address().String(),
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		tx := &sdk.TxResponse{}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"confirmations": "1",
				"updated_at":    time.Now(),
				"return_tx":     "",
				// "signers":       []string{x.privateKey.PublicKey().RawString()},
				"status": models.StatusConfirmed,
			},
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.poktHeight = 100
		app.Config.Pocket.Confirmations = 10
		app.Config.Ethereum.ChainID = "31337"

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		tx := &sdk.TxResponse{}

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
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		success := x.HandleBurn(nil)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.True(t, success)
	})

	t.Run("Validation failure and update failed", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), errors.New("error")).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Validation successful and mint confirmed and signing failed", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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
				// "signers":       []string{x.privateKey.PublicKey().RawString()},
				"status": models.StatusConfirmed,
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 10
		app.Config.Ethereum.ChainID = "31337"

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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(burn)

		assert.True(t, success)
	})

}

/*
func TestBurnSignerSyncInvalidMints(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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

		tx := &sdk.TxResponse{
			Tx: "abcd",
			TxResult: sdk.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: sdk.StdTx{
				Msg: sdk.Msg{
					Type: "pos/Send",
					Value: sdk.Value{
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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

		tx := &sdk.TxResponse{
			Tx: "abcd",
			TxResult: sdk.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: sdk.StdTx{
				Msg: sdk.Msg{
					Type: "pos/Send",
					Value: sdk.Value{
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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

	mockContract := ethMocks.NewMockWrappedPocketContract(t)
	mockEthClient := ethMocks.NewMockEthereumClient(t)
	mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
	mockDB := appMocks.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestBurnSigner(t, mockContract, mockEthClient, mockPoktClient)

	mockPoktClient.EXPECT().GetLatestBlockHeight().Return(200, nil)
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
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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

		tx := &sdk.TxResponse{
			Tx: "abcd",
			TxResult: sdk.TxResult{
				Code:        0,
				MessageType: "send",
			},
			StdTx: sdk.StdTx{
				Msg: sdk.Msg{
					Type: "pos/Send",
					Value: sdk.Value{
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
				// "signers":       []string{x.privateKey.PublicKey().RawString()},
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

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
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
			// "signers":       bson.M{"$nin": []string{strings.ToLower(x.privateKey.PublicKey().RawString())}},
		}

		app.Config.Pocket.Confirmations = 0

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		app.Config.Pocket.Mnemonic = ""
		app.Config.Ethereum.RPCURL = ""

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnSigner(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid Multisig keys", func(t *testing.T) {

		app.Config.BurnSigner.Enabled = true
		app.Config.Pocket.Mnemonic = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
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
		app.Config.Pocket.Mnemonic = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = ""
		app.Config.Pocket.MultisigAddress = ""
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
		app.Config.Pocket.Mnemonic = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = ""
		app.Config.Pocket.MultisigAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

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
		app.Config.Pocket.Mnemonic = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"
		app.Config.Pocket.MultisigAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

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
		app.Config.Pocket.Mnemonic = "8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82"
		app.Config.BurnSigner.IntervalMillis = 1
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"
		app.Config.Pocket.MultisigAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

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

*/
