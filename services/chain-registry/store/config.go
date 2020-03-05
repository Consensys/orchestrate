package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(TypeViperKey, typeEnv)
	viper.SetDefault(TypeViperKey, typeDefault)
	_ = viper.BindEnv(InitViperKey, initEnv)
	viper.SetDefault(InitViperKey, initDefault)
}

func Flags(f *pflag.FlagSet) {
	InitRegistry(f)
}

const (
	typeFlag     = "chain-registry-type"
	TypeViperKey = "chain-registry.type"
	typeDefault  = postgresOpt
	typeEnv      = "CHAIN_REGISTRY_TYPE"
)

// Type register flag for the Chain Registry to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Chain Registry (one of %q)
Environment variable: %q`, []string{postgresOpt}, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(TypeViperKey, f.Lookup(typeFlag))
}

var (
	initFlag     = "chain-registry-init"
	InitViperKey = "chain-registry.init"
	initDefault  []string
	initEnv      = "CHAIN_REGISTRY_INIT"
)

// Init register flag for the Chain Registry to define initialization state
func InitRegistry(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Initialize Chain Registry
Environment variable: %q`, initEnv)
	f.StringSlice(initFlag, initDefault, desc)
	_ = viper.BindPFlag(InitViperKey, f.Lookup(initFlag))
}
