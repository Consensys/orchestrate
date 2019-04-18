package hashicorp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	vault "github.com/hashicorp/vault/api"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"
)

func init() {
	viper.SetDefault(vaultSecretPathViperKey, vaultSecretPathDefault)
	_ = viper.BindEnv(vaultSecretPathViperKey, vaultSecretPathEnv)

	viper.SetDefault(vaultRateLimitViperKey, vaultRateLimitDefault)
	_ = viper.BindEnv(vaultRateLimitViperKey, vaultRateLimitEnv)

	viper.SetDefault(vaultBurstLimitViperKey, vaultBurstLimitDefault)
	_ = viper.BindEnv(vaultBurstLimitViperKey, vaultBurstLimitEnv)

	viper.SetDefault(vaultAddressViperKey, vaultAddressDefault)
	_ = viper.BindEnv(vaultAddressViperKey, vaultAddressEnv)

	viper.SetDefault(vaultCACertViperKey, vaultCACertDefault)
	_ = viper.BindEnv(vaultCACertViperKey, vaultCACertEnv)

	viper.SetDefault(vaultCAPathViperKey, vaultCAPathDefault)
	_ = viper.BindEnv(vaultCAPathViperKey, vaultCAPathEnv)

	viper.SetDefault(vaultClientCertViperKey, vaultClientCertDefault)
	_ = viper.BindEnv(vaultClientCertViperKey, vaultClientCertEnv)

	viper.SetDefault(vaultClientKeyViperKey, vaultClientKeyDefault)
	_ = viper.BindEnv(vaultClientKeyViperKey, vaultClientKeyEnv)

	viper.SetDefault(vaultClientTimeoutViperKey, vaultClientTimeoutDefault)
	_ = viper.BindEnv(vaultClientTimeoutViperKey, vaultClientTimeoutEnv)

	viper.SetDefault(vaultMaxRetriesViperKey, vaultMaxRetriesDefault)
	_ = viper.BindEnv(vaultMaxRetriesViperKey, vaultMaxRetriesEnv)

	viper.SetDefault(vaultSkipVerifyViperKey, vaultSkipVerifyDefault)
	_ = viper.BindEnv(vaultSkipVerifyViperKey, vaultSkipVerifyEnv)

	viper.SetDefault(vaultTLSServerNameViperKey, vaultTLSServerNameDefault)
	_ = viper.BindEnv(vaultTLSServerNameViperKey, vaultTLSServerNameEnv)

	log.Infof("Using vault_address: %v \n", os.Getenv("VAULT_ADDR"))
	// log.Infof("Using vault_agent_address: %v \n", os.Getenv("VAULT_AGENT_ADDR"))
}

