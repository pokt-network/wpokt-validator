package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTxConfig(t *testing.T) {
	bech32Prefix := "pokt"
	txConfig := NewTxConfig(bech32Prefix)

	assert.NotNil(t, txConfig)
}

func TestNewProtoCodec(t *testing.T) {
	bech32Prefix := "pokt"
	protoCodec := NewProtoCodec(bech32Prefix)

	assert.NotNil(t, protoCodec)
}

func TestNewTxDecoder(t *testing.T) {
	bech32Prefix := "pokt"
	txDecoder := NewTxDecoder(bech32Prefix)

	assert.NotNil(t, txDecoder)
}
