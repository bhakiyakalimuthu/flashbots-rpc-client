package common

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
)

type Signer interface {
	SignPayload(payload json.RawMessage) (*string, error)
}
type signer struct {
	logger     *zap.Logger
	privateKey string
}

func NewSigner() *signer {
	return &signer{
		logger:     NewLogger(),
		privateKey: os.Getenv("SIGNER_PRIVATE_KEY"),
	}
}

func (s *signer) SignPayload(payload json.RawMessage) (*string, error) {
	key, err := crypto.HexToECDSA(s.privateKey)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	hashedBody := crypto.Keccak256Hash(body).Hex()
	sig, err := crypto.Sign(accounts.TextHash([]byte(hashedBody)), key)
	if err != nil {
		return nil, err
	}

	signature := crypto.PubkeyToAddress(key.PublicKey).Hex() + ":" + hexutil.Encode(sig)
	return &signature, nil
}
