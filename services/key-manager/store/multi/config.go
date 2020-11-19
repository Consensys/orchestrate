package multi

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(TypeViperKey, typeDefault)
	_ = viper.BindEnv(TypeViperKey, typeEnv)
}

const (
	hashicorpVaultType = "hashicorp-vault"
	azureKeyVaultType  = "azure-key-vault"
	ukcVaultType       = "ukc-key-vault"
)

var availableTypes = []string{hashicorpVaultType, azureKeyVaultType, ukcVaultType}

const (
	typeFlag     = "key-manager-type"
	TypeViperKey = "key-manager.vault.type"
	typeDefault  = hashicorpVaultType
	typeEnv      = "KEY_MANAGER_TYPE"
)

// Type register flag for the Key manager to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Key Manager Vault (one of %q)
Environment variable: %q`, availableTypes, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(TypeViperKey, f.Lookup(typeFlag))
}

type Config struct {
	Type string
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Type: vipr.GetString(TypeViperKey),
	}
}

func Flags(f *pflag.FlagSet) {
	Type(f)
}
