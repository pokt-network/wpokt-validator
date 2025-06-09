package eth

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dan13ram/wpokt-validator/app"
	appMocks "github.com/dan13ram/wpokt-validator/app/mocks"
	"github.com/dan13ram/wpokt-validator/common"
	cosmosMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	ethMocks "github.com/dan13ram/wpokt-validator/eth/client/mocks"
	"github.com/dan13ram/wpokt-validator/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cosmosUtil "github.com/dan13ram/wpokt-validator/cosmos/util"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func NewTestMintSigner(t *testing.T, mockWrappedPocketContract *ethMocks.MockWrappedPocketContract,
	mockMintControllerContract *ethMocks.MockMintControllerContract,
	mockEthClient *ethMocks.MockEthereumClient, mockPoktClient *cosmosMocks.MockCosmosClient) *MintSignerRunner {
	pk, _ := crypto.HexToECDSA("1395eeb9c36ef43e9e05692c9ee34034c00a9bef301135a96d082b2a65fd1680")
	address := crypto.PubkeyToAddress(pk.PublicKey).Hex()

	x := &MintSignerRunner{
		address:         strings.ToLower(address),
		privateKey:      pk,
		vaultAddress:    "vaultAddress",
		wpoktAddress:    "wpoktAddress",
		validatorCount:  3,
		signerThreshold: 2,
		domain: eth.DomainData{
			Name:              "Test",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: common.HexToAddress(""),
		},
		wpoktContract:          mockWrappedPocketContract,
		mintControllerContract: mockMintControllerContract,
		ethClient:              mockEthClient,
		cosmosClient:           mockPoktClient,
		cosmosHeight:           100,
		minimumAmount:          math.NewInt(10000),
		maximumAmount:          math.NewInt(1000000),
	}
	return x
}

func TestMintSignerStatus(t *testing.T) {
	mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
	mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
	mockEthClient := ethMocks.NewMockEthereumClient(t)
	mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
	x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

	status := x.Status()
	assert.Equal(t, status.EthBlockNumber, "")
	assert.Equal(t, status.PoktHeight, "100")
}

func TestMintSignerUpdateBlocks(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetLatestBlockHeight().Return(int64(200), nil)

		x.UpdateBlocks()

		assert.Equal(t, x.cosmosHeight, int64(200))
	})

	t.Run("With Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockPoktClient.EXPECT().GetLatestBlockHeight().Return(0, errors.New("error"))

		x.UpdateBlocks()

		assert.Equal(t, x.cosmosHeight, int64(100))
	})

}

func TestMintSignerUpdateValidatorCount(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockMintControllerContract.EXPECT().ValidatorCount(mock.Anything).Return(big.NewInt(5), nil)

		x.UpdateValidatorCount()

		assert.Equal(t, x.validatorCount, int64(5))
	})

	t.Run("With Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockMintControllerContract.EXPECT().ValidatorCount(mock.Anything).Return(big.NewInt(5), errors.New("error"))

		x.UpdateValidatorCount()

		assert.Equal(t, x.validatorCount, int64(3))
	})

}

func TestMintSignerUpdateDomainData(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		domain := eth.DomainData{
			Name:              "New Domain",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: common.HexToAddress(""),
		}

		mockMintControllerContract.EXPECT().Eip712Domain(mock.Anything).Return(domain, nil)

		x.UpdateDomainData()

		assert.Equal(t, x.domain.Name, "New Domain")
	})

	t.Run("With Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		domain := eth.DomainData{
			Name:              "New Domain",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: common.HexToAddress(""),
		}

		mockMintControllerContract.EXPECT().Eip712Domain(mock.Anything).Return(domain, errors.New("error"))

		x.UpdateDomainData()

		assert.Equal(t, x.domain.Name, "Test")
	})

}

func TestMintSignerUpdateMaxMintLimit(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockMintControllerContract.EXPECT().MaxMintLimit(mock.Anything).Return(big.NewInt(500000), nil)

		x.UpdateMaxMintLimit()

		assert.Equal(t, x.maximumAmount, math.NewInt(500000))
	})

	t.Run("With Error", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockMintControllerContract.EXPECT().MaxMintLimit(mock.Anything).Return(big.NewInt(500000), errors.New("error"))

		x.UpdateMaxMintLimit()

		assert.Equal(t, x.maximumAmount, math.NewInt(1000000))
	})

}

