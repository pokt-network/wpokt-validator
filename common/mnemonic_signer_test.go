package common

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNewMnemonicSigner(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic)
	assert.NoError(t, err)
	assert.NotNil(t, signer)

	assert.NotNil(t, signer.ethPrivKey)
	assert.NotNil(t, signer.cosmosPrivKey)
	assert.NotEqual(t, common.Address{}, signer.ethAddress)
	assert.NotNil(t, signer.cosmosPubKey)
}

func TestMnemonicSigner_EthSign(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic)
	assert.NoError(t, err)

	data := []byte("test data")
	sig, err := signer.EthSign(data)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	if sig[64] != 27 && sig[64] != 28 {
		t.Fatalf("invalid Ethereum signature")
	}

	sig[64] -= 27

	hash := crypto.Keccak256(data)
	pubKey, err := crypto.SigToPub(hash, sig)
	assert.NoError(t, err)

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	assert.Equal(t, signer.EthAddress(), recoveredAddr)
}

func TestMnemonicSigner_CosmosSign(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic)
	assert.NoError(t, err)

	data := []byte("test data")
	sig, err := signer.CosmosSign(data)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	assert.True(t, signer.cosmosPubKey.VerifySignature(data, sig))
}

func TestMnemonicSigner_EthAddress(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic)
	assert.NoError(t, err)

	assert.NotEqual(t, common.Address{}, signer.EthAddress())
}

func TestMnemonicSigner_CosmosPublicKey(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic)
	assert.NoError(t, err)

	assert.NotNil(t, signer.CosmosPublicKey())
}

func TestMnemonicSigner_Destroy(t *testing.T) {

	mnemonic := "test test test test test test test test test test test junk"
	signer, err := NewMnemonicSigner(mnemonic)
	assert.NoError(t, err)

	signer.Destroy()
	// Nothing to assert here since the Destroy method does nothing
}
