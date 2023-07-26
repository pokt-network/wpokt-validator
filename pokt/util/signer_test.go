package util

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

func TestUpdateStatusAndConfirmationsForInvalidMint(t *testing.T) {
	testCases := []struct {
		name                  string
		doc                   models.InvalidMint
		currentHeight         int64
		requiredConfirmations int64
		expectedDoc           models.InvalidMint
		expectedErr           bool
	}{
		{
			name: "Status Pending, Confirmations 0 and requiredConfirmations 0",
			doc: models.InvalidMint{
				Status:        models.StatusPending,
				Confirmations: "0",
				Height:        "100",
			},
			currentHeight:         150,
			requiredConfirmations: 0,
			expectedDoc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "0",
				Height:        "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Pending, Confirmations 0",
			doc: models.InvalidMint{
				Status:        models.StatusPending,
				Confirmations: "0",
				Height:        "100",
			},
			currentHeight:         150,
			requiredConfirmations: 50,
			expectedDoc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "50",
				Height:        "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Pending, Confirmations < 0",
			doc: models.InvalidMint{
				Status:        models.StatusPending,
				Confirmations: "-1",
				Height:        "100",
			},
			currentHeight:         150,
			requiredConfirmations: 50,
			expectedDoc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "50",
				Height:        "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Confirmed",
			doc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "10",
				Height:        "100",
			},
			currentHeight:         150,
			requiredConfirmations: 50,
			expectedDoc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "10",
				Height:        "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Pending, Confirmations >= Config.Pocket.Confirmations",
			doc: models.InvalidMint{
				Status:        models.StatusPending,
				Confirmations: "0",
				Height:        "100",
			},
			currentHeight:         200,
			requiredConfirmations: 100,
			expectedDoc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "100",
				Height:        "100",
			},
			expectedErr: false,
		},
		{
			name: "Invalid Height",
			doc: models.InvalidMint{
				Status:        models.StatusPending,
				Confirmations: "0",
				Height:        "height",
			},
			currentHeight:         200,
			requiredConfirmations: 100,
			expectedDoc: models.InvalidMint{
				Status:        models.StatusConfirmed,
				Confirmations: "100",
				Height:        "100",
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.Confirmations = tc.requiredConfirmations

			result, err := UpdateStatusAndConfirmationsForInvalidMint(tc.doc, tc.currentHeight)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDoc, result)
			}

		})
	}
}

func TestUpdateStatusAndConfirmationsForBurn(t *testing.T) {
	testCases := []struct {
		name                  string
		doc                   models.Burn
		blockNumber           int64
		requiredConfirmations int64
		expectedDoc           models.Burn
		expectedErr           bool
	}{
		{
			name: "Status Pending, Confirmations 0 and requiredConfirmations 0",
			doc: models.Burn{
				Status:        models.StatusPending,
				Confirmations: "0",
				BlockNumber:   "100",
			},
			blockNumber:           150,
			requiredConfirmations: 0,
			expectedDoc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "0",
				BlockNumber:   "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Pending, Confirmations 0",
			doc: models.Burn{
				Status:        models.StatusPending,
				Confirmations: "0",
				BlockNumber:   "100",
			},
			blockNumber:           150,
			requiredConfirmations: 50,
			expectedDoc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "50",
				BlockNumber:   "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Pending, Confirmations < 0",
			doc: models.Burn{
				Status:        models.StatusPending,
				Confirmations: "-1",
				BlockNumber:   "100",
			},
			blockNumber:           150,
			requiredConfirmations: 50,
			expectedDoc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "50",
				BlockNumber:   "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Confirmed",
			doc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "10",
				BlockNumber:   "100",
			},
			blockNumber:           150,
			requiredConfirmations: 50,
			expectedDoc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "10",
				BlockNumber:   "100",
			},
			expectedErr: false,
		},
		{
			name: "Status Pending, Confirmations >= Config.Ethereum.Confirmations",
			doc: models.Burn{
				Status:        models.StatusPending,
				Confirmations: "0",
				BlockNumber:   "100",
			},
			blockNumber:           200,
			requiredConfirmations: 100,
			expectedDoc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "100",
				BlockNumber:   "100",
			},
			expectedErr: false,
		},
		{
			name: "Invalid Block Number",
			doc: models.Burn{
				Status:        models.StatusPending,
				Confirmations: "0",
				BlockNumber:   "number",
			},
			blockNumber:           200,
			requiredConfirmations: 100,
			expectedDoc: models.Burn{
				Status:        models.StatusConfirmed,
				Confirmations: "100",
				BlockNumber:   "100",
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Ethereum.Confirmations = tc.requiredConfirmations

			result, err := UpdateStatusAndConfirmationsForBurn(tc.doc, tc.blockNumber)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDoc, result)
			}

		})
	}
}

