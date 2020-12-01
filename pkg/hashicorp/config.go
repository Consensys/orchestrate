package hashicorp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config object that be converted into an api.Config later
type Config struct {
	TokenFilePath string
	MountPoint    string
	SecretPath    string
	RateLimit     float64
	BurstLimit    int
	Address       string
	CACert        string
	CAPath        string
	ClientCert    string
	ClientKey     string
	ClientTimeout time.Duration
	MaxRetries    int
	SkipVerify    bool
	TLSServerName string
}

func init() {
	viper.SetDefault(vaultTokenFilePathViperKey, vaultTokenFilePathDefault)
	_ = viper.BindEnv(vaultTokenFilePathViperKey, vaultTokenFilePathEnv)

	viper.SetDefault(vaultMountPointViperKey, vaultMountPointDefault)
	_ = viper.BindEnv(vaultMountPointViperKey, vaultMountPointEnv)

	viper.SetDefault(vaultRateLimitViperKey, vaultRateLimitDefault)
	_ = viper.BindEnv(vaultRateLimitViperKey, vaultRateLimitEnv)

	viper.SetDefault(vaultBurstLimitViperKey, vaultBurstLimitDefault)
	_ = viper.BindEnv(vaultBurstLimitViperKey, vaultBurstLimitEnv)

	viper.SetDefault(vaultAddrViperKey, vaultAddrDefault)
	_ = viper.BindEnv(vaultAddrViperKey, vaultAddrEnv)

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
}

const (
	vaultTokenFilePathEnv = "VAULT_TOKEN_FILE"
	vaultMountPointEnv    = "VAULT_MOUNT_POINT"
	vaultRateLimitEnv     = "VAULT_RATE_LIMIT"
	vaultBurstLimitEnv    = "VAULT_BURST_LIMIT"
	vaultAddrEnv          = "VAULT_ADDR"
	vaultCACertEnv        = "VAULT_CACERT"
	vaultCAPathEnv        = "VAULT_CAPATH"
	vaultClientCertEnv    = "VAULT_CLIENT_CERT"
	vaultClientKeyEnv     = "VAULT_CLIENT_KEY"
	vaultClientTimeoutEnv = "VAULT_CLIENT_TIMEOUT"
	vaultMaxRetriesEnv    = "VAULT_MAX_RETRIES"
	vaultSkipVerifyEnv    = "VAULT_SKIP_VERIFY"
	vaultTLSServerNameEnv = "VAULT_TLS_SERVER_NAME"

	vaultTokenFilePathFlag = "vault-token-file"
	vaultMountPointFlag    = "vault-mount-point"
	vaultRateLimitFlag     = "vault-rate-limit"
	vaultBurstLimitFlag    = "vault-burst-limit"
	vaultAddrFlag          = "vault-addr"
	vaultCACertFlag        = "vault-cacert"
	vaultCAPathFlag        = "vault-capath"
	vaultClientCertFlag    = "vault-client-cert"
	vaultClientKeyFlag     = "vault-client-key"
	vaultClientTimeoutFlag = "vault-client-timeout"
	vaultMaxRetriesFlag    = "vault-max-retries"
	vaultSkipVerifyFlag    = "vault-skip-verify"
	vaultTLSServerNameFlag = "vault-tls-server-name"

	vaultTokenFilePathViperKey = "vault.token.file"
	vaultMountPointViperKey    = "vault.mount.point"
	vaultRateLimitViperKey     = "vault.rate.limit"
	vaultBurstLimitViperKey    = "vault.burst.limit"
	vaultAddrViperKey          = "vault.addr"
	vaultCACertViperKey        = "vault.cacert"
	vaultCAPathViperKey        = "vault.capath"
	vaultClientCertViperKey    = "vault.client.cert"
	vaultClientKeyViperKey     = "vault.client.key"
	vaultClientTimeoutViperKey = "vault.client.timeout"
	vaultMaxRetriesViperKey    = "vault.max.retries"
	vaultSkipVerifyViperKey    = "vault.skip.verify"
	vaultTLSServerNameViperKey = "vault.tls.server.name"

	// No need to redefine the default here
	vaultTokenFilePathDefault = "/vault/token/.vault-token"
	vaultMountPointDefault    = "orchestrate"
	vaultRateLimitDefault     = float64(0)
	vaultBurstLimitDefault    = int(0)
	vaultAddrDefault          = "https://127.0.0.1:8200"
	vaultCACertDefault        = ""
	vaultCAPathDefault        = ""
	vaultClientCertDefault    = ""
	vaultClientKeyDefault     = ""
	vaultClientTimeoutDefault = time.Second * 60
	vaultMaxRetriesDefault    = int(0)
	vaultSkipVerifyDefault    = false
	vaultTLSServerNameDefault = ""
)

// InitFlags register flags for HashiCorp Vault
func InitFlags(f *pflag.FlagSet) {
	vaultAddr(f)
	vaultBurstLimit(f)
	vaultCACert(f)
	vaultCAPath(f)
	vaultClientCert(f)
	vaultClientKey(f)
	vaultClientTimeout(f)
	vaultMaxRetries(f)
	vaultMountPoint(f)
	vaultRateLimit(f)
	vaultSkipVerify(f)
	vaultTLSServerName(f)
	vaultTokenFilePath(f)
}

func vaultTokenFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the token file path.
Parameter ignored if the token has been passed by VAULT_TOKEN
Environment variable: %q `, vaultTokenFilePathEnv)
	f.String(vaultTokenFilePathFlag, vaultTokenFilePathDefault, desc)
	_ = viper.BindPFlag(vaultTokenFilePathViperKey, f.Lookup(vaultTokenFilePathFlag))
}

func vaultMountPoint(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the mount point used. Should not start with a //
Environment variable: %q `, vaultMountPointEnv)
	f.String(vaultMountPointFlag, vaultMountPointDefault, desc)
	_ = viper.BindPFlag(vaultMountPointViperKey, f.Lookup(vaultMountPointFlag))
}

func vaultRateLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp query rate limit
Environment variable: %q`, vaultRateLimitEnv)
	f.Float64(vaultRateLimitFlag, vaultRateLimitDefault, desc)
	_ = viper.BindPFlag(vaultRateLimitViperKey, f.Lookup(vaultRateLimitFlag))
}

func vaultBurstLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp query burst limit
Environment variable: %q`, vaultRateLimitEnv)
	f.Int(vaultBurstLimitFlag, vaultBurstLimitDefault, desc)
	_ = viper.BindPFlag(vaultBurstLimitViperKey, f.Lookup(vaultBurstLimitFlag))
}

func vaultAddr(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp URL of the remote hashicorp vault
Environment variable: %q`, vaultAddrEnv)
	f.String(vaultAddrFlag, vaultAddrDefault, desc)
	_ = viper.BindPFlag(vaultAddrViperKey, f.Lookup(vaultAddrFlag))
}

func vaultCACert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp CA certificate
Environment variable: %q`, vaultCACertEnv)
	f.String(vaultCACertFlag, vaultCACertDefault, desc)
	_ = viper.BindPFlag(vaultCACertViperKey, f.Lookup(vaultCACertFlag))
}

func vaultCAPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path toward the CA certificate
Environment variable: %q`, vaultCAPathEnv)
	f.String(vaultCAPathFlag, vaultCAPathDefault, desc)
	_ = viper.BindPFlag(vaultCAPathViperKey, f.Lookup(vaultCAPathFlag))
}

func vaultClientCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Certificate of the client
Environment variable: %q`, vaultClientCertEnv)
	f.String(vaultClientCertFlag, vaultClientCertDefault, desc)
	_ = viper.BindPFlag(vaultClientCertViperKey, f.Lookup(vaultClientCertFlag))
}

func vaultClientKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp client key
Environment variable: %q`, vaultClientKeyEnv)
	f.String(vaultClientKeyFlag, vaultClientKeyDefault, desc)
	_ = viper.BindPFlag(vaultClientKeyViperKey, f.Lookup(vaultClientKeyFlag))
}

func vaultClientTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp clean timeout of the client
Environment variable: %q`, vaultClientTimeoutEnv)
	f.Duration(vaultClientTimeoutFlag, vaultClientTimeoutDefault, desc)
	_ = viper.BindPFlag(vaultClientTimeoutViperKey, f.Lookup(vaultClientTimeoutFlag))
}

func vaultMaxRetries(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp max retry for a request
Environment variable: %q`, vaultMaxRetriesEnv)
	f.Int(vaultMaxRetriesFlag, vaultMaxRetriesDefault, desc)
	_ = viper.BindPFlag(vaultMaxRetriesViperKey, f.Lookup(vaultMaxRetriesFlag))
}

func vaultSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp skip verification
Environment variable: %q`, vaultSkipVerifyEnv)
	f.Bool(vaultSkipVerifyFlag, vaultSkipVerifyDefault, desc)
	_ = viper.BindPFlag(vaultSkipVerifyViperKey, f.Lookup(vaultSkipVerifyFlag))
}

func vaultTLSServerName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp TLS server name
Environment variable: %q`, vaultTLSServerNameEnv)
	f.String(vaultTLSServerNameFlag, vaultTLSServerNameDefault, desc)
	_ = viper.BindPFlag(vaultTLSServerNameViperKey, f.Lookup(vaultTLSServerNameFlag))
}

// ConfigFromViper returns a local config object that be converted into an api.Config
func ConfigFromViper() *Config {
	return &Config{
		Address:       viper.GetString(vaultAddrViperKey),
		BurstLimit:    viper.GetInt(vaultBurstLimitViperKey),
		CACert:        viper.GetString(vaultCACertViperKey),
		CAPath:        viper.GetString(vaultCAPathViperKey),
		ClientCert:    viper.GetString(vaultClientCertViperKey),
		ClientKey:     viper.GetString(vaultClientKeyViperKey),
		ClientTimeout: viper.GetDuration(vaultClientTimeoutViperKey),
		MaxRetries:    viper.GetInt(vaultMaxRetriesViperKey),
		MountPoint:    viper.GetString(vaultMountPointViperKey),
		RateLimit:     viper.GetFloat64(vaultRateLimitViperKey),
		SkipVerify:    viper.GetBool(vaultSkipVerifyViperKey),
		TLSServerName: viper.GetString(vaultTLSServerNameViperKey),
		TokenFilePath: viper.GetString(vaultTokenFilePathViperKey),
	}
}

// ToVaultConfig extracts an api.Config object from self
func ToVaultConfig(c *Config) *api.Config {
	// Create Vault Configuration
	config := api.DefaultConfig()
	config.Address = c.Address
	config.HttpClient = cleanhttp.DefaultClient()
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

	// Configure TLS
	tlsConfig := &api.TLSConfig{
		CACert:        c.CACert,
		CAPath:        c.CAPath,
		ClientCert:    c.ClientCert,
		ClientKey:     c.ClientKey,
		TLSServerName: c.TLSServerName,
		Insecure:      c.SkipVerify,
	}

	_ = config.ConfigureTLS(tlsConfig)

	config.Limiter = rate.NewLimiter(rate.Limit(c.RateLimit), c.BurstLimit)
	config.MaxRetries = c.MaxRetries
	config.Timeout = c.ClientTimeout

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
