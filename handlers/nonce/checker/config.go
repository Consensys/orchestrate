package noncechecker

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(MaxRecoveryViperKey, maxRecoveryDefault)
	_ = viper.BindEnv(MaxRecoveryViperKey, maxRecoveryEnv)
}

// Register Redis flags
func Flags(f *pflag.FlagSet) {
	MaxRecovery(f)
}

const (
	maxRecoveryFlag     = "checker-max-recovery"
	MaxRecoveryViperKey = "checker.max.recovery"
	maxRecoveryDefault  = 5
	maxRecoveryEnv      = "NONCE_CHECKER_MAX_RECOVERY"
)

// MaxRecovery register a flag for Redis server MaxRecovery
func MaxRecovery(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of times tx-sender tries to recover an envelope with invalid nonce before outputing an error
Environment variable: %q`, maxRecoveryEnv)
	f.Int(maxRecoveryFlag, maxRecoveryDefault, desc)
	_ = viper.BindPFlag(MaxRecoveryViperKey, f.Lookup(maxRecoveryFlag))
}

type Configuration struct {
	MaxRecovery int
}

func NewConfig() *Configuration {
	return &Configuration{
		MaxRecovery: viper.GetInt(MaxRecoveryViperKey),
	}
}
