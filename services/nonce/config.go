package nonce

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(typeViperKey, typeDefault)
	_ = viper.BindEnv(typeViperKey, typeEnv)
}

const (
	typeFlag     = "nonce-manager-type"
	typeViperKey = "nonce.manager.type"
	typeDefault  = "redis"
	typeEnv      = "NONCE_MANAGER_TYPE"
)

// Type register flag for Nonce Cooldown
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Nonce (one of %q)
Environment variable: %q`, []string{"in-memory", "redis"}, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(typeViperKey, f.Lookup(typeFlag))
}
