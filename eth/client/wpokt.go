package client

import (
	"math/big"

	"github.com/dan13ram/wpokt-validator/eth/autogen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type WrappedPocketContract interface {
	GetUserNonce(opts *bind.CallOpts, user common.Address) (*big.Int, error)
	FilterMinted(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) (WrappedPocketMintedIterator, error)
	FilterBurnAndBridge(opts *bind.FilterOpts, amount []*big.Int, poktAddress []common.Address, from []common.Address) (WrappedPocketBurnAndBridgeIterator, error)
	ParseBurnAndBridge(log types.Log) (*autogen.WrappedPocketBurnAndBridge, error)
}

type WrappedPocketBurnAndBridgeIterator interface {
	Next() bool
	Event() *autogen.WrappedPocketBurnAndBridge
	Close() error
	Error() error
}

type WrappedPocketBurnAndBridgeIteratorImpl struct {
	iterator *autogen.WrappedPocketBurnAndBridgeIterator
}

func (x *WrappedPocketBurnAndBridgeIteratorImpl) Next() bool {
	return x.iterator.Next()
}

func (x *WrappedPocketBurnAndBridgeIteratorImpl) Event() *autogen.WrappedPocketBurnAndBridge {
	return x.iterator.Event
}

func (x *WrappedPocketBurnAndBridgeIteratorImpl) Close() error {
	return x.iterator.Close()
}

func (x *WrappedPocketBurnAndBridgeIteratorImpl) Error() error {
	return x.iterator.Error()
}

type WrappedPocketMintedIterator interface {
	Next() bool
	Event() *autogen.WrappedPocketMinted
	Close() error
	Error() error
}

type WrappedPocketMintedIteratorImpl struct {
	iterator *autogen.WrappedPocketMintedIterator
}

func (x *WrappedPocketMintedIteratorImpl) Next() bool {
	return x.iterator.Next()
}

func (x *WrappedPocketMintedIteratorImpl) Event() *autogen.WrappedPocketMinted {
	return x.iterator.Event
}

func (x *WrappedPocketMintedIteratorImpl) Close() error {
	return x.iterator.Close()
}

func (x *WrappedPocketMintedIteratorImpl) Error() error {
	return x.iterator.Error()
}

type WrappedPocketContractImpl struct {
	contract *autogen.WrappedPocket
}

func (x *WrappedPocketContractImpl) ParseBurnAndBridge(log types.Log) (*autogen.WrappedPocketBurnAndBridge, error) {
	return x.contract.ParseBurnAndBridge(log)
}

func (x *WrappedPocketContractImpl) FilterBurnAndBridge(opts *bind.FilterOpts, amount []*big.Int, poktAddress []common.Address, from []common.Address) (WrappedPocketBurnAndBridgeIterator, error) {
	iterator, err := x.contract.FilterBurnAndBridge(opts, amount, poktAddress, from)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketBurnAndBridgeIteratorImpl{iterator: iterator}, nil
}

func (x *WrappedPocketContractImpl) FilterMinted(opts *bind.FilterOpts, recipient []common.Address, amount []*big.Int, nonce []*big.Int) (WrappedPocketMintedIterator, error) {
	iterator, err := x.contract.FilterMinted(opts, recipient, amount, nonce)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketMintedIteratorImpl{iterator: iterator}, nil
}

func (x *WrappedPocketContractImpl) GetUserNonce(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	return x.contract.GetUserNonce(opts, user)
}

func NewWrappedPocketContract(contract *autogen.WrappedPocket) WrappedPocketContract {
	return &WrappedPocketContractImpl{contract: contract}
}
