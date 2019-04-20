package creditor

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(creditorAddressViperKey, creditorAddressEnv)
	viper.SetDefault(creditorAddressViperKey, creditorAddressDefault)
}

var (
	creditorAddressFlag     = "faucet-creditor"
	creditorAddressViperKey = "faucet.creditors"
	creditorAddressDefault  = []string{}
	creditorAddressEnv      = "FAUCET_CREDITOR_ADDRESS"
)

// FaucetAddress register flag for Faucet address
func FaucetAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Faucet on each chain (format <chainID>:<Address>)
Environment variable: %q`, creditorAddressEnv)
	f.StringSlice(creditorAddressFlag, creditorAddressDefault, desc)
	viper.BindPFlag(creditorAddressViperKey, f.Lookup(creditorAddressFlag))
}
