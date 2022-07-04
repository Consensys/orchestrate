package key

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(APIKeyViperKey, apiKeyDefault)
	_ = viper.BindEnv(APIKeyViperKey, apiKeyEnv)
}

// Provision trusted certificate of the authentication service (base64 encoded)
const (
	apiKeyFlag     = "auth-api-key"
	APIKeyViperKey = "auth.api-key"
	apiKeyDefault  = ""
	apiKeyEnv      = "AUTH_API_KEY"
)

func Flags(f *pflag.FlagSet) {
	apiKey(f)
}

// apiKey register flag for Authentication with API Key
func apiKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key used for authentication (it should be used only for Orchestrate internal authentication).
Environment variable: %q`, apiKeyEnv)
	f.String(apiKeyFlag, apiKeyDefault, desc)
	_ = viper.BindPFlag(APIKeyViperKey, f.Lookup(apiKeyFlag))
}
