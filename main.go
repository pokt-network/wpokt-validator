package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/dan13ram/wpokt-backend/pocket"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if len(os.Args) < 2 {
		log.Fatal("Please provide config file as parameter")
	}
	absConfigPath, _ := filepath.Abs(os.Args[1])

	app.InitConfig(absConfigPath)
	app.InitLogger()

	{
		res, err := pocket.Height()
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("Height: ", res.Height)
	}

	{
		res, err := pocket.AccountTxs()
		if err != nil {
			log.Error(err)
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

	// Gracefully shut down server
	gracefulStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go waitForExitSignals(gracefulStop, done)
	<-done
	log.Info("Server shutting down")
}

func waitForExitSignals(gracefulStop chan os.Signal, done chan bool) {
	sig := <-gracefulStop
	log.Debug("Got signal:", sig)
	log.Debug("Sending done signal to main")
	done <- true
}
