package pokt

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestBurnExecutor(t *testing.T, mockClient *pokt.MockPocketClient) *BurnExecutorRunner {
	x := &BurnExecutorRunner{
		vaultAddress: "vaultaddress",
		wpoktAddress: "wpoktaddress",
		client:       mockClient,
	}
	return x
}

func TestBurnExecutorStatus(t *testing.T) {
	mockClient := pokt.NewMockPocketClient(t)
	x := NewTestBurnExecutor(t, mockClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "")
}

func TestBurnExecutorHandleInvalidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		success := x.HandleInvalidMint(nil)

		assert.False(t, success)
	})

	t.Run("Invalid status", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{}

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error submitting signed transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSigned,
		}

		p := rpc.SendRawTxParams{
			Addr:        x.vaultAddress,
			RawHexBytes: doc.ReturnTx,
		}

		mockClient.EXPECT().SubmitRawTx(p).Return(nil, errors.New("error"))

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error while submitting signed transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSigned,
		}

		p := rpc.SendRawTxParams{
			Addr:        x.vaultAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res := &pokt.SubmitRawTxResponse{
			TransactionHash: "hash",
		}

		mockClient.EXPECT().SubmitRawTx(p).Return(res, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(errors.New("error"))

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error while submitting signed transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSigned,
		}

		p := rpc.SendRawTxParams{
			Addr:        x.vaultAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res := &pokt.SubmitRawTxResponse{
			TransactionHash: "hash",
		}

		mockClient.EXPECT().SubmitRawTx(p).Return(res, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		update := bson.M{
			"$set": bson.M{
				"status":         models.StatusSubmitted,
				"return_tx_hash": res.TransactionHash,
				"updated_at":     time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(doc)

		assert.True(t, success)
	})

	t.Run("Error fetching submitted transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		mockClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 10,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":         models.StatusConfirmed,
				"updated_at":     time.Now(),
				"return_tx_hash": "",
				"return_tx":      "",
				"signers":        []string{},
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(doc)

		assert.True(t, success)
	})

	t.Run("Submitted transaction successful but update failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(errors.New("error")).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(doc)

		assert.True(t, success)
	})

}

func TestBurnExecutorHandleBurn(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		success := x.HandleBurn(nil)

		assert.False(t, success)
	})

	t.Run("Invalid status", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{}

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error submitting signed transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSigned,
		}

		p := rpc.SendRawTxParams{
			Addr:        x.vaultAddress,
			RawHexBytes: doc.ReturnTx,
		}

		mockClient.EXPECT().SubmitRawTx(p).Return(nil, errors.New("error"))

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error while submitting signed transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSigned,
		}

		p := rpc.SendRawTxParams{
			Addr:        x.vaultAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res := &pokt.SubmitRawTxResponse{
			TransactionHash: "hash",
		}

		mockClient.EXPECT().SubmitRawTx(p).Return(res, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(errors.New("error"))

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error while submitting signed transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSigned,
		}

		p := rpc.SendRawTxParams{
			Addr:        x.vaultAddress,
			RawHexBytes: doc.ReturnTx,
		}

		res := &pokt.SubmitRawTxResponse{
			TransactionHash: "hash",
		}

		mockClient.EXPECT().SubmitRawTx(p).Return(res, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		update := bson.M{
			"$set": bson.M{
				"status":         models.StatusSubmitted,
				"return_tx_hash": res.TransactionHash,
				"updated_at":     time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(doc)

		assert.True(t, success)
	})

	t.Run("Error fetching submitted transaction", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		mockClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 10,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":         models.StatusConfirmed,
				"updated_at":     time.Now(),
				"return_tx_hash": "",
				"return_tx":      "",
				"signers":        []string{},
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(doc)

		assert.True(t, success)
	})

	t.Run("Submitted transaction successful but update failed", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(errors.New("error")).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction successful", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(doc)

		assert.True(t, success)
	})

}

func TestBurnExecutorSyncInvalidMints(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filter := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"vault_address": x.vaultAddress,
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filter, mock.Anything).Return(nil)

		success := x.SyncInvalidMints()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"vault_address": x.vaultAddress,
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
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"vault_address": x.vaultAddress,
		}

		doc := &models.InvalidMint{
			Id:     &primitive.NilObjectID,
			Status: models.StatusSubmitted,
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*doc,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filterUpdate := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"vault_address": x.vaultAddress,
		}

		doc := &models.InvalidMint{
			Id:     &primitive.NilObjectID,
			Status: models.StatusSubmitted,
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*doc,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filterUpdate := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncInvalidMints()

		assert.True(t, success)
	})

}

