package util

import (
	"testing"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestParseTransferEvents(t *testing.T) {
	denom := "upokt"
	tests := []struct {
		name              string
		events            []abci.Event
		expectedSender    string
		expectedRecipient string
		expectedAmount    sdk.Coin
		expectedErr       string
	}{
		{
			name: "Transfer",
			events: []abci.Event{
				{Type: "transfer", Attributes: []abci.EventAttribute{
					{Key: ("sender"), Value: ("pokt1abcd")},
					{Key: ("recipient"), Value: ("pokt1efgh")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expectedSender:    "pokt1abcd",
			expectedRecipient: "pokt1efgh",
			expectedAmount:    sdk.NewCoin(denom, math.NewInt(100)),
			expectedErr:       "",
		},
		{
			name: "No Sender",
			events: []abci.Event{
				{Type: "transfer", Attributes: []abci.EventAttribute{
					{Key: ("recipient"), Value: ("pokt1efgh")},
					{Key: ("amount"), Value: ("100invalid")},
				}},
			},
			expectedSender:    "pokt1abcd",
			expectedRecipient: "pokt1efgh",
			expectedAmount:    sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr:       ("no attribute found with key: sender"),
		},
		{
			name: "No Recipient",
			events: []abci.Event{
				{Type: "transfer", Attributes: []abci.EventAttribute{
					{Key: ("sender"), Value: ("pokt1abcd")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expectedSender:    "pokt1abcd",
			expectedRecipient: "pokt1efgh",
			expectedAmount:    sdk.NewCoin(denom, math.NewInt(100)),
			expectedErr:       ("no attribute found with key: recipient"),
		},
		{
			name: "Different Recipient",
			events: []abci.Event{
				{Type: "transfer", Attributes: []abci.EventAttribute{
					{Key: ("sender"), Value: ("pokt1abcd")},
					{Key: ("recipient"), Value: ("pokt1abcd")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expectedSender:    "",
			expectedRecipient: "",
			expectedAmount:    sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr:       "",
		},
		{
			name: "Invalid Amount",
			events: []abci.Event{
				{Type: "transfer", Attributes: []abci.EventAttribute{
					{Key: ("sender"), Value: ("pokt1abcd")},
					{Key: ("recipient"), Value: ("pokt1efgh")},
					{Key: ("amount"), Value: ("invalid")},
				}},
			},
			expectedSender: "pokt1abcd",
			expectedAmount: sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr:    "invalid decimal coin expression: invalid",
		},
		{
			name: "Invalid Denom",
			events: []abci.Event{
				{Type: "transfer", Attributes: []abci.EventAttribute{
					{Key: ("sender"), Value: ("pokt1abcd")},
					{Key: ("recipient"), Value: ("pokt1efgh")},
					{Key: ("amount"), Value: ("100invalid")},
				}},
			},
			expectedSender:    "",
			expectedRecipient: "",
			expectedAmount:    sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfers, err := ParseTransferEvents(tt.events, "pokt1efgh", denom)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				if tt.expectedSender != "" {
					assert.Equal(t, tt.expectedSender, transfers[0].Sender)
					assert.Equal(t, tt.expectedRecipient, transfers[0].Receiver)
					assert.Equal(t, tt.expectedAmount, transfers[0].Amount)
				}
			}
		})
	}

	t.Run("No Transfer Events", func(t *testing.T) {
		transfers, err := ParseTransferEvents([]abci.Event{}, "pokt1efgh", denom)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(transfers))
	})

	t.Run("Multiple Transfer Events", func(t *testing.T) {
		transfers, err := ParseTransferEvents([]abci.Event{
			{Type: "transfer", Attributes: []abci.EventAttribute{
				{Key: ("sender"), Value: ("pokt1abcd")},
				{Key: ("recipient"), Value: ("pokt1efgh")},
				{Key: ("amount"), Value: ("100upokt")},
			}},
			{Type: "transfer", Attributes: []abci.EventAttribute{
				{Key: ("sender"), Value: ("pokt1abcd")},
				{Key: ("recipient"), Value: ("pokt1efgh")},
				{Key: ("amount"), Value: ("100upokt")},
			}},
		}, "pokt1efgh", denom)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(transfers))
	})
}
