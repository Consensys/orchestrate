package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(dbUserViperKey, dbUserEnv)
	viper.BindEnv(dbPasswordViperKey, dbPasswordEnv)
	viper.BindEnv(dbDatabaseViperKey, dbDatabaseEnv)
	viper.BindEnv(dbHostViperKey, dbHostEnv)
	viper.BindEnv(dbPortViperKey, dbPortEnv)
	viper.BindEnv(dbPoolSizeViperKey, dbPoolSizeEnv)
}

// PGFlags register flags for Postgres database
func PGFlags(f *pflag.FlagSet) {
	DBUser(f, "postgres")
	DBPassword(f, "postgres")
	DBDatabase(f, "postgres")
	DBHost(f, "127.0.0.1")
	DBPort(f, 5432)
	DBPoolSize(f)
}

var (
	dbUserFlag     = "db-user"
	dbUserViperKey = "db.user"
	dbUserEnv      = "DB_USER"
)

// DBUser register flag for db user
func DBUser(f *pflag.FlagSet, defaultUser string) {
	desc := fmt.Sprintf(`Database User.
Environment variable: %q`, dbUserEnv)
	f.String(dbUserFlag, defaultUser, desc)
	viper.BindPFlag(dbUserViperKey, f.Lookup(dbUserFlag))
	viper.SetDefault(dbUserViperKey, defaultUser)
}

var (
	dbPasswordFlag     = "db-password"
	dbPasswordViperKey = "db.password"
	dbPasswordEnv      = "DB_PASSWORD"
)

// DBPassword register flag for db password
func DBPassword(f *pflag.FlagSet, defaultPassword string) {
	desc := fmt.Sprintf(`Database User password
Environment variable: %q`, dbPasswordEnv)
	f.String(dbPasswordFlag, defaultPassword, desc)
	viper.SetDefault(dbPasswordViperKey, defaultPassword)
	viper.BindPFlag(dbPasswordViperKey, f.Lookup(dbPasswordFlag))
}

var (
	dbDatabaseFlag     = "db-database"
	dbDatabaseViperKey = "db.database"
	dbDatabaseEnv      = "DB_DATABASE"
)

// DBDatabase register flag for db database name
func DBDatabase(f *pflag.FlagSet, defaultDatabase string) {
	desc := fmt.Sprintf(`Target Database name
Environment variable: %q`, dbDatabaseEnv)
	f.String(dbDatabaseFlag, defaultDatabase, desc)
	viper.SetDefault(dbDatabaseViperKey, defaultDatabase)
	viper.BindPFlag(dbDatabaseViperKey, f.Lookup(dbDatabaseFlag))
}

var (
	dbHostFlag     = "db-host"
	dbHostViperKey = "db.host"
	dbHostEnv      = "DB_HOST"
)

// DBHost register flag for database host
func DBHost(f *pflag.FlagSet, defaultHost string) {
	desc := fmt.Sprintf(`Database host
Environment variable: %q`, dbHostEnv)
	f.String(dbHostFlag, defaultHost, desc)
	viper.SetDefault(dbHostViperKey, defaultHost)
	viper.BindPFlag(dbHostViperKey, f.Lookup(dbHostFlag))
}

var (
	dbPortFlag     = "db-port"
	dbPortViperKey = "db.port"
	dbPortEnv      = "DB_PORT"
)

// DBPort register flag for database port
func DBPort(f *pflag.FlagSet, defaultPort int) {
	desc := fmt.Sprintf(`Database port
Environment variable: %q`, dbPortEnv)
	f.Int(dbPortFlag, defaultPort, desc)
	viper.SetDefault(dbPortViperKey, defaultPort)
	viper.BindPFlag(dbPortViperKey, f.Lookup(dbPortFlag))
}

var (
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
	viper.SetDefault(dbPoolSizeViperKey, dbPoolSizeDefault)
	viper.BindPFlag(dbPoolSizeViperKey, f.Lookup(dbPoolSizeFlag))
}
