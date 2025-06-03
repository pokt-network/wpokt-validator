package util

import (
	"github.com/cosmos/cosmos-sdk/codec"

	testutil "github.com/cosmos/cosmos-sdk/codec/testutil"

	std "github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/client"
)

func NewProtoCodec(bech32Prefix string) *codec.ProtoCodec {
	codecOptions := testutil.CodecOptions{
		AccAddressPrefix: bech32Prefix,
	}

	reg := codecOptions.NewInterfaceRegistry()

	std.RegisterInterfaces(reg)
	authtypes.RegisterInterfaces(reg)
	banktypes.RegisterInterfaces(reg)
	// TODO: add more modules' interfaces

	codec := codec.NewProtoCodec(reg)

	return codec
}

func NewTxDecoder(bech32Prefix string) sdk.TxDecoder {
	codec := NewProtoCodec(bech32Prefix)

	return authtx.DefaultTxDecoder(codec)
}

func newTxConfig(bech32Prefix string) client.TxConfig {

	codec := NewProtoCodec(bech32Prefix)
	txConfig := authtx.NewTxConfig(codec, authtx.DefaultSignModes)

	return txConfig
}

var NewTxConfig = newTxConfig
