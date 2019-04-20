package cooldown

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(faucetCooldownViperKey, faucetCooldownEnv)
	viper.SetDefault(faucetCooldownViperKey, faucetCooldownDefault)
}

var (
	faucetCooldownFlag     = "faucet-cooldown"
	faucetCooldownViperKey = "faucet.ctrl.cooldown"
	faucetCooldownDefault  = 60 * time.Second
	faucetCooldownEnv      = "FAUCET_COOLDOWN_TIME"
)

// FaucetCooldown register flag for Faucet Cooldown
func FaucetCooldown(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Faucet minimum to wait before crediting an address again
Environment variable: %q`, faucetCooldownEnv)
	f.Duration(faucetCooldownFlag, faucetCooldownDefault, desc)
	viper.BindPFlag(faucetCooldownViperKey, f.Lookup(faucetCooldownFlag))
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
		Delay:   viper.GetDuration(faucetCooldownViperKey),
		Stripes: 100,
	}
}
