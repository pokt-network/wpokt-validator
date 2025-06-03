package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytes32FromAddressHex(t *testing.T) {
	addr := "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"
	expected := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 171, 88, 1, 167, 211, 152, 53, 27, 139, 225, 28, 67, 158, 5, 197, 179, 37, 154, 236, 155}
	result, err := Bytes32FromAddressHex(addr)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBytes32FromAddressHex_Error(t *testing.T) {
	addr := "error"
	result, err := Bytes32FromAddressHex(addr)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestBytes32FromAddressHex_InvalidLength(t *testing.T) {
	addr := "0xAb5801a7D398351b8bE11C439e05C5b3259aec"
	result, err := Bytes32FromAddressHex(addr)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAddressLength, err)
	assert.Empty(t, result)
}

func TestIsZeroHex(t *testing.T) {
	assert.True(t, isZeroHex("0x0000000000000000000000000000000000000000"))
	assert.False(t, isZeroHex("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"))
}

func TestBytesFromHex(t *testing.T) {
	hexStr := "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"
	expected := []byte{171, 88, 1, 167, 211, 152, 53, 27, 139, 225, 28, 67, 158, 5, 197, 179, 37, 154, 236, 155}
	result, err := BytesFromHex(hexStr)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBytesFromHex_EmptyString(t *testing.T) {
	hexStr := ""
	result, err := BytesFromHex(hexStr)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBytesFromAddressHex(t *testing.T) {
	addr := "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"
	expected := []byte{171, 88, 1, 167, 211, 152, 53, 27, 139, 225, 28, 67, 158, 5, 197, 179, 37, 154, 236, 155}
	result, err := BytesFromAddressHex(addr)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBytesFromAddressHex_InvalidLength(t *testing.T) {
	addr := "0xAb5801a7D398351b8bE11C439e05C5b3259aec"
	result, err := BytesFromAddressHex(addr)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAddressLength, err)
	assert.Empty(t, result)
}

func TestIsValidAddressHex(t *testing.T) {
	assert.True(t, IsValidAddressHex("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"))
	assert.False(t, IsValidAddressHex("0xAb5801a7D398351b8bE11C439e05C5b3259ae"))
}

func TestEnsure0xPrefix(t *testing.T) {
	assert.Equal(t, "0xab5801a7d398351b8be11c439e05c5b3259aec9b", Ensure0xPrefix("ab5801a7d398351b8be11c439e05c5b3259aec9b"))
	assert.Equal(t, "0xab5801a7d398351b8be11c439e05c5b3259aec9b", Ensure0xPrefix("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"))
}

func TestAddressHexFromBytes(t *testing.T) {
	address := []byte{171, 88, 1, 167, 211, 152, 53, 27, 139, 225, 28, 67, 158, 5, 197, 179, 37, 154, 236, 155}
	expected := "0xab5801a7d398351b8be11c439e05c5b3259aec9b"
	result, err := AddressHexFromBytes(address)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestAddressHexFromBytes_InvalidLength(t *testing.T) {
	address := []byte{171, 88, 1, 167, 211, 152, 53, 27}
	result, err := AddressHexFromBytes(address)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAddressLength, err)
	assert.Empty(t, result)
}

func TestHexFromBytes(t *testing.T) {
	bytes := []byte{171, 88, 1, 167, 211, 152, 53, 27, 139, 225, 28, 67, 158, 5, 197, 179, 37, 154, 236, 155}
	expected := "0xab5801a7d398351b8be11c439e05c5b3259aec9b"
	result := HexFromBytes(bytes)
	assert.Equal(t, expected, result)
}
