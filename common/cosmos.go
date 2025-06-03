package common

import (
	"encoding/hex"
	"errors"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

var (
	defaultAlgo               = hd.Secp256k1
	ErrInvalidPublicKeyLength = errors.New("invalid public key length")
)

func CosmosPrivateKeyFromMnemonic(mnemonic string) (crypto.PrivKey, error) {
	derivedPriv, err := defaultAlgo.Derive()(mnemonic, DefaultBIP39Passphrase, DefaultCosmosHDPath)
	if err != nil {
		return nil, err
	}

	privKey := defaultAlgo.Generate()(derivedPriv)

	return privKey, nil
}

func CosmosPublicKeyFromMnemonic(mnemonic string) (crypto.PubKey, error) {
	derivedPriv, err := defaultAlgo.Derive()(mnemonic, DefaultBIP39Passphrase, DefaultCosmosHDPath)
	if err != nil {
		return nil, err
	}

	privKey := defaultAlgo.Generate()(derivedPriv)

	pubKey := privKey.PubKey()

	return pubKey, nil
}

func CosmosPublicKeyHexFromMnemonic(mnemonic string) (string, error) {
	derivedPriv, err := defaultAlgo.Derive()(mnemonic, DefaultBIP39Passphrase, DefaultCosmosHDPath)
	if err != nil {
		return "", err
	}

	privKey := defaultAlgo.Generate()(derivedPriv)

	pubKey := privKey.PubKey()

	return hex.EncodeToString(pubKey.Bytes()), nil
}

func CosmosPublicKeyFromHex(pubKeyHex string) (crypto.PubKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}

	if len(pubKeyBytes) != CosmosPublicKeyLength {
		return nil, ErrInvalidPublicKeyLength
	}

	pubKey := &secp256k1.PubKey{}
	err = pubKey.UnmarshalAmino(pubKeyBytes)

	return pubKey, err
}

func IsValidCosmosPublicKey(s string) bool {
	_, err := CosmosPublicKeyFromHex(s)
	return err == nil
}

func BytesFromBech32(bech32Prefix string, address string) (addr []byte, err error) {
	addrCdc := addresscodec.NewBech32Codec(bech32Prefix)
	return addrCdc.StringToBytes(address)
}

func Bech32FromBytes(bech32Prefix string, bs []byte) (string, error) {
	if len(bs) != AddressLength {
		return "", ErrInvalidAddressLength
	}
	return bech32.ConvertAndEncode(bech32Prefix, bs)
}

func AddressBytesFromBech32(bech32Prefix string, address string) ([]byte, error) {
	bytes, err := BytesFromBech32(bech32Prefix, address)
	if err != nil {
		return nil, err
	}
	if len(bytes) != AddressLength {
		return nil, ErrInvalidAddressLength
	}
	return bytes, nil
}

func IsValidBech32Address(bech32Prefix string, address string) bool {
	_, err := AddressBytesFromBech32(bech32Prefix, address)
	return err == nil
}
