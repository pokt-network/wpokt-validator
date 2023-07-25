package ethereum

import (
	"testing"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/models"
)

func TestUpdateStatusAndConfirmationsForMint(t *testing.T) {
	testCases := []struct {
		name                  string
		initialMint           models.Mint
		poktHeight            int64
		requiredConfirmations int64
		expectedStatus        string
		expectedConfs         string
	}{
		{
			name: "Status is pending and confirmations are 0",
			initialMint: models.Mint{
				Status:        models.StatusPending,
				Confirmations: "0",
				Height:        "1000",
			},
			poktHeight:            1000,
			requiredConfirmations: 5,
			expectedStatus:        models.StatusPending,
			expectedConfs:         "0",
		},
		{
			name: "Status is pending and confirmations are 0 but less than required confirmations",
			initialMint: models.Mint{
				Status:        models.StatusPending,
				Confirmations: "0",
				Height:        "1000",
			},
			poktHeight:            1005,
			requiredConfirmations: 10,
			expectedStatus:        models.StatusPending,
			expectedConfs:         "5",
		},
		{
			name: "Status is pending and confirmations are greater than 0",
			initialMint: models.Mint{
				Status:        models.StatusPending,
				Confirmations: "10",
				Height:        "1000",
			},
			poktHeight:            1020,
			requiredConfirmations: 10,
			expectedStatus:        models.StatusConfirmed,
			expectedConfs:         "20",
		},
		{
			name: "Status is confirmed",
			initialMint: models.Mint{
				Status:        models.StatusPending,
				Confirmations: "3",
				Height:        "1000",
			},
			poktHeight:            1024,
			requiredConfirmations: 2,
			expectedStatus:        models.StatusConfirmed,
			expectedConfs:         "24",
		},
		{
			name: "Status is confirmed and required confirmations are 0",
			initialMint: models.Mint{
				Status:        models.StatusConfirmed,
				Confirmations: "3",
				Height:        "1000",
			},
			poktHeight:            1024,
			requiredConfirmations: 0,
			expectedStatus:        models.StatusConfirmed,
			expectedConfs:         "3",
		},
		{
			name: "Status is pending and required confirmations are 0",
			initialMint: models.Mint{
				Status:        models.StatusPending,
				Confirmations: "3",
				Height:        "1000",
			},
			poktHeight:            1024,
			requiredConfirmations: 0,
			expectedStatus:        models.StatusConfirmed,
			expectedConfs:         "3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Config.Pocket.Confirmations = tc.requiredConfirmations

			result, err := updateStatusAndConfirmationsForMint(tc.initialMint, tc.poktHeight)
			if err != nil {
				t.Errorf("Error should be nil, got: %v", err)
			}

			if result.Status != tc.expectedStatus {
				t.Errorf("Expected status %s, got %s", tc.expectedStatus, result.Status)
			}

			if result.Confirmations != tc.expectedConfs {
				t.Errorf("Expected confirmations %s, got %s", tc.expectedConfs, result.Confirmations)
			}
		})
	}
}
