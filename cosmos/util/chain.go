package util

import (
	"math/big"

	"github.com/dan13ram/wpokt-validator/models"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// hash the string chainID to get a uint64 chainDomain
func getChainDomain(chainID string) uint32 {
	chainHash := ethcrypto.Keccak256([]byte(chainID))
	chainDomain := new(big.Int).SetBytes(chainHash).Uint64()
	return uint32(chainDomain)
}

func ParseChain(config models.PocketConfig) models.Chain {
	return models.Chain{
		ChainID:     config.ChainID,
		ChainDomain: getChainDomain(config.ChainID),
		ChainName:   config.ChainName,
		ChainType:   models.ChainTypeCosmos,
	}
}
