package tessera

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(URLsViperKey, urlsEnv)
	viper.SetDefault(URLsViperKey, urlsDefault)
}

var (
	urlsFlag     = "tessera-url"
	URLsViperKey = "tessera.urls"
	urlsEnv      = "TESSERA_URL"
	urlsDefault  = map[string]string{}
)

func InitFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Tessera URLs (endpoints)
Environment variable: %q`, urlsEnv)
	f.StringToString(urlsFlag, urlsDefault, desc)
	_ = viper.BindPFlag(URLsViperKey, f.Lookup(urlsFlag))
}
