package pocket

type HeightResponse struct {
	Height int64 `json:"height"`
}

type SubmitRawTxResponse struct {
	TransactionHash string `json:"txhash"`
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

type StdTxParams struct {
	Memo string `json:"memo"`
	Msg  TxMsg  `json:"msg"`
}

type TxResponse struct {
	Hash   string      `json:"hash"`
	Height int64       `json:"height"`
	Index  int64       `json:"index"`
	StdTx  StdTxParams `json:"stdTx"`
}

type AccountTxsResponse struct {
	PageCount uint32        `json:"page_count"`
	TotalTxs  uint32        `json:"total_txs"`
	Txs       []*TxResponse `json:"txs"`
}

type Header struct {
	ChainID string `json:"chain_id"`
	Height  string `json:"height"`
}

type Block struct {
	Header Header `json:"header"`
}

type BlockResponse struct {
	Block Block `json:"block"`
}
