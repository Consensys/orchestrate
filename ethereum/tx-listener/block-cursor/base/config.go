package base

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(blockBackoffViperKey, blockBackoffDefault)
	_ = viper.BindEnv(blockBackoffViperKey, blockBackoffEnv)
	viper.SetDefault(blockLimitViperKey, blockLimitDefault)
	_ = viper.BindEnv(blockLimitViperKey, blockLimitEnv)
}

type Config struct {
	// How long to wait after failing to retrieve a new mined block
	Backoff time.Duration

	// Limit is a count of blocks that can be pre-fetched and buffered
	Limit uint64
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	BlockBackoff(f)
	BlockLimit(f)
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
	_ = viper.BindPFlag(blockBackoffViperKey, f.Lookup(blockBackoffFlag))
}

var (
	blockLimitFlag     = "listener-block-limit"
	blockLimitViperKey = "listener.block.limit"
	blockLimitDefault  = int64(40)
	blockLimitEnv      = "LISTENER_BLOCK_LIMIT"
)

// BlockLimit register flag for Listener Block limit
func BlockLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Limit number of blocks that can be prefetched while listening
Environment variable: %q`, blockLimitEnv)
	f.Int64(blockLimitFlag, blockLimitDefault, desc)
	_ = viper.BindPFlag(blockLimitViperKey, f.Lookup(blockLimitFlag))
}

func NewConfig() *Config {
	return &Config{
		Backoff: viper.GetDuration(blockBackoffViperKey),
		Limit:   uint64(viper.GetInt64(blockLimitViperKey)),
	}
}
