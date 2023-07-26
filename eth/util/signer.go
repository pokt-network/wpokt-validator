package util

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func sortAddresses(addresses []string) []string {
	for i, address := range addresses {
		addresses[i] = common.HexToAddress(address).Hex()
	}
	sort.Slice(addresses, func(i, j int) bool {
		return common.HexToAddress(addresses[i]).Big().Cmp(common.HexToAddress(addresses[j]).Big()) == -1
	})
	return addresses
}

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

type DomainData struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

func signTypedData(
	domainData DomainData,
	mint autogen.MintControllerMintData,
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
	if err != nil {
		return nil, err
	}
	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return signature, nil
}

func UpdateStatusAndConfirmationsForMint(mint models.Mint, poktHeight int64) (models.Mint, error) {
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

func SignMint(
	mint models.Mint,
	data autogen.MintControllerMintData,
	domain DomainData,
	privateKey *ecdsa.PrivateKey,
	numSigners int,
) (models.Mint, error) {
	signature, err := signTypedData(domain, data, privateKey)
	if err != nil {
		return mint, err
	}

	signatureEncoded := "0x" + hex.EncodeToString(signature)
	if mint.Signatures == nil {
		mint.Signatures = []string{}
	}
	if mint.Signers == nil {
		mint.Signers = []string{}
	}
	signatures := append(mint.Signatures, signatureEncoded)
	signers := append(mint.Signers, crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	sortedSigners := sortAddresses(signers)

	sortedSignatures := make([]string, len(signatures))

	for i, signature := range signatures {
		signer := signers[i]
		index := -1
		for j, validator := range sortedSigners {
			if validator == signer {
				index = j
				break
			}
		}
		sortedSignatures[index] = signature
	}

	if len(sortedSignatures) == numSigners {
		mint.Status = models.StatusSigned
	}

	mint.Signatures = sortedSignatures
	mint.Signers = sortedSigners
	return mint, nil
}
