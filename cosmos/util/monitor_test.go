package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
)

func TestCreateMint(t *testing.T) {
	app.Config.Pocket.ChainID = "0001"
	testCases := []struct {
		name            string
		tx              *sdk.TxResponse
		result          *ValidateTxResult
		memo            models.MintMemo
		wpoktAddress    string
		vaultAddress    string
		expectedMint    models.Mint
		expectedErr     bool
		expectedUpdated time.Duration
	}{
		{
			name: "Valid Mint",
			tx: &sdk.TxResponse{
				Height: 12345,
				TxHash: "0x1234567890abcdef",
			},
			result: &ValidateTxResult{
				Memo: models.MintMemo{
					Address: "0x1234567890",
					ChainID: "0001",
				},
				TxStatus:      models.TransactionStatusPending,
				Confirmations: uint64(0),
				Tx:            nil,
				TxHash:        "0x1234567890abcdef",
				Amount:        sdk.NewCoin("upokt", math.NewInt(100)),
				SenderAddress: "0xabcdef",
				NeedsRefund:   false,
			},
			memo: models.MintMemo{
				Address: "0x1234567890",
				ChainID: "0001",
			},
			wpoktAddress: "0x9876543210",
			vaultAddress: "0xabc123def",
			expectedMint: models.Mint{
				Height:           "12345",
				Confirmations:    "0",
				TransactionHash:  "0x1234567890abcdef",
				SenderAddress:    "0xabcdef",
				SenderChainID:    app.Config.Pocket.ChainID,
				RecipientAddress: "0x1234567890",
				RecipientChainID: "0001",
				WPOKTAddress:     "0x9876543210",
				VaultAddress:     "0xabc123def",
				Amount:           "100",
				Memo: &models.MintMemo{
					Address: "0x1234567890",
					ChainID: "0001",
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

			result := CreateMint(tc.tx, tc.result, tc.wpoktAddress, tc.vaultAddress)

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
		tx                  *sdk.TxResponse
		result              *ValidateTxResult
		vaultAddress        string
		expectedInvalidMint models.InvalidMint
		expectedErr         bool
		expectedUpdated     time.Duration
	}{
		{
			name: "Valid Invalid Mint",
			tx: &sdk.TxResponse{
				Height: 12345,
				TxHash: "0x1234567890abcdef",
			},
			result: &ValidateTxResult{
				Memo: models.MintMemo{
					Address: "0x1234567890",
					ChainID: "0001",
				},
				TxStatus:      models.TransactionStatusPending,
				Confirmations: uint64(0),
				Tx: &tx.Tx{
					Body: &tx.TxBody{
						Memo: "Invalid mint memo",
					},
				},
				TxHash:        "0x1234567890abcdef",
				Amount:        sdk.NewCoin("upokt", math.NewInt(100)),
				SenderAddress: "0xabcdef",
				NeedsRefund:   true,
			},
			vaultAddress: "0xabc123def",
			expectedInvalidMint: models.InvalidMint{
				Height:                "12345",
				Confirmations:         "0",
				TransactionHash:       "0x1234567890abcdef",
				SenderAddress:         "0xabcdef",
				SenderChainID:         app.Config.Pocket.ChainID,
				Memo:                  "Invalid mint memo",
				Amount:                "100",
				VaultAddress:          "0xabc123def",
				CreatedAt:             time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				UpdatedAt:             time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				Status:                models.StatusPending,
				Signatures:            []models.Signature{},
				Sequence:              nil,
				ReturnTransactionHash: "",
				ReturnTransactionBody: "",
			},
			expectedErr:     false,
			expectedUpdated: 2 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.ChainID = "0001"

			result := CreateInvalidMint(tc.tx, tc.result, tc.vaultAddress)

			assert.WithinDuration(t, time.Now(), result.CreatedAt, tc.expectedUpdated)
			assert.WithinDuration(t, time.Now(), result.UpdatedAt, tc.expectedUpdated)

			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedInvalidMint, result)
		})
	}
}

func TestCreateFailedMint(t *testing.T) {
	testCases := []struct {
		name                string
		tx                  *sdk.TxResponse
		result              *ValidateTxResult
		vaultAddress        string
		expectedInvalidMint models.InvalidMint
		expectedErr         bool
		expectedUpdated     time.Duration
	}{
		{
			name: "Valid Invalid Mint",
			tx: &sdk.TxResponse{
				Height: 12345,
				TxHash: "0x1234567890abcdef",
			},
			result: &ValidateTxResult{
				Memo: models.MintMemo{
					Address: "0x1234567890",
					ChainID: "0001",
				},
				TxStatus:      models.TransactionStatusFailed,
				Confirmations: uint64(0),
				Tx: &tx.Tx{
					Body: &tx.TxBody{
						Memo: "Invalid mint memo",
					},
				},
				TxHash:        "0x1234567890abcdef",
				Amount:        sdk.NewCoin("upokt", math.NewInt(100)),
				SenderAddress: "0xabcdef",
				NeedsRefund:   true,
			},
			vaultAddress: "0xabc123def",
			expectedInvalidMint: models.InvalidMint{
				Height:          "12345",
				Confirmations:   "0",
				TransactionHash: "0x1234567890abcdef",
				SenderAddress:   "0xabcdef",
				SenderChainID:   app.Config.Pocket.ChainID,
				Memo:            "Invalid mint memo",
				Amount:          "100",
				VaultAddress:    "0xabc123def",
				CreatedAt:       time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				UpdatedAt:       time.Time{}, // We'll use assert.WithinDuration to check if within an acceptable range
				Status:          models.StatusFailed,

				Signatures:            []models.Signature{},
				Sequence:              nil,
				ReturnTransactionHash: "",
				ReturnTransactionBody: "",
			},
			expectedErr:     false,
			expectedUpdated: 2 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.ChainID = "0001"

			result := CreateFailedMint(tc.tx, tc.result, tc.vaultAddress)

			assert.WithinDuration(t, time.Now(), result.CreatedAt, tc.expectedUpdated)
			assert.WithinDuration(t, time.Now(), result.UpdatedAt, tc.expectedUpdated)

			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}

			assert.Equal(t, tc.expectedInvalidMint, result)
		})
	}
}
