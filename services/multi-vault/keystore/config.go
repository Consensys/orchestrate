package keystore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(secretPkeyViperKey, secretPkeyDefault)
	_ = viper.BindEnv(secretPkeyViperKey, secretPkeyEnv)
}

// InitFlags initialize flags
func InitFlags(f *pflag.FlagSet) {
	SecretPkeys(f)
}

var (
	secretPkeyFlag     = "secret-pkey"
	secretPkeyViperKey = "secret.pkeys"
	secretPkeyDefault  []string
	secretPkeyEnv      = "SECRET_PKEY"
)

// SecretPkeys register flag for Vault accounts
func SecretPkeys(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Private keys to pre-register in key store. Warning - Do not use in production. 
Environment variable: %q`, secretPkeyEnv)
	f.StringSlice(secretPkeyFlag, secretPkeyDefault, desc)
	_ = viper.BindPFlag(secretPkeyViperKey, f.Lookup(secretPkeyFlag))
}
