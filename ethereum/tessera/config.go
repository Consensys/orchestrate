package tessera

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(tesseraEndpointsViperKey, tesseraEndpointsEnv)
	viper.SetDefault(tesseraEndpointsViperKey, tesseraEndpointsDefault)
}

var (
	tesseraEndpointsFlag     = "tessera-endpoints"
	tesseraEndpointsViperKey = "tessera.endpoints"
	tesseraEndpointsEnv      = "TESSERA_ENDPOINTS"
	tesseraEndpointsDefault  = map[string]string{}
)

func InitFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Tessera endpoints
Environment variable: %q`, tesseraEndpointsEnv)
	f.StringToString(tesseraEndpointsFlag, tesseraEndpointsDefault, desc)
	_ = viper.BindPFlag(tesseraEndpointsViperKey, f.Lookup(tesseraEndpointsFlag))
}
