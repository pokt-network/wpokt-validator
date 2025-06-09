package common

import "encoding/asn1"

const (
	AddressLength          = 20
	CosmosPublicKeyLength  = 33
	DefaultBIP39Passphrase = ""
	DefaultCosmosHDPath    = "m/44'/118'/0'/0/0"
	DefaultETHHDPath       = "m/44'/60'/0'/0/0"
	ZeroAddress            = "0x0000000000000000000000000000000000000000"
)

var oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
