package listener

import (
	"time"
)

// Config configuration of a TxListener
type Config struct {
	EthClient struct {
		Retry struct {
			// We use an exponential backoff retry strategy when fetching from an Eth Client
			// See https://github.com/cenkalti/backoff/blob/master/exponential.go
			InitialInterval     time.Duration
			RandomizationFactor float64
			Multiplier          float64
			MaxInterval         time.Duration
			MaxElapsedTime      time.Duration
		}
	}

	BlockCursor struct {
		// How long to wait after failing to retrieve a new mined block
		Backoff time.Duration

		// Limit is a count of blocks that can be pre-fetched and buffered
		Limit uint64

		Tracker struct {
			// Depth under which a block is considered final
			Depth uint64
		}
	}

	TxListener struct {
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
func NewConfig() Config {
	c := Config{}
	c.BlockCursor.Backoff = time.Second
	c.BlockCursor.Limit = 20
	c.BlockCursor.Tracker.Depth = 0

	c.EthClient.Retry.InitialInterval = 500 * time.Millisecond
	c.EthClient.Retry.RandomizationFactor = 0.5
	c.EthClient.Retry.Multiplier = 1.5
	c.EthClient.Retry.MaxInterval = 10 * time.Second
	c.EthClient.Retry.MaxElapsedTime = 30 * time.Second
	c.TxListener.Return.Blocks = false
	c.TxListener.Return.Errors = false

	return c
}
