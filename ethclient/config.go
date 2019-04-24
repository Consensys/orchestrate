package ethclient

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(urlViperKey, urlEnv)
	viper.SetDefault(urlViperKey, urlDefault)
}

var (
	urlFlag     = "eth-client"
	urlViperKey = "eth.clients"
	urlDefault  = []string{}
	urlEnv      = "ETH_CLIENT_URL"
)

// URLs register flag for Ethereum client urls
func URLs(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Ethereum client url
Environment variable: %q`, urlEnv)
	f.StringSlice(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(urlViperKey, f.Lookup(urlFlag))
}
