package cosmos

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dan13ram/wpokt-validator/app"
	appMocks "github.com/dan13ram/wpokt-validator/app/mocks"
	cosmosMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestMintMonitor(t *testing.T, mockClient *cosmosMocks.MockCosmosClient) *MintMonitorRunner {
	x := &MintMonitorRunner{
		vaultAddress:  "vaultaddress",
		wpoktAddress:  "wpoktaddress",
		startHeight:   0,
		currentHeight: 0,
		client:        mockClient,
		minimumAmount: math.NewInt(10000),
		maximumAmount: math.NewInt(100000),
	}
	app.Config.Pocket.TxFee = 10000
	return x
}

func TestMintMonitorStatus(t *testing.T) {
	mockClient := cosmosMocks.NewMockCosmosClient(t)
	x := NewTestMintMonitor(t, mockClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "0")
}

func TestMintMonitorUpdateCurrentHeight(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintMonitor(t, mockClient)

		mockClient.EXPECT().GetLatestBlockHeight().Return(200, nil)

		x.UpdateCurrentHeight()

		assert.Equal(t, x.currentHeight, int64(200))
	})

	t.Run("With Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintMonitor(t, mockClient)

		mockClient.EXPECT().GetLatestBlockHeight().Return(200, errors.New("error"))

		x.UpdateCurrentHeight()

		assert.Equal(t, x.currentHeight, int64(0))
	})

}

func TestMintMonitorHandleFailedMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		success := x.HandleFailedMint(nil, nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), nil)

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       false,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleFailedMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), mongo.CommandError{Code: 11000})

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       false,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleFailedMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error"))

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       false,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleFailedMint(&sdk.TxResponse{}, result)

		assert.False(t, success)
	})

}

func TestMintMonitorHandleInvalidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		success := x.HandleInvalidMint(nil, nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), nil)

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		success := x.HandleInvalidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), mongo.CommandError{Code: 11000})

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		success := x.HandleInvalidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error"))

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		success := x.HandleInvalidMint(&sdk.TxResponse{}, result)

		assert.False(t, success)
	})

	t.Run("With Mint Disabled", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		app.Config.Pocket.MintDisabled = true

		defer func() {
			app.Config.Pocket.MintDisabled = false
		}()

		mockDB.EXPECT().FindOne(models.CollectionMints, mock.Anything, mock.Anything).Return(errors.New("not found"))

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error"))

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		success := x.HandleInvalidMint(&sdk.TxResponse{}, result)

		assert.False(t, success)
	})

	t.Run("With Mint Disabled and duplicate key error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		app.Config.Pocket.MintDisabled = true

		defer func() {
			app.Config.Pocket.MintDisabled = false
		}()

		mockDB.EXPECT().FindOne(models.CollectionMints, mock.Anything, mock.Anything).Return(nil)

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		success := x.HandleInvalidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

}

func TestMintMonitorHandleValidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		success := x.HandleValidMint(nil, nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(errors.New("not found"))

		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(primitive.NewObjectID(), nil)

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleValidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(errors.New("not found"))

		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(primitive.NewObjectID(), mongo.CommandError{Code: 11000})

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleValidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(errors.New("not found"))

		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error"))

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleValidMint(&sdk.TxResponse{}, result)

		assert.False(t, success)
	})

	t.Run("Mint Disabled", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		app.Config.Pocket.MintDisabled = true

		defer func() {
			app.Config.Pocket.MintDisabled = false
		}()

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleValidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

	t.Run("With Duplicate in Invalid Mints", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(nil)

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		success := x.HandleValidMint(&sdk.TxResponse{}, result)

		assert.True(t, success)
	})

}

func TestMintMonitorInitStartHeight(t *testing.T) {

	t.Run("Last Health Pokt Height is valid", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		lastHealth := models.ServiceHealth{
			PoktHeight: "10",
		}

		x.InitStartHeight(lastHealth)

		assert.Equal(t, int64(x.startHeight), int64(10))
	})

	t.Run("Last Health Pokt Height is invalid", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		lastHealth := models.ServiceHealth{
			PoktHeight: "invalid",
		}

		x.InitStartHeight(lastHealth)

		assert.Equal(t, int64(x.startHeight), int64(0))
	})

}

