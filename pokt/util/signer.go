package util

import (
	"encoding/hex"
	"strconv"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
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

func buildMultiSigTxAndSign(
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

func decodeTx(txHex string, chainID string) (authTypes.StdTx, []byte, error) {
	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return authTypes.StdTx{}, nil, err
	}

	t, err := txDecoder(bz, -1)
	if err != nil {
		return authTypes.StdTx{}, nil, err
	}

	tx := t.(authTypes.StdTx)

	bytesToSign, err := authTypes.StdSignBytes(chainID, tx.GetEntropy(), tx.GetFee(), tx.GetMsg(), tx.GetMemo())
	if err != nil {
		return authTypes.StdTx{}, nil, err
	}

	return tx, bytesToSign, nil
}

func signMultisigTx(
	txHex string,
	chainID string,
	signerKey crypto.PrivateKey,
	multisigKey crypto.PublicKeyMultiSig,
) ([]byte, error) {

	tx, bytesToSign, err := decodeTx(txHex, chainID)
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

	ms, err = ms.AddSignature(
		sigBytes,
		signerKey.PublicKey(),
		multisigKey.Keys(),
	)

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

func UpdateStatusAndConfirmationsForInvalidMint(doc models.InvalidMint, currentHeight int64) (models.InvalidMint, error) {
	status := doc.Status
	confirmations, err := strconv.ParseInt(doc.Confirmations, 10, 64)
	if err != nil || confirmations < 0 {
		confirmations = 0
	}

	if status == models.StatusPending || confirmations == 0 {
		status = models.StatusPending
		if app.Config.Pocket.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			mintHeight, err := strconv.ParseInt(doc.Height, 10, 64)
			if err != nil {
				return doc, err
			}
			confirmations = currentHeight - mintHeight
			if confirmations >= app.Config.Pocket.Confirmations {
				status = models.StatusConfirmed
			}
		}
	}

	doc.Status = status
	doc.Confirmations = strconv.FormatInt(confirmations, 10)

	return doc, nil
}

func SignInvalidMint(
	doc models.InvalidMint,
	privateKey crypto.PrivateKey,
	multisigPubKey crypto.PublicKeyMultiSig,
	numSigners int,
) (models.InvalidMint, error) {
	returnTx := doc.ReturnTx
	signers := doc.Signers

	if signers == nil || returnTx == "" || len(signers) == 0 {
		signers = []string{}
		returnTx = ""
	}

	if returnTx == "" {
		amountWithTxFee, err := strconv.ParseInt(doc.Amount, 10, 64)
		if err != nil {
			return doc, err
		}
		amount := amountWithTxFee - app.Config.Pocket.TxFee
		memo := doc.TransactionHash

		returnTxBytes, err := buildMultiSigTxAndSign(
			doc.SenderAddress,
			memo,
			app.Config.Pocket.ChainId,
			amount,
			app.Config.Pocket.TxFee,
			privateKey,
			multisigPubKey,
		)
		if err != nil {
			return doc, err
		}
		returnTx = hex.EncodeToString(returnTxBytes)
	} else {
		returnTxBytes, err := signMultisigTx(
			returnTx,
			app.Config.Pocket.ChainId,
			privateKey,
			multisigPubKey,
		)
		if err != nil {
			return doc, err
		}
		returnTx = hex.EncodeToString(returnTxBytes)
	}

	signers = append(signers, privateKey.PublicKey().RawString())

	if len(signers) == numSigners {
		doc.Status = models.StatusSigned
	}

	doc.ReturnTx = returnTx
	doc.Signers = signers

	return doc, nil
}

func UpdateStatusAndConfirmationsForBurn(doc models.Burn, blockNumber int64) (models.Burn, error) {
	status := doc.Status
	confirmations, err := strconv.ParseInt(doc.Confirmations, 10, 64)
	if err != nil || confirmations < 0 {
		confirmations = 0
	}

	if status == models.StatusPending || confirmations == 0 {
		status = models.StatusPending
		if app.Config.Ethereum.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			burnBlockNumber, err := strconv.ParseInt(doc.BlockNumber, 10, 64)
			if err != nil {
				return doc, err
			}

			confirmations = blockNumber - burnBlockNumber
			if confirmations >= app.Config.Ethereum.Confirmations {
				status = models.StatusConfirmed
			}
		}
	}

	doc.Status = status
	doc.Confirmations = strconv.FormatInt(confirmations, 10)
	return doc, nil
}

func SignBurn(
	doc models.Burn,
	privateKey crypto.PrivateKey,
	multisigPubKey crypto.PublicKeyMultiSig,
	numSigners int,
) (models.Burn, error) {

	signers := doc.Signers
	returnTx := doc.ReturnTx

	if signers == nil || returnTx == "" || len(signers) == 0 {
		signers = []string{}
		returnTx = ""
	}

	if returnTx == "" {
		amountWithTxFee, err := strconv.ParseInt(doc.Amount, 10, 64)
		if err != nil {
			return doc, err
		}
		amount := amountWithTxFee - app.Config.Pocket.TxFee
		memo := doc.TransactionHash

		returnTxBytes, err := buildMultiSigTxAndSign(
			doc.RecipientAddress,
			memo,
			app.Config.Pocket.ChainId,
			amount,
			app.Config.Pocket.TxFee,
			privateKey,
			multisigPubKey,
		)
		if err != nil {
			return doc, err
		}
		returnTx = hex.EncodeToString(returnTxBytes)
	} else {
		returnTxBytes, err := signMultisigTx(
			returnTx,
			app.Config.Pocket.ChainId,
			privateKey,
			multisigPubKey,
		)
		if err != nil {
			return doc, err
		}
		returnTx = hex.EncodeToString(returnTxBytes)

	}

	signers = append(signers, privateKey.PublicKey().RawString())

	if len(signers) == numSigners {
		doc.Status = models.StatusSigned
	}

	doc.Signers = signers
	doc.ReturnTx = returnTx

	return doc, nil
}
