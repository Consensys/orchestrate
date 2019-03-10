package pg

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DBFlags register database flags
func DBFlags(f *pflag.FlagSet) {
	DBUser(f)
	DBPassword(f)
	DBDatabase(f)
	DBHost(f)
	DBPort(f)
	DBPoolSize(f)
}

var (
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
	viper.BindPFlag(dbUserViperKey, f.Lookup(dbUserFlag))
	viper.BindEnv(dbUserViperKey, dbUserEnv)
}

var (
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
	viper.BindPFlag(dbPasswordViperKey, f.Lookup(dbPasswordFlag))
	viper.BindEnv(dbPasswordViperKey, dbPasswordEnv)
}

var (
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
	viper.BindPFlag(dbDatabaseViperKey, f.Lookup(dbDatabaseFlag))
	viper.BindEnv(dbDatabaseViperKey, dbDatabaseEnv)
}

var (
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
	viper.BindPFlag(dbHostViperKey, f.Lookup(dbHostFlag))
	viper.BindEnv(dbHostViperKey, dbHostEnv)
}

var (
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
	viper.BindPFlag(dbPortViperKey, f.Lookup(dbPortFlag))
	viper.BindEnv(dbPortViperKey, dbPortEnv)
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
	viper.BindPFlag(dbPoolSizeViperKey, f.Lookup(dbPoolSizeFlag))
	viper.BindEnv(dbPoolSizeViperKey, dbPoolSizeEnv)
}
