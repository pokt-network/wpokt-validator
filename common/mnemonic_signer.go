package common

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/types"
)

// Struct Definition
type MnemonicSigner struct {
	ethAddress    common.Address
	cosmosPubKey  types.PubKey
	ethPrivKey    *ecdsa.PrivateKey
	cosmosPrivKey types.PrivKey
}

var _ Signer = &MnemonicSigner{}

// Constructor Function
func NewMnemonicSigner(mnemonic string) (*MnemonicSigner, error) {

	ethPrivKey, err := EthereumPrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create ethereum private key: %w", err)
	}

	publicKeyECDSA, _ := ethPrivKey.Public().(*ecdsa.PublicKey) // impossible to get an error since the private key is not nil

	ethAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	cosmosPrivKey, err := CosmosPrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create cosmos private key: %w", err)
	}

	cosmosPubKey := cosmosPrivKey.PubKey()

	return &MnemonicSigner{
		ethPrivKey:    ethPrivKey,
		ethAddress:    ethAddress,
		cosmosPrivKey: cosmosPrivKey,
		cosmosPubKey:  cosmosPubKey,
	}, nil
}

// Destructor Function
func (s *MnemonicSigner) Destroy() {
	// nothing to do
}

// Method Implementations
func (s *MnemonicSigner) EthSign(data []byte) ([]byte, error) {
	digest := data
	if len(digest) != 32 {
		digest = crypto.Keccak256(data)
	}
	hash := common.BytesToHash(digest)
	signature, err := crypto.Sign(hash[:], s.ethPrivKey)
	if err != nil {
		return nil, err
	}

	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return signature, nil
}

func (s *MnemonicSigner) CosmosSign(data []byte) ([]byte, error) {
	return s.cosmosPrivKey.Sign(data[:])
}

func (s *MnemonicSigner) EthAddress() common.Address {
	return s.ethAddress
}

func (s *MnemonicSigner) CosmosPublicKey() types.PubKey {
	return s.cosmosPubKey
}
