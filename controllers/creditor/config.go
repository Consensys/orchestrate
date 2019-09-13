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
	creditorAddressFlag     = "faucet-creditor"
	creditorAddressViperKey = "faucet.creditors"
	creditorAddressDefault  []string
	creditorAddressEnv      = "FAUCET_CREDITOR_ADDRESS"
)

// FaucetAddress register flag for Faucet address
func FaucetAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Faucet on each chain (format <chainID>@<Address>)
Environment variable: %q`, creditorAddressEnv)
	f.StringSlice(creditorAddressFlag, creditorAddressDefault, desc)
	_ = viper.BindPFlag(creditorAddressViperKey, f.Lookup(creditorAddressFlag))
}