var (
	vaultSecretPathEnv    = "VAULT_SECRET_PATH"
	vaultRateLimitEnv     = "VAULT_RATE_LIMIT"
	vaultBurstLimitEnv    = "VAULT_BURST_LIMIT"
	vaultAddressEnv       = "VAULT_ADDR"
	vaultCACertEnv        = "VAULT_CACERT"
	vaultCAPathEnv        = "VAULT_CAPATH"
	vaultClientCertEnv    = "VAULT_CLIENT_CERT"
	vaultClientKeyEnv     = "VAULT_CLIENT_KEY"
	vaultClientTimeoutEnv = "VAULT_CLIENT_TIMEOUT"
	vaultMaxRetriesEnv    = "VAULT_MAX_RETRIES"
	vaultSkipVerifyEnv    = "VAULT_SKIP_VERIFY"
	vaultTLSServerNameEnv = "VAULT_TLS_SERVER_NAME"

	vaultSecretPathFlag    = "vault-secret-path"
	vaultRateLimitFlag     = "vault-rate-limit"
	vaultBurstLimitFlag    = "vault-burst-limit"
	vaultAddressFlag       = "vault-addr"
	vaultCACertFlag        = "vault-cacert"
	vaultCAPathFlag        = "vault-capath"
	vaultClientCertFlag    = "vault-client-cert"
	vaultClientKeyFlag     = "vault-client-key"
	vaultClientTimeoutFlag = "vault-client-timeout"
	vaultMaxRetriesFlag    = "vault-max-retries"
	vaultSkipVerifyFlag    = "vault-skip-verify"
	vaultTLSServerNameFlag = "vault-tls-server-name"

	vaultSecretPathViperKey    = "vault.secret.path"
	vaultRateLimitViperKey     = "vault.rate.limit"
	vaultBurstLimitViperKey    = "vault.burst.limit"
	vaultAddressViperKey       = "vault.addr"
	vaultCACertViperKey        = "vault.cacert"
	vaultCAPathViperKey        = "vault.capath"
	vaultClientCertViperKey    = "vault.client.cert"
	vaultClientKeyViperKey     = "vault.client.key"
	vaultClientTimeoutViperKey = "vault.client.timeout"
	vaultMaxRetriesViperKey    = "vault.max.retries"
	vaultSkipVerifyViperKey    = "vault.skip.verify"
	vaultTLSServerNameViperKey = "vault.tls.server.name"

	// No need to redefine the default here
	vaultSecretPathDefault    = "/secret"
	vaultRateLimitDefault     float64
	vaultBurstLimitDefault    int
	vaultAddressDefault       = "https://127.0.0.1:8200"
	vaultCACertDefault        string
	vaultCAPathDefault        string
	vaultClientCertDefault    string
	vaultClientKeyDefault     string
	vaultClientTimeoutDefault = time.Second * 60
	vaultMaxRetriesDefault    int
	vaultSkipVerifyDefault    bool
	vaultTLSServerNameDefault string
)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultSecretPath(f)
	VaultRateLimit(f)
	VaultBurstLimit(f)
	VaultAddress(f)
	VaultCACert(f)
	VaultCAPath(f)
	VaultClientCert(f)
	VaultClientKey(f)
	VaultClientTimeout(f)
	VaultMaxRetries(f)
	VaultSkipVerify(f)
	VaultTLSServerName(f)
}

// VaultSecretPath registers a flag for the path used by vault secret engine
func VaultSecretPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultSecretPathEnv)
	f.String(vaultSecretPathFlag, vaultSecretPathDefault, desc)
	viper.BindPFlag(vaultSecretPathViperKey, f.Lookup(vaultSecretPathFlag))
}

// VaultRateLimit registers a flag for the path used by vault secret engine
func VaultRateLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultRateLimitEnv)
	f.Float64(vaultRateLimitFlag, vaultRateLimitDefault, desc)
	viper.BindPFlag(vaultRateLimitViperKey, f.Lookup(vaultRateLimitFlag))
}

// VaultBurstLimit registers a flag for the path used by vault secret engine
func VaultBurstLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultRateLimitEnv)
	f.Int(vaultBurstLimitFlag, vaultBurstLimitDefault, desc)
	viper.BindPFlag(vaultBurstLimitViperKey, f.Lookup(vaultBurstLimitFlag))
}

// VaultAddress registers a flag for the path used by vault secret engine
func VaultAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultAddressEnv)
	f.String(vaultAddressFlag, vaultAddressDefault, desc)
	viper.BindPFlag(vaultAddressViperKey, f.Lookup(vaultAddressFlag))
}

// VaultCACert registers a flag for the path used by vault secret engine
func VaultCACert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultCACertEnv)
	f.String(vaultCACertFlag, vaultCACertDefault, desc)
	viper.BindPFlag(vaultCACertViperKey, f.Lookup(vaultCACertFlag))
}

// VaultCAPath registers a flag for the path used by vault secret engine
func VaultCAPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultCAPathEnv)
	f.String(vaultCAPathFlag, vaultCAPathDefault, desc)
	viper.BindPFlag(vaultCAPathViperKey, f.Lookup(vaultCAPathFlag))
}

// VaultClientCert registers a flag for the path used by vault secret engine
func VaultClientCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultClientCertEnv)
	f.String(vaultClientCertFlag, vaultClientCertDefault, desc)
	viper.BindPFlag(vaultClientCertViperKey, f.Lookup(vaultClientCertFlag))
}

