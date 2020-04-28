package postgres

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tls/certificate"
)

func init() {
	viper.SetDefault(DBUserViperKey, dbUserDefault)
	_ = viper.BindEnv(DBUserViperKey, dbUserEnv)
	viper.SetDefault(DBPasswordViperKey, dbPasswordDefault)
	_ = viper.BindEnv(DBPasswordViperKey, dbPasswordEnv)
	viper.SetDefault(DBDatabaseViperKey, dbDatabaseDefault)
	_ = viper.BindEnv(DBDatabaseViperKey, dbDatabaseEnv)
	viper.SetDefault(DBHostViperKey, dbHostDefault)
	_ = viper.BindEnv(DBHostViperKey, dbHostEnv)
	viper.SetDefault(DBPortViperKey, dbPortDefault)
	_ = viper.BindEnv(DBPortViperKey, dbPortEnv)
	viper.SetDefault(DBPoolSizeViperKey, dbPoolSizeDefault)
	_ = viper.BindEnv(DBPoolSizeViperKey, dbPoolSizeEnv)
	viper.SetDefault(DBTLSCertViperKey, dbTLSCertDefault)
	_ = viper.BindEnv(DBTLSCertViperKey, dbTLSCertEnv)
	viper.SetDefault(DBTLSKeyViperKey, dbTLSKeyDefault)
	_ = viper.BindEnv(DBTLSKeyViperKey, dbTLSKeyEnv)
	viper.SetDefault(DBTLSCAViperKey, dbTLSCADefault)
	_ = viper.BindEnv(DBTLSCAViperKey, dbTLSCAEnv)
}

// PGFlags register flags for Postgres database
func PGFlags(f *pflag.FlagSet) {
	DBUser(f)
	DBPassword(f)
	DBDatabase(f)
	DBHost(f)
	DBPort(f)
	DBPoolSize(f)
	DBTLSCert(f)
	DBTLSKey(f)
	DBTLSCA(f)
}

const (
	dbUserFlag     = "db-user"
	DBUserViperKey = "db.user"
	dbUserDefault  = "postgres"
	dbUserEnv      = "DB_USER"
)

// DBUser register flag for db user
func DBUser(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database User.
Environment variable: %q`, dbUserEnv)
	f.String(dbUserFlag, dbUserDefault, desc)
	_ = viper.BindPFlag(DBUserViperKey, f.Lookup(dbUserFlag))
}

const (
	dbPasswordFlag     = "db-password"
	DBPasswordViperKey = "db.password"
	dbPasswordDefault  = "postgres"
	dbPasswordEnv      = "DB_PASSWORD"
)

// DBPassword register flag for db password
func DBPassword(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database User password
Environment variable: %q`, dbPasswordEnv)
	f.String(dbPasswordFlag, dbPasswordDefault, desc)
	_ = viper.BindPFlag(DBPasswordViperKey, f.Lookup(dbPasswordFlag))
}

const (
	dbDatabaseFlag     = "db-database"
	DBDatabaseViperKey = "db.database"
	dbDatabaseDefault  = "postgres"
	dbDatabaseEnv      = "DB_DATABASE"
)

// DBDatabase register flag for db database name
func DBDatabase(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Target Database name
Environment variable: %q`, dbDatabaseEnv)
	f.String(dbDatabaseFlag, dbDatabaseDefault, desc)
	_ = viper.BindPFlag(DBDatabaseViperKey, f.Lookup(dbDatabaseFlag))
}

const (
	dbHostFlag     = "db-host"
	DBHostViperKey = "db.host"
	dbHostDefault  = "127.0.0.1"
	dbHostEnv      = "DB_HOST"
)

// DBHost register flag for database host
func DBHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database host
Environment variable: %q`, dbHostEnv)
	f.String(dbHostFlag, dbHostDefault, desc)
	_ = viper.BindPFlag(DBHostViperKey, f.Lookup(dbHostFlag))
}

const (
	dbPortFlag     = "db-port"
	DBPortViperKey = "db.port"
	dbPortDefault  = 5432
	dbPortEnv      = "DB_PORT"
)

