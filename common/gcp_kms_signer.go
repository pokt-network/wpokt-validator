package common

import (
	"context"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	dcrecSecp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	cryptotypes "github.com/cometbft/cometbft/crypto"

	gax "github.com/googleapis/gax-go/v2"

	"github.com/cosmos/cosmos-sdk/crypto/types"
)

// Interface Definition
type Signer interface {
	EthSign(data []byte) ([]byte, error)
	CosmosSign(data []byte) ([]byte, error)
	EthAddress() common.Address
	CosmosPublicKey() types.PubKey
	Destroy()
}

type GCPKeyManagementClient interface {
	Close() error
	GetPublicKey(ctx context.Context, req *kmspb.GetPublicKeyRequest, opts ...gax.CallOption) (*kmspb.PublicKey, error)
	AsymmetricSign(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error)
	GetCryptoKeyVersion(ctx context.Context, req *kmspb.GetCryptoKeyVersionRequest, opts ...gax.CallOption) (*kmspb.CryptoKeyVersion, error)
}

// Struct Definition
type GcpKmsSigner struct {
	client          GCPKeyManagementClient
	keyName         string
	ethAddress      common.Address
	cosmosPubKey    types.PubKey
	secp256k1PubKey *dcrecSecp256k1.PublicKey
}

var _ Signer = &GcpKmsSigner{}

var NewGCPKeyManagementClient = func(ctx context.Context) (GCPKeyManagementClient, error) {
	return kms.NewKeyManagementClient(ctx)
}

// Constructor Function
func NewGcpKmsSigner(keyName string) (Signer, error) {
	client, err := NewGCPKeyManagementClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS client: %w", err)
	}

	// verify key algorithm
	keyVersionDetails, err := resolveKeyVersionDetails(client, keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get key version details: %w", err)
	}

	if keyVersionDetails.Algorithm != kmspb.CryptoKeyVersion_EC_SIGN_SECP256K1_SHA256 {
		return nil, fmt.Errorf("key algorithm is not EC_SIGN_P256_SHA256")
	}

	// resolve public key
	pubKeyBytes, err := resolvePubKeyBytes(client, keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve public key: %w", err)
	}

	ethPublicKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %w", err)
	}

	ethAddress := getEthAddr(pubKeyBytes)

	if ethAddress != crypto.PubkeyToAddress(*ethPublicKey) {
		return nil, fmt.Errorf("ethereum address mismatch")
	}

	secp256k1PubKey, err := getSecp256k1PubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get secp256k1 public key: %w", err)
	}

	cosmosPubKey := &secp256k1.PubKey{Key: secp256k1PubKey.SerializeCompressed()}

	return &GcpKmsSigner{
		client:          client,
		keyName:         keyName,
		ethAddress:      ethAddress,
		cosmosPubKey:    cosmosPubKey,
		secp256k1PubKey: secp256k1PubKey,
	}, nil
}

// Destructor Function
func (s *GcpKmsSigner) Destroy() {
	s.client.Close()
}

// Method Implementations
func (s *GcpKmsSigner) EthSign(data []byte) ([]byte, error) {
	digest := data
	if len(digest) != 32 {
		digest = crypto.Keccak256(data)
	}
	hash := common.BytesToHash(digest)
	return ethSignHash(hash, s.client, s.keyName, s.ethAddress)
}

func (s *GcpKmsSigner) CosmosSign(data []byte) ([]byte, error) {
	digest := data
	if len(digest) != 32 {
		digest = cryptotypes.Sha256(data)
	}
	hash := common.BytesToHash(digest)
	return cosmosSignHash(s.client, s.keyName, hash, s.secp256k1PubKey)
}

func (s *GcpKmsSigner) EthAddress() common.Address {
	return s.ethAddress
}

func (s *GcpKmsSigner) CosmosPublicKey() types.PubKey {
	return s.cosmosPubKey
}

func resolvePubKeyBytes(client GCPKeyManagementClient, keyName string) ([]byte, error) {
	publicKeyResp, err := client.GetPublicKey(context.Background(), &kmspb.GetPublicKeyRequest{Name: keyName})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	publicKeyPem := publicKeyResp.Pem

	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil {
		return nil, fmt.Errorf("public key %q PEM empty: %.130q", keyName, publicKeyPem)
	}

	var info struct {
		AlgID pkix.AlgorithmIdentifier
		Key   asn1.BitString
	}
	_, err = asn1.Unmarshal(block.Bytes, &info)
	if err != nil {
		return nil, fmt.Errorf("public key %q PEM block %q: %w", keyName, block.Type, err)
	}

	wantAlg := asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	if gotAlg := info.AlgID.Algorithm; !gotAlg.Equal(wantAlg) {
		return nil, fmt.Errorf("public key %q ASN.1 algorithm %s instead of %s", keyName, gotAlg, wantAlg)
	}

	return info.Key.Bytes, nil
}

// PubKeyAddr returns the Ethereum address for (uncompressed-)key bytes.
func getEthAddr(bytes []byte) common.Address {
	digest := crypto.Keccak256(bytes[1:])
	var addr common.Address
	copy(addr[:], digest[12:])
	return addr
}

