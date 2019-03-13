package listener

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("listener.block.limit", 20)
	viper.SetDefault("listener.block.backoff", time.Second)
	viper.SetDefault("listener.tracker.depth", 0)
	viper.SetDefault("ethclient.retry.initinterval", 500*time.Millisecond)
	viper.SetDefault("ethclient.retry.randomfactor", 0.5)
	viper.SetDefault("ethclient.retry.multiplier", 1.5)
	viper.SetDefault("ethclient.retry.maxinterval", 10*time.Second)
	viper.SetDefault("ethclient.retry.maxelapsedtime", 60*time.Second)
}

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
}

// NewConfig creates a new default config
func NewConfig() Config {
	return Config{}
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	BlockBackoff(f)
	BlockLimit(f)
	TrackerDepth(f)
}

var (
	blockBackoffFlag     = "listener-block-backoff"
	blockBackoffViperKey = "listener.block.backoff"
	blockBackoffDefault  = time.Second
	blockBackoffEnv      = "LISTENER_BLOCK_BACKOFF"
)

// BlockBackoff register flag for Listener Block backoff
func BlockBackoff(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Backoff time to wait before retrying after failing to find a mined block
Environment variable: %q`, blockBackoffEnv)
	f.Duration(blockBackoffFlag, blockBackoffDefault, desc)
	viper.SetDefault(blockBackoffViperKey, blockBackoffDefault)
	viper.BindPFlag(blockBackoffViperKey, f.Lookup(blockBackoffFlag))
	viper.BindEnv(blockBackoffViperKey, blockBackoffEnv)
}

var (
	blockLimitFlag     = "listener-block-limit"
	blockLimitViperKey = "listener.block.limit"
	blockLimitDefault  = int64(40)
	blockLimitEnv      = "LISTENER_BLOCK_LIMIT"
)

// BlockLimit register flag for Listener Block limit
func BlockLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Limit number of block that can be prefetched while listening
Environment variable: %q`, blockLimitEnv)
	f.Int64(blockLimitFlag, blockLimitDefault, desc)
	viper.SetDefault(blockLimitViperKey, blockLimitDefault)
	viper.BindPFlag(blockLimitViperKey, f.Lookup(blockLimitFlag))
	viper.BindEnv(blockLimitViperKey, blockLimitEnv)
}

var (
	trackerDepthFlag     = "listener-tracker-depth"
	trackerDepthViperKey = "listener.tracker.depth"
	trackerDepthDefault  = int64(5)
	trackerDepthEnv      = "LISTENER_TRACKER_DEPTH"
)

// TrackerDepth register flag for Listener Tracker Depth
func TrackerDepth(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Depth at which we consider a block final (to avoid falling into a re-org)
Environment variable: %q`, trackerDepthEnv)
	f.Int64(trackerDepthFlag, trackerDepthDefault, desc)
	viper.SetDefault(trackerDepthViperKey, trackerDepthDefault)
	viper.BindPFlag(trackerDepthViperKey, f.Lookup(trackerDepthFlag))
	viper.BindEnv(trackerDepthViperKey, trackerDepthEnv)
}
