package redis

import (
	"fmt"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/tls"
	"github.com/ConsenSys/orchestrate/pkg/tls/certificate"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(HostViperKey, hostDefault)
	_ = viper.BindEnv(HostViperKey, hostEnv)
	viper.SetDefault(PortViperKey, portDefault)
	_ = viper.BindEnv(PortViperKey, portEnv)
	viper.SetDefault(UsernameViperKey, usernameDefault)
	_ = viper.BindEnv(UsernameViperKey, usernameEnv)
	viper.SetDefault(PasswordViperKey, passwordDefault)
	_ = viper.BindEnv(PasswordViperKey, passwordEnv)
	viper.SetDefault(DatabaseViperKey, databaseDefault)
	_ = viper.BindEnv(DatabaseViperKey, databaseEnv)
	viper.SetDefault(TLSEnableViperKey, tlsEnableDefault)
	_ = viper.BindEnv(TLSEnableViperKey, tlsEnableEnv)
	viper.SetDefault(TLSCertViperKey, tlsCertDefault)
	_ = viper.BindEnv(TLSCertViperKey, tlsCertEnv)
	viper.SetDefault(TLSKeyViperKey, tlsKeyDefault)
	_ = viper.BindEnv(TLSKeyViperKey, tlsKeyEnv)
	viper.SetDefault(TLSCAViperKey, tlsCADefault)
	_ = viper.BindEnv(TLSCAViperKey, tlsCAEnv)
	viper.SetDefault(TLSSkipVerifyViperKey, tlsSkipVerifyDefault)
	_ = viper.BindEnv(TLSSkipVerifyViperKey, tlsSkipVerifyEnv)
}

const (
	hostFlag     = "redis-host"
	HostViperKey = "redis.host"
	hostDefault  = "localhost"
	hostEnv      = "REDIS_HOST"
)

const (
	portFlag     = "redis-port"
	PortViperKey = "redis.port"
	portDefault  = "6379"
	portEnv      = "REDIS_PORT"
)

// URL register a flag for Redis server URL
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Host (address) of Redis server to connect to.
Environment variable: %q`, hostEnv)
	f.String(hostFlag, hostDefault, desc)
	_ = viper.BindPFlag(HostViperKey, f.Lookup(hostFlag))

	desc = fmt.Sprintf(`Port (address) of Redis server to connect to.
Environment variable: %q`, portEnv)
	f.String(portFlag, portDefault, desc)
	_ = viper.BindPFlag(PortViperKey, f.Lookup(portFlag))
}

const (
	usernameFlag     = "redis-user"
	UsernameViperKey = "redis.user"
	usernameDefault  = ""
	usernameEnv      = "REDIS_USER"
)

// UsernameFlag register flag for db user
func UsernameFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis Username.
Environment variable: %q`, usernameEnv)
	f.String(usernameFlag, usernameDefault, desc)
	_ = viper.BindPFlag(UsernameViperKey, f.Lookup(usernameFlag))
}

const (
	passwordFlag     = "redis-password"
	PasswordViperKey = "redis.password"
	passwordDefault  = ""
	passwordEnv      = "REDIS_PASSWORD"
)

// PasswordFlag register flag for db password
func PasswordFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis Username password
Environment variable: %q`, passwordEnv)
	f.String(passwordFlag, passwordDefault, desc)
	_ = viper.BindPFlag(PasswordViperKey, f.Lookup(passwordFlag))
}

const (
	databaseFlag     = "redis-database"
	DatabaseViperKey = "redis.database"
	databaseDefault  = -1
	databaseEnv      = "REDIS_DATABASE"
)

// DatabaseFlag register flag for db database name
func DatabaseFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Target Redis database name
Environment variable: %q`, databaseEnv)
	f.Int(databaseFlag, databaseDefault, desc)
	_ = viper.BindPFlag(DatabaseViperKey, f.Lookup(databaseFlag))
}

const (
	tlsEnableFlag     = "redis-tls-enable"
	TLSEnableViperKey = "redis.tls.enable"
	tlsEnableDefault  = false
	tlsEnableEnv      = "REDIS_TLS_ENABLE"
)

func TLSEnableFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable TLS/SSL to connect to Redis
Environment variable: %q`, tlsEnableEnv)
	f.Bool(tlsEnableFlag, tlsEnableDefault, desc)
	_ = viper.BindPFlag(TLSEnableViperKey, f.Lookup(tlsEnableFlag))
}

const (
	tlsCertFlag     = "redis-tls-cert"
	TLSCertViperKey = "redis.tls.cert"
	tlsCertDefault  = ""
	tlsCertEnv      = "REDIS_TLS_CERT"
)

// RedisTLSCert register flag for TLS certificate used to connect to the database
func TLSCertFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Certificate to connect to Redis
Environment variable: %q`, tlsCertEnv)
	f.String(tlsCertFlag, tlsCertDefault, desc)
	_ = viper.BindPFlag(TLSCertViperKey, f.Lookup(tlsCertFlag))
}

const (
	tlsKeyFlag     = "redis-tls-key"
	TLSKeyViperKey = "redis.tls.key"
	tlsKeyDefault  = ""
	tlsKeyEnv      = "REDIS_TLS_KEY"
)

// RedisTLSKey register flag for database TLS private key
func TLSKeyFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Private Key to connect to Redis
Environment variable: %q`, tlsKeyEnv)
	f.String(tlsKeyFlag, tlsKeyDefault, desc)
	_ = viper.BindPFlag(TLSKeyViperKey, f.Lookup(tlsKeyFlag))
}

const (
	tlsCAFlag     = "redis-tls-ca"
	TLSCAViperKey = "redis.tls.ca"
	tlsCADefault  = ""
	tlsCAEnv      = "REDIS_TLS_CA"
)

// RedisTLSCert register flag for database pool size
func TLSCAFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Trusted Certificate Authority
Environment variable: %q`, tlsCAEnv)
	f.String(tlsCAFlag, tlsCADefault, desc)
	_ = viper.BindPFlag(TLSCAViperKey, f.Lookup(tlsCAFlag))
}

const (
	tlsSkipVerifyFlag     = "redis-tls-skip-verify"
	TLSSkipVerifyViperKey = "redis.tls.skip-verify"
	tlsSkipVerifyDefault  = false
	tlsSkipVerifyEnv      = "REDIS_TLS_SKIP_VERIFY"
)

// RedisTLSCert register flag for database pool size
func SkipVerifyFlag(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Skip service certificate verification
Environment variable: %q`, tlsSkipVerifyEnv)
	f.Bool(tlsSkipVerifyFlag, tlsSkipVerifyDefault, desc)
	_ = viper.BindPFlag(TLSSkipVerifyViperKey, f.Lookup(tlsSkipVerifyFlag))
}

// Register Redis flags
func Flags(f *pflag.FlagSet) {
	URL(f)
	UsernameFlag(f)
	DatabaseFlag(f)
	PasswordFlag(f)
	TLSEnableFlag(f)
	TLSCertFlag(f)
	TLSKeyFlag(f)
	TLSCAFlag(f)
	SkipVerifyFlag(f)
}

type Config struct {
	Host       string
	Port       string
	Expiration int
	User       string
	Password   string
	Database   int
	TLS        *tls.Option
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		Host:       vipr.GetString(HostViperKey),
		Port:       vipr.GetString(PortViperKey),
		User:       vipr.GetString(UsernameViperKey),
		Password:   vipr.GetString(PasswordViperKey),
		Database:   vipr.GetInt(DatabaseViperKey),
		Expiration: int(2 * time.Minute),
	}

	if vipr.GetBool(TLSEnableViperKey) {
		cfg.TLS = &tls.Option{
			ServerName: cfg.Host,
		}

		if vipr.GetString(TLSCertViperKey) != "" {
			cfg.TLS.Certificates = []*certificate.KeyPair{
				{
					Cert: []byte(vipr.GetString(TLSCertViperKey)),
					Key:  []byte(vipr.GetString(TLSKeyViperKey)),
				},
			}
		}

		if vipr.GetString(TLSCAViperKey) != "" {
			cfg.TLS.CAs = [][]byte{
				[]byte(vipr.GetString(TLSCAViperKey)),
			}
		}

		if vipr.GetBool(TLSSkipVerifyViperKey) {
			cfg.TLS.InsecureSkipVerify = true
		}
	}

	return cfg
}

func (c *Config) URL() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}
