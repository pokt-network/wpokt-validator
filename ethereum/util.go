package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

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

func getDomain(chainId int64, verifyingContract string) apitypes.TypedDataDomain {
	return apitypes.TypedDataDomain{
		Name:              "MintController",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(chainId),
		VerifyingContract: verifyingContract,
	}
}

const primaryType = "MintData"

type DomainData struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

func SignTypedData(
	domainData DomainData,
	mint autogen.MintControllerMintData, key *ecdsa.PrivateKey,
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
