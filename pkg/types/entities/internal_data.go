package entities

import "time"

type InternalData struct {
	OneTimeKey        bool          `json:"oneTimeKey,omitempty"`
	HasBeenRetried    bool          `json:"hasBeenRetried,omitempty"`
	ChainID           string        `json:"chainID"`
	Priority          string        `json:"priority"`
	ParentJobUUID     string        `json:"parentJobUUID,omitempty"`
	GasPriceIncrement float64       `json:"gasPriceIncrement,omitempty"`
	GasPriceLimit     float64       `json:"gasPriceLimit,omitempty"`
	RetryInterval     time.Duration `json:"retryInterval"`
	ExpectedNonce     string        `json:"expectedNonce,omitempty"` // Using string because 0 is a valid
}
