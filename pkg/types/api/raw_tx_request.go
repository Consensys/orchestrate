package api

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RawTransactionRequest struct {
	ChainName string               `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Params    RawTransactionParams `json:"params" validate:"required"`
}
type RawTransactionParams struct {
	Raw         hexutil.Bytes       `json:"raw" validate:"required" example:"0xfe378324abcde723" swaggertype:"string"`
	RetryPolicy IntervalRetryParams `json:"retryPolicy,omitempty"`
}
