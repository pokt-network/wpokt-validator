package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMemo(t *testing.T) {
	tests := []struct {
		name        string
		txMemo      string
		expectedErr string
	}{
		{
			name:        "Invalid JSON",
			txMemo:      `invalid json`,
			expectedErr: "failed to unmarshal memo",
		},
		{
			name:        "Invalid Ethereum Address",
			txMemo:      `{"address": "invalid_address", "chain_id": "1"}`,
			expectedErr: "invalid address",
		},
		{
			name:        "Zero Ethereum Address",
			txMemo:      `{"address": "0x0000000000000000000000000000000000000000", "chain_id": "1"}`,
			expectedErr: "zero address",
		},
		{
			name:        "Invalid Chain ID",
			txMemo:      `{"address": "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B", "chain_id": "invalid_chain_id"}`,
			expectedErr: "invalid chain id",
		},
		{
			name:        "Unsupported Chain ID",
			txMemo:      `{"address": "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B", "chain_id": "999"}`,
			expectedErr: "unsupported chain id",
		},
		{
			name:        "Valid Memo",
			txMemo:      `{"address": "0xAb5801a7D398351b8bE11C439e05C5b3259aec9B", "chain_id": "1"}`,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memo, err := ValidateMemo(tt.txMemo)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strings.Trim(strings.ToLower("0xAb5801a7D398351b8bE11C439e05C5b3259aec9B"), " "), memo.Address)
				assert.Equal(t, strings.Trim(strings.ToLower("1"), " "), memo.ChainID)
			}
		})
	}
}
