package cosmos

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	appMocks "github.com/dan13ram/wpokt-validator/app/mocks"
	cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
	cosmosMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	// eth "github.com/dan13ram/wpokt-validator/eth/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	ethMocks "github.com/dan13ram/wpokt-validator/eth/client/mocks"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	mockCosmosClient *cosmosMocks.MockCosmosClient,
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
		cosmosClient:           mockCosmosClient,
		cosmosHeight:           0,
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
	mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
	x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "0")
	assert.Equal(t, status.PoktHeight, "0")
}

func TestBurnSignerUpdateBlocks(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mockCosmosClient.EXPECT().GetLatestBlockHeight().Return(200, nil)
		mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

		x.UpdateBlocks()

		assert.Equal(t, x.cosmosHeight, int64(200))
		assert.Equal(t, x.ethBlockNumber, int64(200))
	})

	t.Run("With Error in Pokt Client", func(t *testing.T) {
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mockCosmosClient.EXPECT().GetLatestBlockHeight().Return(200, errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.cosmosHeight, int64(0))
		assert.Equal(t, x.ethBlockNumber, int64(0))
	})

	t.Run("With Error in Eth Client", func(t *testing.T) {
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mockCosmosClient.EXPECT().GetLatestBlockHeight().Return(200, nil)
		mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.cosmosHeight, int64(200))
		assert.Equal(t, x.ethBlockNumber, int64(0))
	})

}

func TestBurnSignerValidateInvalidMint(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mint := &models.InvalidMint{}

		mockCosmosClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid transaction code", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mint := &models.InvalidMint{}

		txResponse := &sdk.TxResponse{}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg from address", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mint := &models.InvalidMint{}

		txResponse := &sdk.TxResponse{}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Invalid transaction msg amount", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
		}

		txResponse := &sdk.TxResponse{}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("invalid", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Memo mismatch", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
		}

		app.Config.Ethereum.ChainID = "31337"

		txResponse := &sdk.TxResponse{}
		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Memo is a valid mint memo", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "100",
			Memo:          fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
		}

		app.Config.Ethereum.ChainID = "31337"

		txResponse := &sdk.TxResponse{}
		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: mint.Memo}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   false,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		valid, err := x.ValidateInvalidMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Successful case", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          `invalid`,
		}

		app.Config.Ethereum.ChainID = "31337"

		txResponse := &sdk.TxResponse{}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		valid, err := x.ValidateInvalidMint(mint)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestBurnSignerValidateBurn(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		burn := &models.Burn{}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(nil, errors.New("error"))

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Invalid log index", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		burn := &models.Burn{
			LogIndex: "0",
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{
				{},
			},
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(nil, errors.New("error"))

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Amount mismatch", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Amount mismatch", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Sender mismatch", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Recipient mismatch", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())

		burn := &models.Burn{
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
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
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Successful case", func(t *testing.T) {

		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())

		burn := &models.Burn{
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		valid, err := x.ValidateBurn(burn)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestBurnSignerHandleInvalidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		success := x.HandleInvalidMint(nil)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		app.Config.Pocket.Confirmations = 1

		invalidMint := &models.InvalidMint{
			Confirmations: "invalid",
			Height:        "invalid",
		}

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Error validating", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100
		app.Config.Pocket.Confirmations = 0

		invalidMint := &models.InvalidMint{
			Confirmations: "1",
			Height:        "99",
		}

		mockCosmosClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation failure and update successful", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100
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

		txResponse := &sdk.TxResponse{}

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

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

	t.Run("Validation failure and update failed", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100
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

		txResponse := &sdk.TxResponse{}

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

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), errors.New("error")).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation successful and invalid mint confirmed and signing failed", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Pocket.CoinDenom = "upokt"
		app.Config.Ethereum.ChainID = "31337"

		invalidMint := &models.InvalidMint{
			SenderAddress: "abcd",
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		txResponse := &sdk.TxResponse{}

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		success := x.HandleInvalidMint(invalidMint)

		assert.False(t, success)
	})

	t.Run("Validation successful and invalid mint confirmed and update successful", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100
		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000
		app.Config.Pocket.CoinDenom = "upokt"
		app.Config.Pocket.Bech32Prefix = "pokt"

		address, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x1245").Bytes())

		invalidMint := &models.InvalidMint{
			SenderAddress: address,
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		txResponse := &sdk.TxResponse{}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		sequence := uint64(2)
		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: address,
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock(mock.Anything).Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		// mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		oldLockWriteSequence := LockWriteSequence
		LockWriteSequence = func() (string, error) {
			return "lockId", nil
		}
		defer func() { LockWriteSequence = oldLockWriteSequence }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

	t.Run("Validation successful and invalid mint pending and update successful", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100
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

		txResponse := &sdk.TxResponse{}

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

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: "abcd",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		success := x.HandleInvalidMint(invalidMint)

		assert.True(t, success)
	})

}

