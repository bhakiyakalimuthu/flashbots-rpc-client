package common

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

const (
	JSONRPCVersion = "2.0"
)

type JSONRPCMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *JSONError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func NewJSONRPCMessage(method string, params json.RawMessage) JSONRPCMessage {
	return JSONRPCMessage{
		ID:      genID(),
		Method:  method,
		Params:  params,
		Version: JSONRPCVersion,
	}
}

func genID() json.RawMessage {
	// id := atomic.AddUint32(&c.idCounter, 1)
	return strconv.AppendUint(nil, uint64(rand.Int()), 10)
}

type CallBundleArgs struct {
	Txs              []string `json:"txs"`              // Array[String], A list of signed transactions to execute in an atomic bundle
	BlockNumber      string   `json:"blockNumber"`      // String, a hex encoded block number for which this bundle is valid on
	StateBlockNumber string   `json:"stateBlockNumber"` // String, either a hex encoded number or a block tag for which state to base this simulation on. Can use "latest"
	Timestamp        *uint64  `json:"timestamp"`        // (Optional) Number, the timestamp to use for this bundle simulation, in seconds since the unix epoch
}

type SendBundleArgs struct {
	Txs               []string  // Array[String], A list of signed transactions to execute in an atomic bundle
	BlockNumber       string    // String, a hex encoded block number for which this bundle is valid on
	MinTimestamp      time.Time // (Optional) Number, the minimum timestamp for which this bundle is valid, in seconds since the unix epoch
	MaxTimestamp      time.Time // (Optional) Number, the maximum timestamp for which this bundle is valid, in seconds since the unix epoch
	RevertingTxHashes []string  // (Optional) Array[String], A list of tx hashes that are allowed to revert
}

type UserStatsArgs struct {
	BlockNumber string // String, a hex encoded recent block number, in order to prevent replay attacks. Must be within 20 blocks of the current chain tip.
}

type BundleStatsArgs struct {
	BundleHash  string `json:"bundleHash"`  // String, returned by the flashbots api when calling eth_sendBundle
	BlockNumber string `json:"blockNumber"` // String, the block number the bundle was targeting (hex encoded)
}

type SendPrivateTxArgs struct {
	Tx             string `json:"tx"`             // String, raw signed transaction
	MaxBlockNumber string `json:"maxBlockNumber"` // Hex-encoded number string, optional. Highest block number in which the transaction should be included.
	Preferences    struct {
		Fast bool `json:"fast"` // optional. "fast" left for backwards compatibility; may be removed in a future version
	} `json:"preferences"`
}

type CancelPrivateTxArgs struct {
	TxHash string `json:"txHash"` // String, transaction hash of private tx to be cancelled
}

type CallBundleResponse struct {
	Results           []TxSimulationResponse `json:"results"`
	CoinbaseDiff      string                 `json:"coinbaseDiff"`
	GasFees           *string                `json:"gasFees"`
	EthSentToCoinbase string                 `json:"ethSentToCoinbase"`
	BundleGasPrice    string                 `json:"bundleGasPrice"`
	TotalGasUsed      int64                  `json:"totalGasUsed"`
	StateBlockNumber  int64                  `json:"stateBlockNumber"`
	BundleHash        string                 `json:"bundleHash"`
}

type TxSimulationResponse struct {
	TxHash            string  `json:"txHash"`
	GasUsed           int64   `json:"gasUsed"`
	GasPrice          *string `json:"gasPrice"`
	GasFees           *string `json:"gasFees"`
	FromAddress       *string `json:"fromAddress"`
	ToAddress         *string `json:"toAddress"`
	CoinbaseDiff      *string `json:"coinbaseDiff"`
	EthSentToCoinbase *string `json:"ethSentToCoinbase"`
	Error             *string `json:"error"`
	Revert            *string `json:"revert"`
	Value             *string `json:"value"`
}

type SendBundleResponse struct {
	BundleHash string `json:"bundleHash"`
}

type BundleStatsResponse struct {
	IsHighPriority bool   `json:"isHighPriority"`
	IsSentToMiners bool   `json:"isSentToMiners"`
	IsSimulated    bool   `json:"isSimulated"`
	SentToMinersAt string `json:"sentToMinersAt,omitempty"`
	SimulatedAt    string `json:"simulatedAt"`
	SubmittedAt    string `json:"submittedAt"`
}

type UserStatsResponse struct {
	AllTimeGasSimulated  string `json:"all_time_gas_simulated"`
	AllTimeMinerPayments string `json:"all_time_miner_payments"`
	IsHighPriority       bool   `json:"is_high_priority"`
	Last1dGasSimulated   string `json:"last_1d_gas_simulated"`
	Last1dMinerPayments  string `json:"last_1d_miner_payments"`
	Last7dGasSimulated   string `json:"last_7d_gas_simulated"`
	Last7dMinerPayments  string `json:"last_7d_miner_payments"`
}

type UserStatsResponseV2 struct {
	AllTimeGasSimulated      string `json:"allTimeGasSimulated"`
	AllTimeValidatorPayments string `json:"allTimeValidatorPayments"`
	IsHighPriority           bool   `json:"isHighPriority"`
	Last1dGasSimulated       string `json:"last1DGasSimulated"`
	Last1dValidatorPayments  string `json:"last1dValidatorPayments"`
	Last7dGasSimulated       string `json:"last7DGasSimulated"`
	Last7dMinerPayments      string `json:"last7DMinerPayments"`
}
type BundleStatsResponseV2 struct {
	IsHighPriority         bool   `json:"isHighPriority"`
	SimulatedAt            string `json:"simulatedAt"`
	ReceivedAt             string `json:"receivedAt"`
	ConsideredByBuildersAt []struct {
		Pubkey    string `json:"pubkey"`
		Timestamp string `json:"timestamp"`
	} `json:"consideredByBuildersAt"`
	SealedByBuildersAt []struct {
		Pubkey    string `json:"pubkey"`
		Timestamp string `json:"timestamp"`
	} `json:"sealedByBuildersAt"`
}
