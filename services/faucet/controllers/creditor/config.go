package creditor

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(creditorAddressViperKey, creditorAddressDefault)
	_ = viper.BindEnv(creditorAddressViperKey, creditorAddressEnv)
}

var (
	creditorAddressFlag     = "faucet-creditor-address"
	creditorAddressViperKey = "faucet.creditor.address"
	creditorAddressDefault  []string
	creditorAddressEnv      = "FAUCET_CREDITOR_ADDRESS"
)

// FaucetAddress register flag for Faucet address
func FaucetAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Addresses of Faucet on each chain (format <address>@<chainID>)
Environment variable: %q`, creditorAddressEnv)
	f.StringSlice(creditorAddressFlag, creditorAddressDefault, desc)
	_ = viper.BindPFlag(creditorAddressViperKey, f.Lookup(creditorAddressFlag))
}
