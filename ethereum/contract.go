package ethereum

import (
	"math/big"

	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type BurnAndBridgeIterator interface {
	Next() bool
	Event() *autogen.WrappedPocketBurnAndBridge
}

type BurnAndBridgeIteratorImpl struct {
	*autogen.WrappedPocketBurnAndBridgeIterator
}

func (i *BurnAndBridgeIteratorImpl) Event() *autogen.WrappedPocketBurnAndBridge {
	return i.WrappedPocketBurnAndBridgeIterator.Event
}

func (i *BurnAndBridgeIteratorImpl) Next() bool {
	return i.WrappedPocketBurnAndBridgeIterator.Next()
}

type WrappedPocketContract interface {
	FilterBurnAndBridge(opts *bind.FilterOpts, _amount []*big.Int, _from []common.Address, _poktAddress []common.Address) (BurnAndBridgeIterator, error)
}

type WrappedPocketContractImpl struct {
	*autogen.WrappedPocket
}

func (c *WrappedPocketContractImpl) FilterBurnAndBridge(opts *bind.FilterOpts, _amount []*big.Int, _from []common.Address, _poktAddress []common.Address) (BurnAndBridgeIterator, error) {
	iterator, err := c.WrappedPocket.FilterBurnAndBridge(opts, _amount, _from, _poktAddress)
	return &BurnAndBridgeIteratorImpl{iterator}, err
}