func TestMintSignerFindNonce(t *testing.T) {

	t.Run("Nonce already set", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{
			Nonce: "10",
		}

		gotNonce, err := x.FindNonce(mint)

		assert.Equal(t, gotNonce, big.NewInt(10))
		assert.Nil(t, err)

	})

	t.Run("No pending mints", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		nonce := big.NewInt(10)

		mockWrappedPocketContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(nonce, nil)

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
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		nonce := big.NewInt(5)

		mockWrappedPocketContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(nonce, nil)

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
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{
			Nonce: "invalid",
		}

		gotNonce, err := x.FindNonce(mint)

		assert.NotNil(t, err)
		assert.Nil(t, gotNonce)

	})

	t.Run("Error finding current nonce", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		mockWrappedPocketContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(big.NewInt(5), errors.New("error"))

		gotNonce, err := x.FindNonce(mint)

		assert.NotNil(t, err)
		assert.Nil(t, gotNonce)

	})

	t.Run("Error finding pending mints", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		app.DB = mockDB

		mint := &models.Mint{}

		mockWrappedPocketContract.EXPECT().GetUserNonce(mock.Anything, common.HexToAddress("")).Return(big.NewInt(5), nil)
		mockDB.EXPECT().FindMany(models.CollectionMints, mock.Anything, mock.Anything).Return(errors.New("error"))

		gotNonce, err := x.FindNonce(mint)

		assert.NotNil(t, err)
		assert.Nil(t, gotNonce)

	})

}

func TestValidateMint(t *testing.T) {

	t.Run("Error fetching transaction", func(t *testing.T) {

		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.NotNil(t, err)

	})

	t.Run("Nil transaction", func(t *testing.T) {

		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		mockPoktClient.EXPECT().GetTx("").Return(nil, nil)

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Error(t, err)

	})

	t.Run("Failed transaction", func(t *testing.T) {

		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		tx := &sdk.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,
			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid: false,
			}
		}

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Needs refund transaction", func(t *testing.T) {

		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mint := &models.Mint{}

		tx := &sdk.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:     true,
				NeedsRefund: true,
			}
		}

		valid, err := x.ValidateMint(mint)

		assert.False(t, valid)
		assert.Nil(t, err)

	})

	t.Run("Valid transaction and mint", func(t *testing.T) {

		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			RecipientChainID: "31337",
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,
			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
		}

		valid, err := x.ValidateMint(mint)

		assert.True(t, valid)
		assert.Nil(t, err)

	})

}

