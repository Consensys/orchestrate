package entities

import "time"

type InternalData struct {
	OneTimeKey        bool          `json:"oneTimeKey,omitempty"`
	ChainID           string        `json:"chainID"`
	Priority          string        `json:"priority"`
	RetryInterval     time.Duration `json:"retryInterval"`
	GasPriceIncrement float64       `json:"gasPriceIncrement,omitempty" `
	GasPriceLimit     float64       `json:"gasPriceLimit,omitempty"`
	ParentJobUUID     string        `json:"parentJobUUID,omitempty"`
}
