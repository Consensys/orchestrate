package base

import (
	"github.com/spf13/pflag"
	cursor "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/block-cursor/base"
	tiptracker "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/tip-tracker/base"
)

// Config configuration of a TxListener
type Config struct {
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

	BlockCursor cursor.Config

	TipTracker tiptracker.Config
}

// NewConfig creates a new default config
func NewConfig() *Config {
	config := &Config{
		BlockCursor: *(cursor.NewConfig()),
		TipTracker:  *(tiptracker.NewConfig()),
	}
	config.TxListener.Return.Blocks = false
	config.TxListener.Return.Errors = false

	return config
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	cursor.InitFlags(f)
	tiptracker.InitFlags(f)
}
