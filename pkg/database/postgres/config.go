package postgres

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(dbUserViperKey, dbUserDefault)
	_ = viper.BindEnv(dbUserViperKey, dbUserEnv)
	viper.SetDefault(dbPasswordViperKey, dbPasswordDefault)
	_ = viper.BindEnv(dbPasswordViperKey, dbPasswordEnv)
	viper.SetDefault(dbDatabaseViperKey, dbDatabaseDefault)
	_ = viper.BindEnv(dbDatabaseViperKey, dbDatabaseEnv)
	viper.SetDefault(dbHostViperKey, dbHostDefault)
	_ = viper.BindEnv(dbHostViperKey, dbHostEnv)
	viper.SetDefault(dbPortViperKey, dbPortDefault)
	_ = viper.BindEnv(dbPortViperKey, dbPortEnv)
	viper.SetDefault(dbPoolSizeViperKey, dbPoolSizeDefault)
	_ = viper.BindEnv(dbPoolSizeViperKey, dbPoolSizeEnv)
}

// PGFlags register flags for Postgres database
func PGFlags(f *pflag.FlagSet) {
	DBUser(f)
	DBPassword(f)
	DBDatabase(f)
	DBHost(f)
	DBPort(f)
	DBPoolSize(f)
}

const (
	dbUserFlag     = "db-user"
	dbUserViperKey = "db.user"
	dbUserDefault  = "postgres"
	dbUserEnv      = "DB_USER"
)

// DBUser register flag for db user
func DBUser(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database User.
Environment variable: %q`, dbUserEnv)
	f.String(dbUserFlag, dbUserDefault, desc)
	_ = viper.BindPFlag(dbUserViperKey, f.Lookup(dbUserFlag))
}

const (
	dbPasswordFlag     = "db-password"
	dbPasswordViperKey = "db.password"
	dbPasswordDefault  = "postgres"
	dbPasswordEnv      = "DB_PASSWORD"
)

// DBPassword register flag for db password
func DBPassword(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database User password
Environment variable: %q`, dbPasswordEnv)
	f.String(dbPasswordFlag, dbPasswordDefault, desc)
	_ = viper.BindPFlag(dbPasswordViperKey, f.Lookup(dbPasswordFlag))
}

const (
	dbDatabaseFlag     = "db-database"
	dbDatabaseViperKey = "db.database"
	dbDatabaseDefault  = "postgres"
	dbDatabaseEnv      = "DB_DATABASE"
)

// DBDatabase register flag for db database name
func DBDatabase(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Target Database name
Environment variable: %q`, dbDatabaseEnv)
	f.String(dbDatabaseFlag, dbDatabaseDefault, desc)
	_ = viper.BindPFlag(dbDatabaseViperKey, f.Lookup(dbDatabaseFlag))
}

const (
	dbHostFlag     = "db-host"
	dbHostViperKey = "db.host"
	dbHostDefault  = "127.0.0.1"
	dbHostEnv      = "DB_HOST"
)

// DBHost register flag for database host
func DBHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database host
Environment variable: %q`, dbHostEnv)
	f.String(dbHostFlag, dbHostDefault, desc)
	_ = viper.BindPFlag(dbHostViperKey, f.Lookup(dbHostFlag))
}

const (
	dbPortFlag     = "db-port"
	dbPortViperKey = "db.port"
	dbPortDefault  = 5432
	dbPortEnv      = "DB_PORT"
)

// DBPort register flag for database port
func DBPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Database port
Environment variable: %q`, dbPortEnv)
	f.Int(dbPortFlag, dbPortDefault, desc)
	_ = viper.BindPFlag(dbPortViperKey, f.Lookup(dbPortFlag))
}

const (
	dbPoolSizeFlag     = "db-poolsize"
	dbPoolSizeViperKey = "db.poolsize"
	dbPoolSizeDefault  = 0
	dbPoolSizeEnv      = "DB_POOLSIZE"
)

// DBPoolSize register flag for database pool size
func DBPoolSize(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of connections on database
Environment variable: %q`, dbPoolSizeEnv)
	f.Int(dbPoolSizeFlag, dbPoolSizeDefault, desc)
	_ = viper.BindPFlag(dbPoolSizeViperKey, f.Lookup(dbPoolSizeFlag))
}

// NewOptions creates new postgres options
func NewOptions(vipr *viper.Viper) *pg.Options {
	return &pg.Options{
		Addr:     fmt.Sprintf("%v:%v", vipr.GetString(dbHostViperKey), vipr.GetString(dbPortViperKey)),
		User:     vipr.GetString(dbUserViperKey),
		Password: vipr.GetString(dbPasswordViperKey),
		Database: vipr.GetString(dbDatabaseViperKey),
		PoolSize: vipr.GetInt(dbPoolSizeViperKey),
	}
}
