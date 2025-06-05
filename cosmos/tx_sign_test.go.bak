package cosmos

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cosmos/cosmos-sdk/client"

	clientMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	"github.com/dan13ram/wpokt-validator/models"

	"github.com/cosmos/cosmos-sdk/x/auth/signing"

	ethcommon "github.com/ethereum/go-ethereum/common"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/dan13ram/wpokt-validator/common"
)

func TestCosmosSignTx(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, nil
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)

	txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)
	txBuilder.EXPECT().GetTx().Return(tx)

	var encoder sdk.TxEncoder = func(tx sdk.Tx) ([]byte, error) {
		return []byte("encoded tx"), nil
	}

	txConfig.EXPECT().TxJSONEncoder().Return(encoder)

	mockClient.EXPECT().GetAccount(multisigAddr).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		0,
		[]models.Signature{},
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	txBuilder.AssertExpectations(t)
	txConfig.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "encoded tx", txBody)
	assert.NotEmpty(t, signatures)
}

func TestCosmosSignTx_InvalidSigner(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{
		{
			Signer:    "0xsigner",
			Signature: "0xsignature",
		},
	}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error parsing signer")
}

func TestCosmosSignTx_AlreadySigned(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signerAddr := ethcommon.BytesToAddress(signerKey.CosmosPublicKey().Address().Bytes())

	signatures := []models.Signature{
		{
			Signer:    signerAddr.Hex(),
			Signature: "0xsignature",
		},
	}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Equal(t, ErrAlreadySigned, err)
}

func TestCosmosSignTx_InvalidMultisigAddress(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: "multisigAddr",
	}

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error parsing multisig address")
}

func TestCosmosSignTx_ErrorNewTx(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "", assert.AnError
	}

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error creating tx body")
}

func TestCosmosSignTx_WrapError(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}
	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return nil, nil, assert.AnError
	}

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error wrapping tx builder")
}

func TestCosmosSignTx_ErrorGetSigners(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	tx.EXPECT().GetSigners().Return(nil, assert.AnError)
	txBuilder.EXPECT().GetTx().Return(tx)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error getting signers")
}

func TestCosmosSignTx_MultisigIsNotSigner(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	signers := [][]byte{}
	assert.False(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)
	txBuilder.EXPECT().GetTx().Return(tx)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "multisig is not a signer")
}

func TestCosmosSignTx_AccountError(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		t.Errorf("utilSignWithPrivKey should not be called")
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, nil
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)

	txBuilder.EXPECT().GetTx().Return(tx)

	mockClient.EXPECT().GetAccount(multisigAddr).Return(nil, assert.AnError)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error getting account")
}

func TestCosmosSignTx_ErrorSigning(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, assert.AnError
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)

	txBuilder.EXPECT().GetTx().Return(tx)

	mockClient.EXPECT().GetAccount(multisigAddr).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error signing tx")
}

func TestCosmosSignTx_ErrorGettingSignatures(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	anotherSignerAddr := ethcommon.BytesToAddress([]byte("anotherSigner"))

	signatures := []models.Signature{{
		Signer:    anotherSignerAddr.Hex(),
		Signature: "0xsignature",
	}}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, nil
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	txBuilder.EXPECT().GetTx().Return(tx)
	tx.EXPECT().GetSigners().Return(signers, nil)

	tx.EXPECT().GetSignaturesV2().Return(nil, assert.AnError)

	mockClient.EXPECT().GetAccount(multisigAddr).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error getting signatures")
}

func TestCosmosSignTx_ErrorSettingSignatures(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, nil
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)

	txBuilder.EXPECT().GetTx().Return(tx)
	txBuilder.EXPECT().SetSignatures(mock.Anything).Return(assert.AnError)
	mockClient.EXPECT().GetAccount(multisigAddr).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error setting signatures")
}

func TestCosmosSignTx_ErrorEncoding(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, nil
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)

	txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)
	txBuilder.EXPECT().GetTx().Return(tx)

	var encoder sdk.TxEncoder = func(tx sdk.Tx) ([]byte, error) {
		return []byte("encoded tx"), assert.AnError
	}

	txConfig.EXPECT().TxJSONEncoder().Return(encoder)

	mockClient.EXPECT().GetAccount(multisigAddr).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.Error(t, err)
	assert.Equal(t, "", txBody)
	assert.Nil(t, signatures)
	assert.Contains(t, err.Error(), "error encoding tx")
}

func TestCosmosSignTx_TransactionBodyNonEmpty(t *testing.T) {
	mockClient := clientMocks.NewMockCosmosClient(t)

	signerKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")
	multisigPk := multisig.NewLegacyAminoPubKey(1, []cryptotypes.PubKey{signerKey.CosmosPublicKey()})
	multisigAddr, _ := common.Bech32FromBytes("pokt", multisigPk.Address().Bytes())

	recipientAddr := ethcommon.BytesToAddress([]byte("recipient"))

	signatures := []models.Signature{}

	config := models.PocketConfig{
		ChainID:         "chain-id",
		CoinDenom:       "upokt",
		Bech32Prefix:    "pokt",
		MultisigAddress: multisigAddr,
	}

	txBuilder := clientMocks.NewMockTxBuilder(t)
	txConfig := clientMocks.NewMockTxConfig(t)
	tx := clientMocks.NewMockTx(t)

	utilNewSendTx = func(string, []byte, []byte, sdk.Coin, string, sdk.Coin) (string, error) {
		t.Errorf("utilNewSendTx should not be called")
		return "txBody", nil
	}

	utilWrapTxBuilder = func(prefix string, txBody string) (client.TxBuilder, client.TxConfig, error) {
		assert.Equal(t, "pokt", prefix)
		assert.Equal(t, "txAlreadyBody", txBody)
		return txBuilder, txConfig, nil
	}

	utilSignWithPrivKey = func(context.Context, signing.SignerData, client.TxBuilder, common.Signer, client.TxConfig, uint64) (signingtypes.SignatureV2, []byte, error) {
		return signingtypes.SignatureV2{
			PubKey: signerKey.CosmosPublicKey(),
			Data: &signingtypes.SingleSignatureData{
				SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
				Signature: []byte("signature"),
			},
		}, nil, nil
	}

	signers := [][]byte{multisigPk.Address().Bytes()}

	assert.True(t, isTxSigner(multisigPk.Address().Bytes(), signers))

	tx.EXPECT().GetSigners().Return(signers, nil)

	txBuilder.EXPECT().SetSignatures(mock.Anything).Return(nil)
	txBuilder.EXPECT().GetTx().Return(tx)

	var encoder sdk.TxEncoder = func(tx sdk.Tx) ([]byte, error) {
		return []byte("encoded tx"), nil
	}

	txConfig.EXPECT().TxJSONEncoder().Return(encoder)

	mockClient.EXPECT().GetAccount(multisigAddr).Return(&authtypes.BaseAccount{AccountNumber: 1, Sequence: 1}, nil)

	amount, _ := sdk.ParseCoinNormalized("100upokt")

	sequence := uint64(0)

	txBody, signatures, err := SignTx(
		signerKey,
		config,
		mockClient,
		sequence,
		signatures,
		"txAlreadyBody",
		recipientAddr[:],
		amount,
		"memo",
	)

	assert.NoError(t, err)
	assert.Equal(t, "encoded tx", txBody)
	assert.NotEmpty(t, signatures)
}