func TestMintSignerHandleMint(t *testing.T) {

	t.Run("Nil mint", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		success := x.HandleMint(nil)

		assert.False(t, success)
	})

	t.Run("Invalid mint amount", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "invalid",
			RecipientChainID: "31337",
		}

		success := x.HandleMint(mint)

		assert.False(t, success)
	})

	t.Run("Error finding nonce", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "invalid",
			RecipientChainID: "31337",
		}

		success := x.HandleMint(mint)

		assert.False(t, success)
	})

	t.Run("Error updating confirmations", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 1

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "invalid",
		}

		success := x.HandleMint(mint)

		assert.False(t, success)
	})

	t.Run("Error validating mint", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
		}

		mockPoktClient.EXPECT().GetTx("").Return(nil, errors.New("error"))

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.False(t, success)
	})

	t.Run("Validating mint returned false", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{}

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
			}).Return(primitive.NewObjectID(), nil)

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid: false,
			}
		}

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.True(t, success)
	})

	t.Run("Validating mint returned true, mint pending", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 10

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{}
		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
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
			}).Return(primitive.NewObjectID(), nil)

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusPending, mint.Status)
		assert.Equal(t, mint.Confirmations, "1")

		assert.True(t, success)
	})

	t.Run("Validating mint returned true, mint confirmed, signing failed", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		x.domain = eth.DomainData{
			ChainId: big.NewInt(1),
		}

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
		}

		app.Config.Ethereum.ChainID = "31337"

		tx := &sdk.TxResponse{}
		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
		}

		mockPoktClient.EXPECT().GetTx("").Return(tx, nil)

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.False(t, success)
	})

	t.Run("Error updating mint", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
		}

		app.Config.Ethereum.ChainID = "31337"

		bech32Prefix := "pokt"

		multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
		multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
		senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		txResponse := &sdk.TxResponse{
			TxHash: "0x123",
			Height: 90,
			Code:   0,
			Events: []abci.Event{
				{
					Type: "transfer",
					Attributes: []abci.EventAttribute{
						{Key: "sender", Value: senderBech32},
						{Key: "recipient", Value: multisigBech32},
						{Key: "amount", Value: "20000upokt"},
					},
				},
			},
		}

		tx := &tx.Tx{
			Body: &tx.TxBody{
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}
		txValue, err := tx.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		txResponse.Tx = &codectypes.Any{Value: txValue}

		mockPoktClient.EXPECT().GetTx("").Return(txResponse, nil)

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
			}).Return(primitive.NewObjectID(), errors.New("error"))

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
		}

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		address := common.HexToAddress("0x1234").Hex()

		app.Config.Pocket.Confirmations = 0

		mint := &models.Mint{
			SenderAddress:    "abcd",
			RecipientAddress: address,
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
		}

		app.Config.Ethereum.ChainID = "31337"
		bech32Prefix := "pokt"

		multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
		multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
		senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		txResponse := &sdk.TxResponse{
			TxHash: "0x123",
			Height: 90,
			Code:   0,
			Events: []abci.Event{
				{
					Type: "transfer",
					Attributes: []abci.EventAttribute{
						{Key: "sender", Value: senderBech32},
						{Key: "recipient", Value: multisigBech32},
						{Key: "amount", Value: "20000upokt"},
					},
				},
			},
		}

		tx := &tx.Tx{
			Body: &tx.TxBody{
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}
		txValue, err := tx.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		txResponse.Tx = &codectypes.Any{Value: txValue}

		mockPoktClient.EXPECT().GetTx("").Return(txResponse, nil)

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
			}).Return(primitive.NewObjectID(), nil)

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
		}

		success := x.HandleMint(mint)

		assert.Equal(t, models.StatusConfirmed, mint.Status)
		assert.Equal(t, mint.Confirmations, "0")

		assert.True(t, success)
	})

}

func TestMintSignerSyncTxs(t *testing.T) {

	t.Run("Error finding mints", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

		mockDB.EXPECT().FindMany(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		success := x.SyncTxs()

		assert.False(t, success)

	})

	t.Run("No mints to handle", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

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
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

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
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

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
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainID = "31337"

		bech32Prefix := "pokt"

		multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
		multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
		senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		txResponse := &sdk.TxResponse{
			TxHash: "0x123",
			Height: 90,
			Code:   0,
			Events: []abci.Event{
				{
					Type: "transfer",
					Attributes: []abci.EventAttribute{
						{Key: "sender", Value: senderBech32},
						{Key: "recipient", Value: multisigBech32},
						{Key: "amount", Value: "20000upokt"},
					},
				},
			},
		}

		tx := &tx.Tx{
			Body: &tx.TxBody{
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}
		txValue, err := tx.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		txResponse.Tx = &codectypes.Any{Value: txValue}

		mockPoktClient.EXPECT().GetTx("").Return(txResponse, nil)

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
			}).Return(primitive.NewObjectID(), nil)

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(errors.New("error"))

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
		}
		success := x.SyncTxs()

		assert.False(t, success)
	})

	t.Run("Successful case", func(t *testing.T) {
		mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
		mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
		mockEthClient := ethMocks.NewMockEthereumClient(t)
		mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
		mockDB := appMocks.NewMockDatabase(t)
		app.DB = mockDB
		x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

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
			Amount:           "20000",
			Nonce:            "1",
			RecipientChainID: "31337",
			Height:           "99",
			Confirmations:    "invalid",
		}

		app.Config.Ethereum.ChainID = "31337"

		bech32Prefix := "pokt"

		multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
		multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
		senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		txResponse := &sdk.TxResponse{
			TxHash: "0x123",
			Height: 90,
			Code:   0,
			Events: []abci.Event{
				{
					Type: "transfer",
					Attributes: []abci.EventAttribute{
						{Key: "sender", Value: senderBech32},
						{Key: "recipient", Value: multisigBech32},
						{Key: "amount", Value: "20000upokt"},
					},
				},
			},
		}

		tx := &tx.Tx{
			Body: &tx.TxBody{
				Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
			},
		}
		txValue, err := tx.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		txResponse.Tx = &codectypes.Any{Value: txValue}

		mockPoktClient.EXPECT().GetTx("").Return(txResponse, nil)

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
			}).Return(primitive.NewObjectID(), nil)

		mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
		mockDB.EXPECT().Unlock("lockId").Return(nil)

		oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
		defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
		cosmosUtilValidateTxToCosmosMultisig = func(
			txResponse *sdk.TxResponse,
			config models.CosmosConfig,

			minAmount math.Int,
			maxAmount math.Int,
		) *cosmosUtil.ValidateTxResult {
			return &cosmosUtil.ValidateTxResult{
				TxValid:       true,
				NeedsRefund:   false,
				SenderAddress: "abcd",
				Memo: models.MintMemo{
					Address: address,
					ChainID: "31337",
				},
				Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
			}
		}
		success := x.SyncTxs()

		assert.True(t, success)
	})

}

