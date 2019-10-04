package hashicorp

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config object that be converted into an api.Config later
type Config struct {
	TokenFilePath string
	MountPoint    string
	KVVersion     string
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

	viper.SetDefault(vaultKVVersionViperKey, vaultKVVersionDefault)
	_ = viper.BindEnv(vaultKVVersionViperKey, vaultKVVersionEnv)

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
}

var (
	vaultTokenFilePathEnv = "VAULT_TOKEN_FILEPATH"
	vaultMountPointEnv    = "VAULT_MOUNT_POINT"
	vaultKVVersionEnv     = "VAULT_KV_VERSION"
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

	vaultTokenFilePathFlag = "vault-token-filepath"
	vaultMountPointFlag    = "vault-mount-point"
	vaultKVVersionFlag     = "vault-kv-version"
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

	vaultTokenFilePathViperKey = "vault.token.filepath"
	vaultMountPointViperKey    = "vault.mount.point"
	vaultKVVersionViperKey     = "vault.kv.version"
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
	vaultTokenFilePathDefault = "/vault/token/.vault-token"
	vaultMountPointDefault    = "secret"
	vaultKVVersionDefault     = "v2" // Could be "v1"
	vaultSecretPathDefault    = "default"
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
	vaultAddress(f)
	vaultBurstLimit(f)
	vaultCACert(f)
	vaultCAPath(f)
	vaultClientCert(f)
	vaultClientKey(f)
	vaultClientTimeout(f)
	vaultKVVersion(f)
	vaultMaxRetries(f)
	vaultMountPoint(f)
	vaultRateLimit(f)
	vaultSecretPath(f)
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
	desc := fmt.Sprintf(`Specifies the mount point used. Should not start with a \"//\"
Environment variable: %q `, vaultMountPointEnv)
	f.String(vaultMountPointFlag, vaultMountPointDefault, desc)
	_ = viper.BindPFlag(vaultMountPointViperKey, f.Lookup(vaultMountPointFlag))
}

func vaultKVVersion(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Determine which version of the kv secret engine we will be using
Can be "v1" or "v2".
Environment variable: %q `, vaultKVVersionEnv)
	f.String(vaultKVVersionFlag, vaultKVVersionDefault, desc)
	_ = viper.BindPFlag(vaultKVVersionViperKey, f.Lookup(vaultKVVersionFlag))
}

func vaultSecretPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultSecretPathEnv)
	f.String(vaultSecretPathFlag, vaultSecretPathDefault, desc)
	_ = viper.BindPFlag(vaultSecretPathViperKey, f.Lookup(vaultSecretPathFlag))
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

func vaultAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp address of the remote hashicorp vault
Environment variable: %q`, vaultAddressEnv)
	f.String(vaultAddressFlag, vaultAddressDefault, desc)
	_ = viper.BindPFlag(vaultAddressViperKey, f.Lookup(vaultAddressFlag))
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
		Address:       viper.GetString(vaultAddressViperKey),
		BurstLimit:    viper.GetInt(vaultBurstLimitViperKey),
		CACert:        viper.GetString(vaultCACertViperKey),
		CAPath:        viper.GetString(vaultCAPathViperKey),
		ClientCert:    viper.GetString(vaultClientCertViperKey),
		ClientKey:     viper.GetString(vaultClientKeyViperKey),
		ClientTimeout: viper.GetDuration(vaultClientTimeoutViperKey),
		KVVersion:     viper.GetString(vaultKVVersionViperKey),
		MaxRetries:    viper.GetInt(vaultMaxRetriesViperKey),
		MountPoint:    viper.GetString(vaultMountPointViperKey),
		RateLimit:     viper.GetFloat64(vaultRateLimitViperKey),
		SecretPath:    viper.GetString(vaultSecretPathViperKey),
		SkipVerify:    viper.GetBool(vaultSkipVerifyViperKey),
		TLSServerName: viper.GetString(vaultTLSServerNameViperKey),
		TokenFilePath: viper.GetString(vaultTokenFilePathViperKey),
	}
}
