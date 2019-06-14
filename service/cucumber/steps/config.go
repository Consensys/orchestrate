package steps

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(cucumberTimeoutViperKey, cucumberTimeoutDefault)
	_ = viper.BindEnv(cucumberTimeoutViperKey, cucumberTimeoutEnv)
	viper.SetDefault(cucumberMiningTimeoutViperKey, cucumberMiningTimeoutDefault)
	_ = viper.BindEnv(cucumberMiningTimeoutViperKey, cucumberMiningTimeoutEnv)
}

// InitFlags register Step flags
func InitFlags(f *pflag.FlagSet) {
	Timeout(f)
	MiningTimeout(f)
}

var (
	cucumberTimeoutFlag     = "cucumber-steps-timeout"
	cucumberTimeoutViperKey = "cucumber.steps.timeout"
	cucumberTimeoutDefault  = 5
	cucumberTimeoutEnv      = "CUCUMBER_STEPS_TIMEOUT"
)

// Timeout register flag for Timeout Option
func Timeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration for waiting envelopes to be processed by a step method before timeout
Environment variable: %q`, cucumberTimeoutEnv)
	f.Int(cucumberTimeoutFlag, cucumberTimeoutDefault, desc)
	_ = viper.BindPFlag(cucumberTimeoutViperKey, f.Lookup(cucumberTimeoutFlag))
}

var (
	cucumberMiningTimeoutFlag     = "cucumber-steps-miningtimeout"
	cucumberMiningTimeoutViperKey = "cucumber.steps.miningtimeout"
	cucumberMiningTimeoutDefault  = 10
	cucumberMiningTimeoutEnv      = "CUCUMBER_STEPS_MININGTIMEOUT"
)

// MiningTimeout register flag for MiningTimeout Option
func MiningTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration for waiting envelopes to be processed by a blockchain before timeout
Environment variable: %q`, cucumberMiningTimeoutEnv)
	f.Int(cucumberMiningTimeoutFlag, cucumberMiningTimeoutDefault, desc)
	_ = viper.BindPFlag(cucumberMiningTimeoutViperKey, f.Lookup(cucumberMiningTimeoutFlag))
}
