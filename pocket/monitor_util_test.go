package pocket

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
	pocket "github.com/dan13ram/wpokt-backend/pocket/client"
)

func TestCreateMint(t *testing.T) {
	app.Config.Pocket.ChainId = "0001"
	testCases := []struct {
		name            string
		tx              *pocket.TxResponse
		memo            models.MintMemo
		wpoktAddress    string
		vaultAddress    string
		expectedMint    models.Mint
		expectedErr     bool
		expectedUpdated time.Duration
	}{
		{
			name: "Valid Mint",
			tx: &pocket.TxResponse{
				Height: 12345,
				Hash:   "0x1234567890abcdef",
				StdTx: pocket.StdTx{
					Msg: pocket.Msg{
						Value: pocket.Value{
							FromAddress: "0xabcdef",
							Amount:      "100",
						},
					},
				},
			},
			memo: models.MintMemo{
				Address: "0x1234567890",
				ChainId: "0001",
			},
			wpoktAddress: "0x9876543210",
			vaultAddress: "0xabc123def",
			expectedMint: models.Mint{
				Height:           "12345",
				Confirmations:    "0",
				TransactionHash:  "0x1234567890abcdef",
				SenderAddress:    "0xabcdef",
				SenderChainId:    app.Config.Pocket.ChainId,
				RecipientAddress: "0x1234567890",
				RecipientChainId: "0001",
				WPOKTAddress:     "0x9876543210",
				VaultAddress:     "0xabc123def",
				Amount:           "100",
				Memo: &models.MintMemo{
					Address: "0x1234567890",
					ChainId: "0001",
				},
				CreatedAt:           time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				UpdatedAt:           time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				Status:              models.StatusPending,
				Data:                nil,
				MintTransactionHash: "",
				Signers:             []string{},
				Signatures:          []string{},
			},
			expectedErr:     false,
			expectedUpdated: 2 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result := createMint(tc.tx, tc.memo, tc.wpoktAddress, tc.vaultAddress)

			assert.WithinDuration(t, time.Now(), result.CreatedAt, tc.expectedUpdated)
			assert.WithinDuration(t, time.Now(), result.UpdatedAt, tc.expectedUpdated)

			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedMint, result)
		})
	}
}

func TestCreateInvalidMint(t *testing.T) {
	testCases := []struct {
		name                string
		tx                  *pocket.TxResponse
		vaultAddress        string
		expectedInvalidMint models.InvalidMint
		expectedErr         bool
		expectedUpdated     time.Duration
	}{
		{
			name: "Valid Invalid Mint",
			tx: &pocket.TxResponse{
				Height: 12345,
				Hash:   "0x1234567890abcdef",
				StdTx: pocket.StdTx{
					Msg: pocket.Msg{
						Value: pocket.Value{
							FromAddress: "0xabcdef",
							Amount:      "100",
						},
					},
					Memo: "Invalid mint memo",
				},
			},
			vaultAddress: "0xabc123def",
			expectedInvalidMint: models.InvalidMint{
				Height:          "12345",
				Confirmations:   "0",
				TransactionHash: "0x1234567890abcdef",
				SenderAddress:   "0xabcdef",
				SenderChainId:   app.Config.Pocket.ChainId,
				Memo:            "Invalid mint memo",
				Amount:          "100",
				VaultAddress:    "0xabc123def",
				CreatedAt:       time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				UpdatedAt:       time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				Status:          models.StatusPending,
				Signers:         []string{},
				ReturnTx:        "",
				ReturnTxHash:    "",
			},
			expectedErr:     false,
			expectedUpdated: 2 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.ChainId = "0001"

			result := createInvalidMint(tc.tx, tc.vaultAddress)

			assert.WithinDuration(t, time.Now(), result.CreatedAt, tc.expectedUpdated)
			assert.WithinDuration(t, time.Now(), result.UpdatedAt, tc.expectedUpdated)

			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedInvalidMint, result)
		})
	}
}

func TestValidateMemo(t *testing.T) {
	app.Config.Ethereum.ChainId = "1"
	address := common.HexToAddress("0x1234567890")

	testCases := []struct {
		name          string
		txMemo        string
		expectedMemo  models.MintMemo
		expectedValid bool
	}{
		{
			name:   "Valid Memo",
			txMemo: fmt.Sprintf(`{"address": "%s","chain_id": "0001"}`, strings.ToLower(address.Hex())),
			expectedMemo: models.MintMemo{
				Address: address.Hex(),
				ChainId: app.Config.Ethereum.ChainId,
			},
			expectedValid: true,
		},
		{
			name: "Invalid JSON Memo",
			txMemo: `{
				"address": "0x1234567890",
				"chain_id": "0001",
				"extraField": "invalid"
			}`,
			expectedMemo:  models.MintMemo{},
			expectedValid: false,
		},
		{
			name: "Invalid ChainID",
			txMemo: `{
				"address": "0x1234567890",
				"chain_id": "0002"
			}`,
			expectedMemo: models.MintMemo{
				Address: "0x1234567890",
				ChainId: "0002",
			},
			expectedValid: false,
		},
		{
			name: "Invalid Address",
			txMemo: `{
				"address": "0xinvalid",
				"chain_id": "0001"
			}`,
			expectedMemo: models.MintMemo{
				Address: "0xinvalid",
				ChainId: "0001",
			},
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result, valid := validateMemo(tc.txMemo)
			assert.Equal(t, tc.expectedValid, valid)

			if valid {
				assert.Equal(t, tc.expectedMemo, result)
			}
		})
	}
}
