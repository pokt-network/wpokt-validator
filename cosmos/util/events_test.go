package util

import (
	"testing"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestParseMessageSenderEvent(t *testing.T) {
	tests := []struct {
		name        string
		events      []abci.Event
		expected    string
		expectedErr string
	}{
		{
			name: "Sender Found",
			events: []abci.Event{
				{Type: "message", Attributes: []abci.EventAttribute{
					{Key: ("sender"), Value: ("pokt1abcd")},
				}},
			},
			expected:    "pokt1abcd",
			expectedErr: "",
		},
		{
			name: "Sender Not Found",
			events: []abci.Event{
				{Type: "message", Attributes: []abci.EventAttribute{
					{Key: ("not_sender"), Value: ("value")},
				}},
			},
			expected:    "",
			expectedErr: "no sender found in message events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := ParseMessageSenderEvent(tt.events)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, sender)
			}
		})
	}
}

func TestParseCoinsReceivedEvents(t *testing.T) {
	denom := "upokt"
	receiver := "pokt1abcd"
	tests := []struct {
		name        string
		events      []abci.Event
		expected    sdk.Coin
		expectedErr string
	}{
		{
			name: "Coins Received",
			events: []abci.Event{
				{Type: "coin_received", Attributes: []abci.EventAttribute{
					{Key: ("receiver"), Value: (receiver)},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expected:    sdk.NewCoin(denom, math.NewInt(100)),
			expectedErr: "",
		},
		{
			name: "Invalid Denom",
			events: []abci.Event{
				{Type: "coin_received", Attributes: []abci.EventAttribute{
					{Key: ("receiver"), Value: (receiver)},
					{Key: ("amount"), Value: ("100invalid")},
				}},
			},
			expected:    sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr: "invalid coin denom",
		},
		{
			name: "Invalid Amount",
			events: []abci.Event{
				{Type: "coin_received", Attributes: []abci.EventAttribute{
					{Key: ("receiver"), Value: (receiver)},
					{Key: ("amount"), Value: ("invalid")},
				}},
			},
			expected:    sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr: "unable to parse coin amount",
		},
		{
			name: "Receiver Not Found",
			events: []abci.Event{
				{Type: "coin_received", Attributes: []abci.EventAttribute{
					{Key: ("receiver"), Value: ("another")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expected:    sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coins, err := ParseCoinsReceivedEvents(denom, receiver, tt.events)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, coins)
			}
		})
	}
}

func TestParseCoinsSpentEvents(t *testing.T) {
	denom := "upokt"
	tests := []struct {
		name        string
		events      []abci.Event
		expectedSp  string
		expectedAmt sdk.Coin
		expectedErr string
	}{
		{
			name: "Coins Spent",
			events: []abci.Event{
				{Type: "coin_spent", Attributes: []abci.EventAttribute{
					{Key: ("spender"), Value: ("pokt1abcd")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expectedSp:  "pokt1abcd",
			expectedAmt: sdk.NewCoin(denom, math.NewInt(100)),
			expectedErr: "",
		},
		{
			name: "Invalid Denom",
			events: []abci.Event{
				{Type: "coin_spent", Attributes: []abci.EventAttribute{
					{Key: ("spender"), Value: ("pokt1abcd")},
					{Key: ("amount"), Value: ("100invalid")},
				}},
			},
			expectedSp:  "pokt1abcd",
			expectedAmt: sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr: "invalid coin denom",
		},
		{
			name: "Invalid Amount",
			events: []abci.Event{
				{Type: "coin_spent", Attributes: []abci.EventAttribute{
					{Key: ("spender"), Value: ("pokt1abcd")},
					{Key: ("amount"), Value: ("invalid")},
				}},
			},
			expectedSp:  "pokt1abcd",
			expectedAmt: sdk.NewCoin(denom, math.NewInt(0)),
			expectedErr: "invalid decimal coin expression: invalid",
		},
		{
			name: "Multiple Spenders",
			events: []abci.Event{
				{Type: "coin_spent", Attributes: []abci.EventAttribute{
					{Key: ("spender"), Value: ("pokt1abcd")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
				{Type: "coin_spent", Attributes: []abci.EventAttribute{
					{Key: ("spender"), Value: ("pokt1efgh")},
					{Key: ("amount"), Value: ("100upokt")},
				}},
			},
			expectedSp:  "pokt1abcd",
			expectedAmt: sdk.NewCoin(denom, math.NewInt(100)),
			expectedErr: "multiple spenders found in coin spent events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spender, amount, err := ParseCoinsSpentEvents(denom, tt.events)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSp, spender)
				assert.Equal(t, tt.expectedAmt, amount)
			}
		})
	}
}