// DBPort register flag for database port
func DBPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database port
Environment variable: %q`, dbPortEnv)
	f.Int(dbPortFlag, dbPortDefault, desc)
	_ = viper.BindPFlag(DBPortViperKey, f.Lookup(dbPortFlag))
}

const (
	dbPoolSizeFlag     = "db-poolsize"
	DBPoolSizeViperKey = "db.poolsize"
	dbPoolSizeDefault  = 0
	dbPoolSizeEnv      = "DB_POOLSIZE"
)

// DBPoolSize register flag for database pool size
func DBPoolSize(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of connections on database
Environment variable: %q`, dbPoolSizeEnv)
	f.Int(dbPoolSizeFlag, dbPoolSizeDefault, desc)
	_ = viper.BindPFlag(DBPoolSizeViperKey, f.Lookup(dbPoolSizeFlag))
}

const (
	dbTLSCertFlag     = "db-tls-cert"
	DBTLSCertViperKey = "db.tls.cert"
	dbTLSCertDefault  = ""
	dbTLSCertEnv      = "DB_TLS_CERT"
)

// DBTLSCert register flag for TLS certificate used to connect to the database
func DBTLSCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Certificate to connect to database
Environment variable: %q`, dbTLSCertEnv)
	f.String(dbTLSCertFlag, dbTLSCertDefault, desc)
	_ = viper.BindPFlag(DBTLSCertViperKey, f.Lookup(dbTLSCertFlag))
}

const (
	dbTLSKeyFlag     = "db-tls-key"
	DBTLSKeyViperKey = "db.tls.key"
	dbTLSKeyDefault  = ""
	dbTLSKeyEnv      = "DB_TLS_KEY"
)

// DBTLSKey register flag for database TLS private key
func DBTLSKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`TLS Private Key to connect to database
Environment variable: %q`, dbTLSKeyEnv)
	f.String(dbTLSKeyFlag, dbTLSKeyDefault, desc)
	_ = viper.BindPFlag(DBTLSKeyViperKey, f.Lookup(dbTLSKeyFlag))
}

const (
	dbTLSCAFlag     = "db-tls-ca"
	DBTLSCAViperKey = "db.tls.ca"
	dbTLSCADefault  = ""
	dbTLSCAEnv      = "DB_TLS_CA"
)

// DBTLSCert register flag for database pool size
func DBTLSCA(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Trusted Certificate Authority
Environment variable: %q`, dbTLSCAEnv)
	f.String(dbTLSCAFlag, dbTLSCADefault, desc)
	_ = viper.BindPFlag(DBTLSCAViperKey, f.Lookup(dbTLSCAFlag))
}

type Config struct {
	Addr            string
	User            string
	Password        string
	Database        string
	PoolSize        int
	TLS             *tls.Option
	ApplicationName string
}

func (cfg *Config) PGOptions() (*pg.Options, error) {
	opt := &pg.Options{
		Addr:            cfg.Addr,
		User:            cfg.User,
		Password:        cfg.Password,
		Database:        cfg.Database,
		PoolSize:        cfg.PoolSize,
		ApplicationName: cfg.ApplicationName,
	}

	if cfg.TLS != nil {
		tlsConfig, err := tls.NewConfig(cfg.TLS)
		if err != nil {
			return nil, err
		}

		opt.TLSConfig = tlsConfig
	}

	return opt, nil
}

// NewONewConfigNewConfigptions creates new postgres options
func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		Addr:     fmt.Sprintf("%v:%v", vipr.GetString(DBHostViperKey), vipr.GetString(DBPortViperKey)),
		User:     vipr.GetString(DBUserViperKey),
		Password: vipr.GetString(DBPasswordViperKey),
		Database: vipr.GetString(DBDatabaseViperKey),
		PoolSize: vipr.GetInt(DBPoolSizeViperKey),
	}

	if vipr.GetString(DBTLSCertViperKey) != "" {
		cfg.TLS = &tls.Option{
			ServerName: vipr.GetString(DBHostViperKey),
			Certificates: []*certificate.KeyPair{
				&certificate.KeyPair{
					Cert: []byte(vipr.GetString(DBTLSCertViperKey)),
					Key:  []byte(vipr.GetString(DBTLSKeyViperKey)),
				},
			},
		}

		if vipr.GetString(DBTLSCAViperKey) != "" {
			cfg.TLS.CAs = [][]byte{
				[]byte(vipr.GetString(DBTLSCAViperKey)),
			}
		}
	}

	return cfg
}

func Copy(opts *pg.Options) *pg.Options {
	if opts == nil {
		return nil
	}
	o := (*opts)
	return &o
}
