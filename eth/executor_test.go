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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestMintExecutor(t *testing.T, mockContract *ethMocks.MockWrappedPocketContract, mockClient *ethMocks.MockEthereumClient) *MintExecutorRunner {
	mintControllerAbi, _ := autogen.MintControllerMetaData.GetAbi()
	x := &MintExecutorRunner{
		startBlockNumber:   0,
		currentBlockNumber: 100,
		wpoktContract:      mockContract,
		mintControllerAbi:  mintControllerAbi,
		client:             mockClient,
		vaultAddress:       "vaultAddress",
		wpoktAddress:       "wpoktAddress",
	}
	return x
}

func TestMintExecutorStatus(t *testing.T) {
	mockContract := ethMocks.NewMockWrappedPocketContract(t)
	mockClient := ethMocks.NewMockEthereumClient(t)
	x := NewTestMintExecutor(t, mockContract, mockClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "0")
	assert.Equal(t, status.PoktHeight, "")
}

func TestMintExecutorUpdateCurrentBlockNumber(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		x := NewTestMintExecutor(t, mockContract, mockClient)

		mockClient.EXPECT().GetBlockNumber().Return(uint64(200), nil)

		x.UpdateCurrentBlockNumber()

		assert.Equal(t, x.currentBlockNumber, int64(200))
	})

	t.Run("With Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		x := NewTestMintExecutor(t, mockContract, mockClient)

		mockClient.EXPECT().GetBlockNumber().Return(uint64(200), errors.New("error"))

		x.UpdateCurrentBlockNumber()

		assert.Equal(t, x.currentBlockNumber, int64(100))
	})

}

func TestMintExecutorHandleMintEvent(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintExecutor(t, mockContract, mockClient)

		success := x.HandleMintEvent(nil)

		assert.False(t, success)
	})

	t.Run("No Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintExecutor(t, mockContract, mockClient)

		event := &autogen.WrappedPocketMinted{}

		filter := bson.M{
			"wpokt_address":     x.wpoktAddress,
			"vault_address":     x.vaultAddress,
			"recipient_address": strings.ToLower(event.Recipient.Hex()),
			"amount":            event.Amount.String(),
			"nonce":             event.Nonce.String(),
			"status": bson.M{
				"$in": []string{models.StatusConfirmed, models.StatusSigned, models.StatusSuccess},
			},
		}

		update := bson.M{
			"$set": bson.M{
				"status":       models.StatusSuccess,
				"mint_tx_hash": strings.ToLower(event.Raw.TxHash.String()),
				"updated_at":   time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(_ string, _ interface{}, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleMintEvent(&autogen.WrappedPocketMinted{})

		assert.True(t, success)
	})

	t.Run("With Error", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintExecutor(t, mockContract, mockClient)

		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), errors.New("error"))

		success := x.HandleMintEvent(&autogen.WrappedPocketMinted{})

		assert.False(t, success)
	})

}

func TestMintExecutorInitStartBlockNumber(t *testing.T) {

	t.Run("Last Health Eth Block Number is valid", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintExecutor(t, mockContract, mockClient)

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
		x := NewTestMintExecutor(t, mockContract, mockClient)

		lastHealth := models.ServiceHealth{
			EthBlockNumber: "invalid",
		}

		x.InitStartBlockNumber(lastHealth)

		assert.Equal(t, x.startBlockNumber, int64(100))
	})

}

func TestMintExecutorSyncBlocks(t *testing.T) {

	t.Run("Successful Case", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(100))
			}).Once()
		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncBlocks(1, 100)
		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(100))
			}).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", errors.New("error"))

		success := x.SyncBlocks(1, 100)
		assert.False(t, success)
	})

	t.Run("Error unlocking", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(100))
			}).Once()
		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncBlocks(1, 100)
		assert.False(t, success)
	})

	t.Run("Error in Filtering", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(nil, errors.New("some error")).Once()

		success := x.SyncBlocks(1, 100)
		assert.False(t, success)
	})

	t.Run("Error in Handling Events", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(nil)
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})

	t.Run("Some events were removed", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{}).Once()
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{
			Raw: types.Log{Removed: true},
		}).Once()
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{}).Once()
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Times(3)
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil).Times(2)

		assert.True(t, x.SyncBlocks(1, 100))
	})

	t.Run("Error in Handling First Event", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(nil).Once()
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{}).Once()
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Times(2)
		mockFilter.EXPECT().Next().Return(false).Once()
		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).Once()
		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		assert.False(t, x.SyncBlocks(1, 100))
	})

	t.Run("Error During Filtering Iteration", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Error().Return(errors.New("iteration error"))
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})

	t.Run("Error After Filtering Iteration", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(nil).Once()
		mockFilter.EXPECT().Error().Return(nil).Once()
		mockFilter.EXPECT().Error().Return(errors.New("iteration error")).Once()
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).Once()

		assert.False(t, x.SyncBlocks(1, 100))
	})
}

