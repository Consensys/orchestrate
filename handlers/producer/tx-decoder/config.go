package txdecoder

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(disableExternalTxViperKey, disableExternalTxDefault)
	_ = viper.BindEnv(disableExternalTxViperKey, disableExternalTxEnv)
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	disableExternalTx(f)
}

var (
	disableExternalTxFlag     = "disable-external-tx"
	disableExternalTxViperKey = "disable.external.tx"
	disableExternalTxDefault  = false
	disableExternalTxEnv      = "DISABLE_EXTERNAL_TX"
)

// disableExternalTx register flag for Listener Start Default
func disableExternalTx(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Boolean: skip all tx that are not sent directly by core-stack
Environment variable: %q`, disableExternalTxEnv)
	f.Bool(disableExternalTxFlag, disableExternalTxDefault, desc)
	_ = viper.BindPFlag(disableExternalTxViperKey, f.Lookup(disableExternalTxFlag))
}

// ExternalTxDisabled returns the config field "disable.external.tx"
func ExternalTxDisabled() bool {
	return viper.GetBool(disableExternalTxViperKey)
}
