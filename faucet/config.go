package faucet

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(typeViperKey, typeEnv)
	viper.SetDefault(typeViperKey, typeDefault)
}

var (
	typeFlag     = "faucet"
	typeViperKey = "faucet.type"
	typeDefault  = "sarama"
	typeEnv      = "FAUCET"
)

// Type register flag for Faucet Cooldown
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Faucet (one of %q)
Environment variable: %q`, []string{"mock", "sarama"}, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	viper.BindPFlag(typeViperKey, f.Lookup(typeFlag))
}

// Config is Cooldown configuration object
type Config struct {
	// Cooldown Delay
	Delay time.Duration

	// Cooldown uses an underlying SripeMutext
	Stripes int
}

// NewConfig creates new configuration
func NewConfig() *Config {
	return &Config{
		Delay:   viper.GetDuration(typeViperKey),
		Stripes: 100,
	}
}