func TestMintExecutorSyncTxs(t *testing.T) {

	t.Run("Start & Current Block Number are equal", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintExecutor(t, mockContract, mockClient)
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
		x := NewTestMintExecutor(t, mockContract, mockClient)
		x.currentBlockNumber = 100
		x.startBlockNumber = 101

		success := x.SyncTxs()

		assert.True(t, success)
	})

	t.Run("Start Block Number is less than Current Block Number but diff is less than MAX_QUERY_BLOCKS", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		x.startBlockNumber = 1
		final := eth.MAX_QUERY_BLOCKS / 2
		x.currentBlockNumber = final

		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).
			Run(func(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) {
				assert.Equal(t, opts.Start, uint64(1))
				assert.Equal(t, *opts.End, uint64(final))
			}).Once()
		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil).Once()
		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncTxs()

		assert.True(t, success)

		assert.Equal(t, x.currentBlockNumber, x.startBlockNumber)
		assert.Equal(t, x.startBlockNumber, int64(final))
	})

	t.Run("Start Block Number is less than Current Block Number but diff is greater than MAX_QUERY_BLOCKS", func(t *testing.T) {
		mockContract := ethMocks.NewMockWrappedPocketContract(t)
		mockClient := ethMocks.NewMockEthereumClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
		mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{})
		mockFilter.EXPECT().Error().Return(nil)
		mockFilter.EXPECT().Close().Return(nil)
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		mockFilter.EXPECT().Next().Return(true).Once()
		mockFilter.EXPECT().Next().Return(false).Once()
		app.DB = mockDB

		x := NewTestMintExecutor(t, mockContract, mockClient)
		x.currentBlockNumber = 2*eth.MAX_QUERY_BLOCKS - 1
		x.startBlockNumber = 1

		mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
			Return(mockFilter, nil).Times(2)
		mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil)
		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncTxs()

		assert.True(t, success)

		assert.Equal(t, x.currentBlockNumber, x.startBlockNumber)
	})

}

func TestNewMintExecutor(t *testing.T) {

	t.Run("Disabled", func(t *testing.T) {

		app.Config.MintExecutor.Enabled = false

		service := NewMintExecutor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, app.EmptyServiceName)

	})

	t.Run("Invalid ETH RPC", func(t *testing.T) {

		app.Config.MintExecutor.Enabled = true
		app.Config.Ethereum.RPCURL = ""

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewMintExecutor(&sync.WaitGroup{}, models.ServiceHealth{})
		})

	})

	t.Run("Interval is 0", func(t *testing.T) {

		app.Config.MintExecutor.Enabled = true
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"

		service := NewMintExecutor(&sync.WaitGroup{}, models.ServiceHealth{})

		assert.Nil(t, service)
	})

	t.Run("Valid", func(t *testing.T) {

		app.Config.MintExecutor.Enabled = true
		app.Config.MintExecutor.IntervalMillis = 1
		app.Config.Ethereum.RPCURL = "https://eth.llamarpc.com"

		service := NewMintExecutor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, MintExecutorName)

	})

}

func TestMintExecutorRun(t *testing.T) {

	mockContract := ethMocks.NewMockWrappedPocketContract(t)
	mockClient := ethMocks.NewMockEthereumClient(t)
	mockDB := appMocks.NewMockDatabase(t)
	mockFilter := ethMocks.NewMockWrappedPocketMintedIterator(t)
	mockFilter.EXPECT().Event().Return(&autogen.WrappedPocketMinted{})
	mockFilter.EXPECT().Error().Return(nil)
	mockFilter.EXPECT().Close().Return(nil)
	mockFilter.EXPECT().Next().Return(true).Once()
	mockFilter.EXPECT().Next().Return(false).Once()

	app.DB = mockDB
	x := NewTestMintExecutor(t, mockContract, mockClient)
	x.currentBlockNumber = 100
	x.startBlockNumber = 1

	mockClient.EXPECT().GetBlockNumber().Return(uint64(100), nil)
	mockContract.EXPECT().FilterMinted(mock.Anything, []common.Address{}, []*big.Int{}, []*big.Int{}).
		Return(mockFilter, nil).
		Run(func(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) {
			assert.Equal(t, opts.Start, uint64(1))
			assert.Equal(t, *opts.End, uint64(100))
		}).Once()
	mockDB.EXPECT().UpdateOne(models.CollectionMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil).Once()
	mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
	mockDB.EXPECT().Unlock("lockId").Return(nil)

	x.Run()

}
