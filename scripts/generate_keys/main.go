package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/dan13ram/wpokt-validator/common"

	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
)

var (
	DefaultAlgo = hd.Secp256k1
)

const (
	DefaultEntropySize     = 256
	DefaultBIP39Passphrase = ""
	DefaultCosmosHDPath    = "m/44'/118'/0'/0/0"
	DefaultBech32Prefix    = "pokt"

	DefaultETHHDPath = "m/44'/60'/0'/0/0"
)

func NewMnemonic() (string, error) {
	// Default number of words (24): This generates a mnemonic directly from the
	// number of words by reading system entropy.
	entropy, err := bip39.NewEntropy(DefaultEntropySize)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func NewAccount(mnemonic string) (crypto.PrivKey, error) {

	// create master key and derive first key for keyring
	derivedPriv, err := DefaultAlgo.Derive()(mnemonic, DefaultBIP39Passphrase, DefaultCosmosHDPath)
	if err != nil {
		return nil, err
	}

	privKey := DefaultAlgo.Generate()(derivedPriv)

	return privKey, nil
}

func main() {

	var mnemonic string
	var insecure bool
	var err error
	flag.StringVar(&mnemonic, "mnemonic", "", "24 word mnemonic")
	flag.BoolVar(&insecure, "insecure", false, "insecure mode")
	flag.Parse()

	if mnemonic == "" {
		mnemonic, err = NewMnemonic()
		if err != nil {
			fmt.Printf("Error generating mnemonic: %v\n", err)
			return
		}
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		fmt.Println("mnemonic is invalid")
		return
	}

	fields := strings.Fields(mnemonic)
	if !insecure && len(fields) != 24 {
		fmt.Println("mnemonic is invalid, must be 24 words")
		return
	}

	fmt.Println("mnemonic: ", mnemonic)
	fmt.Println("num words: ", len(fields))
	fmt.Println()

	{
		// cosmos
		privKey, err := NewAccount(mnemonic)
		if err != nil {
			fmt.Printf("Error generating account: %v\n", err)
			return
		}

		pubKey := privKey.PubKey()

		address := pubKey.Address()

		fmt.Println("cosmos private key: ", hex.EncodeToString(privKey.Bytes()))
		fmt.Println("cosmos public key: ", hex.EncodeToString(pubKey.Bytes()))

		fmt.Println("commos address: ", hex.EncodeToString(address.Bytes()))

		bech32, err := common.Bech32FromBytes(DefaultBech32Prefix, address.Bytes())
		if err != nil {
			fmt.Printf("Error converting address to bech32: %v\n", err)
			return
		}

		fmt.Println("cosmos bech32 address: ", bech32)
	}

	fmt.Println()

	{
		// ethereum
		wallet, err := hdwallet.NewFromMnemonic(mnemonic)
		if err != nil {
			fmt.Printf("Error creating wallet: %v\n", err)
			return
		}

		path := hdwallet.MustParseDerivationPath(common.DefaultETHHDPath)
		account, err := wallet.Derive(path, false)
		if err != nil {
			fmt.Printf("Error deriving account: %v\n", err)
			return
		}

		privKey, err := wallet.PrivateKeyHex(account)
		if err != nil {
			fmt.Printf("Error getting private key: %v\n", err)
			return
		}
		pubKey, err := wallet.PublicKeyHex(account)
		if err != nil {
			fmt.Printf("Error getting public key: %v\n", err)
			return
		}

		address, err := wallet.AddressHex(account)
		if err != nil {
			fmt.Printf("Error getting address: %v\n", err)
			return
		}

		fmt.Println("eth private key: ", privKey)
		fmt.Println("eth public key: ", pubKey)
		fmt.Println("eth address: ", address)
	}

}
