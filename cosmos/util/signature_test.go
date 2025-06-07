package util

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/dan13ram/wpokt-validator/common"
	clientMocks "github.com/dan13ram/wpokt-validator/cosmos/client/mocks"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func TestSignWithPrivKey(t *testing.T) {
	// Generate a new private key
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, msg, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, sigV2)
	assert.Equal(t, privKey.CosmosPublicKey(), sigV2.PubKey)
	assert.Equal(t, signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, sigV2.Data.(*signingtypes.SingleSignatureData).SignMode)
	assert.NotEmpty(t, sigV2.Data.(*signingtypes.SingleSignatureData).Signature)
	assert.Equal(t, uint64(1), sigV2.Sequence)
	pub := privKey.CosmosPublicKey()
	assert.Equal(t, true, pub.VerifySignature(msg, sigV2.Data.(*signingtypes.SingleSignatureData).Signature))
}

func TestSignWithPrivKey_ErrorSignBytes(t *testing.T) {
	// Generate a new private key
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	pubKey := &secp256k1.PubKey{}

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        pubKey,
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, _, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, sigV2)
}

type mockSigner struct {
}

func (m *mockSigner) CosmosSign(msg []byte) ([]byte, error) {
	return nil, errors.New("error signing")
}

func (m *mockSigner) CosmosPublicKey() types.PubKey {
	return nil
}

func (m *mockSigner) EthSign(msg []byte) ([]byte, error) {
	return nil, errors.New("error signing")
}

func (m *mockSigner) EthAddress() ethcommon.Address {
	return ethcommon.Address{}
}

func (m *mockSigner) Destroy() {
}

func TestSignWithPrivKey_ErrorSigning(t *testing.T) {
	// Generate a new private key
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	mockPrivKey := &mockSigner{}

	// Call the SignWithPrivKey function
	sigV2, _, err := SignWithPrivKey(ctx, signerData, txBuilder, mockPrivKey, txConfig, 1)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, sigV2)
}

func TestValidateSignature(t *testing.T) {
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, _, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)
	assert.NoError(t, err)

	config := models.PocketConfig{
		ChainID:      "poktroll",
		Bech32Prefix: "pokt",
	}

	err = ValidateSignature(config, &sigV2, 1, 1, txConfig, txBuilder)
	assert.NoError(t, err)
}

func TestValidateSignature_VerificationFailure(t *testing.T) {
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, _, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)
	assert.NoError(t, err)

	config := models.PocketConfig{
		ChainID:      "poktroll-different",
		Bech32Prefix: "pokt",
	}

	err = ValidateSignature(config, &sigV2, 1, 1, txConfig, txBuilder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't verify signature for address")
}

func TestValidateSignature_AnyError(t *testing.T) {
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, _, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)
	assert.NoError(t, err)

	config := models.PocketConfig{
		ChainID:      "poktroll-different",
		Bech32Prefix: "pokt",
	}

	sigV2.PubKey = nil

	err = ValidateSignature(config, &sigV2, 1, 1, txConfig, txBuilder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating any pubkey")
}

func TestValidateSignature_TxError(t *testing.T) {
	privKey, _ := common.NewMnemonicSigner("test test test test test test test test test test test junk")

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := clientMocks.NewMockTxBuilder(t)

	tx := clientMocks.NewMockTx(t)

	txBuilder.EXPECT().GetTx().Return(tx)

	// Call the SignWithPrivKey function
	sigV2 := signingtypes.SignatureV2{
		PubKey: privKey.CosmosPublicKey(),
	}

	config := models.PocketConfig{
		ChainID:      "poktroll",
		Bech32Prefix: "pokt",
	}

	err := ValidateSignature(config, &sigV2, 1, 1, txConfig, txBuilder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected Tx to be signing.V2AdaptableTx")
}

func TestSignWithPrivKey_WithGCPKMS(t *testing.T) {
	keyName := os.Getenv("GCP_KMS_KEY_NAME")
	if keyName == "" {
		t.Skip("GCP KMS key name not set")
	}
	credentails := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentails == "" {
		t.Skip("GCP credentials not set")
	}

	privKey, _ := common.NewGcpKmsSigner(keyName)

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, msg, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, sigV2)
	assert.Equal(t, privKey.CosmosPublicKey(), sigV2.PubKey)
	assert.Equal(t, signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, sigV2.Data.(*signingtypes.SingleSignatureData).SignMode)
	assert.NotEmpty(t, sigV2.Data.(*signingtypes.SingleSignatureData).Signature)
	assert.Equal(t, uint64(1), sigV2.Sequence)
	pub := privKey.CosmosPublicKey()
	assert.Equal(t, true, pub.VerifySignature(msg, sigV2.Data.(*signingtypes.SingleSignatureData).Signature))
}

func TestValidateSignature_WithGCPKMS(t *testing.T) {
	keyName := os.Getenv("GCP_KMS_KEY_NAME")
	if keyName == "" {
		t.Skip("GCP KMS key name not set")
	}
	credentails := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentails == "" {
		t.Skip("GCP credentials not set")
	}

	privKey, _ := common.NewGcpKmsSigner(keyName)

	// Create a new TxConfig
	txConfig := NewTxConfig("pokt")

	// Create a new TxBuilder
	txBuilder := txConfig.NewTxBuilder()

	// Create dummy signer data
	signerData := authsigning.SignerData{
		ChainID:       "poktroll",
		AccountNumber: 1,
		Sequence:      1,
		PubKey:        privKey.CosmosPublicKey(),
		Address:       sdk.AccAddress(privKey.CosmosPublicKey().Address()).String(),
	}

	// Create a new context
	ctx := context.Background()

	// Call the SignWithPrivKey function
	sigV2, _, err := SignWithPrivKey(ctx, signerData, txBuilder, privKey, txConfig, 1)
	assert.NoError(t, err)

	config := models.PocketConfig{
		ChainID:      "poktroll",
		Bech32Prefix: "pokt",
	}

	err = ValidateSignature(config, &sigV2, 1, 1, txConfig, txBuilder)
	assert.NoError(t, err)
}