func TestBurnSignerHandleBurn(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		success := x.HandleBurn(nil)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		app.Config.Ethereum.Confirmations = 1

		burn := &models.Burn{
			Confirmations: "invalid",
			BlockNumber:   "invalid",
		}

		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Error validating", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

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

	t.Run("Validation successful and burn confirmed and signing failed", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())

		burn := &models.Burn{
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, errors.New("error")
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock("lock-id").Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		sequence := uint64(2)

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		success := x.HandleBurn(burn)

		assert.False(t, success)
	})

	t.Run("Validation successful and burn confirmed and update successful", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())

		burn := &models.Burn{
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock(mock.Anything).Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		// mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		oldLockWriteSequence := LockWriteSequence
		LockWriteSequence = func() (string, error) {
			return "lockId", nil
		}
		defer func() { LockWriteSequence = oldLockWriteSequence }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		sequence := uint64(2)

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
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

	t.Run("Validation successful and burn pending and update successful", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 10
		app.Config.Ethereum.ChainID = "31337"
		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())

		burn := &models.Burn{
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

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

func TestBurnSignerSyncInvalidMints(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filter := bson.M{
			"$and": []bson.M{
				{
					"vault_address": x.vaultAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filter, mock.Anything).Return(nil)

		success := x.SyncInvalidMints()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"vault_address": x.vaultAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
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
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100

		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000
		app.Config.Pocket.CoinDenom = "upokt"
		app.Config.Pocket.Bech32Prefix = "pokt"

		address, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x1245").Bytes())

		id := primitive.NewObjectID()
		invalidMint := &models.InvalidMint{
			Id:            &id,
			SenderAddress: address,
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		txResponse := &sdk.TxResponse{}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		sequence := uint64(2)
		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: address,
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock("lock-id").Return(nil).Once()

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"vault_address": x.vaultAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}
		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*invalidMint,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		x.cosmosHeight = 100

		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000
		app.Config.Pocket.CoinDenom = "upokt"
		app.Config.Pocket.Bech32Prefix = "pokt"

		address, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x1245").Bytes())

		id := primitive.NewObjectID()
		invalidMint := &models.InvalidMint{
			Id:            &id,
			SenderAddress: address,
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		txResponse := &sdk.TxResponse{}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		sequence := uint64(2)
		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: address,
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock(mock.Anything).Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		oldLockWriteSequence := LockWriteSequence
		LockWriteSequence = func() (string, error) {
			return "lockId", nil
		}
		defer func() { LockWriteSequence = oldLockWriteSequence }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"vault_address": x.vaultAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}
		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*invalidMint,
				}
			})

		success := x.SyncInvalidMints()

		assert.True(t, success)
	})
}

func TestBurnSignerSyncBurns(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filter := bson.M{
			"$and": []bson.M{
				{
					"wpokt_address": x.wpoktAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filter, mock.Anything).Return(nil)

		success := x.SyncBurns()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"wpokt_address": x.wpoktAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
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
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"wpokt_address": x.wpoktAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())
		burn := &models.Burn{
			Id:               &primitive.NilObjectID,
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock("lock-id").Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		sequence := uint64(2)

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*burn,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
		mockMintController := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"wpokt_address": x.wpoktAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())
		burn := &models.Burn{
			Id:               &primitive.NilObjectID,
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock("lock-id").Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		sequence := uint64(2)

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*burn,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncBurns()

		assert.True(t, success)
	})

}