func TestMintSignerRun(t *testing.T) {

	mockWrappedPocketContract := ethMocks.NewMockWrappedPocketContract(t)
	mockMintControllerContract := ethMocks.NewMockMintControllerContract(t)
	mockEthClient := ethMocks.NewMockEthereumClient(t)
	mockPoktClient := cosmosMocks.NewMockCosmosClient(t)
	mockDB := appMocks.NewMockDatabase(t)
	app.DB = mockDB
	x := NewTestMintSigner(t, mockWrappedPocketContract, mockMintControllerContract, mockEthClient, mockPoktClient)

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
		Amount:           "20000",
		Nonce:            "1",
		RecipientChainID: "31337",
		Height:           "99",
		Confirmations:    "invalid",
	}

	app.Config.Ethereum.ChainID = "31337"

	bech32Prefix := "pokt"

	multisigAddress := ethcommon.BytesToAddress([]byte("pokt1multisig"))
	multisigBech32, err := common.Bech32FromBytes(bech32Prefix, multisigAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	senderAddress := ethcommon.BytesToAddress([]byte("pokt1sender"))
	senderBech32, err := common.Bech32FromBytes(bech32Prefix, senderAddress.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	txResponse := &sdk.TxResponse{
		TxHash: "0x123",
		Height: 90,
		Code:   0,
		Events: []abci.Event{
			{
				Type: "transfer",
				Attributes: []abci.EventAttribute{
					{Key: "sender", Value: senderBech32},
					{Key: "recipient", Value: multisigBech32},
					{Key: "amount", Value: "20000upokt"},
				},
			},
		},
	}

	tx := &tx.Tx{
		Body: &tx.TxBody{
			Memo: fmt.Sprintf(`{ "address": "%s", "chain_id": "31337" }`, address),
		},
	}
	txValue, err := tx.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	txResponse.Tx = &codectypes.Any{Value: txValue}

	mockPoktClient.EXPECT().GetTx("").Return(txResponse, nil)

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
		}).Return(primitive.NewObjectID(), nil)

	mockDB.EXPECT().XLock(mock.Anything).Return("lockId", nil)
	mockDB.EXPECT().Unlock("lockId").Return(nil)

	mockPoktClient.EXPECT().GetLatestBlockHeight().Return(int64(200), nil)

	mockMintControllerContract.EXPECT().ValidatorCount(mock.Anything).Return(big.NewInt(3), nil)

	mockMintControllerContract.EXPECT().MaxMintLimit(mock.Anything).Return(big.NewInt(1000000), nil)

	oldCosmosUtilValidateTxToCosmosMultisig := cosmosUtilValidateTxToCosmosMultisig
	defer func() { cosmosUtilValidateTxToCosmosMultisig = oldCosmosUtilValidateTxToCosmosMultisig }()
	cosmosUtilValidateTxToCosmosMultisig = func(
		txResponse *sdk.TxResponse,
		config models.CosmosConfig,

		minAmount math.Int,
		maxAmount math.Int,
	) *cosmosUtil.ValidateTxResult {
		return &cosmosUtil.ValidateTxResult{
			TxValid:       true,
			NeedsRefund:   false,
			SenderAddress: "abcd",
			Memo: models.MintMemo{
				Address: address,
				ChainID: "31337",
			},
			Amount: sdk.NewCoin("upokt", math.NewInt(20000)),
		}
	}

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
