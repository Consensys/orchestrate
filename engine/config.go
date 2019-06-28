package engine

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
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
func (c *Config) Validate() *err.Error {
	return errors.ConfigError(c.validate()).SetComponent(component)
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

var (
	engineSlotsFlag     = "engine-slots"
	engineSlotsViperKey = "engine.slots"
	engineSlotsDefault  = uint(20)
	engineSlotsEnv      = "ENGINE_SLOTS"
)

// Slots register flag for Kafka server addresses
func Slots(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of messages the engine can treat concurrently.
Environment variable: %q`, engineSlotsEnv)
	f.Uint(engineSlotsFlag, engineSlotsDefault, desc)
	_ = viper.BindPFlag(engineSlotsViperKey, f.Lookup(engineSlotsFlag))
}
