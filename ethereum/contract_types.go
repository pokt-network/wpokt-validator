package ethereum

import (
	"math/big"

	"github.com/dan13ram/wpokt-backend/ethereum/autogen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const MAX_QUERY_BLOCKS uint64 = 100000
const ZERO_ADDRESS string = "0x0000000000000000000000000000000000000000"

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

type TransferIterator interface {
	Next() bool
	Event() *autogen.WrappedPocketTransfer
}

type TransferIteratorImpl struct {
	*autogen.WrappedPocketTransferIterator
}

func (i *TransferIteratorImpl) Event() *autogen.WrappedPocketTransfer {
	return i.WrappedPocketTransferIterator.Event
}

func (i *TransferIteratorImpl) Next() bool {
	return i.WrappedPocketTransferIterator.Next()
}

type MintedIterator interface {
	Next() bool
	Event() *autogen.WrappedPocketMinted
}

type MintedIteratorImpl struct {
	*autogen.WrappedPocketMintedIterator
}

func (i *MintedIteratorImpl) Event() *autogen.WrappedPocketMinted {
	return i.WrappedPocketMintedIterator.Event
}

func (i *MintedIteratorImpl) Next() bool {
	return i.WrappedPocketMintedIterator.Next()
}

type WrappedPocketContract interface {
	FilterBurnAndBridge(opts *bind.FilterOpts, _amount []*big.Int, _from []common.Address, _poktAddress []common.Address) (BurnAndBridgeIterator, error)
	FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (TransferIterator, error)
	FilterMinted(opts *bind.FilterOpts, _recipient []common.Address, _amount []*big.Int, _nonce []*big.Int) (MintedIterator, error)
}

type WrappedPocketContractImpl struct {
	*autogen.WrappedPocket
}

func (c *WrappedPocketContractImpl) FilterBurnAndBridge(opts *bind.FilterOpts, _amount []*big.Int, _from []common.Address, _poktAddress []common.Address) (BurnAndBridgeIterator, error) {
	iterator, err := c.WrappedPocket.FilterBurnAndBridge(opts, _amount, _from, _poktAddress)
	return &BurnAndBridgeIteratorImpl{iterator}, err
}

func (c *WrappedPocketContractImpl) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (TransferIterator, error) {
	iterator, err := c.WrappedPocket.FilterTransfer(opts, _from, _to)
	return &TransferIteratorImpl{iterator}, err
}

func (c *WrappedPocketContractImpl) FilterMinted(opts *bind.FilterOpts, _recipient []common.Address, _amount []*big.Int, _nonce []*big.Int) (MintedIterator, error) {
	iterator, err := c.WrappedPocket.FilterMinted(opts, _recipient, _amount, _nonce)
	return &MintedIteratorImpl{iterator}, err
}
