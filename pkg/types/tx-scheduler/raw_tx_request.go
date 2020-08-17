package txschedulertypes

type RawTransactionRequest struct {
	ChainName string               `json:"chain" validate:"required" example:"myChain"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Params    RawTransactionParams `json:"params" validate:"required"`
}
type RawTransactionParams struct {
	Raw         string              `json:"raw" validate:"required,isHex" example:"0xfe378324abcde723..."`
	RetryPolicy IntervalRetryParams `json:"retryPolicy,omitempty"`
}
