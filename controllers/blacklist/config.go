package blacklist

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(faucetBlacklistViperKey, faucetBlacklistEnv)
	viper.SetDefault(faucetBlacklistViperKey, faucetBlacklistDefault)
}

var (
	faucetBlacklistFlag     = "faucet-blacklist"
	faucetBlacklistViperKey = "faucet.ctrl.blacklist"
	faucetBlacklistDefault  = []string{}
	faucetBlacklistEnv      = "FAUCET_BLACKLIST"
)

// FaucetBlacklist register flag for Faucet address
func FaucetBlacklist(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Blacklisted address (format <chainID>-<Address>)
Environment variable: %q`, faucetBlacklistEnv)
	f.StringSlice(faucetBlacklistFlag, faucetBlacklistDefault, desc)
	viper.BindPFlag(faucetBlacklistViperKey, f.Lookup(faucetBlacklistFlag))
}
