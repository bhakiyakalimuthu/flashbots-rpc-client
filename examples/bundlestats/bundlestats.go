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

	// create user stats argument
	arg := []common.UserStatsArgs{{
		BlockNumber: blockNum,
	}}
	// create flashbots client
	c := client.NewFlashbotsClient("https://relay-goerli.flashbots.net")

	// call bundle
	res, err := c.UserStats(context.Background(), arg)
	if err != nil {
		fmt.Printf("send bundle failed %v", err)
		return
	}
	l.Info("sendBundle response", zap.Any("response", res))
}
