package pocket

import (
	"encoding/hex"

	pokt "github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

func BuildMultiSigTxAndSign(
	signerAddr string,
	toAddr string,
	memo string,
	chainID string,
	amount int64,
	fees int64,
	kb keys.Keybase,
	pk crypto.PublicKeyMultiSig,
) ([]byte, error) {
	fa, err := sdk.AddressFromHex(pk.Address().String())
	if err != nil {
		return nil, err
	}

	sa, err := sdk.AddressFromHex(signerAddr)
	if err != nil {
		return nil, err
	}

	ta, err := sdk.AddressFromHex(toAddr)
	if err != nil {
		return nil, err
	}

	protoMsg := nodeTypes.MsgSend{
		FromAddress: fa,
		ToAddress:   ta,
		Amount:      sdk.NewInt(amount),
	}

	passphrase := "PASSPHRASE"
	legacyCodec := false

	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(pokt.Codec()),
		auth.DefaultTxDecoder(pokt.Codec()),
		chainID,
		memo,
		nil).WithKeybase(kb)

	return txBuilder.BuildAndSignMultisigTransaction(sa, pk, &protoMsg, passphrase, fees, legacyCodec)
}

func SignMultisigTx(
	signerAddr string,
	txHex string,
	chainID string,
	kb keys.Keybase,
	pk crypto.PublicKeyMultiSig,
) ([]byte, error) {
	sa, err := sdk.AddressFromHex(signerAddr)
	if err != nil {
		return nil, err
	}

	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}

	passphrase := "PASSPHRASE"
	legacyCodec := false

	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(pokt.Codec()),
		auth.DefaultTxDecoder(pokt.Codec()),
		chainID,
		"",
		nil).WithKeybase(kb)

	return txBuilder.SignMultisigTransaction(sa, pk.Keys(), passphrase, bz, legacyCodec)
}
