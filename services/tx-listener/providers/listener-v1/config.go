package listenerv1

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
)

func init() {
	viper.SetDefault(BlockBackoffViperKey, blockBackoffDefault)
	_ = viper.BindEnv(BlockBackoffViperKey, blockBackoffEnv)
	viper.SetDefault(BlockLimitViperKey, blockLimitDefault)
	_ = viper.BindEnv(BlockLimitViperKey, blockLimitEnv)
	viper.SetDefault(DepthViperKey, depthDefault)
	_ = viper.BindEnv(DepthViperKey, depthEnv)
	viper.SetDefault(StartDefaultViperKey, startDefaultDefault)
	_ = viper.BindEnv(StartDefaultViperKey, startDefaultEnv)
	viper.SetDefault(StartViperKey, startDefault)
	_ = viper.BindEnv(StartViperKey, startEnv)
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	BlockBackoff(f)
	BlockLimit(f)
	Depth(f)
	Start(f)
	StartDefault(f)
}

const (
	blockBackoffFlag     = "listener-block-backoff"
	BlockBackoffViperKey = "listener.block.backoff"
	blockBackoffDefault  = time.Second
	blockBackoffEnv      = "LISTENER_BLOCK_BACKOFF"
)

// BlockBackoff register flag for Listener Block backoff
func BlockBackoff(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Backoff time to wait before retrying after failing to find a mined block
Environment variable: %q`, blockBackoffEnv)
	f.Duration(blockBackoffFlag, blockBackoffDefault, desc)
	_ = viper.BindPFlag(BlockBackoffViperKey, f.Lookup(blockBackoffFlag))
}

const (
	blockLimitFlag     = "listener-block-limit"
	BlockLimitViperKey = "listener.block.limit"
	blockLimitDefault  = int64(40)
	blockLimitEnv      = "LISTENER_BLOCK_LIMIT"
)

// BlockLimit register flag for Listener Block limit
func BlockLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Limit number of blocks that can be prefetched while listening
Environment variable: %q`, blockLimitEnv)
	f.Int64(blockLimitFlag, blockLimitDefault, desc)
	_ = viper.BindPFlag(BlockLimitViperKey, f.Lookup(blockLimitFlag))
}

const (
	depthFlag     = "listener-tracker-depth"
	DepthViperKey = "listener.tracker.depth"
	depthDefault  = int64(0)
	depthEnv      = "LISTENER_TRACKER_DEPTH"
)

// Depth register flag for Listener Tracker Depth
func Depth(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Depth at which we consider a block final (to avoid falling into a re-org)
Environment variable: %q`, depthEnv)
	f.Int64(depthFlag, depthDefault, desc)
	_ = viper.BindPFlag(DepthViperKey, f.Lookup(depthFlag))
}

const (
	startDefaultFlag     = "listener-start-default"
	StartDefaultViperKey = "listener.start-default"
	startDefaultDefault  = "oldest"
	startDefaultEnv      = "LISTENER_START_DEFAULT"
)

// StartDefault register flag for Listener Start Default
func StartDefault(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Default block position listener should start listening from (one of 'latest', 'oldest', 'genesis')
Environment variable: %q`, startDefaultEnv)
	f.String(startDefaultFlag, startDefaultDefault, desc)
	_ = viper.BindPFlag(StartDefaultViperKey, f.Lookup(startDefaultFlag))
}

var (
	startFlag     = "listener-start"
	StartViperKey = "listener.start"
	startDefault  []string
	startEnv      = "LISTENER_START"
)

// Start register flag for Listener Start Position
func Start(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Position listener should start listening from (format <chainID>:<blockNumber>-<txIndex> or <chainID>:<blockNumber>) (e.g. 42:2348721-5 or 3:latest)
Environment variable: %q`, startEnv)
	f.StringSlice(startFlag, startDefault, desc)
	_ = viper.BindPFlag(StartViperKey, f.Lookup(startFlag))
}

// Position is an helpful type for storing a starting position
type Position struct {
	BlockNumber int64
	TxIndex     int64
}

type StartConfig struct {
	Positions map[string]Position
	Default   Position
}

type Config struct {
	// EthClient URLS
	URLs []string

	// How long to wait after failing to retrieve a new mined block
	Backoff time.Duration

	// Limit is a count of blocks that can be pre-fetched and buffered
	Limit uint64

	//
	Depth uint64

	// Start positions
	Start *StartConfig
}

// NewConfig create a new configuration
func NewStartConfig() (*StartConfig, error) {
	blockNumber, err := ParseBlock(viper.GetString(StartDefaultViperKey))
	if err != nil {
		return nil, err
	}

	conf := &StartConfig{}
	conf.Default = Position{BlockNumber: blockNumber}

	conf.Positions = make(map[string]Position)
	for _, position := range viper.GetStringSlice(StartViperKey) {
		chain, pos, err := ParsePosition(position)
		if err != nil {
			return nil, err
		}
		conf.Positions[chain] = *pos
	}

	return conf, nil
}

func NewConfig() (*Config, error) {
	start, err := NewStartConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		URLs:    viper.GetStringSlice(rpc.URLViperKey),
		Backoff: viper.GetDuration(BlockBackoffViperKey),
		Limit:   uint64(viper.GetInt64(BlockLimitViperKey)),
		Depth:   uint64(viper.GetInt64(DepthViperKey)),
		Start:   start,
	}, nil
}
