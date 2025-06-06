package util

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	eth "github.com/dan13ram/wpokt-validator/eth/client"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const primaryType = "MintData"

var typesStandard = apitypes.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"MintData": {
		{
			Name: "recipient",
			Type: "address",
		},
		{
			Name: "amount",
			Type: "uint256",
		},
		{
			Name: "nonce",
			Type: "uint256",
		},
	},
}

func signTypedData(
	domainData eth.DomainData,
	mint *autogen.MintControllerMintData,
	key *ecdsa.PrivateKey,
) ([]byte, error) {

	message := apitypes.TypedDataMessage{
		"recipient": mint.Recipient.String(),
		"amount":    mint.Amount.String(),
		"nonce":     mint.Nonce.String(),
	}

	domain := apitypes.TypedDataDomain{
		Name:              domainData.Name,
		Version:           domainData.Version,
		ChainId:           math.NewHexOrDecimal256(domainData.ChainId.Int64()),
		VerifyingContract: domainData.VerifyingContract.String(),
	}

	typedData := apitypes.TypedData{
		Types:       typesStandard,
		PrimaryType: primaryType,
		Domain:      domain,
		Message:     message,
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(rawData)

	signature, err := crypto.Sign(sighash, key)
	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return signature, err
}

func UpdateStatusAndConfirmationsForMint(mint *models.Mint, poktHeight int64) (*models.Mint, error) {
	status := mint.Status
	confirmations, err := strconv.ParseInt(mint.Confirmations, 10, 64)
	if err != nil || confirmations < 0 {
		confirmations = 0
	}

	if status == models.StatusPending || confirmations == 0 {
		status = models.StatusPending
		if app.Config.Pocket.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			mintHeight, err := strconv.ParseInt(mint.Height, 10, 64)
			if err != nil {
				return mint, err
			}
			confirmations = poktHeight - mintHeight
			if confirmations >= app.Config.Pocket.Confirmations {
				status = models.StatusConfirmed
			}
		}
	}

	mint.Status = status
	mint.Confirmations = strconv.FormatInt(confirmations, 10)
	return mint, nil
}

func sortSignersAndSignatures(signers, signatures []string) ([]string, []string) {
	type SignerSignaturePair struct {
		Signer    string
		Signature string
	}

	// Pair up signers and signatures
	pairs := make([]SignerSignaturePair, len(signers))
	for i := range signers {
		pairs[i] = SignerSignaturePair{
			Signer:    signers[i],
			Signature: signatures[i],
		}
	}

	// Sort pairs based on signer
	sort.Slice(pairs, func(i, j int) bool {
		return common.HexToAddress(pairs[i].Signer).Big().Cmp(common.HexToAddress(pairs[j].Signer).Big()) == -1
	})

	// Extract sorted signers and signatures
	for i := range pairs {
		signers[i] = pairs[i].Signer
		signatures[i] = pairs[i].Signature
	}

	return signers, signatures
}

func SignMint(
	mint *models.Mint,
	data *autogen.MintControllerMintData,
	domain eth.DomainData,
	privateKey *ecdsa.PrivateKey,
	signerThreshold int,
) (*models.Mint, error) {
	signature, err := signTypedData(domain, data, privateKey)
	if err != nil {
		return mint, err
	}

	signatureEncoded := "0x" + hex.EncodeToString(signature)
	signatures := mint.Signatures
	signers := mint.Signers
	if signatures == nil || signers == nil || len(signatures) != len(signers) || len(signatures) == 0 {
		signatures = []string{}
		signers = []string{}
	}
	signatures = append(signatures, signatureEncoded)
	signers = append(signers, strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex()))

	sortedSigners, sortedSignatures := sortSignersAndSignatures(signers, signatures)

	if len(sortedSignatures) >= signerThreshold {
		mint.Status = models.StatusSigned
	}

	mint.Signatures = sortedSignatures
	mint.Signers = sortedSigners
	return mint, nil
}
