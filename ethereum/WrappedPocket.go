// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// WrappedPocketMetaData contains all meta data concerning the WrappedPocket contract.
var WrappedPocketMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BatchMintLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BlockBurn\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeBasisDust\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeCollectorZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxBasis\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"UserNonce\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"poktAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"BurnAndBridge\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeCollector\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeCollected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newFeeBasis\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeCollector\",\"type\":\"address\"}],\"name\":\"FeeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BASIS_POINTS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_FEE_BASIS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amount\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"nonce\",\"type\":\"uint256[]\"}],\"name\":\"batchMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"poktAddress\",\"type\":\"address\"}],\"name\":\"burnAndBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeBasis\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeFlag\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newCollector\",\"type\":\"address\"}],\"name\":\"setFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// WrappedPocketABI is the input ABI used to generate the binding from.
// Deprecated: Use WrappedPocketMetaData.ABI instead.
var WrappedPocketABI = WrappedPocketMetaData.ABI

// WrappedPocket is an auto generated Go binding around an Ethereum contract.
type WrappedPocket struct {
	WrappedPocketCaller     // Read-only binding to the contract
	WrappedPocketTransactor // Write-only binding to the contract
	WrappedPocketFilterer   // Log filterer for contract events
}

// WrappedPocketCaller is an auto generated read-only Go binding around an Ethereum contract.
type WrappedPocketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedPocketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WrappedPocketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedPocketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WrappedPocketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedPocketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WrappedPocketSession struct {
	Contract     *WrappedPocket    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WrappedPocketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WrappedPocketCallerSession struct {
	Contract *WrappedPocketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// WrappedPocketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WrappedPocketTransactorSession struct {
	Contract     *WrappedPocketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// WrappedPocketRaw is an auto generated low-level Go binding around an Ethereum contract.
type WrappedPocketRaw struct {
	Contract *WrappedPocket // Generic contract binding to access the raw methods on
}

// WrappedPocketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WrappedPocketCallerRaw struct {
	Contract *WrappedPocketCaller // Generic read-only contract binding to access the raw methods on
}

// WrappedPocketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WrappedPocketTransactorRaw struct {
	Contract *WrappedPocketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWrappedPocket creates a new instance of WrappedPocket, bound to a specific deployed contract.
func NewWrappedPocket(address common.Address, backend bind.ContractBackend) (*WrappedPocket, error) {
	contract, err := bindWrappedPocket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WrappedPocket{WrappedPocketCaller: WrappedPocketCaller{contract: contract}, WrappedPocketTransactor: WrappedPocketTransactor{contract: contract}, WrappedPocketFilterer: WrappedPocketFilterer{contract: contract}}, nil
}

// NewWrappedPocketCaller creates a new read-only instance of WrappedPocket, bound to a specific deployed contract.
func NewWrappedPocketCaller(address common.Address, caller bind.ContractCaller) (*WrappedPocketCaller, error) {
	contract, err := bindWrappedPocket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketCaller{contract: contract}, nil
}

// NewWrappedPocketTransactor creates a new write-only instance of WrappedPocket, bound to a specific deployed contract.
func NewWrappedPocketTransactor(address common.Address, transactor bind.ContractTransactor) (*WrappedPocketTransactor, error) {
	contract, err := bindWrappedPocket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketTransactor{contract: contract}, nil
}

// NewWrappedPocketFilterer creates a new log filterer instance of WrappedPocket, bound to a specific deployed contract.
func NewWrappedPocketFilterer(address common.Address, filterer bind.ContractFilterer) (*WrappedPocketFilterer, error) {
	contract, err := bindWrappedPocket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketFilterer{contract: contract}, nil
}

// bindWrappedPocket binds a generic wrapper to an already deployed contract.
func bindWrappedPocket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WrappedPocketMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappedPocket *WrappedPocketRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedPocket.Contract.WrappedPocketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappedPocket *WrappedPocketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedPocket.Contract.WrappedPocketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappedPocket *WrappedPocketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedPocket.Contract.WrappedPocketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappedPocket *WrappedPocketCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedPocket.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappedPocket *WrappedPocketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedPocket.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappedPocket *WrappedPocketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedPocket.Contract.contract.Transact(opts, method, params...)
}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) BASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) BASISPOINTS() (*big.Int, error) {
	return _WrappedPocket.Contract.BASISPOINTS(&_WrappedPocket.CallOpts)
}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) BASISPOINTS() (*big.Int, error) {
	return _WrappedPocket.Contract.BASISPOINTS(&_WrappedPocket.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _WrappedPocket.Contract.DEFAULTADMINROLE(&_WrappedPocket.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _WrappedPocket.Contract.DEFAULTADMINROLE(&_WrappedPocket.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_WrappedPocket *WrappedPocketSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _WrappedPocket.Contract.DOMAINSEPARATOR(&_WrappedPocket.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _WrappedPocket.Contract.DOMAINSEPARATOR(&_WrappedPocket.CallOpts)
}

// MAXFEEBASIS is a free data retrieval call binding the contract method 0x8312ebd1.
//
// Solidity: function MAX_FEE_BASIS() view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) MAXFEEBASIS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "MAX_FEE_BASIS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXFEEBASIS is a free data retrieval call binding the contract method 0x8312ebd1.
//
// Solidity: function MAX_FEE_BASIS() view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) MAXFEEBASIS() (*big.Int, error) {
	return _WrappedPocket.Contract.MAXFEEBASIS(&_WrappedPocket.CallOpts)
}

// MAXFEEBASIS is a free data retrieval call binding the contract method 0x8312ebd1.
//
// Solidity: function MAX_FEE_BASIS() view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) MAXFEEBASIS() (*big.Int, error) {
	return _WrappedPocket.Contract.MAXFEEBASIS(&_WrappedPocket.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketSession) MINTERROLE() ([32]byte, error) {
	return _WrappedPocket.Contract.MINTERROLE(&_WrappedPocket.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCallerSession) MINTERROLE() ([32]byte, error) {
	return _WrappedPocket.Contract.MINTERROLE(&_WrappedPocket.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketSession) PAUSERROLE() ([32]byte, error) {
	return _WrappedPocket.Contract.PAUSERROLE(&_WrappedPocket.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_WrappedPocket *WrappedPocketCallerSession) PAUSERROLE() ([32]byte, error) {
	return _WrappedPocket.Contract.PAUSERROLE(&_WrappedPocket.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.Allowance(&_WrappedPocket.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.Allowance(&_WrappedPocket.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.BalanceOf(&_WrappedPocket.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.BalanceOf(&_WrappedPocket.CallOpts, account)
}

// Burn is a free data retrieval call binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 ) pure returns()
func (_WrappedPocket *WrappedPocketCaller) Burn(opts *bind.CallOpts, arg0 *big.Int) error {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "burn", arg0)

	if err != nil {
		return err
	}

	return err

}

// Burn is a free data retrieval call binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 ) pure returns()
func (_WrappedPocket *WrappedPocketSession) Burn(arg0 *big.Int) error {
	return _WrappedPocket.Contract.Burn(&_WrappedPocket.CallOpts, arg0)
}

// Burn is a free data retrieval call binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 ) pure returns()
func (_WrappedPocket *WrappedPocketCallerSession) Burn(arg0 *big.Int) error {
	return _WrappedPocket.Contract.Burn(&_WrappedPocket.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedPocket *WrappedPocketCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedPocket *WrappedPocketSession) Decimals() (uint8, error) {
	return _WrappedPocket.Contract.Decimals(&_WrappedPocket.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedPocket *WrappedPocketCallerSession) Decimals() (uint8, error) {
	return _WrappedPocket.Contract.Decimals(&_WrappedPocket.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_WrappedPocket *WrappedPocketCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_WrappedPocket *WrappedPocketSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _WrappedPocket.Contract.Eip712Domain(&_WrappedPocket.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_WrappedPocket *WrappedPocketCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _WrappedPocket.Contract.Eip712Domain(&_WrappedPocket.CallOpts)
}

// FeeBasis is a free data retrieval call binding the contract method 0x5a94ee46.
//
// Solidity: function feeBasis() view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) FeeBasis(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "feeBasis")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeBasis is a free data retrieval call binding the contract method 0x5a94ee46.
//
// Solidity: function feeBasis() view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) FeeBasis() (*big.Int, error) {
	return _WrappedPocket.Contract.FeeBasis(&_WrappedPocket.CallOpts)
}

// FeeBasis is a free data retrieval call binding the contract method 0x5a94ee46.
//
// Solidity: function feeBasis() view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) FeeBasis() (*big.Int, error) {
	return _WrappedPocket.Contract.FeeBasis(&_WrappedPocket.CallOpts)
}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_WrappedPocket *WrappedPocketCaller) FeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "feeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_WrappedPocket *WrappedPocketSession) FeeCollector() (common.Address, error) {
	return _WrappedPocket.Contract.FeeCollector(&_WrappedPocket.CallOpts)
}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_WrappedPocket *WrappedPocketCallerSession) FeeCollector() (common.Address, error) {
	return _WrappedPocket.Contract.FeeCollector(&_WrappedPocket.CallOpts)
}

// FeeFlag is a free data retrieval call binding the contract method 0x118c6897.
//
// Solidity: function feeFlag() view returns(bool)
func (_WrappedPocket *WrappedPocketCaller) FeeFlag(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "feeFlag")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FeeFlag is a free data retrieval call binding the contract method 0x118c6897.
//
// Solidity: function feeFlag() view returns(bool)
func (_WrappedPocket *WrappedPocketSession) FeeFlag() (bool, error) {
	return _WrappedPocket.Contract.FeeFlag(&_WrappedPocket.CallOpts)
}

// FeeFlag is a free data retrieval call binding the contract method 0x118c6897.
//
// Solidity: function feeFlag() view returns(bool)
func (_WrappedPocket *WrappedPocketCallerSession) FeeFlag() (bool, error) {
	return _WrappedPocket.Contract.FeeFlag(&_WrappedPocket.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedPocket *WrappedPocketCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedPocket *WrappedPocketSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _WrappedPocket.Contract.GetRoleAdmin(&_WrappedPocket.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedPocket *WrappedPocketCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _WrappedPocket.Contract.GetRoleAdmin(&_WrappedPocket.CallOpts, role)
}

// GetUserNonce is a free data retrieval call binding the contract method 0x6834e3a8.
//
// Solidity: function getUserNonce(address user) view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) GetUserNonce(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "getUserNonce", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserNonce is a free data retrieval call binding the contract method 0x6834e3a8.
//
// Solidity: function getUserNonce(address user) view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) GetUserNonce(user common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.GetUserNonce(&_WrappedPocket.CallOpts, user)
}

// GetUserNonce is a free data retrieval call binding the contract method 0x6834e3a8.
//
// Solidity: function getUserNonce(address user) view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) GetUserNonce(user common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.GetUserNonce(&_WrappedPocket.CallOpts, user)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedPocket *WrappedPocketCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedPocket *WrappedPocketSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _WrappedPocket.Contract.HasRole(&_WrappedPocket.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedPocket *WrappedPocketCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _WrappedPocket.Contract.HasRole(&_WrappedPocket.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedPocket *WrappedPocketCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedPocket *WrappedPocketSession) Name() (string, error) {
	return _WrappedPocket.Contract.Name(&_WrappedPocket.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedPocket *WrappedPocketCallerSession) Name() (string, error) {
	return _WrappedPocket.Contract.Name(&_WrappedPocket.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) Nonces(owner common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.Nonces(&_WrappedPocket.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _WrappedPocket.Contract.Nonces(&_WrappedPocket.CallOpts, owner)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_WrappedPocket *WrappedPocketCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_WrappedPocket *WrappedPocketSession) Paused() (bool, error) {
	return _WrappedPocket.Contract.Paused(&_WrappedPocket.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_WrappedPocket *WrappedPocketCallerSession) Paused() (bool, error) {
	return _WrappedPocket.Contract.Paused(&_WrappedPocket.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedPocket *WrappedPocketCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedPocket *WrappedPocketSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WrappedPocket.Contract.SupportsInterface(&_WrappedPocket.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedPocket *WrappedPocketCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WrappedPocket.Contract.SupportsInterface(&_WrappedPocket.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedPocket *WrappedPocketCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedPocket *WrappedPocketSession) Symbol() (string, error) {
	return _WrappedPocket.Contract.Symbol(&_WrappedPocket.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedPocket *WrappedPocketCallerSession) Symbol() (string, error) {
	return _WrappedPocket.Contract.Symbol(&_WrappedPocket.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedPocket *WrappedPocketCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedPocket.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedPocket *WrappedPocketSession) TotalSupply() (*big.Int, error) {
	return _WrappedPocket.Contract.TotalSupply(&_WrappedPocket.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedPocket *WrappedPocketCallerSession) TotalSupply() (*big.Int, error) {
	return _WrappedPocket.Contract.TotalSupply(&_WrappedPocket.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Approve(&_WrappedPocket.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Approve(&_WrappedPocket.TransactOpts, spender, amount)
}

// BatchMint is a paid mutator transaction binding the contract method 0xd559f05b.
//
// Solidity: function batchMint(address[] to, uint256[] amount, uint256[] nonce) returns()
func (_WrappedPocket *WrappedPocketTransactor) BatchMint(opts *bind.TransactOpts, to []common.Address, amount []*big.Int, nonce []*big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "batchMint", to, amount, nonce)
}

// BatchMint is a paid mutator transaction binding the contract method 0xd559f05b.
//
// Solidity: function batchMint(address[] to, uint256[] amount, uint256[] nonce) returns()
func (_WrappedPocket *WrappedPocketSession) BatchMint(to []common.Address, amount []*big.Int, nonce []*big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.BatchMint(&_WrappedPocket.TransactOpts, to, amount, nonce)
}

// BatchMint is a paid mutator transaction binding the contract method 0xd559f05b.
//
// Solidity: function batchMint(address[] to, uint256[] amount, uint256[] nonce) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) BatchMint(to []common.Address, amount []*big.Int, nonce []*big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.BatchMint(&_WrappedPocket.TransactOpts, to, amount, nonce)
}

// BurnAndBridge is a paid mutator transaction binding the contract method 0x8402eb6b.
//
// Solidity: function burnAndBridge(uint256 amount, address poktAddress) returns()
func (_WrappedPocket *WrappedPocketTransactor) BurnAndBridge(opts *bind.TransactOpts, amount *big.Int, poktAddress common.Address) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "burnAndBridge", amount, poktAddress)
}

// BurnAndBridge is a paid mutator transaction binding the contract method 0x8402eb6b.
//
// Solidity: function burnAndBridge(uint256 amount, address poktAddress) returns()
func (_WrappedPocket *WrappedPocketSession) BurnAndBridge(amount *big.Int, poktAddress common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.BurnAndBridge(&_WrappedPocket.TransactOpts, amount, poktAddress)
}

// BurnAndBridge is a paid mutator transaction binding the contract method 0x8402eb6b.
//
// Solidity: function burnAndBridge(uint256 amount, address poktAddress) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) BurnAndBridge(amount *big.Int, poktAddress common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.BurnAndBridge(&_WrappedPocket.TransactOpts, amount, poktAddress)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_WrappedPocket *WrappedPocketTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_WrappedPocket *WrappedPocketSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.BurnFrom(&_WrappedPocket.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.BurnFrom(&_WrappedPocket.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_WrappedPocket *WrappedPocketTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_WrappedPocket *WrappedPocketSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.DecreaseAllowance(&_WrappedPocket.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_WrappedPocket *WrappedPocketTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.DecreaseAllowance(&_WrappedPocket.TransactOpts, spender, subtractedValue)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.GrantRole(&_WrappedPocket.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.GrantRole(&_WrappedPocket.TransactOpts, role, account)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_WrappedPocket *WrappedPocketTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_WrappedPocket *WrappedPocketSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.IncreaseAllowance(&_WrappedPocket.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_WrappedPocket *WrappedPocketTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.IncreaseAllowance(&_WrappedPocket.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 amount, uint256 nonce) returns()
func (_WrappedPocket *WrappedPocketTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int, nonce *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "mint", to, amount, nonce)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 amount, uint256 nonce) returns()
func (_WrappedPocket *WrappedPocketSession) Mint(to common.Address, amount *big.Int, nonce *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Mint(&_WrappedPocket.TransactOpts, to, amount, nonce)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 amount, uint256 nonce) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) Mint(to common.Address, amount *big.Int, nonce *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Mint(&_WrappedPocket.TransactOpts, to, amount, nonce)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_WrappedPocket *WrappedPocketTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_WrappedPocket *WrappedPocketSession) Pause() (*types.Transaction, error) {
	return _WrappedPocket.Contract.Pause(&_WrappedPocket.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_WrappedPocket *WrappedPocketTransactorSession) Pause() (*types.Transaction, error) {
	return _WrappedPocket.Contract.Pause(&_WrappedPocket.TransactOpts)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_WrappedPocket *WrappedPocketTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_WrappedPocket *WrappedPocketSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Permit(&_WrappedPocket.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Permit(&_WrappedPocket.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.RenounceRole(&_WrappedPocket.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.RenounceRole(&_WrappedPocket.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.RevokeRole(&_WrappedPocket.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.RevokeRole(&_WrappedPocket.TransactOpts, role, account)
}

// SetFee is a paid mutator transaction binding the contract method 0xe11ee7c8.
//
// Solidity: function setFee(bool flag, uint256 newFee, address newCollector) returns()
func (_WrappedPocket *WrappedPocketTransactor) SetFee(opts *bind.TransactOpts, flag bool, newFee *big.Int, newCollector common.Address) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "setFee", flag, newFee, newCollector)
}

// SetFee is a paid mutator transaction binding the contract method 0xe11ee7c8.
//
// Solidity: function setFee(bool flag, uint256 newFee, address newCollector) returns()
func (_WrappedPocket *WrappedPocketSession) SetFee(flag bool, newFee *big.Int, newCollector common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.SetFee(&_WrappedPocket.TransactOpts, flag, newFee, newCollector)
}

// SetFee is a paid mutator transaction binding the contract method 0xe11ee7c8.
//
// Solidity: function setFee(bool flag, uint256 newFee, address newCollector) returns()
func (_WrappedPocket *WrappedPocketTransactorSession) SetFee(flag bool, newFee *big.Int, newCollector common.Address) (*types.Transaction, error) {
	return _WrappedPocket.Contract.SetFee(&_WrappedPocket.TransactOpts, flag, newFee, newCollector)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Transfer(&_WrappedPocket.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.Transfer(&_WrappedPocket.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.TransferFrom(&_WrappedPocket.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_WrappedPocket *WrappedPocketTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedPocket.Contract.TransferFrom(&_WrappedPocket.TransactOpts, from, to, amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_WrappedPocket *WrappedPocketTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedPocket.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_WrappedPocket *WrappedPocketSession) Unpause() (*types.Transaction, error) {
	return _WrappedPocket.Contract.Unpause(&_WrappedPocket.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_WrappedPocket *WrappedPocketTransactorSession) Unpause() (*types.Transaction, error) {
	return _WrappedPocket.Contract.Unpause(&_WrappedPocket.TransactOpts)
}

// WrappedPocketApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WrappedPocket contract.
type WrappedPocketApprovalIterator struct {
	Event *WrappedPocketApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketApproval represents a Approval event raised by the WrappedPocket contract.
type WrappedPocketApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedPocket *WrappedPocketFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WrappedPocketApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketApprovalIterator{contract: _WrappedPocket.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedPocket *WrappedPocketFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WrappedPocketApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketApproval)
				if err := _WrappedPocket.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedPocket *WrappedPocketFilterer) ParseApproval(log types.Log) (*WrappedPocketApproval, error) {
	event := new(WrappedPocketApproval)
	if err := _WrappedPocket.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketBurnAndBridgeIterator is returned from FilterBurnAndBridge and is used to iterate over the raw logs and unpacked data for BurnAndBridge events raised by the WrappedPocket contract.
type WrappedPocketBurnAndBridgeIterator struct {
	Event *WrappedPocketBurnAndBridge // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketBurnAndBridgeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketBurnAndBridge)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketBurnAndBridge)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketBurnAndBridgeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketBurnAndBridgeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketBurnAndBridge represents a BurnAndBridge event raised by the WrappedPocket contract.
type WrappedPocketBurnAndBridge struct {
	Amount      *big.Int
	PoktAddress common.Address
	From        common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBurnAndBridge is a free log retrieval operation binding the contract event 0xac314bfa2d21af3d3e0937c97bc805574e3b1abb58a457ef02280c6f5e3faa75.
//
// Solidity: event BurnAndBridge(uint256 indexed amount, address indexed poktAddress, address indexed from)
func (_WrappedPocket *WrappedPocketFilterer) FilterBurnAndBridge(opts *bind.FilterOpts, amount []*big.Int, poktAddress []common.Address, from []common.Address) (*WrappedPocketBurnAndBridgeIterator, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var poktAddressRule []interface{}
	for _, poktAddressItem := range poktAddress {
		poktAddressRule = append(poktAddressRule, poktAddressItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "BurnAndBridge", amountRule, poktAddressRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketBurnAndBridgeIterator{contract: _WrappedPocket.contract, event: "BurnAndBridge", logs: logs, sub: sub}, nil
}

// WatchBurnAndBridge is a free log subscription operation binding the contract event 0xac314bfa2d21af3d3e0937c97bc805574e3b1abb58a457ef02280c6f5e3faa75.
//
// Solidity: event BurnAndBridge(uint256 indexed amount, address indexed poktAddress, address indexed from)
func (_WrappedPocket *WrappedPocketFilterer) WatchBurnAndBridge(opts *bind.WatchOpts, sink chan<- *WrappedPocketBurnAndBridge, amount []*big.Int, poktAddress []common.Address, from []common.Address) (event.Subscription, error) {

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var poktAddressRule []interface{}
	for _, poktAddressItem := range poktAddress {
		poktAddressRule = append(poktAddressRule, poktAddressItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "BurnAndBridge", amountRule, poktAddressRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketBurnAndBridge)
				if err := _WrappedPocket.contract.UnpackLog(event, "BurnAndBridge", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBurnAndBridge is a log parse operation binding the contract event 0xac314bfa2d21af3d3e0937c97bc805574e3b1abb58a457ef02280c6f5e3faa75.
//
// Solidity: event BurnAndBridge(uint256 indexed amount, address indexed poktAddress, address indexed from)
func (_WrappedPocket *WrappedPocketFilterer) ParseBurnAndBridge(log types.Log) (*WrappedPocketBurnAndBridge, error) {
	event := new(WrappedPocketBurnAndBridge)
	if err := _WrappedPocket.contract.UnpackLog(event, "BurnAndBridge", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the WrappedPocket contract.
type WrappedPocketEIP712DomainChangedIterator struct {
	Event *WrappedPocketEIP712DomainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketEIP712DomainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketEIP712DomainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketEIP712DomainChanged represents a EIP712DomainChanged event raised by the WrappedPocket contract.
type WrappedPocketEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_WrappedPocket *WrappedPocketFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*WrappedPocketEIP712DomainChangedIterator, error) {

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &WrappedPocketEIP712DomainChangedIterator{contract: _WrappedPocket.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_WrappedPocket *WrappedPocketFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *WrappedPocketEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketEIP712DomainChanged)
				if err := _WrappedPocket.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_WrappedPocket *WrappedPocketFilterer) ParseEIP712DomainChanged(log types.Log) (*WrappedPocketEIP712DomainChanged, error) {
	event := new(WrappedPocketEIP712DomainChanged)
	if err := _WrappedPocket.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketFeeCollectedIterator is returned from FilterFeeCollected and is used to iterate over the raw logs and unpacked data for FeeCollected events raised by the WrappedPocket contract.
type WrappedPocketFeeCollectedIterator struct {
	Event *WrappedPocketFeeCollected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketFeeCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketFeeCollected)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketFeeCollected)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketFeeCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketFeeCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketFeeCollected represents a FeeCollected event raised by the WrappedPocket contract.
type WrappedPocketFeeCollected struct {
	FeeCollector common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFeeCollected is a free log retrieval operation binding the contract event 0x06c5efeff5c320943d265dc4e5f1af95ad523555ce0c1957e367dda5514572df.
//
// Solidity: event FeeCollected(address indexed feeCollector, uint256 indexed amount)
func (_WrappedPocket *WrappedPocketFilterer) FilterFeeCollected(opts *bind.FilterOpts, feeCollector []common.Address, amount []*big.Int) (*WrappedPocketFeeCollectedIterator, error) {

	var feeCollectorRule []interface{}
	for _, feeCollectorItem := range feeCollector {
		feeCollectorRule = append(feeCollectorRule, feeCollectorItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "FeeCollected", feeCollectorRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketFeeCollectedIterator{contract: _WrappedPocket.contract, event: "FeeCollected", logs: logs, sub: sub}, nil
}

// WatchFeeCollected is a free log subscription operation binding the contract event 0x06c5efeff5c320943d265dc4e5f1af95ad523555ce0c1957e367dda5514572df.
//
// Solidity: event FeeCollected(address indexed feeCollector, uint256 indexed amount)
func (_WrappedPocket *WrappedPocketFilterer) WatchFeeCollected(opts *bind.WatchOpts, sink chan<- *WrappedPocketFeeCollected, feeCollector []common.Address, amount []*big.Int) (event.Subscription, error) {

	var feeCollectorRule []interface{}
	for _, feeCollectorItem := range feeCollector {
		feeCollectorRule = append(feeCollectorRule, feeCollectorItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "FeeCollected", feeCollectorRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketFeeCollected)
				if err := _WrappedPocket.contract.UnpackLog(event, "FeeCollected", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeCollected is a log parse operation binding the contract event 0x06c5efeff5c320943d265dc4e5f1af95ad523555ce0c1957e367dda5514572df.
//
// Solidity: event FeeCollected(address indexed feeCollector, uint256 indexed amount)
func (_WrappedPocket *WrappedPocketFilterer) ParseFeeCollected(log types.Log) (*WrappedPocketFeeCollected, error) {
	event := new(WrappedPocketFeeCollected)
	if err := _WrappedPocket.contract.UnpackLog(event, "FeeCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketFeeSetIterator is returned from FilterFeeSet and is used to iterate over the raw logs and unpacked data for FeeSet events raised by the WrappedPocket contract.
type WrappedPocketFeeSetIterator struct {
	Event *WrappedPocketFeeSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketFeeSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketFeeSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketFeeSet represents a FeeSet event raised by the WrappedPocket contract.
type WrappedPocketFeeSet struct {
	Flag         bool
	NewFeeBasis  *big.Int
	FeeCollector common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFeeSet is a free log retrieval operation binding the contract event 0xd67a9beee624d9012cc270d0c9f894e26cc7f615c2762f9081862e653ea4c7b8.
//
// Solidity: event FeeSet(bool indexed flag, uint256 indexed newFeeBasis, address indexed feeCollector)
func (_WrappedPocket *WrappedPocketFilterer) FilterFeeSet(opts *bind.FilterOpts, flag []bool, newFeeBasis []*big.Int, feeCollector []common.Address) (*WrappedPocketFeeSetIterator, error) {

	var flagRule []interface{}
	for _, flagItem := range flag {
		flagRule = append(flagRule, flagItem)
	}
	var newFeeBasisRule []interface{}
	for _, newFeeBasisItem := range newFeeBasis {
		newFeeBasisRule = append(newFeeBasisRule, newFeeBasisItem)
	}
	var feeCollectorRule []interface{}
	for _, feeCollectorItem := range feeCollector {
		feeCollectorRule = append(feeCollectorRule, feeCollectorItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "FeeSet", flagRule, newFeeBasisRule, feeCollectorRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketFeeSetIterator{contract: _WrappedPocket.contract, event: "FeeSet", logs: logs, sub: sub}, nil
}

// WatchFeeSet is a free log subscription operation binding the contract event 0xd67a9beee624d9012cc270d0c9f894e26cc7f615c2762f9081862e653ea4c7b8.
//
// Solidity: event FeeSet(bool indexed flag, uint256 indexed newFeeBasis, address indexed feeCollector)
func (_WrappedPocket *WrappedPocketFilterer) WatchFeeSet(opts *bind.WatchOpts, sink chan<- *WrappedPocketFeeSet, flag []bool, newFeeBasis []*big.Int, feeCollector []common.Address) (event.Subscription, error) {

	var flagRule []interface{}
	for _, flagItem := range flag {
		flagRule = append(flagRule, flagItem)
	}
	var newFeeBasisRule []interface{}
	for _, newFeeBasisItem := range newFeeBasis {
		newFeeBasisRule = append(newFeeBasisRule, newFeeBasisItem)
	}
	var feeCollectorRule []interface{}
	for _, feeCollectorItem := range feeCollector {
		feeCollectorRule = append(feeCollectorRule, feeCollectorItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "FeeSet", flagRule, newFeeBasisRule, feeCollectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketFeeSet)
				if err := _WrappedPocket.contract.UnpackLog(event, "FeeSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeSet is a log parse operation binding the contract event 0xd67a9beee624d9012cc270d0c9f894e26cc7f615c2762f9081862e653ea4c7b8.
//
// Solidity: event FeeSet(bool indexed flag, uint256 indexed newFeeBasis, address indexed feeCollector)
func (_WrappedPocket *WrappedPocketFilterer) ParseFeeSet(log types.Log) (*WrappedPocketFeeSet, error) {
	event := new(WrappedPocketFeeSet)
	if err := _WrappedPocket.contract.UnpackLog(event, "FeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the WrappedPocket contract.
type WrappedPocketPausedIterator struct {
	Event *WrappedPocketPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketPaused represents a Paused event raised by the WrappedPocket contract.
type WrappedPocketPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_WrappedPocket *WrappedPocketFilterer) FilterPaused(opts *bind.FilterOpts) (*WrappedPocketPausedIterator, error) {

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &WrappedPocketPausedIterator{contract: _WrappedPocket.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_WrappedPocket *WrappedPocketFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *WrappedPocketPaused) (event.Subscription, error) {

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketPaused)
				if err := _WrappedPocket.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_WrappedPocket *WrappedPocketFilterer) ParsePaused(log types.Log) (*WrappedPocketPaused, error) {
	event := new(WrappedPocketPaused)
	if err := _WrappedPocket.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the WrappedPocket contract.
type WrappedPocketRoleAdminChangedIterator struct {
	Event *WrappedPocketRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketRoleAdminChanged represents a RoleAdminChanged event raised by the WrappedPocket contract.
type WrappedPocketRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedPocket *WrappedPocketFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*WrappedPocketRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketRoleAdminChangedIterator{contract: _WrappedPocket.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedPocket *WrappedPocketFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *WrappedPocketRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketRoleAdminChanged)
				if err := _WrappedPocket.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedPocket *WrappedPocketFilterer) ParseRoleAdminChanged(log types.Log) (*WrappedPocketRoleAdminChanged, error) {
	event := new(WrappedPocketRoleAdminChanged)
	if err := _WrappedPocket.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the WrappedPocket contract.
type WrappedPocketRoleGrantedIterator struct {
	Event *WrappedPocketRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketRoleGranted represents a RoleGranted event raised by the WrappedPocket contract.
type WrappedPocketRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedPocket *WrappedPocketFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*WrappedPocketRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketRoleGrantedIterator{contract: _WrappedPocket.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedPocket *WrappedPocketFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *WrappedPocketRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketRoleGranted)
				if err := _WrappedPocket.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedPocket *WrappedPocketFilterer) ParseRoleGranted(log types.Log) (*WrappedPocketRoleGranted, error) {
	event := new(WrappedPocketRoleGranted)
	if err := _WrappedPocket.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the WrappedPocket contract.
type WrappedPocketRoleRevokedIterator struct {
	Event *WrappedPocketRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketRoleRevoked represents a RoleRevoked event raised by the WrappedPocket contract.
type WrappedPocketRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedPocket *WrappedPocketFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*WrappedPocketRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketRoleRevokedIterator{contract: _WrappedPocket.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedPocket *WrappedPocketFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *WrappedPocketRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketRoleRevoked)
				if err := _WrappedPocket.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedPocket *WrappedPocketFilterer) ParseRoleRevoked(log types.Log) (*WrappedPocketRoleRevoked, error) {
	event := new(WrappedPocketRoleRevoked)
	if err := _WrappedPocket.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WrappedPocket contract.
type WrappedPocketTransferIterator struct {
	Event *WrappedPocketTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketTransfer represents a Transfer event raised by the WrappedPocket contract.
type WrappedPocketTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedPocket *WrappedPocketFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedPocketTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WrappedPocketTransferIterator{contract: _WrappedPocket.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedPocket *WrappedPocketFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WrappedPocketTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketTransfer)
				if err := _WrappedPocket.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedPocket *WrappedPocketFilterer) ParseTransfer(log types.Log) (*WrappedPocketTransfer, error) {
	event := new(WrappedPocketTransfer)
	if err := _WrappedPocket.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedPocketUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the WrappedPocket contract.
type WrappedPocketUnpausedIterator struct {
	Event *WrappedPocketUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedPocketUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedPocketUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedPocketUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedPocketUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedPocketUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedPocketUnpaused represents a Unpaused event raised by the WrappedPocket contract.
type WrappedPocketUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_WrappedPocket *WrappedPocketFilterer) FilterUnpaused(opts *bind.FilterOpts) (*WrappedPocketUnpausedIterator, error) {

	logs, sub, err := _WrappedPocket.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &WrappedPocketUnpausedIterator{contract: _WrappedPocket.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_WrappedPocket *WrappedPocketFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *WrappedPocketUnpaused) (event.Subscription, error) {

	logs, sub, err := _WrappedPocket.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedPocketUnpaused)
				if err := _WrappedPocket.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_WrappedPocket *WrappedPocketFilterer) ParseUnpaused(log types.Log) (*WrappedPocketUnpaused, error) {
	event := new(WrappedPocketUnpaused)
	if err := _WrappedPocket.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
