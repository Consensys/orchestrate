package steps

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(CucumberTimeoutViperKey, cucumberTimeoutDefault)
	_ = viper.BindEnv(CucumberTimeoutViperKey, cucumberTimeoutEnv)

}

// InitFlags register Step flags
func InitFlags(f *pflag.FlagSet) {
	Timeout(f)
}

const (
	cucumberTimeoutFlag     = "cucumber-steps-timeout"
	CucumberTimeoutViperKey = "cucumber.steps.timeout"
	cucumberTimeoutDefault  = 10 * time.Second
	cucumberTimeoutEnv      = "CUCUMBER_STEPS_TIMEOUT"
)

// Timeout register flag for Timeout Option
func Timeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration for waiting envelopes to be processed by a step method before timeout
Environment variable: %q`, cucumberTimeoutEnv)
	f.Duration(cucumberTimeoutFlag, cucumberTimeoutDefault, desc)
	_ = viper.BindPFlag(CucumberTimeoutViperKey, f.Lookup(cucumberTimeoutFlag))
}
