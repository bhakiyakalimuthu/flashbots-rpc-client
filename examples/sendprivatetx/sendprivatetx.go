package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/fbclient"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/util"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/zap"
)

func main() {
	l := common.NewLogger()

	// tx utility for creating rawTx
	txMgr := util.NewTxMgr(os.Getenv("INFURA_GOERLI"))
	rawTx, blockNum := txMgr.CreateTx(context.Background())

	// create send private tx argument
	arg := []common.SendPrivateTxArgs{{
		Tx:             hexutil.Encode(rawTx),
		MaxBlockNumber: blockNum,
		Preferences:    nil,
	}}

	// create flashbots client
	c := fbclient.NewFlashbotsClient("https://relay-goerli.flashbots.net")

	// send private tx
	res, err := c.SendPrivateTransaction(context.Background(), arg)
	if err != nil {
		fmt.Printf("send bundle failed %v", err)
		return
	}
	l.Info("sendBundle response", zap.Any("response", res))
}
