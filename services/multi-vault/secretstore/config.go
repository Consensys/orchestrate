package secretstore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(secretStoreViperKey, secretStoreDefault)
	_ = viper.BindEnv(secretStoreViperKey, secretStoreEnv)
}

const (
	secretStoreFlag     = "secret-store"
	secretStoreViperKey = "secret.store"
	secretStoreDefault  = memoryOpt
	secretStoreEnv      = "SECRET_STORE"
)

// InitFlags initialize flags
func InitFlags(f *pflag.FlagSet) {
	SecFlag(f)
}

// SecFlag register flag for Vault accounts
func SecFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of secret store for private keys (one of %q %q)
Environment variable: %q`, memoryOpt, hashicorpOpt, secretStoreEnv)
	f.String(secretStoreFlag, secretStoreDefault, desc)
	_ = viper.BindPFlag(secretStoreViperKey, f.Lookup(secretStoreFlag))
}
