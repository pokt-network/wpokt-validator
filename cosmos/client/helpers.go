package client

import (
	"fmt"
	"time"

	rpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/dan13ram/wpokt-validator/cosmos/util"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"context"
)

func getBlocksForTxResults(node CosmosHTTPClient, resTxs []*rpctypes.ResultTx) (map[int64]*rpctypes.ResultBlock, error) {
	resBlocks := make(map[int64]*rpctypes.ResultBlock)

	for _, resTx := range resTxs {
		resTx := resTx

		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(context.Background(), &resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func formatTxResults(bech32Prefix string, resTxs []*rpctypes.ResultTx, resBlocks map[int64]*rpctypes.ResultBlock) ([]*sdk.TxResponse, error) {
	var err error
	out := make([]*sdk.TxResponse, len(resTxs))
	for i := range resTxs {
		out[i], err = mkTxResult(bech32Prefix, resTxs[i], resBlocks[resTxs[i].Height])
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

var utilNewTxDecoder = util.NewTxDecoder

func mkTxResult(bech32Prefix string, resTx *rpctypes.ResultTx, resBlock *rpctypes.ResultBlock) (*sdk.TxResponse, error) {
	txDecoder := utilNewTxDecoder(bech32Prefix)
	txb, err := txDecoder(resTx.Tx)
	if err != nil {
		return nil, fmt.Errorf("decoding tx: %w", err)
	}
	p, ok := txb.(AnyTx)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}

	any := p.AsAny()

	return sdk.NewResponseResultTx(resTx, any, resBlock.Block.Time.Format(time.RFC3339)), nil
}

// Deprecated: this interface is used only internally for scenario we are
// deprecating (StdTxConfig support)
type AnyTx interface {
	sdk.Tx
	AsAny() *codectypes.Any
}
