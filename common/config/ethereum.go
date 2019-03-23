package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(ethClientViperKey, ethClientEnv)
	viper.SetDefault(ethClientViperKey, ethClientDefault)
}

var (
	ethClientFlag     = "eth-client"
	ethClientViperKey = "eth.clients"
	ethClientDefault  = []string{}
	ethClientEnv      = "ETH_CLIENT_URL"
)

// EthClientURLs register flag for Ethereum client URLs
func EthClientURLs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Ethereum client URLs.
Environment variable: %q`, ethClientEnv)
	f.StringSlice(ethClientFlag, ethClientDefault, desc)
	viper.BindPFlag(ethClientViperKey, f.Lookup(ethClientFlag))
}
