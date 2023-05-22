package main

import (
	"fmt"

	"github.com/dan13ram/wpokt-backend/pocket"
)

func main() {
	{
		res, err := pocket.Height()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Height: %d\n", res.Height)
	}

	{
		res, err := pocket.AccountTxs()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("AccountTxs:")
		fmt.Printf("PageCount: %d\n", res.PageCount)
		fmt.Printf("TotalTxs: %d\n", res.TotalTxs)
		fmt.Println("Txs:")
		for _, tx := range res.Txs {
			fmt.Printf("[%d]\tHash: %s\n", tx.Index, tx.Hash)
			fmt.Printf("\tHeight: %d\n", tx.Height)
			fmt.Printf("\tFrom: %s\n", tx.StdTx.Msg.Value.FromAddress)
			fmt.Printf("\tTo: %s\n", tx.StdTx.Msg.Value.ToAddress)
			fmt.Printf("\tAmount: %s\n", tx.StdTx.Msg.Value.Amount)
			fmt.Printf("\tMemo: %s\n", tx.StdTx.Memo)
			fmt.Printf("\tFee: %s %s\n", tx.StdTx.Fee[0].Amount, tx.StdTx.Fee[0].Denom)

		}
	}
}
