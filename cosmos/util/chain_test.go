package util

import (
	"testing"

	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"
)

func TestGetChainDomain(t *testing.T) {
	chainID := "test-chain-id"
	expectedDomain := uint32(0x13a153d) // Replace with the expected domain based on your specific chainID

	chainDomain := getChainDomain(chainID)
	assert.Equal(t, expectedDomain, chainDomain)
}

func TestParseChain(t *testing.T) {
	config := models.PocketConfig{
		ChainID:   "cosmoshub-4",
		ChainName: "Cosmos Hub",
	}
	expectedChain := models.Chain{
		ChainID:     "cosmoshub-4",
		ChainDomain: getChainDomain("cosmoshub-4"),
		ChainName:   "Cosmos Hub",
		ChainType:   models.ChainTypeCosmos,
	}

	chain := ParseChain(config)
	assert.Equal(t, expectedChain, chain)
}
