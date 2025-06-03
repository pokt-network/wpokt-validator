package common

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/assert"
)

const (
	testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
)

func TestCosmosPrivateKeyFromMnemonic(t *testing.T) {
	privKey, err := CosmosPrivateKeyFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	assert.NotNil(t, privKey)
}

func TestCosmosPrivateKeyFromMnemonic_Error(t *testing.T) {
	privKey, err := CosmosPrivateKeyFromMnemonic("error")
	assert.Error(t, err)
	assert.Nil(t, privKey)
}

func TestCosmosPublicKeyFromMnemonic(t *testing.T) {
	pubKey, err := CosmosPublicKeyFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)
}

func TestCosmosPublicKeyFromMnemonic_Error(t *testing.T) {
	pubKey, err := CosmosPublicKeyFromMnemonic("error")
	assert.Error(t, err)
	assert.Nil(t, pubKey)
}

func TestCosmosPublicKeyHexFromMnemonic(t *testing.T) {
	pubKeyHex, err := CosmosPublicKeyHexFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	assert.NotEmpty(t, pubKeyHex)
}

func TestCosmosPublicKeyHexFromMnemonic_Error(t *testing.T) {
	pubKeyHex, err := CosmosPublicKeyHexFromMnemonic("error")
	assert.Error(t, err)
	assert.Empty(t, pubKeyHex)
}

func TestCosmosPublicKeyFromHex(t *testing.T) {
	pubKeyHex, _ := CosmosPublicKeyHexFromMnemonic(testMnemonic)
	pubKey, err := CosmosPublicKeyFromHex(pubKeyHex)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)
}

func TestCosmosPublicKeyFromHex_Error(t *testing.T) {
	pubKey, err := CosmosPublicKeyFromHex("error")
	assert.Error(t, err)
	assert.Nil(t, pubKey)
}

func TestCosmosPublicKeyFromHex_InvalidLength(t *testing.T) {
	_, err := CosmosPublicKeyFromHex("abcd")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPublicKeyLength, err)
}

func TestIsValidCosmosPublicKey(t *testing.T) {
	pubKeyHex, _ := CosmosPublicKeyHexFromMnemonic(testMnemonic)
	valid := IsValidCosmosPublicKey(pubKeyHex)
	assert.True(t, valid)
}

func TestIsValidCosmosPublicKey_Invalid(t *testing.T) {
	valid := IsValidCosmosPublicKey("abcd")
	assert.False(t, valid)
}

func TestBytesFromBech32(t *testing.T) {
	bech32Prefix := "cosmos"
	address, _ := Bech32FromBytes(bech32Prefix, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	bytes, err := BytesFromBech32(bech32Prefix, address)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(bytes))
}

func TestBech32FromBytes(t *testing.T) {
	bech32Prefix := "cosmos"
	bytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	address, err := Bech32FromBytes(bech32Prefix, bytes)
	assert.NoError(t, err)
	assert.NotEmpty(t, address)
}

func TestBech32FromBytes_InvalidLength(t *testing.T) {
	bech32Prefix := "cosmos"
	bytes := []byte{1, 2, 3}
	address, err := Bech32FromBytes(bech32Prefix, bytes)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAddressLength, err)
	assert.Empty(t, address)
}

func TestAddressBytesFromBech32(t *testing.T) {
	bech32Prefix := "cosmos"
	address, _ := Bech32FromBytes(bech32Prefix, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	bytes, err := AddressBytesFromBech32(bech32Prefix, address)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(bytes))
}

func TestAddressBytesFromBech32_InvalidLength(t *testing.T) {
	bech32Prefix := "cosmos"
	address, _ := bech32.ConvertAndEncode(bech32Prefix, []byte{1, 2, 3})
	bytes, err := AddressBytesFromBech32(bech32Prefix, address)
	assert.Error(t, err)
	assert.Nil(t, bytes)
}

func TestAddressBytesFromBech32_Error(t *testing.T) {
	bech32Prefix := "cosmos"
	bytes, err := AddressBytesFromBech32(bech32Prefix, "error")
	assert.Error(t, err)
	assert.Nil(t, bytes)
}

func TestIsValidBech32Address(t *testing.T) {
	bech32Prefix := "cosmos"
	address, _ := Bech32FromBytes(bech32Prefix, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	valid := IsValidBech32Address(bech32Prefix, address)
	assert.True(t, valid)
}

func TestIsValidBech32Address_Invalid(t *testing.T) {
	bech32Prefix := "cosmos"
	address, _ := bech32.ConvertAndEncode(bech32Prefix, []byte{1, 2, 3})
	valid := IsValidBech32Address(bech32Prefix, address)
	assert.False(t, valid)
}
