package eth

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"testing"

	"github.com/dan13ram/wpokt-validator/app"
	appMocks "github.com/dan13ram/wpokt-validator/app/mocks"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	ethMocks "github.com/dan13ram/wpokt-validator/eth/client/mocks"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestBurnMonitor(t *testing.T, mockContract *ethMocks.MockWrappedPocketContract, mockClient *ethMocks.MockEthereumClient) *BurnMonitorRunner {
	x := &BurnMonitorRunner{
		startBlockNumber:   0,
		currentBlockNumber: 100,
		wpoktContract:      mockContract,
		client:             mockClient,
		minimumAmount:      big.NewInt(10000),
	}
	return x
}

func TestBurnMonitorStatus(t *testing.T) {
	mockContract := ethMocks.NewMockWrappedPocketContract(t)
	mockClient := ethMocks.NewMockEthereumClient(t)
	x := NewTestBurnMonitor(t, mockContract, mockClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "0")
	assert.Equal(t, status.PoktHeight, "")
}

func TestBurnMonitorUpdateCurrentBlockNumber(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		mockClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

		x.UpdateCurrentBlockNumber()

		assert.Equal(t, x.currentBlockNumber, int64(200))
	})

	t.Run("With Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		mockClient.EXPECT().GetBlockNumber().Return(uint64(200), errors.New("error"))

		x.UpdateCurrentBlockNumber()

		assert.Equal(t, x.currentBlockNumber, int64(100))
	})

}

func TestBurnMonitorHandleBurnEvent(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		success := x.HandleBurnEvent(nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil)

		success := x.HandleBurnEvent(&autogen.WrappedPocketBurnAndBridge{})

		assert.True(t, success)
	})

	t.Run("With Duplicate Key Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), mongo.CommandError{Code: 11000})

		success := x.HandleBurnEvent(&autogen.WrappedPocketBurnAndBridge{})

		assert.True(t, success)
	})

	t.Run("With Other Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), errors.New("error"))

		success := x.HandleBurnEvent(&autogen.WrappedPocketBurnAndBridge{})

		assert.False(t, success)
	})

}

func TestBurnMonitorInitStartBlockNumber(t *testing.T) {

	t.Run("Last Health Eth Block Number is valid", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		lastHealth := models.ServiceHealth{
			EthBlockNumber: "10",
		}

		x.InitStartBlockNumber(lastHealth)

		assert.Equal(t, x.startBlockNumber, int64(10))
	})

	t.Run("Last Health Eth Block Number is invalid", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)

		lastHealth := models.ServiceHealth{
			EthBlockNumber: "invalid",
		}

		x.InitStartBlockNumber(lastHealth)

		assert.Equal(t, x.startBlockNumber, int64(100))
	})

}

func TestBurnMonitorSyncBlocks(t *testing.T) {

	t.Run("Successful Case", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
		})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, amount []*big.Int, to []common.Address, from []common.Address) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(100))
			}).Once()
		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil).Once()

		success := x.SyncBlocks(1, 100)
		assert.True(t, success)
	})

	t.Run("Amount less than minimumAmount", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20),
		})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, amount []*big.Int, to []common.Address, from []common.Address) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(100))
			}).Once()

		success := x.SyncBlocks(1, 100)
		assert.True(t, success)
	})

	t.Run("Error in Filtering", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(nil, errors.New("some error")).Once()

		success := x.SyncBlocks(1, 100)
		assert.False(t, success)
	})

	t.Run("Error in Handling Events", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(nil)
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})

	t.Run("Some events were removed", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
		}).Once()
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Raw: types.Log{Removed: true},
		}).Once()
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
		}).Once()
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Times(3)
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).Once()

		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil).Times(2)

		assert.True(t, x.SyncBlocks(1, 100))
	})

	t.Run("Error in Handling First Event", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(nil).Once()
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
		}).Once()
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Times(2)
		mockFilter.EXPECT().Next().Return(false).Once()
		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})

	t.Run("Error During Filtering Iteration", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Error().Return(errors.New("iteration error"))
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})

	t.Run("Error After Filtering Iteration", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(nil)
		mockFilter.EXPECT().Error().Return(nil).Once()
		mockFilter.EXPECT().Error().Return(errors.New("iteration error")).Once()
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})
}

