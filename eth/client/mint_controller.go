package client

import (
	"math/big"

	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type DomainData struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

type MintControllerContract interface {
	ValidatorCount(opts *bind.CallOpts) (*big.Int, error)
	Eip712Domain(opts *bind.CallOpts) (DomainData, error)
}

type MintControllerContractImpl struct {
	contract *autogen.MintController
}

func (x *MintControllerContractImpl) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	return x.contract.ValidatorCount(opts)
}

func (x *MintControllerContractImpl) Eip712Domain(opts *bind.CallOpts) (DomainData, error) {
	return x.contract.Eip712Domain(opts)
}

func NewMintControllerContract(contract *autogen.MintController) MintControllerContract {
	return &MintControllerContractImpl{contract: contract}
}