func ethSignHash(hash common.Hash, client GCPKeyManagementClient, keyName string, ethAddress common.Address) ([]byte, error) {
	// Resolve a signature
	req := kmspb.AsymmetricSignRequest{
		Name: keyName,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: hash[:],
			},
		},
	}
	resp, err := client.AsymmetricSign(context.Background(), &req)
	if err != nil {
		return nil, fmt.Errorf("asymmetric sign operation: %w", err)
	}

	// Parse signature
	var params struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(resp.Signature, &params)
	if err != nil {
		return nil, fmt.Errorf("asymmetric signature encoding: %w", err)
	}

	var rLen, sLen int // byte size
	if params.R != nil {
		rLen = (params.R.BitLen() + 7) / 8
	}
	if params.S != nil {
		sLen = (params.S.BitLen() + 7) / 8
	}
	if rLen == 0 || rLen > 32 || sLen == 0 || sLen > 32 {
		return nil, fmt.Errorf("asymmetric signature with %d-byte r and %d-byte s denied on size", rLen, sLen)
	}

	// Need uncompressed signature with "recovery ID" at end:
	// https://bitcointalk.org/index.php?topic=5249677.0
	// https://ethereum.stackexchange.com/a/53182/39582
	var sig [66]byte // + 1-byte header + 1-byte tailer
	params.R.FillBytes(sig[33-rLen : 33])
	params.S.FillBytes(sig[65-sLen : 65])

	// Brute force try includes KMS verification
	var recoverErr error
	var finalSig []byte
	for recoveryID := byte(0); recoveryID < 2; recoveryID++ {
		sig[0] = recoveryID + 27 // BitCoin header
		btcsig := sig[:65]       // Exclude Ethereum 'v' parameter
		pubKey, _, err := btcecdsa.RecoverCompact(btcsig, hash[:])
		if err != nil {
			recoverErr = err
			continue
		}

		if getEthAddr(pubKey.SerializeUncompressed()) == ethAddress {
			// Sign the transaction
			sig[65] = recoveryID // Ethereum 'v' parameter

			finalSig = sig[1:] // Exclude BitCoin header
			break

		}
	}

	if recoverErr != nil {
		return nil, fmt.Errorf("asymmetric signature address recovery failed: %w", recoverErr)
	}

	if finalSig == nil {
		return nil, fmt.Errorf("signature address mismatch")
	}

	recoveredPubKey, err := crypto.Ecrecover(hash[:], finalSig)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %w", err)
	}

	recoveredPubKeyECDSA, err := crypto.UnmarshalPubkey(recoveredPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal recoverred public key: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKeyECDSA)
	if recoveredAddr != ethAddress {
		return nil, fmt.Errorf("recovered address mismatch")
	}

	if finalSig[64] < 4 {
		finalSig[64] += 27
	}

	return finalSig, nil
}

func cosmosSignHash(client GCPKeyManagementClient, keyName string, hash [32]byte, pubKey *dcrecSecp256k1.PublicKey) ([]byte, error) {
	// Sign the hash using KMS
	req := &kmspb.AsymmetricSignRequest{
		Name: keyName,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: hash[:],
			},
		},
	}

	resp, err := client.AsymmetricSign(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	signature := resp.Signature

	// Extract r and s values from the signature
	var params struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(signature, &params)
	if err != nil {
		return nil, fmt.Errorf("asymmetric signature encoding: %w", err)
	}

	rBytes := params.R.Bytes()
	sBytes := params.S.Bytes()

	// Ensure r and s are 32 bytes each
	rPadded := make([]byte, 32)
	sPadded := make([]byte, 32)
	copy(rPadded[32-len(rBytes):], rBytes)
	copy(sPadded[32-len(sBytes):], sBytes)

	finalSig := append(rPadded, sPadded...)

	sig, err := signatureFromBytes(finalSig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	if !sig.Verify(hash[:], pubKey) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return finalSig, nil
}

func signatureFromBytes(sigStr []byte) (*btcecdsa.Signature, error) {
	if len(sigStr) != 64 {
		return nil, fmt.Errorf("signature length is not 64 bytes")
	}

	var r dcrecSecp256k1.ModNScalar
	r.SetByteSlice(sigStr[:32])
	var s dcrecSecp256k1.ModNScalar
	s.SetByteSlice(sigStr[32:64])
	if s.IsOverHalfOrder() {
		return nil, fmt.Errorf("signature is not in lower-S form")
	}

	return btcecdsa.NewSignature(&r, &s), nil
}

// func getCosmosPubKey(pubKeyBytes []byte) (*secp256k1.PubKey, error) {
// 	pubkeyObject, err := getSecp256k1PubKey(pubKeyBytes)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	pk := pubkeyObject.SerializeCompressed()
//
// 	return &secp256k1.PubKey{Key: pk}, nil
// }

func getSecp256k1PubKey(pubKeyBytes []byte) (*dcrecSecp256k1.PublicKey, error) {
	pubkeyObject, err := dcrecSecp256k1.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return pubkeyObject, nil
}

func resolveKeyVersionDetails(client GCPKeyManagementClient, keyName string) (*kmspb.CryptoKeyVersion, error) {
	// Request the key version details
	req := &kmspb.GetCryptoKeyVersionRequest{
		Name: keyName,
	}

	resp, err := client.GetCryptoKeyVersion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get key version details: %w", err)
	}

	return resp, nil

}
