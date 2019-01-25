package listener

import (
	"time"
)

// Config configuration of a TxListener
type Config struct {
	BlockListener struct {
		// How long to wait after failing to retrieve a new mined block
		Backoff time.Duration

		// WARNING: Retries are not implemented yet (TODO)
		Retry struct {
			// The total number of times to retry retrieving Receipts from an Ethereum client
			Max int
			// How long to wait for the client to settle between retries
			Backoff time.Duration
		}

		Return struct {
			// If enabled, any errors that occurred while listening for blocks are returned on
			// the Errors channel
			// If set to True you must drain the Errors channel
			Errors bool
		}
	}

	TxListener struct {
		// How many receipts can be retrieved in parallel goroutines
		MaxReceiptCount uint

		// WARNING: Retries are not implemented yet (TODO)
		Retry struct {
			// The total number of times to retry retrieving Receipts from an Ethereum client
			Max int
			// How long to wait for the client to settle between retries
			Backoff time.Duration
		}

		Return struct {
			// If enabled, all mined blocks are returned on the Blocks channel
			// If set to true you must drain the block channel
			Blocks bool

			// If enabled, any errors that occurred while listening for tx are returned on
			// the Errors channel
			// If set to True you must drain the Errors channel
			Errors bool
		}
	}
}

// NewConfig creates a new default config
func NewConfig() *Config {
	c := &Config{}
	c.BlockListener.Backoff = time.Second
	c.BlockListener.Retry.Max = 3
	c.BlockListener.Retry.Backoff = 2 * time.Second
	c.BlockListener.Return.Errors = false

	c.TxListener.MaxReceiptCount = 100
	c.TxListener.Retry.Max = 3
	c.TxListener.Retry.Backoff = 2 * time.Second
	c.TxListener.Return.Blocks = false
	c.TxListener.Return.Errors = false

	return c
}
