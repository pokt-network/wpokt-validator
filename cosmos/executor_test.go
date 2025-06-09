package cosmos

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	multisigtypes "github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dan13ram/wpokt-validator/app"
	appMocks "github.com/dan13ram/wpokt-validator/app/mocks"
	"github.com/dan13ram/wpokt-validator/common"
	cosmosMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	"github.com/dan13ram/wpokt-validator/cosmos/util"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

const (
	mnemonic1 = "test test test test test test test test test test test junk"
	mnemonic2 = "all all all all all all all all all all all all"
	mnemonic3 = "flock image pipe glory position until viable price steak market ring tragic"
)

var signer1, _ = common.NewMnemonicSigner(mnemonic1)
var signer2, _ = common.NewMnemonicSigner(mnemonic2)
var signer3, _ = common.NewMnemonicSigner(mnemonic3)

var pubKey1 = signer1.CosmosPublicKey()
var pubKey2 = signer2.CosmosPublicKey()
var pubKey3 = signer3.CosmosPublicKey()

func init() {
	log.SetOutput(io.Discard)
}

func NewTestBurnExecutor(t *testing.T, mockClient *cosmosMocks.MockCosmosClient) *BurnExecutorRunner {

	pubKeyHex1 := hex.EncodeToString(pubKey1.Bytes())
	pubKeyHex2 := hex.EncodeToString(pubKey2.Bytes())
	pubKeyHex3 := hex.EncodeToString(pubKey3.Bytes())

	pks := []crypto.PubKey{
		pubKey1,
		pubKey2,
		pubKey3,
	}

	sort.Slice(pks, func(i, j int) bool {
		return bytes.Compare(pks[i].Address(), pks[j].Address()) < 0
	})

	multisigPk := multisig.NewLegacyAminoPubKey(2, pks)
	multisigAddressBytes := multisigPk.Address().Bytes()
	multisigAddress, _ := common.Bech32FromBytes("pokt", multisigAddressBytes)

	app.Config.Ethereum.PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	app.Config.Pocket.Mnemonic = mnemonic1
	app.Config.Pocket.MultisigPublicKeys = []string{
		pubKeyHex1,
		pubKeyHex2,
		pubKeyHex3,
	}

	app.Config.Pocket.MultisigAddress = multisigAddress
	app.Config.Pocket.MultisigThreshold = 2
	app.Config.Pocket.Bech32Prefix = "pokt"
	app.Config.Pocket.TxFee = 10000

	signer, err := app.GetPocketSignerAndMultisig()
	assert.Nil(t, err)

	x := &BurnExecutorRunner{
		vaultAddress: "vaultaddress",
		wpoktAddress: "wpoktaddress",
		client:       mockClient,
		signer:       signer,
	}
	return x
}

func TestBurnExecutorStatus(t *testing.T) {
	mockClient := cosmosMocks.NewMockCosmosClient(t)
	x := NewTestBurnExecutor(t, mockClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "")
}

