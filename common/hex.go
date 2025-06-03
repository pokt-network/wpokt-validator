package common

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

var ErrInvalidAddressLength = errors.New("invalid address length")

func Bytes32FromAddressHex(addr string) ([]byte, error) {
	addr = strings.TrimPrefix(addr, "0x")
	address, err := hex.DecodeString(addr)
	if err != nil {
		return []byte{}, err
	}

	if len(address) != AddressLength {
		return []byte{}, ErrInvalidAddressLength
	}

	addressBig := new(big.Int).SetBytes(address)
	addressBytes := addressBig.FillBytes(make([]byte, 32)) // pad address to 32 bytes

	return addressBytes, nil
}

func isZeroHex(s string) bool {
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "0x")
	for _, c := range []byte(s) {
		if c != '0' {
			return false
		}
	}
	return true
}

func BytesFromHex(hexStr string) ([]byte, error) {
	if len(hexStr) == 0 {
		return nil, errors.New("decoding hex string failed: empty hex string")
	}
	hexStr = strings.ToLower(hexStr)
	hexStr = strings.TrimPrefix(hexStr, "0x")
	return hex.DecodeString(hexStr)
}

func BytesFromAddressHex(addr string) ([]byte, error) {
	bytes, err := BytesFromHex(addr)
	if err != nil {
		return nil, err
	}
	if len(bytes) != AddressLength {
		return nil, ErrInvalidAddressLength
	}
	return bytes, nil
}

func IsValidAddressHex(s string) bool {
	_, err := BytesFromAddressHex(s)
	return err == nil
}

func Ensure0xPrefix(str string) string {
	str = strings.ToLower(str)
	if !strings.HasPrefix(str, "0x") {
		return "0x" + str
	}
	return str
}

func AddressHexFromBytes(address []byte) (string, error) {
	if len(address) != AddressLength {
		return "", ErrInvalidAddressLength
	}
	return Ensure0xPrefix(hex.EncodeToString(address)), nil
}

func HexFromBytes(bytes []byte) string {
	return Ensure0xPrefix(hex.EncodeToString(bytes))
}
