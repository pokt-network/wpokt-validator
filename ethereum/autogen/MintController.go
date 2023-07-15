// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package autogen

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

// MintControllerMintData is an auto generated low-level Go binding around an user-defined struct.
type MintControllerMintData struct {
	Recipient common.Address
	Amount    *big.Int
	Nonce     *big.Int
}

// MintControllerMetaData contains all meta data concerning the MintController contract.
var MintControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wPokt\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BelowMinThreshold\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignatureRatio\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonCopper\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OverMintLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"lastMint\",\"type\":\"uint256\"}],\"name\":\"CurrentMintLimit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newCooldown\",\"type\":\"uint256\"}],\"name\":\"MintCooldownSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"NewValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"RemovedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"ratio\",\"type\":\"uint256\"}],\"name\":\"SignerThresholdSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentMintLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMintLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxMintLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mintPerSecond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structMintController.MintData\",\"name\":\"data\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"mintWrappedPocket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newMintPerSecond\",\"type\":\"uint256\"}],\"name\":\"setMintCooldown\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"signatureRatio\",\"type\":\"uint256\"}],\"name\":\"setSignerThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"signerThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wPokt\",\"outputs\":[{\"internalType\":\"contractIWPokt\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// MintControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use MintControllerMetaData.ABI instead.
var MintControllerABI = MintControllerMetaData.ABI

// MintController is an auto generated Go binding around an Ethereum contract.
type MintController struct {
	MintControllerCaller     // Read-only binding to the contract
	MintControllerTransactor // Write-only binding to the contract
	MintControllerFilterer   // Log filterer for contract events
}

// MintControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type MintControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MintControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MintControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MintControllerSession struct {
	Contract     *MintController   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MintControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MintControllerCallerSession struct {
	Contract *MintControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MintControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MintControllerTransactorSession struct {
	Contract     *MintControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MintControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type MintControllerRaw struct {
	Contract *MintController // Generic contract binding to access the raw methods on
}

// MintControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MintControllerCallerRaw struct {
	Contract *MintControllerCaller // Generic read-only contract binding to access the raw methods on
}

// MintControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MintControllerTransactorRaw struct {
	Contract *MintControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMintController creates a new instance of MintController, bound to a specific deployed contract.
func NewMintController(address common.Address, backend bind.ContractBackend) (*MintController, error) {
	contract, err := bindMintController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MintController{MintControllerCaller: MintControllerCaller{contract: contract}, MintControllerTransactor: MintControllerTransactor{contract: contract}, MintControllerFilterer: MintControllerFilterer{contract: contract}}, nil
}

// NewMintControllerCaller creates a new read-only instance of MintController, bound to a specific deployed contract.
func NewMintControllerCaller(address common.Address, caller bind.ContractCaller) (*MintControllerCaller, error) {
	contract, err := bindMintController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MintControllerCaller{contract: contract}, nil
}

// NewMintControllerTransactor creates a new write-only instance of MintController, bound to a specific deployed contract.
func NewMintControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*MintControllerTransactor, error) {
	contract, err := bindMintController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MintControllerTransactor{contract: contract}, nil
}

// NewMintControllerFilterer creates a new log filterer instance of MintController, bound to a specific deployed contract.
func NewMintControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*MintControllerFilterer, error) {
	contract, err := bindMintController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MintControllerFilterer{contract: contract}, nil
}

// bindMintController binds a generic wrapper to an already deployed contract.
func bindMintController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MintControllerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintController *MintControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MintController.Contract.MintControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintController *MintControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintController.Contract.MintControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintController *MintControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintController.Contract.MintControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintController *MintControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MintController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintController *MintControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintController *MintControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintController.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MintController *MintControllerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MintController *MintControllerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MintController.Contract.DEFAULTADMINROLE(&_MintController.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MintController *MintControllerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MintController.Contract.DEFAULTADMINROLE(&_MintController.CallOpts)
}

