package http

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hostnameViperKey, hostnameDefault)
	_ = viper.BindEnv(hostnameViperKey, hostnameEnv)
}

var (
	hostnameFlag     = "http-hostname"
	hostnameViperKey = "http.hostname"
	hostnameDefault  = ":8080"
	hostnameEnv      = "HTTP_HOSTNAME"
)

// Hostname register a flag for Redis server address
func Hostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose HTTP server
Environment variable: %q`, hostnameEnv)
	f.String(hostnameFlag, hostnameDefault, desc)
	_ = viper.BindPFlag(hostnameViperKey, f.Lookup(hostnameFlag))
}
