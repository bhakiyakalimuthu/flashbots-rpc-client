package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/client"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/util"
	"go.uber.org/zap"
)

func main() {
	l := common.NewLogger()

	// tx utility for creating rawTx
	txMgr := util.NewTxMgr(os.Getenv("INFURA_GOERLI"))
	_, blockNum := txMgr.CreateTx(context.Background())

	// create bundle stats argument
	arg := []common.BundleStatsArgs{{
		BundleHash:  "0x05e440b106aefe2b7375f08fab4dbb9554baa384a8df0315ffe9627f3104bea6",
		BlockNumber: blockNum,
	}}

	// create flashbots client
	c := client.NewFlashbotsClient("https://relay-goerli.flashbots.net")

	// call bundle
	res, err := c.BundleStats(context.Background(), arg)
	if err != nil {
		fmt.Printf("send bundle failed %v", err)
		return
	}
	l.Info("sendBundle response", zap.Any("response", res))
}
