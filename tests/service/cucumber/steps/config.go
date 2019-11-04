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
	viper.SetDefault(cucumberChainIDPrimaryViperKey, cucumberChainIDPrimaryDefault)
	_ = viper.BindEnv(cucumberChainIDPrimaryViperKey, cucumberChainIDPrimaryEnv)
	viper.SetDefault(cucumberChainIDSecondaryViperKey, cucumberChainIDSecondaryDefault)
	_ = viper.BindEnv(cucumberChainIDSecondaryViperKey, cucumberChainIDSecondaryEnv)
}

// InitFlags register Step flags
func InitFlags(f *pflag.FlagSet) {
	Timeout(f)
	MiningTimeout(f)
	ChainIDPrimary(f)
	ChainIDSecondary(f)
}

const (
	cucumberTimeoutFlag     = "cucumber-steps-timeout"
	cucumberTimeoutViperKey = "cucumber.steps.timeout"
	cucumberTimeoutDefault  = 15
	cucumberTimeoutEnv      = "CUCUMBER_STEPS_TIMEOUT"
)

// Timeout register flag for Timeout Option
func Timeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration for waiting envelopes to be processed by a step method before timeout
Environment variable: %q`, cucumberTimeoutEnv)
	f.Int(cucumberTimeoutFlag, cucumberTimeoutDefault, desc)
	_ = viper.BindPFlag(cucumberTimeoutViperKey, f.Lookup(cucumberTimeoutFlag))
}

const (
	cucumberMiningTimeoutFlag     = "cucumber-steps-miningtimeout"
	cucumberMiningTimeoutViperKey = "cucumber.steps.miningtimeout"
	cucumberMiningTimeoutDefault  = 30
	cucumberMiningTimeoutEnv      = "CUCUMBER_STEPS_MININGTIMEOUT"
)

// MiningTimeout register flag for MiningTimeout Option
func MiningTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration for waiting envelopes to be processed by a blockchain before timeout
Environment variable: %q`, cucumberMiningTimeoutEnv)
	f.Int(cucumberMiningTimeoutFlag, cucumberMiningTimeoutDefault, desc)
	_ = viper.BindPFlag(cucumberMiningTimeoutViperKey, f.Lookup(cucumberMiningTimeoutFlag))
}

const (
	cucumberChainIDPrimaryFlag     = "cucumber-chainid-primary"
	cucumberChainIDPrimaryViperKey = "cucumber.chainid.primary"
	cucumberChainIDPrimaryDefault  = ""
	cucumberChainIDPrimaryEnv      = "CUCUMBER_CHAINID_PRIMARY"
)

// ChainIDPrimary register flag for ChainIDPrimary Option
func ChainIDPrimary(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`ChainID corresponding to the alias "primary" in the scenario features
Environment variable: %q`, cucumberChainIDPrimaryEnv)
	f.String(cucumberChainIDPrimaryFlag, cucumberChainIDPrimaryDefault, desc)
	_ = viper.BindPFlag(cucumberChainIDPrimaryViperKey, f.Lookup(cucumberChainIDPrimaryFlag))
}

const (
	cucumberChainIDSecondaryFlag     = "cucumber-chainid-secondary"
	cucumberChainIDSecondaryViperKey = "cucumber.chainid.secondary"
	cucumberChainIDSecondaryDefault  = ""
	cucumberChainIDSecondaryEnv      = "CUCUMBER_CHAINID_SECONDARY"
)

// ChainIDSecondary register flag for ChainIDSecondary Option
func ChainIDSecondary(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`ChainID corresponding to the alias "secondary" in the scenario features
Environment variable: %q`, cucumberChainIDSecondaryEnv)
	f.String(cucumberChainIDSecondaryFlag, cucumberChainIDSecondaryDefault, desc)
	_ = viper.BindPFlag(cucumberChainIDSecondaryViperKey, f.Lookup(cucumberChainIDSecondaryFlag))
}