func TestBurnMonitorSyncTxs(t *testing.T) {

	t.Run("Start & Current Block Number are equal", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)
		x.currentBlockNumber = 100
		x.startBlockNumber = 100

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Start Block Number is greater than Current Block Number", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnMonitor(t, mockContract, mockClient)
		x.currentBlockNumber = 100
		x.startBlockNumber = 101

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Start Block Number is less than Current Block Number but diff is less than MAX_QUERY_BLOCKS", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
		})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		x.startBlockNumber = 1
		final := eth.MAX_QUERY_BLOCKS / 2
		x.currentBlockNumber = final

		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, amount []*big.Int, to []common.Address, from []common.Address) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(final))
			}).Once()
		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil).Once()

		success := x.SyncTxs()

		assert.True(t, success)

		assert.Equal(t, x.currentBlockNumber, x.startBlockNumber)
		assert.Equal(t, x.startBlockNumber, int64(final))
	})

	t.Run("Start Block Number is less than Current Block Number but diff is greater than MAX_QUERY_BLOCKS", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
			Amount: big.NewInt(20000),
		})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestBurnMonitor(t, mockContract, mockClient)
		x.startBlockNumber = 1
		x.currentBlockNumber = eth.MAX_QUERY_BLOCKS + 2

		mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
			Return(mockFilter, nil).Times(2)
		mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil)

		success := x.SyncTxs()

		assert.True(t, success)

		assert.Equal(t, x.currentBlockNumber, x.startBlockNumber)
	})

}

func TestNewBurnMonitor(t *testing.T) {

	t.Run("Disabled", func(t *testing.T) {

		app.Config.BurnMonitor.Enabled = false

		service := NewBurnMonitor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, app.EmptyServiceName)

	})

	t.Run("Invalid ETH RPC", func(t *testing.T) {

		app.Config.BurnMonitor.Enabled = true
		app.Config.Ethereum.RPCURL = ""

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnMonitor(&sync.WaitGroup{}, models.ServiceHealth{})
		})

	})

	t.Run("Interval is 0", func(t *testing.T) {

		app.Config.BurnMonitor.Enabled = true
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"

		service := NewBurnMonitor(&sync.WaitGroup{}, models.ServiceHealth{})

		assert.Nil(t, service)
	})

	t.Run("Valid", func(t *testing.T) {

		app.Config.BurnMonitor.Enabled = true
		app.Config.BurnMonitor.IntervalMillis = 1
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"

		service := NewBurnMonitor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, BurnMonitorName)

	})

}

func TestBurnMonitorRun(t *testing.T) {

	mockContract := ethMocks.NewMockWrappedPocketContract(t)
	mockClient := ethMocks.NewMockEthereumClient(t)
	mockDB := appMocks.NewMockDatabase(t)
	mockFilter := ethMocks.NewMockWrappedPocketBurnAndBridgeIterator(t)
	mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketBurnAndBridge{
		Amount: big.NewInt(20000),
	})
	mockFilter.EXPECT().Error().Return(nil)
	mockFilter.EXPECT().Close().Return(nil)
	mockFilter.EXPECT().Next().Return(true).Once()
	mockFilter.EXPECT().Next().Return(false).Once()

	app.DB = mockDB
	x := NewTestBurnMonitor(t, mockContract, mockClient)
	x.currentBlockNumber = 100
	x.startBlockNumber = 1

	mockClient.EXPECT().GetBlockNumber().Return(uint64(100), nil)
	mockContract.EXPECT().FilterBurnAndBridge(mock.Anything, []*big.Int{}, []common.Address{}, []common.Address{}).
		Return(mockFilter, nil).
		Run(func(opts *bind.FilterOpts, amount []*big.Int, to []common.Address, from []common.Address) {
			assert.Equal(t, opts.Start, uint64(1))
			assert.Equal(t, *opts.End, uint64(100))
		}).Once()
	mockDB.EXPECT().InsertOne(models.CollectionBurns, mock.Anything).Return(primitive.NewObjectID(), nil).Once()

	x.Run()

}