// VaultClientKey registers a flag for the path used by vault secret engine
func VaultClientKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultClientKeyEnv)
	f.String(vaultClientKeyFlag, vaultClientKeyDefault, desc)
	viper.BindPFlag(vaultClientKeyViperKey, f.Lookup(vaultClientKeyFlag))
}

// VaultClientTimeout registers a flag for the path used by vault secret engine
func VaultClientTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultClientTimeoutEnv)
	f.Duration(vaultClientTimeoutFlag, vaultClientTimeoutDefault, desc)
	viper.BindPFlag(vaultClientTimeoutViperKey, f.Lookup(vaultClientTimeoutFlag))
}

// VaultMaxRetries registers a flag for the path used by vault secret engine
func VaultMaxRetries(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultMaxRetriesEnv)
	f.Int(vaultMaxRetriesFlag, vaultMaxRetriesDefault, desc)
	viper.BindPFlag(vaultMaxRetriesViperKey, f.Lookup(vaultMaxRetriesFlag))
}

// VaultSkipVerify registers a flag for vault client
func VaultSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultSkipVerifyEnv)
	f.Bool(vaultSkipVerifyFlag, vaultSkipVerifyDefault, desc)
	viper.BindPFlag(vaultSkipVerifyViperKey, f.Lookup(vaultSkipVerifyFlag))
}

// VaultTLSServerName registers a flag for vault client
func VaultTLSServerName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultTLSServerNameEnv)
	f.String(vaultTLSServerNameFlag, vaultTLSServerNameDefault, desc)
	viper.BindPFlag(vaultTLSServerNameViperKey, f.Lookup(vaultTLSServerNameFlag))
}

// TODO: update Hashicorp creation

// NewConfig override the environment variable
// defined by the SDK with the parameters passed by Viper
func NewConfig() *vault.Config {
	// Create Vault Configuration
	config := &vault.Config{
		Address:    viper.GetString(vaultAddressViperKey),
		HttpClient: cleanhttp.DefaultClient(),
	}
	config.HttpClient.Timeout = time.Second * 60

	// Create Transport
	transport := config.HttpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		config.Error = err
		return config
	}

	// Replicate ReadEnvironment behaviour

	// Configure TLS
	tlsConfig := &vault.TLSConfig{
		CACert:        viper.GetString(vaultCACertViperKey),
		CAPath:        viper.GetString(vaultCAPathViperKey),
		ClientCert:    viper.GetString(vaultClientCertViperKey),
		ClientKey:     viper.GetString(vaultClientKeyViperKey),
		TLSServerName: viper.GetString(vaultTLSServerNameViperKey),
		Insecure: 	   viper.GetBool(vaultSkipVerifyViperKey),
	}

	config.ConfigureTLS(tlsConfig)

	rateLimit := viper.GetFloat64(vaultRateLimitViperKey)
	burstLimit := viper.GetInt(vaultBurstLimitViperKey)
	config.Limiter = rate.NewLimiter(rate.Limit(rateLimit), burstLimit)

	config.MaxRetries = viper.GetInt(vaultMaxRetriesViperKey)

	config.Timeout = viper.GetDuration(vaultClientTimeoutViperKey)

	// Ensure redirects are not automatically followed
	// Note that this is sane for the API client as it has its own
	// redirect handling logic (and thus also for command/meta),
	// but in e.g. http_test actual redirect handling is necessary
	config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Returning this value causes the Go net library to not close the
		// response body and to nil out the error. Otherwise retry clients may
		// try three times on every redirect because it sees an error from this
		// function (to prevent redirects) passing through to it.
		return http.ErrUseLastResponse
	}

	config.Backoff = retryablehttp.LinearJitterBackoff

	return config
}

// GetSecretPath returns the secret path set in deployment by vault
func GetSecretPath() string {
	return viper.GetString(vaultSecretPathViperKey)
}
