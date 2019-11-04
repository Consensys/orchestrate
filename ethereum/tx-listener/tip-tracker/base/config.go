package base

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(depthViperKey, depthDefault)
	_ = viper.BindEnv(depthViperKey, depthEnv)
}

// Config is a configuration object for tracker
type Config struct {
	Depth uint64
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	Depth(f)
}

const (
	depthFlag     = "listener-tracker-depth"
	depthViperKey = "listener.tracker.depth"
	depthDefault  = int64(0)
	depthEnv      = "LISTENER_TRACKER_DEPTH"
)

// Depth register flag for Listener Tracker Depth
func Depth(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Depth at which we consider a block final (to avoid falling into a re-org)
Environment variable: %q`, depthEnv)
	f.Int64(depthFlag, depthDefault, desc)
	_ = viper.BindPFlag(depthViperKey, f.Lookup(depthFlag))
}

func NewConfig() *Config {
	return &Config{
		Depth: uint64(viper.GetInt64(depthViperKey)),
	}
}
