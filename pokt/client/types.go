package client

type HeightResponse struct {
	Height int64 `json:"height"`
}

type SubmitRawTxResponse struct {
	TransactionHash string `json:"txhash"`
}

type Value struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
}

type Msg struct {
	Value Value  `json:"value"`
	Type  string `json:"type"`
}

type StdTx struct {
	Memo string `json:"memo"`
	Msg  Msg    `json:"msg"`
}

type TxResponse struct {
	Hash     string   `json:"hash"`
	Height   int64    `json:"height"`
	Index    int64    `json:"index"`
	StdTx    StdTx    `json:"stdTx"`
	TxResult TxResult `json:"tx_result"`
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

type TxResult struct {
	Code        int64  `json:"code"`
	Codespace   string `json:"codespace"`
	Data        string `json:"data"`
	Events      string `json:"events"`
	Info        string `json:"info"`
	Log         string `json:"log"`
	MessageType string `json:"message_type"`
	Recipient   string `json:"recipient"`
	Signer      string `json:"signer"`
}