// CurrentMintLimit is a free data retrieval call binding the contract method 0xa40c3eab.
//
// Solidity: function currentMintLimit() view returns(uint256)
func (_MintController *MintControllerCaller) CurrentMintLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "currentMintLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentMintLimit is a free data retrieval call binding the contract method 0xa40c3eab.
//
// Solidity: function currentMintLimit() view returns(uint256)
func (_MintController *MintControllerSession) CurrentMintLimit() (*big.Int, error) {
	return _MintController.Contract.CurrentMintLimit(&_MintController.CallOpts)
}

// CurrentMintLimit is a free data retrieval call binding the contract method 0xa40c3eab.
//
// Solidity: function currentMintLimit() view returns(uint256)
func (_MintController *MintControllerCallerSession) CurrentMintLimit() (*big.Int, error) {
	return _MintController.Contract.CurrentMintLimit(&_MintController.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_MintController *MintControllerCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "eip712Domain")

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
func (_MintController *MintControllerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _MintController.Contract.Eip712Domain(&_MintController.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_MintController *MintControllerCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _MintController.Contract.Eip712Domain(&_MintController.CallOpts)
}

// LastMint is a free data retrieval call binding the contract method 0x586fc5b5.
//
// Solidity: function lastMint() view returns(uint256)
func (_MintController *MintControllerCaller) LastMint(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "lastMint")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastMint is a free data retrieval call binding the contract method 0x586fc5b5.
//
// Solidity: function lastMint() view returns(uint256)
func (_MintController *MintControllerSession) LastMint() (*big.Int, error) {
	return _MintController.Contract.LastMint(&_MintController.CallOpts)
}

// LastMint is a free data retrieval call binding the contract method 0x586fc5b5.
//
// Solidity: function lastMint() view returns(uint256)
func (_MintController *MintControllerCallerSession) LastMint() (*big.Int, error) {
	return _MintController.Contract.LastMint(&_MintController.CallOpts)
}

// LastMintLimit is a free data retrieval call binding the contract method 0x2f99582a.
//
// Solidity: function lastMintLimit() view returns(uint256)
func (_MintController *MintControllerCaller) LastMintLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "lastMintLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastMintLimit is a free data retrieval call binding the contract method 0x2f99582a.
//
// Solidity: function lastMintLimit() view returns(uint256)
func (_MintController *MintControllerSession) LastMintLimit() (*big.Int, error) {
	return _MintController.Contract.LastMintLimit(&_MintController.CallOpts)
}

// LastMintLimit is a free data retrieval call binding the contract method 0x2f99582a.
//
// Solidity: function lastMintLimit() view returns(uint256)
func (_MintController *MintControllerCallerSession) LastMintLimit() (*big.Int, error) {
	return _MintController.Contract.LastMintLimit(&_MintController.CallOpts)
}

// MaxMintLimit is a free data retrieval call binding the contract method 0x70e2f827.
//
// Solidity: function maxMintLimit() view returns(uint256)
func (_MintController *MintControllerCaller) MaxMintLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "maxMintLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxMintLimit is a free data retrieval call binding the contract method 0x70e2f827.
//
// Solidity: function maxMintLimit() view returns(uint256)
func (_MintController *MintControllerSession) MaxMintLimit() (*big.Int, error) {
	return _MintController.Contract.MaxMintLimit(&_MintController.CallOpts)
}

// MaxMintLimit is a free data retrieval call binding the contract method 0x70e2f827.
//
// Solidity: function maxMintLimit() view returns(uint256)
func (_MintController *MintControllerCallerSession) MaxMintLimit() (*big.Int, error) {
	return _MintController.Contract.MaxMintLimit(&_MintController.CallOpts)
}

// MintPerSecond is a free data retrieval call binding the contract method 0x272c444d.
//
// Solidity: function mintPerSecond() view returns(uint256)
func (_MintController *MintControllerCaller) MintPerSecond(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "mintPerSecond")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MintPerSecond is a free data retrieval call binding the contract method 0x272c444d.
//
// Solidity: function mintPerSecond() view returns(uint256)
func (_MintController *MintControllerSession) MintPerSecond() (*big.Int, error) {
	return _MintController.Contract.MintPerSecond(&_MintController.CallOpts)
}

// MintPerSecond is a free data retrieval call binding the contract method 0x272c444d.
//
// Solidity: function mintPerSecond() view returns(uint256)
func (_MintController *MintControllerCallerSession) MintPerSecond() (*big.Int, error) {
	return _MintController.Contract.MintPerSecond(&_MintController.CallOpts)
}

// SignerThreshold is a free data retrieval call binding the contract method 0xa4a4f390.
//
// Solidity: function signerThreshold() view returns(uint256)
func (_MintController *MintControllerCaller) SignerThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "signerThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SignerThreshold is a free data retrieval call binding the contract method 0xa4a4f390.
//
// Solidity: function signerThreshold() view returns(uint256)
func (_MintController *MintControllerSession) SignerThreshold() (*big.Int, error) {
	return _MintController.Contract.SignerThreshold(&_MintController.CallOpts)
}

// SignerThreshold is a free data retrieval call binding the contract method 0xa4a4f390.
//
// Solidity: function signerThreshold() view returns(uint256)
func (_MintController *MintControllerCallerSession) SignerThreshold() (*big.Int, error) {
	return _MintController.Contract.SignerThreshold(&_MintController.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_MintController *MintControllerCaller) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "validatorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_MintController *MintControllerSession) ValidatorCount() (*big.Int, error) {
	return _MintController.Contract.ValidatorCount(&_MintController.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_MintController *MintControllerCallerSession) ValidatorCount() (*big.Int, error) {
	return _MintController.Contract.ValidatorCount(&_MintController.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_MintController *MintControllerCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_MintController *MintControllerSession) Validators(arg0 common.Address) (bool, error) {
	return _MintController.Contract.Validators(&_MintController.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_MintController *MintControllerCallerSession) Validators(arg0 common.Address) (bool, error) {
	return _MintController.Contract.Validators(&_MintController.CallOpts, arg0)
}

// WPokt is a free data retrieval call binding the contract method 0xd72a828b.
//
// Solidity: function wPokt() view returns(address)
func (_MintController *MintControllerCaller) WPokt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MintController.contract.Call(opts, &out, "wPokt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WPokt is a free data retrieval call binding the contract method 0xd72a828b.
//
// Solidity: function wPokt() view returns(address)
func (_MintController *MintControllerSession) WPokt() (common.Address, error) {
	return _MintController.Contract.WPokt(&_MintController.CallOpts)
}

// WPokt is a free data retrieval call binding the contract method 0xd72a828b.
//
// Solidity: function wPokt() view returns(address)
func (_MintController *MintControllerCallerSession) WPokt() (common.Address, error) {
	return _MintController.Contract.WPokt(&_MintController.CallOpts)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(address validator) returns()
func (_MintController *MintControllerTransactor) AddValidator(opts *bind.TransactOpts, validator common.Address) (*types.Transaction, error) {
	return _MintController.contract.Transact(opts, "addValidator", validator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(address validator) returns()
func (_MintController *MintControllerSession) AddValidator(validator common.Address) (*types.Transaction, error) {
	return _MintController.Contract.AddValidator(&_MintController.TransactOpts, validator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(address validator) returns()
func (_MintController *MintControllerTransactorSession) AddValidator(validator common.Address) (*types.Transaction, error) {
	return _MintController.Contract.AddValidator(&_MintController.TransactOpts, validator)
}

// MintWrappedPocket is a paid mutator transaction binding the contract method 0x0f22ed52.
//
// Solidity: function mintWrappedPocket((address,uint256,uint256) data, bytes[] signatures) returns()
func (_MintController *MintControllerTransactor) MintWrappedPocket(opts *bind.TransactOpts, data MintControllerMintData, signatures [][]byte) (*types.Transaction, error) {
	return _MintController.contract.Transact(opts, "mintWrappedPocket", data, signatures)
}

// MintWrappedPocket is a paid mutator transaction binding the contract method 0x0f22ed52.
//
// Solidity: function mintWrappedPocket((address,uint256,uint256) data, bytes[] signatures) returns()
func (_MintController *MintControllerSession) MintWrappedPocket(data MintControllerMintData, signatures [][]byte) (*types.Transaction, error) {
	return _MintController.Contract.MintWrappedPocket(&_MintController.TransactOpts, data, signatures)
}

// MintWrappedPocket is a paid mutator transaction binding the contract method 0x0f22ed52.
//
// Solidity: function mintWrappedPocket((address,uint256,uint256) data, bytes[] signatures) returns()
func (_MintController *MintControllerTransactorSession) MintWrappedPocket(data MintControllerMintData, signatures [][]byte) (*types.Transaction, error) {
	return _MintController.Contract.MintWrappedPocket(&_MintController.TransactOpts, data, signatures)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address validator) returns()
func (_MintController *MintControllerTransactor) RemoveValidator(opts *bind.TransactOpts, validator common.Address) (*types.Transaction, error) {
	return _MintController.contract.Transact(opts, "removeValidator", validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address validator) returns()
func (_MintController *MintControllerSession) RemoveValidator(validator common.Address) (*types.Transaction, error) {
	return _MintController.Contract.RemoveValidator(&_MintController.TransactOpts, validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address validator) returns()
func (_MintController *MintControllerTransactorSession) RemoveValidator(validator common.Address) (*types.Transaction, error) {
	return _MintController.Contract.RemoveValidator(&_MintController.TransactOpts, validator)
}

// SetMintCooldown is a paid mutator transaction binding the contract method 0x59aa6859.
//
// Solidity: function setMintCooldown(uint256 newLimit, uint256 newMintPerSecond) returns()
func (_MintController *MintControllerTransactor) SetMintCooldown(opts *bind.TransactOpts, newLimit *big.Int, newMintPerSecond *big.Int) (*types.Transaction, error) {
	return _MintController.contract.Transact(opts, "setMintCooldown", newLimit, newMintPerSecond)
}

// SetMintCooldown is a paid mutator transaction binding the contract method 0x59aa6859.
//
// Solidity: function setMintCooldown(uint256 newLimit, uint256 newMintPerSecond) returns()
func (_MintController *MintControllerSession) SetMintCooldown(newLimit *big.Int, newMintPerSecond *big.Int) (*types.Transaction, error) {
	return _MintController.Contract.SetMintCooldown(&_MintController.TransactOpts, newLimit, newMintPerSecond)
}

// SetMintCooldown is a paid mutator transaction binding the contract method 0x59aa6859.
//
// Solidity: function setMintCooldown(uint256 newLimit, uint256 newMintPerSecond) returns()
func (_MintController *MintControllerTransactorSession) SetMintCooldown(newLimit *big.Int, newMintPerSecond *big.Int) (*types.Transaction, error) {
	return _MintController.Contract.SetMintCooldown(&_MintController.TransactOpts, newLimit, newMintPerSecond)
}

// SetSignerThreshold is a paid mutator transaction binding the contract method 0x251b8192.
//
// Solidity: function setSignerThreshold(uint256 signatureRatio) returns()
func (_MintController *MintControllerTransactor) SetSignerThreshold(opts *bind.TransactOpts, signatureRatio *big.Int) (*types.Transaction, error) {
	return _MintController.contract.Transact(opts, "setSignerThreshold", signatureRatio)
}

// SetSignerThreshold is a paid mutator transaction binding the contract method 0x251b8192.
//
// Solidity: function setSignerThreshold(uint256 signatureRatio) returns()
func (_MintController *MintControllerSession) SetSignerThreshold(signatureRatio *big.Int) (*types.Transaction, error) {
	return _MintController.Contract.SetSignerThreshold(&_MintController.TransactOpts, signatureRatio)
}

// SetSignerThreshold is a paid mutator transaction binding the contract method 0x251b8192.
//
// Solidity: function setSignerThreshold(uint256 signatureRatio) returns()
func (_MintController *MintControllerTransactorSession) SetSignerThreshold(signatureRatio *big.Int) (*types.Transaction, error) {
	return _MintController.Contract.SetSignerThreshold(&_MintController.TransactOpts, signatureRatio)
}

// MintControllerCurrentMintLimitIterator is returned from FilterCurrentMintLimit and is used to iterate over the raw logs and unpacked data for CurrentMintLimit events raised by the MintController contract.
type MintControllerCurrentMintLimitIterator struct {
	Event *MintControllerCurrentMintLimit // Event containing the contract specifics and raw log

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
func (it *MintControllerCurrentMintLimitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintControllerCurrentMintLimit)
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
		it.Event = new(MintControllerCurrentMintLimit)
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
func (it *MintControllerCurrentMintLimitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintControllerCurrentMintLimitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintControllerCurrentMintLimit represents a CurrentMintLimit event raised by the MintController contract.
type MintControllerCurrentMintLimit struct {
	Limit    *big.Int
	LastMint *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCurrentMintLimit is a free log retrieval operation binding the contract event 0x37ce17c2014463ff310cd0d994a30795f5c48b6d6c3104ba83677f096383e3e4.
//
// Solidity: event CurrentMintLimit(uint256 indexed limit, uint256 indexed lastMint)
func (_MintController *MintControllerFilterer) FilterCurrentMintLimit(opts *bind.FilterOpts, limit []*big.Int, lastMint []*big.Int) (*MintControllerCurrentMintLimitIterator, error) {

	var limitRule []interface{}
	for _, limitItem := range limit {
		limitRule = append(limitRule, limitItem)
	}
	var lastMintRule []interface{}
	for _, lastMintItem := range lastMint {
		lastMintRule = append(lastMintRule, lastMintItem)
	}

	logs, sub, err := _MintController.contract.FilterLogs(opts, "CurrentMintLimit", limitRule, lastMintRule)
	if err != nil {
		return nil, err
	}
	return &MintControllerCurrentMintLimitIterator{contract: _MintController.contract, event: "CurrentMintLimit", logs: logs, sub: sub}, nil
}

// WatchCurrentMintLimit is a free log subscription operation binding the contract event 0x37ce17c2014463ff310cd0d994a30795f5c48b6d6c3104ba83677f096383e3e4.
//
// Solidity: event CurrentMintLimit(uint256 indexed limit, uint256 indexed lastMint)
func (_MintController *MintControllerFilterer) WatchCurrentMintLimit(opts *bind.WatchOpts, sink chan<- *MintControllerCurrentMintLimit, limit []*big.Int, lastMint []*big.Int) (event.Subscription, error) {

	var limitRule []interface{}
	for _, limitItem := range limit {
		limitRule = append(limitRule, limitItem)
	}
	var lastMintRule []interface{}
	for _, lastMintItem := range lastMint {
		lastMintRule = append(lastMintRule, lastMintItem)
	}

	logs, sub, err := _MintController.contract.WatchLogs(opts, "CurrentMintLimit", limitRule, lastMintRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintControllerCurrentMintLimit)
				if err := _MintController.contract.UnpackLog(event, "CurrentMintLimit", log); err != nil {
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

// ParseCurrentMintLimit is a log parse operation binding the contract event 0x37ce17c2014463ff310cd0d994a30795f5c48b6d6c3104ba83677f096383e3e4.
//
// Solidity: event CurrentMintLimit(uint256 indexed limit, uint256 indexed lastMint)
func (_MintController *MintControllerFilterer) ParseCurrentMintLimit(log types.Log) (*MintControllerCurrentMintLimit, error) {
	event := new(MintControllerCurrentMintLimit)
	if err := _MintController.contract.UnpackLog(event, "CurrentMintLimit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintControllerEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the MintController contract.
type MintControllerEIP712DomainChangedIterator struct {
	Event *MintControllerEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *MintControllerEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintControllerEIP712DomainChanged)
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
		it.Event = new(MintControllerEIP712DomainChanged)
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
func (it *MintControllerEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintControllerEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintControllerEIP712DomainChanged represents a EIP712DomainChanged event raised by the MintController contract.
type MintControllerEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_MintController *MintControllerFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*MintControllerEIP712DomainChangedIterator, error) {

	logs, sub, err := _MintController.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &MintControllerEIP712DomainChangedIterator{contract: _MintController.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_MintController *MintControllerFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *MintControllerEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _MintController.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintControllerEIP712DomainChanged)
				if err := _MintController.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_MintController *MintControllerFilterer) ParseEIP712DomainChanged(log types.Log) (*MintControllerEIP712DomainChanged, error) {
	event := new(MintControllerEIP712DomainChanged)
	if err := _MintController.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintControllerMintCooldownSetIterator is returned from FilterMintCooldownSet and is used to iterate over the raw logs and unpacked data for MintCooldownSet events raised by the MintController contract.
type MintControllerMintCooldownSetIterator struct {
	Event *MintControllerMintCooldownSet // Event containing the contract specifics and raw log

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
func (it *MintControllerMintCooldownSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintControllerMintCooldownSet)
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
		it.Event = new(MintControllerMintCooldownSet)
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
func (it *MintControllerMintCooldownSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintControllerMintCooldownSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintControllerMintCooldownSet represents a MintCooldownSet event raised by the MintController contract.
type MintControllerMintCooldownSet struct {
	NewLimit    *big.Int
	NewCooldown *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMintCooldownSet is a free log retrieval operation binding the contract event 0x2c1471be9f5586272d41bc8d7d0ce89b3b10e4693755e30752510cb88d42fdf5.
//
// Solidity: event MintCooldownSet(uint256 newLimit, uint256 newCooldown)
func (_MintController *MintControllerFilterer) FilterMintCooldownSet(opts *bind.FilterOpts) (*MintControllerMintCooldownSetIterator, error) {

	logs, sub, err := _MintController.contract.FilterLogs(opts, "MintCooldownSet")
	if err != nil {
		return nil, err
	}
	return &MintControllerMintCooldownSetIterator{contract: _MintController.contract, event: "MintCooldownSet", logs: logs, sub: sub}, nil
}

// WatchMintCooldownSet is a free log subscription operation binding the contract event 0x2c1471be9f5586272d41bc8d7d0ce89b3b10e4693755e30752510cb88d42fdf5.
//
// Solidity: event MintCooldownSet(uint256 newLimit, uint256 newCooldown)
func (_MintController *MintControllerFilterer) WatchMintCooldownSet(opts *bind.WatchOpts, sink chan<- *MintControllerMintCooldownSet) (event.Subscription, error) {

	logs, sub, err := _MintController.contract.WatchLogs(opts, "MintCooldownSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintControllerMintCooldownSet)
				if err := _MintController.contract.UnpackLog(event, "MintCooldownSet", log); err != nil {
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

// ParseMintCooldownSet is a log parse operation binding the contract event 0x2c1471be9f5586272d41bc8d7d0ce89b3b10e4693755e30752510cb88d42fdf5.
//
// Solidity: event MintCooldownSet(uint256 newLimit, uint256 newCooldown)
func (_MintController *MintControllerFilterer) ParseMintCooldownSet(log types.Log) (*MintControllerMintCooldownSet, error) {
	event := new(MintControllerMintCooldownSet)
	if err := _MintController.contract.UnpackLog(event, "MintCooldownSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintControllerNewValidatorIterator is returned from FilterNewValidator and is used to iterate over the raw logs and unpacked data for NewValidator events raised by the MintController contract.
type MintControllerNewValidatorIterator struct {
	Event *MintControllerNewValidator // Event containing the contract specifics and raw log

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
func (it *MintControllerNewValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintControllerNewValidator)
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
		it.Event = new(MintControllerNewValidator)
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
func (it *MintControllerNewValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintControllerNewValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintControllerNewValidator represents a NewValidator event raised by the MintController contract.
type MintControllerNewValidator struct {
	Validator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewValidator is a free log retrieval operation binding the contract event 0x29b4645f23b856eccf12b3b38e036c3221ca1b5a9afa2a83aea7ead34e47987c.
//
// Solidity: event NewValidator(address indexed validator)
func (_MintController *MintControllerFilterer) FilterNewValidator(opts *bind.FilterOpts, validator []common.Address) (*MintControllerNewValidatorIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _MintController.contract.FilterLogs(opts, "NewValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return &MintControllerNewValidatorIterator{contract: _MintController.contract, event: "NewValidator", logs: logs, sub: sub}, nil
}

// WatchNewValidator is a free log subscription operation binding the contract event 0x29b4645f23b856eccf12b3b38e036c3221ca1b5a9afa2a83aea7ead34e47987c.
//
// Solidity: event NewValidator(address indexed validator)
func (_MintController *MintControllerFilterer) WatchNewValidator(opts *bind.WatchOpts, sink chan<- *MintControllerNewValidator, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _MintController.contract.WatchLogs(opts, "NewValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintControllerNewValidator)
				if err := _MintController.contract.UnpackLog(event, "NewValidator", log); err != nil {
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

// ParseNewValidator is a log parse operation binding the contract event 0x29b4645f23b856eccf12b3b38e036c3221ca1b5a9afa2a83aea7ead34e47987c.
//
// Solidity: event NewValidator(address indexed validator)
func (_MintController *MintControllerFilterer) ParseNewValidator(log types.Log) (*MintControllerNewValidator, error) {
	event := new(MintControllerNewValidator)
	if err := _MintController.contract.UnpackLog(event, "NewValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintControllerRemovedValidatorIterator is returned from FilterRemovedValidator and is used to iterate over the raw logs and unpacked data for RemovedValidator events raised by the MintController contract.
type MintControllerRemovedValidatorIterator struct {
	Event *MintControllerRemovedValidator // Event containing the contract specifics and raw log

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
func (it *MintControllerRemovedValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintControllerRemovedValidator)
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
		it.Event = new(MintControllerRemovedValidator)
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
func (it *MintControllerRemovedValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintControllerRemovedValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintControllerRemovedValidator represents a RemovedValidator event raised by the MintController contract.
type MintControllerRemovedValidator struct {
	Validator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRemovedValidator is a free log retrieval operation binding the contract event 0xb625c55cf7e37b54fcd18bc4edafdf3f4f9acd59a5ec824c77c795dcb2d65070.
//
// Solidity: event RemovedValidator(address indexed validator)
func (_MintController *MintControllerFilterer) FilterRemovedValidator(opts *bind.FilterOpts, validator []common.Address) (*MintControllerRemovedValidatorIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _MintController.contract.FilterLogs(opts, "RemovedValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return &MintControllerRemovedValidatorIterator{contract: _MintController.contract, event: "RemovedValidator", logs: logs, sub: sub}, nil
}

// WatchRemovedValidator is a free log subscription operation binding the contract event 0xb625c55cf7e37b54fcd18bc4edafdf3f4f9acd59a5ec824c77c795dcb2d65070.
//
// Solidity: event RemovedValidator(address indexed validator)
func (_MintController *MintControllerFilterer) WatchRemovedValidator(opts *bind.WatchOpts, sink chan<- *MintControllerRemovedValidator, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _MintController.contract.WatchLogs(opts, "RemovedValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintControllerRemovedValidator)
				if err := _MintController.contract.UnpackLog(event, "RemovedValidator", log); err != nil {
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

// ParseRemovedValidator is a log parse operation binding the contract event 0xb625c55cf7e37b54fcd18bc4edafdf3f4f9acd59a5ec824c77c795dcb2d65070.
//
// Solidity: event RemovedValidator(address indexed validator)
func (_MintController *MintControllerFilterer) ParseRemovedValidator(log types.Log) (*MintControllerRemovedValidator, error) {
	event := new(MintControllerRemovedValidator)
	if err := _MintController.contract.UnpackLog(event, "RemovedValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintControllerSignerThresholdSetIterator is returned from FilterSignerThresholdSet and is used to iterate over the raw logs and unpacked data for SignerThresholdSet events raised by the MintController contract.
type MintControllerSignerThresholdSetIterator struct {
	Event *MintControllerSignerThresholdSet // Event containing the contract specifics and raw log

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
func (it *MintControllerSignerThresholdSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintControllerSignerThresholdSet)
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
		it.Event = new(MintControllerSignerThresholdSet)
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
func (it *MintControllerSignerThresholdSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintControllerSignerThresholdSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintControllerSignerThresholdSet represents a SignerThresholdSet event raised by the MintController contract.
type MintControllerSignerThresholdSet struct {
	Ratio *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSignerThresholdSet is a free log retrieval operation binding the contract event 0x030f8d880b0bf7b6311c00e028809570118cbccc8302e08c8776ad77d10f505e.
//
// Solidity: event SignerThresholdSet(uint256 indexed ratio)
func (_MintController *MintControllerFilterer) FilterSignerThresholdSet(opts *bind.FilterOpts, ratio []*big.Int) (*MintControllerSignerThresholdSetIterator, error) {

	var ratioRule []interface{}
	for _, ratioItem := range ratio {
		ratioRule = append(ratioRule, ratioItem)
	}

	logs, sub, err := _MintController.contract.FilterLogs(opts, "SignerThresholdSet", ratioRule)
	if err != nil {
		return nil, err
	}
	return &MintControllerSignerThresholdSetIterator{contract: _MintController.contract, event: "SignerThresholdSet", logs: logs, sub: sub}, nil
}

// WatchSignerThresholdSet is a free log subscription operation binding the contract event 0x030f8d880b0bf7b6311c00e028809570118cbccc8302e08c8776ad77d10f505e.
//
// Solidity: event SignerThresholdSet(uint256 indexed ratio)
func (_MintController *MintControllerFilterer) WatchSignerThresholdSet(opts *bind.WatchOpts, sink chan<- *MintControllerSignerThresholdSet, ratio []*big.Int) (event.Subscription, error) {

	var ratioRule []interface{}
	for _, ratioItem := range ratio {
		ratioRule = append(ratioRule, ratioItem)
	}

	logs, sub, err := _MintController.contract.WatchLogs(opts, "SignerThresholdSet", ratioRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintControllerSignerThresholdSet)
				if err := _MintController.contract.UnpackLog(event, "SignerThresholdSet", log); err != nil {
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

// ParseSignerThresholdSet is a log parse operation binding the contract event 0x030f8d880b0bf7b6311c00e028809570118cbccc8302e08c8776ad77d10f505e.
//
// Solidity: event SignerThresholdSet(uint256 indexed ratio)
func (_MintController *MintControllerFilterer) ParseSignerThresholdSet(log types.Log) (*MintControllerSignerThresholdSet, error) {
	event := new(MintControllerSignerThresholdSet)
	if err := _MintController.contract.UnpackLog(event, "SignerThresholdSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
