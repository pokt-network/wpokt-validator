package common

import (
	"crypto/ecdsa"
	"errors"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/ethereum/go-ethereum/crypto"
)

func IsValidEthereumAddress(s string) bool {
	s = strings.ToLower(s)
	_, err := BytesFromAddressHex(s)
	return err == nil && strings.HasPrefix(s, "0x")
}

func EthereumPrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath(DefaultETHHDPath)

	account, _ := wallet.Derive(path, false) // impossible to get an error since the path is hardcoded

	return wallet.PrivateKey(account)
}

func EthereumAddressFromMnemonic(mnemonic string) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	path := hdwallet.MustParseDerivationPath(DefaultETHHDPath)

	account, _ := wallet.Derive(path, false) // impossible to get an error since the path is hardcoded

	return wallet.AddressHex(account)
}

func HexToAddress(hex string) ethcommon.Address {
	return ethcommon.HexToAddress(hex)
}

func EthereumPrivateKeyToAddressHex(privateKey *ecdsa.PrivateKey) (string, error) {
	if privateKey == nil || privateKey.Public() == nil {
		return "", errors.New("private key is nil or public key is nil")
	}
	publicKeyECDSA, _ := privateKey.Public().(*ecdsa.PublicKey) // impossible to get an error since the private key is not nil

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex(), nil
}
