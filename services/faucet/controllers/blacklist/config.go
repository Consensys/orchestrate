package blacklist

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(faucetBlacklistViperKey, faucetBlacklistDefault)
	_ = viper.BindEnv(faucetBlacklistViperKey, faucetBlacklistEnv)
}

var (
	faucetBlacklistFlag     = "faucet-blacklist"
	faucetBlacklistViperKey = "faucet.ctrl.blacklist"
	faucetBlacklistDefault  []string
	faucetBlacklistEnv      = "FAUCET_BLACKLIST"
)

// FaucetBlacklist register flag for Faucet address
func FaucetBlacklist(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Blacklisted addresses (format <address1>@<chainID1> <address2>@<chainID2>)
Environment variable: %q`, faucetBlacklistEnv)
	f.StringSlice(faucetBlacklistFlag, faucetBlacklistDefault, desc)
	_ = viper.BindPFlag(faucetBlacklistViperKey, f.Lookup(faucetBlacklistFlag))
}
