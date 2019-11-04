package cooldown

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(faucetCooldownViperKey, faucetCooldownDefault)
	_ = viper.BindEnv(faucetCooldownViperKey, faucetCooldownEnv)
}

const (
	faucetCooldownFlag     = "faucet-cooldown-time"
	faucetCooldownViperKey = "faucet.ctrl.cooldown.time"
	faucetCooldownDefault  = 60 * time.Second
	faucetCooldownEnv      = "FAUCET_COOLDOWN_TIME"
)

// FaucetCooldown register flag for Faucet Cooldown
func FaucetCooldown(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Faucet minimum time to wait before crediting an address again
Environment variable: %q`, faucetCooldownEnv)
	f.Duration(faucetCooldownFlag, faucetCooldownDefault, desc)
	_ = viper.BindPFlag(faucetCooldownViperKey, f.Lookup(faucetCooldownFlag))
}

// Config is Cooldown configuration object
type Config struct {
	// Cooldown Delay
	Delay time.Duration

	// Cooldown uses an underlying StripeMutex
	Stripes int
}

// NewConfig creates new configuration
func NewConfig() *Config {
	return &Config{
		Delay:   viper.GetDuration(faucetCooldownViperKey),
		Stripes: 100,
	}
}
