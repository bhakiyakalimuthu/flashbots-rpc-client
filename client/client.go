package client

import (
	"context"
	"encoding/json"

	"github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"go.uber.org/zap"
)

const (
	// V1 methods
	_CallBundle = "eth_callBundle"
	_SendBundle = "eth_sendBundle"
	//_CancelBundle    = "eth_cancelBundle"
	_UserStats       = "flashbots_getUserStats"
	_BundleStats     = "flashbots_getBundleStats"
	_SendPrivateTx   = "eth_sendPrivateTransaction"
	_CancelPrivateTx = "eth_cancelPrivateTransaction"

	// V2 methods
	//_UserStatsV2   = "flashbots_getUserStatsV2"
	//_BundleStatsV2 = "flashbots_getBundleStatsV2"
)

type FlashbotsClient struct {
	logger     *zap.Logger
	httpClient *HttpClient
}

func NewFlashbotsClient(url string) *FlashbotsClient {
	logger := common.NewLogger()
	httpClient, err := DialHttpClient(url)
	if err != nil {
		logger.Fatal("failed to dial http client", zap.Error(err))
	}

	return &FlashbotsClient{
		logger:     logger,
		httpClient: httpClient,
	}
}

func NewFlashbotsClientWithSigner(url, signerKey string) *FlashbotsClient {
	logger := common.NewLogger()
	httpClient, err := DialHttpClientWithSingerKey(url, signerKey)
	if err != nil {
		logger.Fatal("failed to dial http client", zap.Error(err))
	}

	return &FlashbotsClient{
		logger:     logger,
		httpClient: httpClient,
	}
}

func (fbc *FlashbotsClient) CallBundle(ctx context.Context, arg interface{}) (*common.CallBundleResponse, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPCMessage(_CallBundle, msg)
	res, err := fbc.httpClient.CallContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var callBundleResponse *common.CallBundleResponse
	if err = json.Unmarshal(res.Result, &callBundleResponse); err != nil {
		return nil, err
	}
	return callBundleResponse, nil
}

func (fbc *FlashbotsClient) BundleStats(ctx context.Context, arg interface{}) (*common.BundleStatsResponse, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPCMessage(_BundleStats, msg)
	res, err := fbc.httpClient.CallContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var bundleStatsResponse *common.BundleStatsResponse
	if err = json.Unmarshal(res.Result, &bundleStatsResponse); err != nil {
		return nil, err
	}
	return bundleStatsResponse, nil
}

func (fbc *FlashbotsClient) UserStats(ctx context.Context, arg interface{}) (*common.UserStatsResponse, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPCMessage(_UserStats, msg)
	res, err := fbc.httpClient.CallContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var userStatsResponse *common.UserStatsResponse
	if err = json.Unmarshal(res.Result, &userStatsResponse); err != nil {
		return nil, err
	}
	return userStatsResponse, nil
}

func (fbc *FlashbotsClient) SendBundle(ctx context.Context, arg interface{}) (*common.SendBundleResponse, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPCMessage(_SendBundle, msg)
	res, err := fbc.httpClient.CallContext(ctx, request)
	if err != nil {
		return nil, err
	}
	var sendBundleResponse *common.SendBundleResponse
	if err = json.Unmarshal(res.Result, &sendBundleResponse); err != nil {
		return nil, err
	}
	return sendBundleResponse, nil
}

func (fbc *FlashbotsClient) SendPrivateTransaction(ctx context.Context, arg interface{}) (*common.SendPrivateTransactionResponse, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPCMessage(_SendPrivateTx, msg)
	res, err := fbc.httpClient.CallContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	var txHash string
	if err = json.Unmarshal(res.Result, &txHash); err != nil {
		return nil, err
	}
	return &common.SendPrivateTransactionResponse{TxHash: txHash}, nil
}

func (fbc *FlashbotsClient) CancelPrivateTransaction(ctx context.Context, arg interface{}) (*common.CancelPrivateTransactionResponse, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		fbc.logger.Error("failed to marshal param", zap.Error(err))
		return nil, err
	}
	msg := json.RawMessage(b)
	request := common.NewJSONRPCMessage(_CancelPrivateTx, msg)
	res, err := fbc.httpClient.CallContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	var isCancelled bool
	if err = json.Unmarshal(res.Result, &isCancelled); err != nil {
		return nil, err
	}
	return &common.CancelPrivateTransactionResponse{IsCancelled: isCancelled}, nil
}
