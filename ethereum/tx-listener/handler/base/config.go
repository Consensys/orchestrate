package base

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(startDefaultViperKey, startDefaultDefault)
	_ = viper.BindEnv(startDefaultViperKey, startDefaultEnv)
	viper.SetDefault(startViperKey, startDefault)
	_ = viper.BindEnv(startViperKey, startEnv)
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	StartDefault(f)
	Start(f)
}

const (
	startDefaultFlag     = "listener-start-default"
	startDefaultViperKey = "listener.start-default"
	startDefaultDefault  = "oldest"
	startDefaultEnv      = "LISTENER_START_DEFAULT"
)

// StartDefault register flag for Listener Start Default
func StartDefault(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Default block position listener should start listening from (one of 'latest', 'oldest', 'genesis')
Environment variable: %q`, startDefaultEnv)
	f.String(startDefaultFlag, startDefaultDefault, desc)
	_ = viper.BindPFlag(startDefaultViperKey, f.Lookup(startDefaultFlag))
}

var (
	startFlag     = "listener-start"
	startViperKey = "listener.start"
	startDefault  []string
	startEnv      = "LISTENER_START"
)

// Start register flag for Listener Start Position
func Start(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Position listener should start listening from (format <chainID>:<blockNumber>-<txIndex> or <chainID>:<blockNumber>) (e.g. 42:2348721-5 or 3:latest)
Environment variable: %q`, startEnv)
	f.StringSlice(startFlag, startDefault, desc)
	_ = viper.BindPFlag(startViperKey, f.Lookup(startFlag))
}

// Position is an helpful type for storing a starting position
type Position struct {
	BlockNumber int64
	TxIndex     int64
}

type Config struct {
	Start struct {
		Positions map[string]Position
		Default   Position
	}
}

// NewConfig create a new configuration
func NewConfig() (*Config, error) {
	blockNumber, err := ParseBlock(viper.GetString(startDefaultViperKey))
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	conf.Start.Default = Position{BlockNumber: blockNumber}

	conf.Start.Positions = make(map[string]Position)
	for _, position := range viper.GetStringSlice(startViperKey) {
		chain, pos, err := ParsePosition(position)
		if err != nil {
			return nil, err
		}
		conf.Start.Positions[chain] = *pos
	}

	return conf, nil
}
