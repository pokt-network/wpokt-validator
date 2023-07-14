package pocket

import (
	"encoding/hex"

	"github.com/dchest/uniuri"
	pokt "github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/libs/rand"
)

const legacyCodec bool = false

var txEncoder sdk.TxEncoder = auth.DefaultTxEncoder(pokt.Codec())
var txDecoder sdk.TxDecoder = auth.DefaultTxDecoder(pokt.Codec())

var passphrase string = uniuri.NewLen(32)

func BuildMultiSigTxAndSign(
	toAddr string,
	memo string,
	chainID string,
	amount int64,
	fees int64,
	signerKey crypto.PrivateKey,
	multisigKey crypto.PublicKeyMultiSig,
) ([]byte, error) {

	fa, err := sdk.AddressFromHex(multisigKey.Address().String())
	if err != nil {
		return nil, err
	}

	ta, err := sdk.AddressFromHex(toAddr)
	if err != nil {
		return nil, err
	}

	m := &nodeTypes.MsgSend{
		FromAddress: fa,
		ToAddress:   ta,
		Amount:      sdk.NewInt(amount),
	}

	entropy := rand.Int64()
	fee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(fees)))

	signBz, err := authTypes.StdSignBytes(chainID, entropy, fee, m, memo)
	if err != nil {
		return nil, err
	}

	sigBytes, err := signerKey.Sign(signBz)
	if err != nil {
		return nil, err
	}

	// sign using multisignature structure
	var ms = crypto.MultiSig(crypto.MultiSignature{})
	ms = ms.NewMultiSignature()

	// loop through all the keys and add signatures
	for i := 0; i < len(multisigKey.Keys()); i++ {
		ms = ms.AddSignatureByIndex(sigBytes, i)
		// when new signatures are added they will replace the old ones
	}

	sig := authTypes.StdSignature{
		PublicKey: multisigKey,
		Signature: ms.Marshal(),
	}

	// create a new standard transaction object
	tx := authTypes.NewTx(m, fee, sig, memo, entropy)

	// encode it using the default encoder
	return txEncoder(tx, -1)
}

func SignMultisigTx(
	txHex string,
	chainID string,
	signerKey crypto.PrivateKey,
	multisigKey crypto.PublicKeyMultiSig,
) ([]byte, error) {

	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}

	t, err := txDecoder(bz, -1)
	if err != nil {
		return nil, err
	}

	tx := t.(authTypes.StdTx)

	bytesToSign, err := authTypes.StdSignBytes(chainID, tx.GetEntropy(), tx.GetFee(), tx.GetMsg(), tx.GetMemo())
	if err != nil {
		return nil, err
	}

	sigBytes, err := signerKey.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}

	var ms = crypto.MultiSig(crypto.MultiSignature{})

	if tx.GetSignature().GetSignature() == nil || len(tx.GetSignature().GetSignature()) == 0 {
		ms = ms.NewMultiSignature()
	} else {
		ms = ms.Unmarshal(tx.GetSignature().GetSignature())
	}

	ms, err = ms.AddSignature(sigBytes,
		signerKey.PublicKey(), multisigKey.Keys())

	if err != nil {
		return nil, err
	}

	sig := authTypes.StdSignature{
		PublicKey: tx.Signature.PublicKey,
		Signature: ms.Marshal(),
	}

	// replace the old multi-signature with the new multi-signature (containing the additional signature)
	tx, err = tx.WithSignature(sig)
	if err != nil {
		return nil, err
	}

	// encode using the standard encoder
	return txEncoder(tx, -1)
}
