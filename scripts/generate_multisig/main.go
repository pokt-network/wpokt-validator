package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	"bytes"
	"sort"

	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/dan13ram/wpokt-validator/common"
)

const DefaultBech32Prefix = "pokt"

func main() {
	var publickeys string
	var threshold int
	flag.StringVar(&publickeys, "publickeys", "", "comma separated list of public keys")
	flag.IntVar(&threshold, "threshold", 0, "threshold for multisig")
	flag.Parse()

	if publickeys == "" {
		fmt.Printf("publickeys is required\n")
		return
	}

	var keys []string
	if publickeys != "" {
		keys = strings.Split(publickeys, ",")
	}

	if threshold <= 0 {
		fmt.Printf("threshold is required\n")
		return
	}

	if len(keys) < threshold {
		fmt.Printf("threshold must be less than or equal to the number of public keys\n")
		return
	}

	var pKeys []crypto.PubKey

	for i, key := range keys {
		if !common.IsValidCosmosPublicKey(key) {
			fmt.Printf("invalid public key %d: %v\n", i, key)
			return
		}
		pKey := &secp256k1.PubKey{}
		pKeyBytes, err := hex.DecodeString(key)
		if err != nil {
			fmt.Printf("error decoding public key %d: %v\n", i, err)
			return
		}
		err = pKey.UnmarshalAmino(pKeyBytes)
		if err != nil {
			fmt.Printf("error unmarshalling public key %d: %v\n", i, err)
			return
		}
		pKeys = append(pKeys, pKey)
		fmt.Printf("public key %d: %v\n", i, key)
		address, _ := common.Bech32FromBytes(DefaultBech32Prefix, pKey.Address().Bytes())
		addressHex := "0x" + hex.EncodeToString(pKey.Address().Bytes())
		fmt.Printf("address hex: %s\n", addressHex)
		fmt.Printf("address bech32: %s\n", address)
	}

	sort.Slice(pKeys, func(i, j int) bool {
		return bytes.Compare(pKeys[i].Address(), pKeys[j].Address()) < 0
	})

	fmt.Printf("threshold: %v\n", threshold)
	pk := multisig.NewLegacyAminoPubKey(threshold, pKeys)

	fmt.Printf("multisig address hex: %v\n", "0x"+strings.ToLower(pk.Address().String()))

	bech32, err := common.Bech32FromBytes(DefaultBech32Prefix, pk.Address().Bytes())
	if err != nil {
		fmt.Printf("error encoding address: %v\n", err)
		return
	}
	fmt.Printf("multisig address bech32: %v\n", bech32)

}
