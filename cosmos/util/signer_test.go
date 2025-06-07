package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
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

			result, err := UpdateStatusAndConfirmationsForInvalidMint(&tc.doc, tc.currentHeight)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDoc, *result)
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

			result, err := UpdateStatusAndConfirmationsForBurn(&tc.doc, tc.blockNumber)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDoc, *result)
			}

		})
	}
}
