package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(TypeViperKey, typeEnv)
	viper.SetDefault(TypeViperKey, typeDefault)
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
Environment variable: %q`, []string{postgresOpt, memoryOpt}, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(TypeViperKey, f.Lookup(typeFlag))
}