func TestBurnSignerRun(t *testing.T) {

	mockDB := appMocks.NewMockDatabase(t)
	app.DB = mockDB
	mockWPOKT := ethMocks.NewMockWrappedPocketContract(t)
	mockMintController := ethMocks.NewMockMintControllerContract(t)
	mockEthClient := ethMocks.NewMockEthereumClient(t)
	mockCosmosClient := cosmosMocks.NewMockCosmosClient(t)
	x := NewTestBurnSigner(t, mockWPOKT, mockMintController, mockEthClient, mockCosmosClient)

	mockCosmosClient.EXPECT().GetLatestBlockHeight().Return(200, nil)
	mockEthClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

	{
		x.cosmosHeight = 100

		app.Config.Pocket.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000
		app.Config.Pocket.CoinDenom = "upokt"
		app.Config.Pocket.Bech32Prefix = "pokt"

		address, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x1245").Bytes())

		id := primitive.NewObjectID()
		invalidMint := &models.InvalidMint{
			Id:            &id,
			SenderAddress: address,
			Amount:        "20000",
			Memo:          "invalid",
			Confirmations: "1",
			Height:        "99",
			Status:        models.StatusPending,
		}

		txResponse := &sdk.TxResponse{}

		filter := bson.M{
			"_id":    invalidMint.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		sequence := uint64(2)
		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockCosmosClient.EXPECT().GetTx("").Return(txResponse, nil)
		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "0xtxhash",
			Amount:        sdk.NewCoin("upokt", math.NewInt(20000)),
			SenderAddress: address,
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock("lock-id").Return(nil).Once()

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"vault_address": x.vaultAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}
		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*invalidMint,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().Unlock("lockId").Return(nil)

	}

	{

		addressHex, _ := common.AddressHexFromBytes(x.signer.Signer.CosmosPublicKey().Address().Bytes())
		filterFind := bson.M{
			"$and": []bson.M{
				{
					"wpokt_address": x.wpoktAddress,
				},
				{"$or": []bson.M{
					{"status": models.StatusPending},
					{"status": models.StatusConfirmed},
				}},
				{"$nor": []bson.M{
					{"signatures": bson.M{
						"$elemMatch": bson.M{"signer": addressHex},
					}},
				}},
			},
		}

		recipient, _ := common.Bech32FromBytes("pokt", common.HexToAddress("0x2345").Bytes())
		burn := &models.Burn{
			Id:               &primitive.NilObjectID,
			Confirmations:    "1",
			BlockNumber:      "99",
			Status:           models.StatusPending,
			LogIndex:         "0",
			Amount:           "20000",
			SenderAddress:    common.HexToAddress("0x1234").Hex(),
			RecipientAddress: recipient,
		}

		x.ethBlockNumber = 100
		app.Config.Ethereum.Confirmations = 0
		app.Config.Ethereum.ChainID = "31337"
		app.Config.Pocket.ChainID = "testnet"
		app.Config.Pocket.TxFee = 10000

		txReceipt := &types.Receipt{
			Logs: []*types.Log{{}},
		}

		event := &autogen.WrappedPocketBurnAndBridge{
			Amount:      big.NewInt(20000),
			From:        common.HexToAddress("0x1234"),
			PoktAddress: common.HexToAddress("0x2345"),
		}

		mockEthClient.EXPECT().GetTransactionReceipt("").Return(txReceipt, nil)
		mockWPOKT.EXPECT().ParseBurnAndBridge(mock.Anything).Return(event, nil)

		oldCosmosSignTx := CosmosSignTx
		CosmosSignTx = func(
			signerKey common.Signer,
			config models.CosmosConfig,
			client cosmos.CosmosClient,
			sequence uint64,
			signatures []models.Signature,
			transactionBody string,
			toAddress []byte,
			amount sdk.Coin,
			memo string,
		) (string, []models.Signature, error) {
			return "encoded tx", []models.Signature{{}}, nil
		}
		defer func() { CosmosSignTx = oldCosmosSignTx }()

		mockDB.EXPECT().Unlock("lock-id").Return(nil)

		oldLockReadSequences := LockReadSequences
		LockReadSequences = func() (string, error) {
			return "lock-id", nil
		}
		defer func() { LockReadSequences = oldLockReadSequences }()

		oldFindMaxSequence := FindMaxSequence
		FindMaxSequence = func() (*uint64, error) {
			return nil, nil
		}
		defer func() { FindMaxSequence = oldFindMaxSequence }()

		sequence := uint64(2)

		mockCosmosClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: sequence}, nil)

		filter := bson.M{
			"_id":    burn.Id,
			"status": bson.M{"$in": []string{models.StatusPending, models.StatusConfirmed}},
		}

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"return_transaction_body": "encoded tx",
				"signatures":              []models.Signature{{}},
				"sequence":                &sequence,
				"updated_at":              time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*burn,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		mockDB.EXPECT().Unlock("lockId").Return(nil)
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
		app.Config.Pocket.MultisigThreshold = 2

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

}