func TestSignBurn(t *testing.T) {
	privateKey1, _ := crypto.NewPrivateKey("8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82")
	privateKey2, _ := crypto.NewPrivateKey("f2c227cd1299f62750e48d3e44c2d29cb3add4c8e9a171ae260b8fdeff49c761ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d055")
	privateKey3, _ := crypto.NewPrivateKey("05339b10520335644fe486e4d39ce33db4d079b5c1d3bceb725e75e4354f5ca7351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b")

	privateKeys := []crypto.PrivateKey{privateKey1, privateKey2, privateKey3}

	address := privateKeys[2].PublicKey().Address().String()

	testCases := []struct {
		name        string
		doc         models.Burn
		expectedDoc models.Burn
		numSigners  int
		expectedErr bool
		privateKey  crypto.PrivateKey
	}{
		{
			name:       "Invalid Amount",
			numSigners: 2,
			doc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "",
				Signers:          []string{},
				RecipientAddress: address,
				Amount:           "amount",
				TransactionHash:  "transaction_hash",
			},
			expectedDoc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "",
				Signers:          []string{},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedErr: true,
			privateKey:  privateKeys[0],
		},
		{
			name:       "Sign without ReturnTx",
			numSigners: 2,
			doc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "",
				Signers:          []string{},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedDoc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "",
				Signers:          []string{},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
		{
			name:       "Sign without ReturnTx 1 signer fully signed",
			numSigners: 1,
			doc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "",
				Signers:          []string{},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedDoc: models.Burn{
				Status:           models.StatusSigned,
				Confirmations:    "0",
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
		{
			name: "Sign with ReturnTx 2 signers fully signed",
			doc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "d7020a470a102f782e6e6f6465732e4d736753656e6412330a14ec1a1c10d8e5290e1a3027fee4dd4f670e5c915d1214e8ae8fdce6b5fc62455f778d6660bc330a47aef11a053930303030120e0a0575706f6b74120531303030301adf010a52f325b8ad0a259d544774206a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad820a259d54477420ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d055128801b2f515f90a40b4dbdd1a36270a23ea1227bd037e0adfd72a1d57aa32d84812ff0d3739c81d804d72ffa8fadbb5e06fae221a9510a6b50a06591ca328aa3c8a44c58ead91f8080a40b4dbdd1a36270a23ea1227bd037e0adfd72a1d57aa32d84812ff0d3739c81d804d72ffa8fadbb5e06fae221a9510a6b50a06591ca328aa3c8a44c58ead91f80822107472616e73616374696f6e5f6861736828b2f398cae4f0f0c86f",
				Signers:          []string{privateKeys[1].PublicKey().RawString()},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			numSigners: 2,
			expectedDoc: models.Burn{
				Status:           models.StatusSigned,
				Confirmations:    "0",
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
		{
			name: "Sign with ReturnTx 3 signers",
			doc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "c1030a470a102f782e6e6f6465732e4d736753656e6412330a140b7a71f5baa23e493be517ac44264b4d5b90b3521214e8ae8fdce6b5fc62455f778d6660bc330a47aef11a053930303030120e0a0575706f6b74120531303030301ac8020a79f325b8ad0a259d544774206a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad820a259d54477420ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d0550a259d54477420351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b12ca01b2f515f90a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d79000a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d79000a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d790022107472616e73616374696f6e5f6861736828ded6cbadc48ee2a1b101",
				Signers:          []string{privateKeys[2].PublicKey().RawString()},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			numSigners: 3,
			expectedDoc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[1],
		},
		{
			name: "Sign with ReturnTx 3 signers fully signed",
			doc: models.Burn{
				Status:           models.StatusPending,
				Confirmations:    "0",
				ReturnTx:         "c1030a470a102f782e6e6f6465732e4d736753656e6412330a140b7a71f5baa23e493be517ac44264b4d5b90b3521214e8ae8fdce6b5fc62455f778d6660bc330a47aef11a053930303030120e0a0575706f6b74120531303030301ac8020a79f325b8ad0a259d544774206a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad820a259d54477420ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d0550a259d54477420351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b12ca01b2f515f90a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d79000a4096a74f64b20a5e274bf1c0cda4dc9f8bbde016a59d9701947529e727ddb9d29d3c3a421e54884993050eceaff9781d268ce3456fdb617086010afded6c23dd050a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d790022107472616e73616374696f6e5f6861736828ded6cbadc48ee2a1b101",
				Signers:          []string{privateKeys[2].PublicKey().RawString(), privateKeys[1].PublicKey().RawString()},
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			numSigners: 3,
			expectedDoc: models.Burn{
				Status:           models.StatusSigned,
				Confirmations:    "0",
				RecipientAddress: address,
				Amount:           "100000",
				TransactionHash:  "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.ChainId = "testnet"
			app.Config.Pocket.Fees = 10000
			pubKeys := []crypto.PublicKey{}
			for i := 0; i < tc.numSigners; i++ {
				pubKeys = append(pubKeys, privateKeys[i].PublicKey())
			}
			multisigPubKey := crypto.PublicKeyMultiSignature{PublicKeys: pubKeys}

			result, err := SignBurn(tc.doc, tc.privateKey, multisigPubKey, tc.numSigners)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, "", result.ReturnTx)

				tx, bytesToSign, err := decodeTx(result.ReturnTx, app.Config.Pocket.ChainId)
				assert.NoError(t, err)
				assert.NotNil(t, tx)
				assert.NotNil(t, bytesToSign)

				fee := types.NewCoins(types.NewCoin(types.DefaultStakeDenom, types.NewInt(app.Config.Pocket.Fees)))
				assert.Equal(t, fee, tx.GetFee())

				assert.Equal(t, tc.doc.TransactionHash, tx.GetMemo())

				assert.Equal(t, strings.ToLower(tc.doc.RecipientAddress), strings.ToLower(tx.GetMsg().(*nodeTypes.MsgSend).ToAddress.String()))

				fa, _ := types.AddressFromHex(multisigPubKey.Address().String())
				assert.Equal(t, fa, tx.GetMsg().(*nodeTypes.MsgSend).FromAddress)

				amountInt, _ := strconv.ParseInt(tc.doc.Amount, 10, 64)
				finalAmount := amountInt - app.Config.Pocket.Fees
				amount := strconv.FormatInt(finalAmount, 10)
				assert.Equal(t, amount, tx.GetMsg().(*nodeTypes.MsgSend).Amount.String())

				assert.NotNil(t, tx.GetSignature())
				var ms = crypto.MultiSig(crypto.MultiSignature{})
				ms = ms.Unmarshal(tx.GetSignature().GetSignature())
				assert.Equal(t, tc.numSigners, ms.NumOfSigs())

				signers := tc.doc.Signers
				signers = append(signers, tc.privateKey.PublicKey().RawString())
				tc.expectedDoc.Signers = signers

				sigs := [][]byte{}
				for i := 0; i < tc.numSigners; i++ {
					sig, _ := ms.GetSignatureByIndex(i)
					sigs = append(sigs, sig)
				}

				if tc.doc.ReturnTx == "" {
					for i := 0; i < tc.numSigners; i++ {
						assert.Equal(t, sigs[i], sigs[0])
					}
				}

				result.ReturnTx = ""
				assert.Equal(t, tc.expectedDoc, result)

				for i := 0; i < len(tc.expectedDoc.Signers); i++ {
					pubKey, _ := crypto.NewPublicKey(tc.expectedDoc.Signers[i])
					var index = -1
					for j := 0; j < tc.numSigners; j++ {
						if pubKey.RawString() == privateKeys[j].PublicKey().RawString() {
							index = j
							break
						}
					}
					assert.True(t, pubKey.VerifyBytes(bytesToSign, sigs[index]))
				}
			}

		})
	}
}

func TestSignInvalidMint(t *testing.T) {
	privateKey1, _ := crypto.NewPrivateKey("8d8da5d374c559b2f80c99c0f4cfb4405b6095487989bb8a5d5a7e579a4e76646a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad82")
	privateKey2, _ := crypto.NewPrivateKey("f2c227cd1299f62750e48d3e44c2d29cb3add4c8e9a171ae260b8fdeff49c761ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d055")
	privateKey3, _ := crypto.NewPrivateKey("05339b10520335644fe486e4d39ce33db4d079b5c1d3bceb725e75e4354f5ca7351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b")

	privateKeys := []crypto.PrivateKey{privateKey1, privateKey2, privateKey3}

	address := privateKeys[2].PublicKey().Address().String()

	testCases := []struct {
		name        string
		doc         models.InvalidMint
		expectedDoc models.InvalidMint
		numSigners  int
		expectedErr bool
		privateKey  crypto.PrivateKey
	}{
		{
			name:       "Invalid Amount",
			numSigners: 2,
			doc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "",
				Signers:         []string{},
				SenderAddress:   address,
				Amount:          "amount",
				TransactionHash: "transaction_hash",
			},
			expectedDoc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "",
				Signers:         []string{},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedErr: true,
			privateKey:  privateKeys[0],
		},
		{
			name:       "Sign without ReturnTx",
			numSigners: 2,
			doc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "",
				Signers:         []string{},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedDoc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "",
				Signers:         []string{},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
		{
			name:       "Sign without ReturnTx 1 signer fully signed",
			numSigners: 1,
			doc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "",
				Signers:         []string{},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedDoc: models.InvalidMint{
				Status:          models.StatusSigned,
				Confirmations:   "0",
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
		{
			name: "Sign with ReturnTx 2 signers fully signed",
			doc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "d7020a470a102f782e6e6f6465732e4d736753656e6412330a14ec1a1c10d8e5290e1a3027fee4dd4f670e5c915d1214e8ae8fdce6b5fc62455f778d6660bc330a47aef11a053930303030120e0a0575706f6b74120531303030301adf010a52f325b8ad0a259d544774206a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad820a259d54477420ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d055128801b2f515f90a40b4dbdd1a36270a23ea1227bd037e0adfd72a1d57aa32d84812ff0d3739c81d804d72ffa8fadbb5e06fae221a9510a6b50a06591ca328aa3c8a44c58ead91f8080a40b4dbdd1a36270a23ea1227bd037e0adfd72a1d57aa32d84812ff0d3739c81d804d72ffa8fadbb5e06fae221a9510a6b50a06591ca328aa3c8a44c58ead91f80822107472616e73616374696f6e5f6861736828b2f398cae4f0f0c86f",
				Signers:         []string{privateKeys[1].PublicKey().RawString()},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			numSigners: 2,
			expectedDoc: models.InvalidMint{
				Status:          models.StatusSigned,
				Confirmations:   "0",
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
		{
			name: "Sign with ReturnTx 3 signers",
			doc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "c1030a470a102f782e6e6f6465732e4d736753656e6412330a140b7a71f5baa23e493be517ac44264b4d5b90b3521214e8ae8fdce6b5fc62455f778d6660bc330a47aef11a053930303030120e0a0575706f6b74120531303030301ac8020a79f325b8ad0a259d544774206a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad820a259d54477420ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d0550a259d54477420351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b12ca01b2f515f90a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d79000a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d79000a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d790022107472616e73616374696f6e5f6861736828ded6cbadc48ee2a1b101",
				Signers:         []string{privateKeys[2].PublicKey().RawString()},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			numSigners: 3,
			expectedDoc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[1],
		},
		{
			name: "Sign with ReturnTx 3 signers fully signed",
			doc: models.InvalidMint{
				Status:          models.StatusPending,
				Confirmations:   "0",
				ReturnTx:        "c1030a470a102f782e6e6f6465732e4d736753656e6412330a140b7a71f5baa23e493be517ac44264b4d5b90b3521214e8ae8fdce6b5fc62455f778d6660bc330a47aef11a053930303030120e0a0575706f6b74120531303030301ac8020a79f325b8ad0a259d544774206a456564a026788cd201a1a324a26d090e8df3dd0f3a233796552bdcaa95ad820a259d54477420ee604c6068452fa886c196afd7dd3a284ce9082d23baae2bfa6fe9cc1cd9d0550a259d54477420351799d14073dca9e5b7d50355b6b3a85d28a6a4b7f67ecb2ac8217732c4070b12ca01b2f515f90a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d79000a4096a74f64b20a5e274bf1c0cda4dc9f8bbde016a59d9701947529e727ddb9d29d3c3a421e54884993050eceaff9781d268ce3456fdb617086010afded6c23dd050a407749d57ad893e0a9a60fb9e7f13111a5d9dd2fc74f43eae83a59a22806e621041bbb37b7908902d490f759da63da9eaa1e5f73e9aea7182d34ff94c1474d790022107472616e73616374696f6e5f6861736828ded6cbadc48ee2a1b101",
				Signers:         []string{privateKeys[2].PublicKey().RawString(), privateKeys[1].PublicKey().RawString()},
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			numSigners: 3,
			expectedDoc: models.InvalidMint{
				Status:          models.StatusSigned,
				Confirmations:   "0",
				SenderAddress:   address,
				Amount:          "100000",
				TransactionHash: "transaction_hash",
			},
			expectedErr: false,
			privateKey:  privateKeys[0],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.ChainId = "testnet"
			app.Config.Pocket.Fees = 10000
			pubKeys := []crypto.PublicKey{}
			for i := 0; i < tc.numSigners; i++ {
				pubKeys = append(pubKeys, privateKeys[i].PublicKey())
			}
			multisigPubKey := crypto.PublicKeyMultiSignature{PublicKeys: pubKeys}

			result, err := SignInvalidMint(tc.doc, tc.privateKey, multisigPubKey, tc.numSigners)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, "", result.ReturnTx)

				tx, bytesToSign, err := decodeTx(result.ReturnTx, app.Config.Pocket.ChainId)
				assert.NoError(t, err)
				assert.NotNil(t, tx)
				assert.NotNil(t, bytesToSign)

				fee := types.NewCoins(types.NewCoin(types.DefaultStakeDenom, types.NewInt(app.Config.Pocket.Fees)))
				assert.Equal(t, fee, tx.GetFee())

				assert.Equal(t, tc.doc.TransactionHash, tx.GetMemo())

				assert.Equal(t, strings.ToLower(tc.doc.SenderAddress), strings.ToLower(tx.GetMsg().(*nodeTypes.MsgSend).ToAddress.String()))

				fa, _ := types.AddressFromHex(multisigPubKey.Address().String())
				assert.Equal(t, fa, tx.GetMsg().(*nodeTypes.MsgSend).FromAddress)

				amountInt, _ := strconv.ParseInt(tc.doc.Amount, 10, 64)
				finalAmount := amountInt - app.Config.Pocket.Fees
				amount := strconv.FormatInt(finalAmount, 10)
				assert.Equal(t, amount, tx.GetMsg().(*nodeTypes.MsgSend).Amount.String())

				assert.NotNil(t, tx.GetSignature())
				var ms = crypto.MultiSig(crypto.MultiSignature{})
				ms = ms.Unmarshal(tx.GetSignature().GetSignature())
				assert.Equal(t, tc.numSigners, ms.NumOfSigs())

				signers := tc.doc.Signers
				signers = append(signers, tc.privateKey.PublicKey().RawString())
				tc.expectedDoc.Signers = signers

				sigs := [][]byte{}
				for i := 0; i < tc.numSigners; i++ {
					sig, _ := ms.GetSignatureByIndex(i)
					sigs = append(sigs, sig)
				}

				if tc.doc.ReturnTx == "" {
					for i := 0; i < tc.numSigners; i++ {
						assert.Equal(t, sigs[i], sigs[0])
					}
				}

				result.ReturnTx = ""
				assert.Equal(t, tc.expectedDoc, result)

				for i := 0; i < len(tc.expectedDoc.Signers); i++ {
					pubKey, _ := crypto.NewPublicKey(tc.expectedDoc.Signers[i])
					var index = -1
					for j := 0; j < tc.numSigners; j++ {
						if pubKey.RawString() == privateKeys[j].PublicKey().RawString() {
							index = j
							break
						}
					}
					assert.True(t, pubKey.VerifyBytes(bytesToSign, sigs[index]))
				}
			}

		})
	}
}
