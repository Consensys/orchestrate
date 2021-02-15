package e2e

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(e2eDataViperKey, e2eDataDefault)
	_ = viper.BindEnv(e2eDataViperKey, e2eDataEnv)
}

// InitFlags register Aliases flags
func InitFlags(f *pflag.FlagSet) {
	TestData(f)
}

var (
	e2eDataFlag     = "e2e-data"
	e2eDataViperKey = "e2e.data"
	e2eDataDefault  = "{}"
	e2eDataEnv      = "TEST_GLOBAL_DATA"
)

// Aliases register flag for aliases
func TestData(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Aliases for cucumber test scenarios (e.g {{"nodes":{"besu":[{"URLs":["http://validator1:8545"]}})
Environment variable: %q`, e2eDataEnv)
	f.String(e2eDataFlag, e2eDataDefault, desc)
	_ = viper.BindPFlag(e2eDataViperKey, f.Lookup(e2eDataFlag))
}
