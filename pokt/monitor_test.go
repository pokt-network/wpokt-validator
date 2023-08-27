package pokt

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestMintMonitor(t *testing.T, mockClient *pokt.MockPocketClient) *MintMonitorRunner {
	x := &MintMonitorRunner{
		vaultAddress:  "vaultaddress",
		wpoktAddress:  "wpoktaddress",
		startHeight:   0,
		currentHeight: 0,
		client:        mockClient,
	}
	return x
}

func TestMintMonitorStatus(t *testing.T) {
	mockClient := pokt.NewMockPocketClient(t)
	x := NewTestMintMonitor(t, mockClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "0")
}

func TestMintMonitorUpdateCurrentHeight(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		x := NewTestMintMonitor(t, mockClient)

		mockClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, nil)

		x.UpdateCurrentHeight()

		assert.Equal(t, x.currentHeight, int64(200))
	})

	t.Run("With Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		x := NewTestMintMonitor(t, mockClient)

		mockClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, errors.New("error"))

		x.UpdateCurrentHeight()

		assert.Equal(t, x.currentHeight, int64(0))
	})

}

func TestMintMonitorHandleFailedMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		success := x.HandleFailedMint(nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil)

		success := x.HandleFailedMint(&pokt.TxResponse{})

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(mongo.CommandError{Code: 11000})

		success := x.HandleFailedMint(&pokt.TxResponse{})

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(errors.New("error"))

		success := x.HandleFailedMint(&pokt.TxResponse{})

		assert.False(t, success)
	})

}

func TestMintMonitorHandleInvalidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		success := x.HandleInvalidMint(nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil)

		success := x.HandleInvalidMint(&pokt.TxResponse{})

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(mongo.CommandError{Code: 11000})

		success := x.HandleInvalidMint(&pokt.TxResponse{})

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(errors.New("error"))

		success := x.HandleInvalidMint(&pokt.TxResponse{})

		assert.False(t, success)
	})

}

func TestMintMonitorHandleValidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		success := x.HandleValidMint(nil, models.MintMemo{})

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(nil)

		success := x.HandleValidMint(&pokt.TxResponse{}, models.MintMemo{})

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(mongo.CommandError{Code: 11000})

		success := x.HandleValidMint(&pokt.TxResponse{}, models.MintMemo{})

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(errors.New("error"))

		success := x.HandleValidMint(&pokt.TxResponse{}, models.MintMemo{})

		assert.False(t, success)
	})

}

func TestMintMonitorInitStartHeight(t *testing.T) {

	t.Run("Last Health Pokt Height is valid", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		lastHealth := models.ServiceHealth{
			PoktHeight: "10",
		}

		x.InitStartHeight(lastHealth)

		assert.Equal(t, x.startHeight, int64(10))
	})

	t.Run("Last Health Pokt Height is invalid", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)

		lastHealth := models.ServiceHealth{
			PoktHeight: "invalid",
		}

		x.InitStartHeight(lastHealth)

		assert.Equal(t, x.startHeight, int64(0))
	})

}

func TestMintMonitorSyncTxs(t *testing.T) {

	t.Run("Start & Current Height are equal", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 100

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Start Height is greater than Current Height", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 101

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Error fetching account txs", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(nil, errors.New("error"))

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("No account txs found", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Invalid tx and insert failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "",
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(errors.New("error")).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("Invalid tx and insert successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "",
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Failed tx and insert successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code: 10,
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Wrong tx recipient and insert successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code:      0,
					Recipient: "random",
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Wrong tx msg type and insert successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code:        0,
					Recipient:   x.vaultAddress,
					MessageType: "invalid",
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("invalid memo and insert failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code:        0,
					Recipient:   x.vaultAddress,
					MessageType: "send",
				},
				StdTx: pokt.StdTx{
					Memo: "invalid",
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(errors.New("error")).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("invalid memo and insert successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code:        0,
					Recipient:   x.vaultAddress,
					MessageType: "send",
				},
				StdTx: pokt.StdTx{
					Memo: "invalid",
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("valid memo and insert failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		app.Config.Ethereum.ChainId = "31337"

		address := common.HexToAddress("0x1c")

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code:        0,
					Recipient:   x.vaultAddress,
					MessageType: "send",
				},
				StdTx: pokt.StdTx{
					Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337"}`, address),
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(errors.New("error")).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.Mint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("valid memo and insert successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintMonitor(t, mockClient)
		x.currentHeight = 100
		x.startHeight = 1

		app.Config.Ethereum.ChainId = "31337"

		address := common.HexToAddress("0x1c")

		txs := []*pokt.TxResponse{
			{
				Tx: "abcd",
				TxResult: pokt.TxResult{
					Code:        0,
					Recipient:   x.vaultAddress,
					MessageType: "send",
				},
				StdTx: pokt.StdTx{
					Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337"}`, address),
				},
			},
		}

		mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil)
		mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(nil).
			Run(func(_ string, doc interface{}) {
				assert.Equal(t, doc.(models.Mint).Status, models.StatusPending)
			})

		success := x.SyncTxs()

		assert.True(t, success)
	})

}

func TestMintMonitorRun(t *testing.T) {

	mockClient := pokt.NewMockPocketClient(t)
	mockDB := app.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestMintMonitor(t, mockClient)
	x.currentHeight = 100
	x.startHeight = 1

	app.Config.Ethereum.ChainId = "31337"

	address := common.HexToAddress("0x1c")

	txs := []*pokt.TxResponse{
		{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				Recipient:   x.vaultAddress,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Memo: "invalid",
			},
		},
		{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				Recipient:   x.vaultAddress,
				MessageType: "send",
			},
			StdTx: pokt.StdTx{
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337"}`, address),
			},
		},
		{
			Tx: "abcd",
			TxResult: pokt.TxResult{
				Code:        0,
				Recipient:   x.vaultAddress,
				MessageType: "invalid",
			},
		},
	}

	mockClient.EXPECT().GetHeight().Return(&pokt.HeightResponse{Height: 200}, nil).Once()
	mockClient.EXPECT().GetAccountTxsByHeight(x.vaultAddress, x.startHeight).Return(txs, nil).Once()
	mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
		Run(func(_ string, doc interface{}) {
			assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusPending)
		}).Once()
	mockDB.EXPECT().InsertOne(models.CollectionMints, mock.Anything).Return(nil).
		Run(func(_ string, doc interface{}) {
			assert.Equal(t, doc.(models.Mint).Status, models.StatusPending)
		}).Once()
	mockDB.EXPECT().InsertOne(models.CollectionInvalidMints, mock.Anything).Return(nil).
		Run(func(_ string, doc interface{}) {
			assert.Equal(t, doc.(models.InvalidMint).Status, models.StatusFailed)
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
		app.Config.Pocket.MultisigPublicKeys = []string{
			"invalid",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid Vault Address", func(t *testing.T) {

		app.Config.MintMonitor.Enabled = true
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
			NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Interval is 0", func(t *testing.T) {

		app.Config.MintMonitor.Enabled = true
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		service := NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})

		assert.Nil(t, service)
	})

	t.Run("Valid", func(t *testing.T) {

		app.Config.MintMonitor.Enabled = true
		app.Config.MintMonitor.IntervalMillis = 1
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"

		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		service := NewMintMonitor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, MintMonitorName)

	})

}
