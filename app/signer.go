package app

import (
	"fmt"
	"strings"

	"bytes"
	"sort"

	"crypto/ecdsa"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	multisigtypes "github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/dan13ram/wpokt-validator/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"

	"encoding/hex"
)

type PocketSigner struct {
	Signer          common.Signer
	Address         string
	SignerIndex     int
	Multisig        multisigtypes.PubKey
	MultisigAddress string
}

func CreatePocketSigner() (common.Signer, error) {
	config := Config.Pocket
	if config.Mnemonic == "" && config.GcpKmsKeyName == "" {
		return nil, fmt.Errorf("both Mnemonic and GcpKmsKeyName are empty")
	}
	if config.Mnemonic != "" {
		return common.NewMnemonicSigner(config.Mnemonic)
	}

	return common.NewGcpKmsSigner(config.GcpKmsKeyName)

}

func GetPocketSignerAndMultisig() (*PocketSigner, error) {
	signer, err := CreatePocketSigner()
	if err != nil {
		return nil, fmt.Errorf("error initializing pokt signer: %w", err)
	}

	cosmosPubKey := signer.CosmosPublicKey()

	hexPubKey := hex.EncodeToString(cosmosPubKey.Bytes())
	log.Debugf("[SIGNER] Pocket public key: %s", hexPubKey)

	poktAddress, err := common.Bech32FromBytes(Config.Pocket.Bech32Prefix, cosmosPubKey.Address().Bytes())
	if err != nil {
		return nil, fmt.Errorf("error getting pokt address: %w", err)
	}

	log.Debugf("[SIGNER] Pocket address: %s", poktAddress)

	var pks []crypto.PubKey
	signerIndex := -1
	for index, pk := range Config.Pocket.MultisigPublicKeys {
		pKey, err := common.CosmosPublicKeyFromHex(pk)
		if err != nil {
			return nil, fmt.Errorf("error parsing multisig public key [%d]: %w", index, err)
		}
		pks = append(pks, pKey)
		if pKey.Equals(cosmosPubKey) {
			signerIndex = index
			log.Debugf("[SIGNER] Found current pocket signer at index %d", index)
		}
	}

	if signerIndex == -1 {
		return nil, fmt.Errorf("could not find current signer in list of multisig public keys")
	}

	if Config.Pocket.MultisigThreshold == 0 || Config.Pocket.MultisigThreshold > uint64(len(Config.Pocket.MultisigPublicKeys)) {
		return nil, fmt.Errorf("multisig threshold is invalid")
	}

	sort.Slice(pks, func(i, j int) bool {
		return bytes.Compare(pks[i].Address(), pks[j].Address()) < 0
	})

	multisigPk := multisig.NewLegacyAminoPubKey(int(Config.Pocket.MultisigThreshold), pks)
	multisigAddressBytes := multisigPk.Address().Bytes()
	multisigAddress, _ := common.Bech32FromBytes(Config.Pocket.Bech32Prefix, multisigAddressBytes)

	log.Debugf("[SIGNER] Pocket multisig address: %s", multisigAddress)

	if !strings.EqualFold(multisigAddress, Config.Pocket.MultisigAddress) {
		return nil, fmt.Errorf("multisig address does not match vault address")
	}

	return &PocketSigner{
		Signer:          signer,
		SignerIndex:     signerIndex,
		Multisig:        multisigPk,
		MultisigAddress: multisigAddress,
		Address:         poktAddress,
	}, nil
}

type EthereumSigner struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

func GetEthereumSigner() (*EthereumSigner, error) {
	ethPK, err := ethCrypto.HexToECDSA(Config.Ethereum.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error initializing ethereum signer: %w", err)
	}

	ethAddress := ethCrypto.PubkeyToAddress(ethPK.PublicKey).Hex()

	return &EthereumSigner{
		PrivateKey: ethPK,
		Address:    ethAddress,
	}, nil
}