func TestBurnExecutorSyncBurns(t *testing.T) {

	t.Run("Error finding", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filter := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"wpokt_address": x.wpoktAddress,
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filter, mock.Anything).Return(nil)

		success := x.SyncBurns()

		assert.True(t, success)
	})

	t.Run("Error locking", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"wpokt_address": x.wpoktAddress,
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
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"wpokt_address": x.wpoktAddress,
		}

		doc := &models.Burn{
			Id:     &primitive.NilObjectID,
			Status: models.StatusSubmitted,
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*doc,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filterUpdate := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockClient := pokt.NewMockPocketClient(t)
		mockDB := app.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"wpokt_address": x.wpoktAddress,
		}

		doc := &models.Burn{
			Id:     &primitive.NilObjectID,
			Status: models.StatusSubmitted,
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*doc,
				}
			})

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filterUpdate := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil)

		success := x.SyncBurns()

		assert.True(t, success)
	})

}

func TestBurnExecutorRun(t *testing.T) {

	mockClient := pokt.NewMockPocketClient(t)
	mockDB := app.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestBurnExecutor(t, mockClient)

	{
		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"vault_address": x.vaultAddress,
		}

		doc := &models.InvalidMint{
			Id:     &primitive.NilObjectID,
			Status: models.StatusSubmitted,
		}

		mockDB.EXPECT().FindMany(models.CollectionInvalidMints, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.InvalidMint)
				*v = []models.InvalidMint{
					*doc,
				}
			}).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil).Once()

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil).Once()

		filterUpdate := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil).Once()
	}

	{
		filterFind := bson.M{
			"status": bson.M{
				"$in": []string{
					string(models.StatusSigned),
					string(models.StatusSubmitted),
				},
			},
			"wpokt_address": x.wpoktAddress,
		}

		doc := &models.Burn{
			Id:     &primitive.NilObjectID,
			Status: models.StatusSubmitted,
		}

		mockDB.EXPECT().FindMany(models.CollectionBurns, filterFind, mock.Anything).Return(nil).
			Run(func(_ string, _ interface{}, result interface{}) {
				v := result.(*[]models.Burn)
				*v = []models.Burn{
					*doc,
				}
			}).Once()

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil).Once()

		tx := &pokt.TxResponse{
			TxResult: pokt.TxResult{
				Code: 0,
			},
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil).Once()

		filterUpdate := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":     models.StatusSuccess,
				"updated_at": time.Now(),
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(nil).Once()
	}

	x.Run()

}

func TestNewBurnExecutor(t *testing.T) {

	t.Run("Disabled", func(t *testing.T) {

		app.Config.BurnExecutor.Enabled = false

		service := NewBurnExecutor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, app.EmptyServiceName)

	})

	t.Run("Invalid Multisig keys", func(t *testing.T) {

		app.Config.BurnExecutor.Enabled = true
		app.Config.Pocket.MultisigPublicKeys = []string{
			"invalid",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnExecutor(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Invalid Vault Address", func(t *testing.T) {

		app.Config.BurnExecutor.Enabled = true
		app.Config.Pocket.VaultAddress = ""
		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() {
			NewBurnExecutor(&sync.WaitGroup{}, models.ServiceHealth{})
		})
	})

	t.Run("Interval is 0", func(t *testing.T) {

		app.Config.BurnExecutor.Enabled = true
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"
		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		service := NewBurnExecutor(&sync.WaitGroup{}, models.ServiceHealth{})

		assert.Nil(t, service)
	})

	t.Run("Valid", func(t *testing.T) {

		app.Config.BurnExecutor.Enabled = true
		app.Config.Pocket.VaultAddress = "E3BB46007E9BF127FD69B02DD5538848A80CADCE"
		app.Config.BurnExecutor.IntervalMillis = 1

		app.Config.Pocket.MultisigPublicKeys = []string{
			"eb0cf2a891382677f03c1b080ec270c693dda7a4c3ee4bcac259ad47c5fe0743",
			"ec69e25c0f2d79e252c1fe0eb8ae07c3a3d8ff7bd616d736f2ded2e9167488b2",
			"abc364918abe9e3966564f60baf74d7ea1c4f3efe92889de066e617989c54283",
		}

		service := NewBurnExecutor(&sync.WaitGroup{}, models.ServiceHealth{})

		health := service.Health()

		assert.NotNil(t, health)
		assert.Equal(t, health.Name, BurnExecutorName)

	})

}
