package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	httpHostnameFlag     = "http-hostname"
	httpHostnameViperKey = "http.hostname"
	httpHostnameDefault  = ":8080"
	httpHostnameEnv      = "HTTP_HOSTNAME"
)

// HTTPHostname register a flag for Redis server address
func HTTPHostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose healthchecks and metrics.
Environment variable: %q`, httpHostnameEnv)
	f.String(httpHostnameFlag, httpHostnameDefault, desc)
	viper.BindPFlag(httpHostnameViperKey, f.Lookup(httpHostnameFlag))
	viper.BindEnv(httpHostnameViperKey, httpHostnameEnv)
}
