package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	// create call bundle argument
	now := uint64(time.Now().Unix())
	arg := []common.CallBundleArgs{{
		Txs:              []string{hexutil.Encode(rawTx)},
		BlockNumber:      blockNum,
		StateBlockNumber: "latest",
		Timestamp:        &now,
	}}

	// create flashbots client
	c := fbclient.NewFlashbotsClient("https://relay-goerli.flashbots.net")

	// call bundle
	res, err := c.CallBundle(context.Background(), arg)
	if err != nil {
		fmt.Printf("call bundle failed %v", err)
		return
	}
	l.Info("callBundle response", zap.Any("response", res))
}
