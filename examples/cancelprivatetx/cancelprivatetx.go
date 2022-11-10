package main

import (
	"context"
	"fmt"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/client"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"go.uber.org/zap"
)

func main() {
	l := common.NewLogger()

	// create send private tx argument
	arg := []common.CancelPrivateTxArgs{{
		TxHash: "0x68f70d7d939d9efae9ca31e8a96dfd074175da483ddcebd71dc2d2ba04f2861b",
	}}

	// create flashbots client
	c := client.NewFlashbotsClient("https://relay-goerli.flashbots.net")

	// cancel private tx
	res, err := c.CancelPrivateTransaction(context.Background(), arg)
	if err != nil {
		fmt.Printf("send bundle failed %v", err)
		return
	}
	l.Info("sendBundle response", zap.Any("response", res))
}