func TestBurnExecutorHandleInvalidMint(t *testing.T) {

	t.Run("Nil event", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		success := x.HandleInvalidMint(nil)

		assert.False(t, success)
	})

	t.Run("Error wrapping tx builder", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSigned,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, assert.AnError
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error GetSignaturesV2", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{{}}, assert.AnError)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error not enough sigs", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{{}}, nil)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in GetAccount", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{{}, {}}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, assert.AnError)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in ValidateSig", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return assert.AnError
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in AddSig", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return assert.AnError
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in SetSigs", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(assert.AnError)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in json encoding", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, assert.AnError
		})

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in tx encoding", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, assert.AnError
		})

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in broadcast", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, nil
		})

		txHash := "0xHash"

		mockClient.EXPECT().BroadcastTx(txBytes).Return(txHash, assert.AnError)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Error in update db", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, nil
		})

		txHash := "0xHash"

		mockClient.EXPECT().BroadcastTx(txBytes).Return(txHash, nil)

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), assert.AnError)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("signed tx submitted successfully", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.InvalidMint{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, nil
		})

		txHash := "0xhash"

		mockClient.EXPECT().BroadcastTx(txBytes).Return(txHash, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		date := time.Now()

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusSubmitted,
				"return_transaction_body": string(txJSON),
				"return_transaction_hash": txHash,
				"updated_at":              date,
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Run(func(_collection string, _filter interface{}, _update interface{}) {
			_update.(bson.M)["$set"].(bson.M)["updated_at"] = date
			assert.Equal(t, update, _update)
		}).Return(primitive.NewObjectID(), nil)

		success := x.HandleInvalidMint(doc)

		assert.True(t, success)
	})
	t.Run("Error fetching submitted transaction", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		mockClient.EXPECT().GetTx("").Return(nil, assert.AnError)

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		tx := &sdk.TxResponse{
			Code: 10,
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"updated_at":              time.Now(),
				"return_transaction_hash": "",
				"return_transaction_body": "",
				"signatures":              []models.Signature{},
				"sequence":                nil,
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(doc)

		assert.True(t, success)
	})

	t.Run("Submitted transaction successful but update failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		tx := &sdk.TxResponse{
			Code: 0,
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

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), assert.AnError).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleInvalidMint(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction successful", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.InvalidMint{
			Status: models.StatusSubmitted,
		}

		tx := &sdk.TxResponse{
			Code: 0,
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

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		success := x.HandleBurn(nil)

		assert.False(t, success)
	})

	t.Run("Error wrapping tx builder", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSigned,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, assert.AnError
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error GetSignaturesV2", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{{}}, assert.AnError)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error not enough sigs", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{{}}, nil)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in GetAccount", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}

		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{{}, {}}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, assert.AnError)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in ValidateSig", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return assert.AnError
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in AddSig", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return assert.AnError
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in SetSigs", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(assert.AnError)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in json encoding", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, assert.AnError
		})

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in tx encoding", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, assert.AnError
		})

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in broadcast", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, nil
		})

		txHash := "0xHash"

		mockClient.EXPECT().BroadcastTx(txBytes).Return(txHash, assert.AnError)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Error in update db", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, nil
		})

		txHash := "0xHash"

		mockClient.EXPECT().BroadcastTx(txBytes).Return(txHash, nil)

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, mock.Anything, mock.Anything).Return(primitive.NewObjectID(), assert.AnError)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("signed tx submitted successfully", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		seq := uint64(1)
		doc := &models.Burn{
			Status:   models.StatusSigned,
			Sequence: &seq,
		}

		txBuilder := cosmosMocks.NewMockTxBuilder(t)
		txConfig := cosmosMocks.NewMockTxConfig(t)

		utilWrapTxBuilder = func(string, string) (client.TxBuilder, client.TxConfig, error) {
			return txBuilder, txConfig, nil
		}
		utilValidateSignature = func(models.CosmosConfig, *signingtypes.SignatureV2, uint64, uint64, client.TxConfig, client.TxBuilder,
		) error {
			return nil
		}
		multisigtypesAddSignatureV2 = func(*signingtypes.MultiSignatureData, signingtypes.SignatureV2, []crypto.PubKey) error {
			return nil
		}
		defer func() {
			utilWrapTxBuilder = util.WrapTxBuilder
			utilValidateSignature = util.ValidateSignature
			multisigtypesAddSignatureV2 = multisigtypes.AddSignatureV2
		}()

		tx := cosmosMocks.NewMockTx(t)
		txBuilder.EXPECT().GetTx().Return(tx)
		tx.EXPECT().GetSignaturesV2().Return([]signingtypes.SignatureV2{
			{PubKey: pubKey1},
			{PubKey: pubKey2},
		}, nil)

		mockClient.EXPECT().GetAccount(mock.Anything).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

		txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)

		txJSON := []byte("encoded tx")

		txConfig.EXPECT().TxJSONEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txJSON, nil
		})
		txBytes := []byte("encoded tx as bytes")

		txConfig.EXPECT().TxEncoder().Return(func(tx sdk.Tx) ([]byte, error) {
			return txBytes, nil
		})

		txHash := "0xhash"

		mockClient.EXPECT().BroadcastTx(txBytes).Return(txHash, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSigned,
		}

		date := time.Now()

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusSubmitted,
				"return_transaction_body": string(txJSON),
				"return_transaction_hash": txHash,
				"updated_at":              date,
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Run(func(_collection string, _filter interface{}, _update interface{}) {
			_update.(bson.M)["$set"].(bson.M)["updated_at"] = date
			assert.Equal(t, update, _update)
		}).Return(primitive.NewObjectID(), nil)

		success := x.HandleBurn(doc)

		assert.True(t, success)
	})
	t.Run("Error fetching submitted transaction", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		mockClient.EXPECT().GetTx("").Return(nil, assert.AnError)

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		tx := &sdk.TxResponse{
			Code: 10,
		}

		mockClient.EXPECT().GetTx("").Return(tx, nil)

		filter := bson.M{
			"_id":    doc.Id,
			"status": models.StatusSubmitted,
		}

		update := bson.M{
			"$set": bson.M{
				"status":                  models.StatusConfirmed,
				"updated_at":              time.Now(),
				"return_transaction_hash": "",
				"return_transaction_body": "",
				"signatures":              []models.Signature{},
				"sequence":                nil,
			},
		}

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(doc)

		assert.True(t, success)
	})

	t.Run("Submitted transaction successful but update failed", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		tx := &sdk.TxResponse{
			Code: 0,
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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), assert.AnError).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		success := x.HandleBurn(doc)

		assert.False(t, success)
	})

	t.Run("Submitted transaction successful", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		doc := &models.Burn{
			Status: models.StatusSubmitted,
		}

		tx := &sdk.TxResponse{
			Code: 0,
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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filter, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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

		tx := &sdk.TxResponse{}

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

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncInvalidMints()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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

		tx := &sdk.TxResponse{}

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

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestBurnExecutor(t, mockClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)

	})

	t.Run("Nothing to handle", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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

		tx := &sdk.TxResponse{}

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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
			Run(func(collection string, filter, gotUpdate interface{}) {
				gotUpdate.(bson.M)["$set"].(bson.M)["updated_at"] = update["$set"].(bson.M)["updated_at"]
				assert.Equal(t, update, gotUpdate)
			}).Once()

		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		success := x.SyncBurns()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
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

		tx := &sdk.TxResponse{}

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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
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

	mockClient := cosmosMocks.NewMockCosmosClient(t)
	mockDB := appMocks.NewMockDatabase(t)
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

		tx := &sdk.TxResponse{}

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

		mockDB.EXPECT().UpdateOne(models.CollectionInvalidMints, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
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

		tx := &sdk.TxResponse{}

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

		mockDB.EXPECT().UpdateOne(models.CollectionBurns, filterUpdate, mock.Anything).Return(primitive.NewObjectID(), nil).
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
		app.Config.Pocket.MultisigAddress = ""
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

}
