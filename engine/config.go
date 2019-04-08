package engine

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(engineSlotsViperKey, engineSlotsDefault)
	viper.BindEnv(engineSlotsViperKey, engineSlotsEnv)
	viper.SetDefault(enginePartitionsViperKey, enginePartitionsDefault)
	viper.BindEnv(enginePartitionsViperKey, enginePartitionsEnv)
	viper.SetDefault(engineTimeoutViperKey, engineTimeoutDefault)
	viper.BindEnv(engineTimeoutViperKey, engineTimeoutEnv)
}

// Config is engine configuration
type Config struct {
	Slots      int64
	Partitions int64
}

// Validate ensure configuration is valid
func (c *Config) Validate() error {
	if c.Slots <= 0 {
		return fmt.Errorf("At least one engine slot is required")
	}

	if c.Partitions <= 0 {
		return fmt.Errorf("At least one partition is required")
	}

	return nil
}

// NewConfig create new engine configuration
func NewConfig() Config {
	config := Config{}
	config.Slots = viper.GetInt64("engine.slots")
	config.Partitions = viper.GetInt64("engine.partitions")
	return config
}

// InitFlags register flags for engine
func InitFlags(f *pflag.FlagSet) {
	Slots(f)
	Partitions(f)
}

var (
	engineSlotsFlag     = "engine-slots"
	engineSlotsViperKey = "engine.slots"
	engineSlotsDefault  = uint(20)
	engineSlotsEnv      = "WORKER_SLOTS"
)

// Slots register flag for Kafka server addresses
func Slots(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of messages the engine can treat concurrently.
Environment variable: %q`, engineSlotsEnv)
	f.Uint(engineSlotsFlag, engineSlotsDefault, desc)
	viper.BindPFlag(engineSlotsViperKey, f.Lookup(engineSlotsFlag))
}

var (
	enginePartitionsFlag     = "engine-partitions"
	enginePartitionsViperKey = "engine.partitions"
	enginePartitionsDefault  = uint(50)
	enginePartitionsEnv      = "WORKER_PARTITIONS"
)

// Partitions register flag for Kafka server addresses
func Partitions(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Number of partitions spawned by engine to treat messages in parallel.
Environment variable: %q`, enginePartitionsEnv)
	f.Uint(enginePartitionsFlag, enginePartitionsDefault, desc)
	viper.BindPFlag(enginePartitionsViperKey, f.Lookup(enginePartitionsFlag))
}

var (
	engineTimeoutFlag     = "engine-timeout"
	engineTimeoutViperKey = "engine.timeout"
	engineTimeoutDefault  = 60 * time.Second
	engineTimeoutEnv      = "WORKER_TIMEOUT"
)

// Timeout register flag for Kafka server addresses
func Timeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum time for a message to be handled a message
Environment variable: %q`, engineTimeoutEnv)
	f.Duration(engineTimeoutFlag, engineTimeoutDefault, desc)
	viper.BindPFlag(engineTimeoutViperKey, f.Lookup(engineTimeoutFlag))
}
