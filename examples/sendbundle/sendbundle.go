package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/client"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/util"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/zap"
)

func main() {
	l := common.NewLogger()

	// tx utility for creating rawTx
	txMgr := util.NewTxMgr(os.Getenv("INFURA_GOERLI"))
	rawTx, blockNum := txMgr.CreateTx(context.Background())

	// create send bundle argument
	arg := []common.SendBundleArgs{{
		Txs:         []string{hexutil.Encode(rawTx)},
		BlockNumber: blockNum,
	}}

	// create flashbots client
	c := client.NewFlashbotsClient("https://relay-goerli.flashbots.net")

	// call bundle
	res, err := c.SendBundle(context.Background(), arg)
	if err != nil {
		fmt.Printf("send bundle failed %v", err)
		return
	}
	l.Info("sendBundle response", zap.Any("response", res))
}
