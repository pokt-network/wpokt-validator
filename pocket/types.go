package pocket

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

type HeightResponse struct {
	Height int64 `json:"height"`
}

type TxSignature struct {
	PubKey    bytes.HexBytes `json:"pub_key"`
	Signature bytes.HexBytes `json:"signature"`
}

type TxMsgValue struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
}

type TxMsg struct {
	Value TxMsgValue `json:"value"`
	Type  string     `json:"type"`
}

type TxFee struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type TxParams struct {
	Memo      string      `json:"memo"`
	Entropy   int64       `json:"entropy"`
	Fee       []*TxFee    `json:"fee"`
	Msg       TxMsg       `json:"msg"`
	Signature TxSignature `json:"signature"`
}

type ResultTx struct {
	Hash     bytes.HexBytes         `json:"hash"`
	Height   int64                  `json:"height"`
	Index    int64                  `json:"index"`
	TxResult abci.ResponseDeliverTx `json:"tx_result"`
	Tx       types.Tx               `json:"tx"`
	Proof    types.TxProof          `json:"proof,omitempty"`
	StdTx    TxParams               `json:"stdTx"`
}

type AccountTxsResponse struct {
	PageCount uint32      `json:"page_count"`
	TotalTxs  uint32      `json:"total_txs"`
	Txs       []*ResultTx `json:"txs"`
}
