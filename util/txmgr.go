package util

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"

	common2 "github.com/bhakiyakalimuthu/flashbots-rpc-client/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// TxMgr Transaction manager is used as test utility to create transaction
type TxMgr interface {
	CreateTx(ctx context.Context) ([]byte, string)
}

type txMgr struct {
	logger     *zap.Logger
	url        string
	privateKey string
}

func NewTxMgr(url string) TxMgr {
	return &txMgr{
		logger:     common2.NewLogger(),
		url:        url,
		privateKey: os.Getenv("SIGNER_PRIVATE_KEY"),
	}
}

func (t *txMgr) CreateTx(ctx context.Context) ([]byte, string) {
	client, err := ethclient.DialContext(ctx, t.url)
	if err != nil {
		t.logger.Fatal("failed to dial ethClient", zap.Error(err))
	}
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		t.logger.Fatal("failed to get block number", zap.Error(err))
	}
	t.logger.Info("show current block", zap.Uint64("currentBlock", currentBlock))

	// transaction addresses
	fromAddress := common.HexToAddress(os.Getenv("WALLET_1"))
	toAddress := common.HexToAddress(os.Getenv("WALLET_2"))

	// Nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		t.logger.Error("failed to get pending nonce", zap.Error(err))
		return nil, ""
	}
	t.logger.Info("transaction nonce", zap.Uint64("nonce", nonce))
	//// Transfer Amount
	//amount := big.NewInt(1000000000000000000) // in wei (1 eth)

	// Gas
	gasLimit := uint64(126000) // in units
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		t.logger.Error("failed to get gas price", zap.Error(err))
		return nil, ""
	}
	t.logger.Info("gasPrice and gasLimit", zap.Uint64("gasLimit", gasLimit), zap.Any("gasPrice", gasPrice))

	// ChainID to differentiate between mainnet & test networks
	id, err := client.NetworkID(ctx)
	if err != nil {
		t.logger.Error("failed to get chain id", zap.Error(err))
		return nil, ""
	}

	t.logger.Info("chainID", zap.Any("id", id))
	// balance

	bal, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		t.logger.Error("failed to get 2-balance", zap.Error(err))
		return nil, ""
	}
	t.logger.Info("2-balance", zap.Any("balanceInWei", bal))
	// 1 ether = 10^18
	fBal := new(big.Float)
	fBal.SetString(bal.String())

	balance := new(big.Float).Quo(fBal, big.NewFloat(math.Pow10(10)))
	t.logger.Info("2-balance", zap.Any("balanceInEther", balance))
	// transfer amount
	amount := new(big.Int)
	amount.SetString("50000000", 10) // 0.05 tokens
	// amount.SetString("1000000000000000000000", 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
	gasFeeCap, gasTipCap := big.NewInt(38694000460), big.NewInt(3869400046)
	//gasTipCap := new(big.Int)
	//gasTipCap.SetString("30", 10) // 0.05 tokens
	//
	//gasFeeCap := new(big.Int)
	//gasFeeCap.SetString("50", 10) // 0.05 tokens
	// create data
	//transferFnSignature := []byte("transfer(address,uint256)")
	//hash := sha3.New256()
	//hash.Write(transferFnSignature)
	//methodID := hash.Sum(nil)[:4]
	//fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb
	//
	//paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	//fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d
	//
	//var data []byte
	//data = append(data, methodID...)
	//data = append(data, paddedAddress...)
	//data = append(data, paddedAmount...)
	dynamicFeeTx := &types.DynamicFeeTx{
		ChainID:   id,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddress,
	}
	tx := types.NewTx(dynamicFeeTx)
	// tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	key, err := crypto.HexToECDSA(strings.TrimPrefix(t.privateKey, "0x"))
	if err != nil {
		t.logger.Error("Error creating tx signing key", zap.Error(err))
	}

	// Sign the transaction using private key
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(id), key)
	if err != nil {
		t.logger.Error("failed to sign tx", zap.Error(err))
		return nil, ""
	}
	t.logger.Info("signed tx", zap.String("tx", signedTx.Hash().Hex()))
	blockNumHex := fmt.Sprintf("0x%x", currentBlock+25)
	t.logger.Info("Transaction hash", zap.String("tx", signedTx.Hash().Hex()))
	t.logger.Info("block num hex ", zap.String("blockNumHex", blockNumHex))

	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.logger.Panic("failed to marshal tx", zap.Error(err))
	}
	return rawTx, blockNumHex
}
