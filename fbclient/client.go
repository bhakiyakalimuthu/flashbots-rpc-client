package fbclient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"github.com/bhakiyakalimuthu/flashbots-rpc-client/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/zap"
)

const (
	// V1 methods
	_CallBundle = "eth_callBundle"
	//_SendBundle      = "eth_sendBundle"
	//_CancelBundle    = "eth_cancelBundle"
	//_UserStats       = "flashbots_getUserStats"
	//_BundleStats     = "flashbots_getBundleStats"
	//_SendPrivateTx   = "eth_sendPrivateTransaction"
	//_CancelPrivateTx = "eth_cancelPrivateTransaction"

	// V2 methods
	//_UserStatsV2   = "flashbots_getUserStatsV2"
	//_BundleStatsV2 = "flashbots_getBundleStatsV2"
)

type flashbotsClient struct {
	logger     *zap.Logger
	signer     common.Signer
	httpClient *rpc.HttpClient
}

func NewFlashbotsClient(url string) *flashbotsClient {
	l := common.NewLogger()
	httpClient, err := rpc.DialHttpClient(url)
	if err != nil {
		l.Fatal("failed to dial http client", zap.Error(err))
	}
	return &flashbotsClient{
		logger:     l,
		signer:     common.NewSigner(),
		httpClient: httpClient,
	}
}

func (fbc *flashbotsClient) CallBundle(ctx context.Context, rawTx []byte, blockNum string) (_CallBundleResponse *common.CallBundleResponse, err error) {
	now := uint64(time.Now().Unix())
	params := []common.CallBundleArgs{{
		Txs:              []string{hexutil.Encode(rawTx)},
		BlockNumber:      blockNum,
		StateBlockNumber: "latest",
		Timestamp:        &now,
	}}
	b, err := json.Marshal(params)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPC(_CallBundle, msg)
	payload, err := json.Marshal(request)
	if err != nil {
		fbc.logger.Error("failed to marshal request", zap.String("method", _CallBundle), zap.Error(err))
		return nil, err
	}
	signature, err := fbc.signer.SignPayload(payload)
	if err != nil {
		return nil, err
	}
	res, err := fbc.httpClient.CallContext(ctx, "eth_callBundle", *signature)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res.Result, &_CallBundleResponse); err != nil {
		return nil, err
	}
	return
}
