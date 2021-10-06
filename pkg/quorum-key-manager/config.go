package quorumkeymanager

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
	viper.SetDefault(MetricsURLViperKey, metricsURLDefault)
	_ = viper.BindEnv(MetricsURLViperKey, metricsURLEnv)
	viper.SetDefault(StoreNameViperKey, StoreNameDefault)
	_ = viper.BindEnv(StoreNameViperKey, storeNameEnv)
	viper.SetDefault(tlsSkipVerifyViperKey, tlsSKipVerifyDefault)
	_ = viper.BindEnv(tlsSkipVerifyViperKey, tlsSkipVerifyEnv)
	viper.SetDefault(AuthAPIKeyViperKey, AuthAPIKeyDefault)
	_ = viper.BindEnv(AuthAPIKeyViperKey, authAPIKeyEnv)
	viper.SetDefault(AuthClientTLSCertViperKey, AuthClientTLSCertDefault)
	_ = viper.BindEnv(AuthClientTLSCertViperKey, authClientTLSCertEnv)
	viper.SetDefault(AuthClientTLSKeyViperKey, AuthClientTLSKeyDefault)
	_ = viper.BindEnv(AuthClientTLSKeyViperKey, authClientTLSKeyEnv)
}

const (
	urlFlag     = "key-manager-url"
	URLViperKey = "key.manager.url"
	urlDefault  = "http://localhost:8081"
	urlEnv      = "KEY_MANAGER_URL"
)

func url(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager HTTP domain.
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

const (
	metricsURLFlag     = "key-manager-metrics-url"
	MetricsURLViperKey = "key.manager.metrics.url"
	metricsURLDefault  = "http://localhost:8082"
	metricsURLEnv      = "KEY_MANAGER_METRICS_URL"
)

func metricsURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager HTTP metrics domain.
Environment variable: %q`, metricsURLEnv)
	f.String(metricsURLFlag, metricsURLDefault, desc)
	_ = viper.BindPFlag(MetricsURLViperKey, f.Lookup(metricsURLFlag))
}

const (
	storeNameFlag     = "key-manager-store-name"
	StoreNameViperKey = "key.manager.store.name"
	StoreNameDefault  = ""
	storeNameEnv      = "KEY_MANAGER_STORE_NAME"
)

func storeName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager ethereum account store name.
Environment variable: %q`, storeNameEnv)
	f.String(storeNameFlag, metricsURLDefault, desc)
	_ = viper.BindPFlag(StoreNameViperKey, f.Lookup(storeNameFlag))
}

const (
	tlsSkipVerifyFlag     = "key-manager-tls-skip-verify"
	tlsSkipVerifyViperKey = "key.manager.tls.skip.verify"
	tlsSKipVerifyDefault  = false
	tlsSkipVerifyEnv      = "KEY_MANAGER_TLS_SKIP_VERIFY"
)

func tlsSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager, disables SSL certificate verification.
Environment variable: %q`, tlsSkipVerifyEnv)
	f.Bool(tlsSkipVerifyFlag, tlsSKipVerifyDefault, desc)
	_ = viper.BindPFlag(tlsSkipVerifyViperKey, f.Lookup(tlsSkipVerifyFlag))
}

const (
	authAPIKeyFlag     = "key-manager-api-key"
	AuthAPIKeyViperKey = "key.manager.api.key"
	AuthAPIKeyDefault  = ""
	authAPIKeyEnv      = "KEY_MANAGER_API_KEY"
)

func authAPIKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager API-KEY authentication.
Environment variable: %q`, authAPIKeyEnv)
	f.String(authAPIKeyFlag, AuthAPIKeyDefault, desc)
	_ = viper.BindPFlag(AuthAPIKeyViperKey, f.Lookup(authAPIKeyFlag))
}

const (
	authClientTLSCertFlag     = "key-manager-client-tls-cert"
	AuthClientTLSCertViperKey = "key.manager.client.tls.cert"
	AuthClientTLSCertDefault  = ""
	authClientTLSCertEnv      = "KEY_MANAGER_CLIENT_TLS_CERT"
)

const (
	authClientTLSKeyFlag     = "key-manager-client-tls-key"
	AuthClientTLSKeyViperKey = "key.manager.client.tls.key"
	AuthClientTLSKeyDefault  = ""
	authClientTLSKeyEnv      = "KEY_MANAGER_CLIENT_TLS_KEY"
)

func authTLSCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication (crt file).
Environment variable: %q`, authClientTLSCertEnv)
	f.String(authClientTLSCertFlag, AuthClientTLSCertDefault, desc)
	_ = viper.BindPFlag(AuthClientTLSCertViperKey, f.Lookup(authClientTLSCertFlag))
}

func authTLSKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication (key file).
Environment variable: %q`, authClientTLSKeyEnv)
	f.String(authClientTLSKeyFlag, AuthClientTLSKeyDefault, desc)
	_ = viper.BindPFlag(AuthClientTLSKeyViperKey, f.Lookup(authClientTLSKeyFlag))
}

func Flags(f *pflag.FlagSet) {
	url(f)
	metricsURL(f)
	storeName(f)
	tlsSkipVerify(f)
	authAPIKey(f)
	authTLSCert(f)
	authTLSKey(f)
}

type Config struct {
	URL        string
	MetricsURL string
	StoreName  string
}

func NewConfigFromViper(vipr *viper.Viper) *Config {
	return &Config{
		URL:        vipr.GetString(URLViperKey),
		MetricsURL: vipr.GetString(MetricsURLViperKey),
		StoreName:  vipr.GetString(StoreNameViperKey),
	}
}
