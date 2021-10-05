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
	viper.SetDefault(AuthTLSSkipVerifyViperKey, AuthTLSSKipVerifyDefault)
	_ = viper.BindEnv(AuthTLSSkipVerifyViperKey, authTSLSkipVerifyEnv)
	viper.SetDefault(AuthAPIKeyViperKey, AuthAPIKeyDefault)
	_ = viper.BindEnv(AuthAPIKeyViperKey, authAPIKeyEnv)
	viper.SetDefault(AuthTLSCertViperKey, AuthTLSCertDefault)
	_ = viper.BindEnv(AuthTLSCertViperKey, authTLSCertEnv)
	viper.SetDefault(AuthTLSKeyViperKey, AuthTLSKeyDefault)
	_ = viper.BindEnv(AuthTLSKeyViperKey, authTLSKeyEnv)
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
	authTLSCertFlag     = "key-manager-tls-cert"
	AuthTLSCertViperKey = "key.manager.tls.cert"
	AuthTLSCertDefault  = ""
	authTLSCertEnv      = "KEY_MANAGER_TLS_CERT"
)

const (
	authTLSKeyFlag     = "key-manager-tls-key"
	AuthTLSKeyViperKey = "key.manager.tls.key"
	AuthTLSKeyDefault  = ""
	authTLSKeyEnv      = "KEY_MANAGER_TLS_KEY"
)

const (
	authTLSSkipVerifyFlag     = "key-manager-tls-skip-verify"
	AuthTLSSkipVerifyViperKey = "key.manager.tls.skip.verify"
	AuthTLSSKipVerifyDefault  = false
	authTSLSkipVerifyEnv      = "KEY_MANAGER_TLS_SKIP_VERIFY"
)

func authTLSCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication (crt file).
Environment variable: %q`, authTLSCertEnv)
	f.String(authTLSCertFlag, AuthTLSCertDefault, desc)
	_ = viper.BindPFlag(AuthTLSCertViperKey, f.Lookup(authTLSCertFlag))
}

func authTLSKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication (key file).
Environment variable: %q`, authTLSKeyEnv)
	f.String(authTLSKeyFlag, AuthTLSKeyDefault, desc)
	_ = viper.BindPFlag(AuthTLSKeyViperKey, f.Lookup(authTLSKeyFlag))
}

func authTLSSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication, disables SSL certificate verification.
Environment variable: %q`, authTSLSkipVerifyEnv)
	f.Bool(authTLSSkipVerifyFlag, AuthTLSSKipVerifyDefault, desc)
	_ = viper.BindPFlag(AuthTLSSkipVerifyViperKey, f.Lookup(authTLSSkipVerifyFlag))
}

func Flags(f *pflag.FlagSet) {
	url(f)
	metricsURL(f)
	storeName(f)
	authTLSSkipVerify(f)
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
