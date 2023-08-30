package client

import (
	"math/big"

	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type MintControllerContract interface {
	ValidatorCount(opts *bind.CallOpts) (*big.Int, error)
}

type MintControllerContractImpl struct {
	contract *autogen.MintController
}

func (x *MintControllerContractImpl) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	return x.contract.ValidatorCount(opts)
}

func NewMintControllerContract(contract *autogen.MintController) MintControllerContract {
	return &MintControllerContractImpl{contract: contract}
}
