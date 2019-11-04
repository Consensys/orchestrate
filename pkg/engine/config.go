package engine

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func init() {
	viper.SetDefault(engineSlotsViperKey, engineSlotsDefault)
	_ = viper.BindEnv(engineSlotsViperKey, engineSlotsEnv)
}

// Config is engine configuration
type Config struct {
	Slots int64
}

func (c *Config) validate() error {
	if c.Slots <= 0 {
		return fmt.Errorf("at least one engine slot is required")
	}
	return nil
}

// Validate ensure configuration is valid
func (c *Config) Validate() error {
	if err := c.validate(); err != nil {
		return errors.ConfigError(err.Error()).SetComponent(component)
	}
	return nil
}

// NewConfig create new engine configuration
func NewConfig() Config {
	config := Config{}
	config.Slots = viper.GetInt64("engine.slots")
	return config
}

// InitFlags register flags for engine
func InitFlags(f *pflag.FlagSet) {
	Slots(f)
}

const (
	engineSlotsFlag     = "engine-slots"
	engineSlotsViperKey = "engine.slots"
	engineSlotsDefault  = uint(20)
	engineSlotsEnv      = "ENGINE_SLOTS"
)

// Slots register flag for Kafka server addresses
func Slots(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of messages that can be treated concurrently.
Environment variable: %q`, engineSlotsEnv)
	f.Uint(engineSlotsFlag, engineSlotsDefault, desc)
	_ = viper.BindPFlag(engineSlotsViperKey, f.Lookup(engineSlotsFlag))
}
