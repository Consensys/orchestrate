package rest

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(HostnameViperKey, hostnameDefault)
	_ = viper.BindEnv(HostnameViperKey, hostnameEnv)

	viper.SetDefault(PortViperKey, portDefault)
	_ = viper.BindEnv(PortViperKey, portEnv)
}

const (
	hostnameFlag     = "rest-hostname"
	HostnameViperKey = "rest.hostname"
	hostnameDefault  = ""
	hostnameEnv      = "REST_HOSTNAME"
)

// Hostname register a flag for server address
func Hostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose REST services
Environment variable: %q`, hostnameEnv)
	f.String(hostnameFlag, hostnameDefault, desc)
	_ = viper.BindPFlag(HostnameViperKey, f.Lookup(hostnameFlag))
}

const (
	portFlag     = "rest-port"
	PortViperKey = "rest.port"
	portDefault  = uint(8081)
	portEnv      = "REST_PORT"
)

// Port register a flag for server port
func Port(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose REST services
Environment variable: %q`, portEnv)
	f.Uint(portFlag, portDefault, desc)
	_ = viper.BindPFlag(PortViperKey, f.Lookup(portFlag))
}

func URL() string {
	return fmt.Sprintf("%s:%d", viper.GetString(HostnameViperKey), viper.GetUint(PortViperKey))
}
