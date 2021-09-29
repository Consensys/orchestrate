package engine

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/consensys/orchestrate/pkg/errors"
)

func init() {
	viper.SetDefault(SlotsViperKey, slotsDefault)
	_ = viper.BindEnv(SlotsViperKey, slotsEnv)
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
	config.Slots = viper.GetInt64(SlotsViperKey)
	return config
}

// InitFlags register flags for engine
func Flags(f *pflag.FlagSet) {
	Slots(f)
}

const (
	slotsFlag     = "engine-slots"
	SlotsViperKey = "engine.slots"
	slotsDefault  = uint(20)
	slotsEnv      = "ENGINE_SLOTS"
)

// Slots register flag for Kafka server addresses
func Slots(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of messages that can be treated concurrently.
Environment variable: %q`, slotsEnv)
	f.Uint(slotsFlag, slotsDefault, desc)
	_ = viper.BindPFlag(SlotsViperKey, f.Lookup(slotsFlag))
}
