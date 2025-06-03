package common

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func TestIsValidEthereumAddress(t *testing.T) {
	assert.True(t, IsValidEthereumAddress("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"))
	assert.False(t, IsValidEthereumAddress("0xAb5801a7D398351b8bE11C439e05C5b3259aec9"))
	assert.False(t, IsValidEthereumAddress("Ab5801a7D398351b8bE11C439e05C5b3259aec9B"))
}

func TestEthereumPrivateKeyFromMnemonic(t *testing.T) {
	privateKey, err := EthereumPrivateKeyFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
}

func TestEthereumPrivateKeyFromMnemonic_Error(t *testing.T) {
	privateKey, err := EthereumPrivateKeyFromMnemonic("Error")
	assert.Error(t, err)
	assert.Nil(t, privateKey)
}

func TestEthereumAddressFromMnemonic(t *testing.T) {
	address, err := EthereumAddressFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.True(t, strings.HasPrefix(address, "0x"))
	assert.Equal(t, "0x9858EfFD232B4033E47d90003D41EC34EcaEda94", address)
}

func TestEthereumAddressFromMnemonic_Error(t *testing.T) {
	address, err := EthereumAddressFromMnemonic("Error")
	assert.Error(t, err)
	assert.Empty(t, address)
}

func TestHexToAddress(t *testing.T) {
	hexStr := "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"
	expected := ethcommon.HexToAddress(hexStr)
	result := HexToAddress(hexStr)
	assert.Equal(t, expected, result)
}

func TestEthereumPrivateKeyToAddressHex(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	addressHex, err := EthereumPrivateKeyToAddressHex(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressHex)
	assert.True(t, strings.HasPrefix(addressHex, "0x"))
}

func TestEthereumPrivateKeyToAddressHex_Error(t *testing.T) {
	addressHex, err := EthereumPrivateKeyToAddressHex(nil)
	assert.Error(t, err)
	assert.Empty(t, addressHex)
	assert.Equal(t, "private key is nil or public key is nil", err.Error())
}
