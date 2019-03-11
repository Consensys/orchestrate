package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ethClientFlag     = "eth-client"
	ethClientViperKey = "eth.clients"
	ethClientDefault  = []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	}
	ethClientEnv = "ETH_CLIENT_URL"
)

// EthClientURLs register flag for Ethereum client URLs
func EthClientURLs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Ethereum client URLs.
Environment variable: %q`, ethClientEnv)
	f.StringSlice(ethClientFlag, ethClientDefault, desc)
	viper.BindPFlag(ethClientViperKey, f.Lookup(ethClientFlag))
	viper.BindEnv(ethClientViperKey, ethClientEnv)
}
