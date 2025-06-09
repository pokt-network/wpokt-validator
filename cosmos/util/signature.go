package util

import (
	"fmt"

	"github.com/dan13ram/wpokt-validator/common"
	"github.com/dan13ram/wpokt-validator/models"
	"google.golang.org/protobuf/types/known/anypb"

	"context"

	"github.com/cosmos/cosmos-sdk/client"

	sdk "github.com/cosmos/cosmos-sdk/types"

	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	txsigning "cosmossdk.io/x/tx/signing"
)

var signMode = signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON

func SignWithPrivKey(
	ctx context.Context,
	signerData authsigning.SignerData,
	txBuilder client.TxBuilder,
	signer common.Signer,
	txConfig client.TxConfig,
	accSeq uint64,
) (sigV2 signingtypes.SignatureV2, msg []byte, err error) {

	// Generate the bytes to be signed.
	msg, err = authsigning.GetSignBytesAdapter(ctx, txConfig.SignModeHandler(), signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return sigV2, msg, err
	}

	// Sign those bytes
	signature, err := signer.CosmosSign(msg)
	if err != nil {
		return sigV2, msg, err
	}

	// Construct the SignatureV2 struct
	sigData := signingtypes.SingleSignatureData{
		SignMode:  signMode,
		Signature: signature,
	}

	sigV2 = signingtypes.SignatureV2{
		PubKey:   signer.CosmosPublicKey(),
		Data:     &sigData,
		Sequence: accSeq,
	}

	return sigV2, msg, nil
}

func ValidateSignature(
	config models.CosmosConfig,
	sig *signingtypes.SignatureV2,
	accountNumber uint64,
	sequence uint64,
	txConfig client.TxConfig,
	txBuilder client.TxBuilder,
) error {
	anyPk, err := codectypes.NewAnyWithValue(sig.PubKey)
	if err != nil {
		return fmt.Errorf("error creating any pubkey: %w", err)
	}
	txSignerData := txsigning.SignerData{
		ChainID:       config.ChainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		Address:       sdk.AccAddress(sig.PubKey.Address()).String(),
		PubKey: &anypb.Any{
			TypeUrl: anyPk.TypeUrl,
			Value:   anyPk.Value,
		},
	}
	builtTx := txBuilder.GetTx()
	adaptableTx, ok := builtTx.(authsigning.V2AdaptableTx)
	if !ok {
		return fmt.Errorf("expected Tx to be signing.V2AdaptableTx")
	}
	txData := adaptableTx.GetSigningTxData()

	err = authsigning.VerifySignature(context.Background(), sig.PubKey, txSignerData, sig.Data, txConfig.SignModeHandler(), txData)
	if err != nil {
		addr, _ := common.Bech32FromBytes(config.Bech32Prefix, sig.PubKey.Address().Bytes())
		return fmt.Errorf("couldn't verify signature for address %s", addr)
	}
	return nil
}