func TestMintMonitorSyncTxs(t *testing.T) {

	t.Run("Start & Current Height are equal", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 100

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Start Height is greater than Current Height", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 101

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Error fetching account txs", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(nil, errors.New("error"))

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("No account txs found", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*sdk.TxResponse{}

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Invalid tx and insert failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*sdk.TxResponse{
			{},
		}

		result := &util.ValidateTxResult{
			Memo:          models.MintMemo{},
			TxValid:       false,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error")).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("Invalid tx and insert successful", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*sdk.TxResponse{
			{},
		}
		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       false,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Failed tx and insert successful", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*sdk.TxResponse{
			{},
		}
		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       false,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("invalid memo and insert failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*sdk.TxResponse{{}}

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error")).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("invalid memo and insert successful", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*sdk.TxResponse{
			{},
		}

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{},

			TxValid:       true,
			Tx:            &tx.Tx{Body: &tx.TxBody{Memo: "invalid"}},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   true,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("valid memo and insert failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		app.Config.Ethereum.ChainID = "31337"

		txs := []*sdk.TxResponse{
			{},
		}

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(errors.New("not found"))
		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(primitive.NewObjectID(), errors.New("error")).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.Mint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("valid memo and insert successful", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		app.Config.Ethereum.ChainID = "31337"

		txs := []*sdk.TxResponse{
			{},
		}

		result := &util.ValidateTxResult{
			Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

			TxValid:       true,
			Tx:            &tx.Tx{},
			TxHash:        "abcd",
			Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
			SenderAddress: "sender",
			NeedsRefund:   false,
		}

		oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
		utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
			return result
		}
		defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

		mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
		mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(errors.New("not found"))
		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.Mint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

}

func TestMintMonitorRun(t *testing.T) {

	mockClient := cosmosMocks.NewMockCosmosClient(t)
	mockDB := appMocks.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestMintMonitor(t, mockClient)
	x.currentHeight = 100
	x.startHeight = 1

	app.Config.Ethereum.ChainID = "31337"

	mockClient.EXPECT().GetLatestBlockHeight().Return(200, nil).Once()

	txs := []*sdk.TxResponse{
		{},
	}

	result := &util.ValidateTxResult{
		Memo: models.MintMemo{ChainID: "31337", Address: "0x1c"},

		TxValid:       true,
		Tx:            &tx.Tx{},
		TxHash:        "abcd",
		Amount:        sdk.NewCoin("pokt", math.NewInt(10000)),
		SenderAddress: "sender",
		NeedsRefund:   false,
	}

	oldValidateTxToCosmosMultisig := utilValidateTxToCosmosMultisig
	utilValidateTxToCosmosMultisig = func(txResponse *sdk.TxResponse, config models.CosmosConfig, minAmount math.Int, maxAmount math.Int) *util.ValidateTxResult {
		return result
	}
	defer func() { utilValidateTxToCosmosMultisig = oldValidateTxToCosmosMultisig }()

	mockClient.EXPECT().GetTxsSentToAddressAfterHeight(x.vaultAddress, uint64(x.startHeight)).Return(txs, nil)
	mockDB.EXPECT().FindOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(errors.New("not found"))
	mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(primitive.NewObjectID(), nil).
		Run(func(_ string, doc interface{}) {
			assert.Equal(t, doc.(models.Mint).Status, models.StatusPending)
		})

	x.Run()

}

func TestNewMintMonitor(t *testing.T) {

	t.Run("Disabled", func(t *testing.T) {

		app.Config.MintMonitor.Enabled = false

		service := NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, app.EmptyServiceName)

	})

	t.Run("Invalid Multisig keys", func(t *testing.T) {

		app.Config.MintMonitor.Enabled = true
		app.Config.Ethereum.RPCURL = ""
		app.Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		app.Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		app.Config.Pocket.MultisigPublicKeys = []string{
			"invalid",
			"0223aa679d6d5344e201e0df9f02ab15a84726eee0dfb4e953c46a9e2cb52349dc",
			"02faaaf0f385bb17381f36dcd86ab2486e8ff8d93440436496665ac007953076c2",
			"02cae233806460db75a941a269490ca5165a620b43241edb8bc72e169f4143a6df",
		}
		app.Config.Pocket.MultisigAddress = "pokt10r5n6x28p9qntchsmhxd4ftq9lk6vzcx3dv4gx"
		app.Config.Pocket.MultisigThreshold = 2
		app.Config.Pocket.Bech32Prefix = "pokt"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid Vault Address", func(t *testing.T) {

		app.Config.MintMonitor.Enabled = true
		app.Config.Ethereum.RPCURL = ""
		app.Config.Pocket.MultisigAddress = ""
		app.Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		app.Config.Pocket.Mnemonic = "test test test test test test test test test test test junk"
		app.Config.Pocket.MultisigPublicKeys = []string{
			"0223aa679d6d5344e201e0df9f02ab15a84726eee0dfb4e953c46a9e2cb52349dc",
			"02faaaf0f385bb17381f36dcd86ab2486e8ff8d93440436496665ac007953076c2",
			"02cae233806460db75a941a269490ca5165a620b43241edb8bc72e169f4143a6df",
		}
		app.Config.Pocket.MultisigThreshold = 2
		app.Config.Pocket.Bech32Prefix = "pokt"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

}
