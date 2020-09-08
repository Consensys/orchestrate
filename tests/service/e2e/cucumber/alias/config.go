package alias

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(cucumberAliasesViperKey, cucumberAliasesDefault)
	_ = viper.BindEnv(cucumberAliasesViperKey, cucumberAliasesEnv)
}

// InitFlags register Aliases flags
func InitFlags(f *pflag.FlagSet) {
	Aliases(f)
}

var (
	cucumberAliasesFlag     = "cucumber-aliases"
	cucumberAliasesViperKey = "cucumber.aliases"
	cucumberAliasesDefault  = "{}"
	cucumberAliasesEnv      = "TEST_GLOBAL_DATA"
)

// Aliases register flag for aliases
func Aliases(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Aliases for cucumber test scenarios (e.g chain.primary:888)
Environment variable: %q`, cucumberAliasesEnv)
	f.String(cucumberAliasesFlag, cucumberAliasesDefault, desc)
	_ = viper.BindPFlag(cucumberAliasesViperKey, f.Lookup(cucumberAliasesFlag))
}
