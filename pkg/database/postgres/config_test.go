package postgres

import (
	"os"
	"strconv"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tlstestutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tls/testutils"
)

func TestPGFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	PGFlags(f)
}

func TestDBUser(t *testing.T) {
	name := "db.user"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBUser(flgs)

	expected := "postgres" //nolint:goconst // reason
	assert.Equal(t, expected, viper.GetString(name), "Default db user should be %q but got %q", expected, viper.GetString(name))

	expected = "env-user"
	_ = os.Setenv("DB_USER", expected)
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db user should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-user=flag-user",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-user"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db user should be %q but got %q", expected, viper.GetString(name))
}

func TestDBPassword(t *testing.T) {
	name := "db.password"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPassword(flgs)

	expected := "postgres"
	assert.Equal(t, expected, viper.GetString(name), "Default db password should be %q but got %q", expected, viper.GetString(name))

	expected = "env-password"
	_ = os.Setenv("DB_PASSWORD", expected)
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db password should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-password=flag-password",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-password"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db password should be %q but got %q", expected, viper.GetString(name))
}

func TestDBDatabase(t *testing.T) {
	name := "db.database"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBDatabase(flgs)

	expected := "postgres"
	assert.Equal(t, expected, viper.GetString(name), "Default db database should be %q but got %q", expected, viper.GetString(name))

	expected = "env-database"
	_ = os.Setenv("DB_DATABASE", expected)
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db database should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-database=flag-database",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-database"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db database should be %q but got %q", expected, viper.GetString(name))
}

func TestDBHost(t *testing.T) {
	name := "db.host"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBHost(flgs)

	expected := "127.1.1.1"
	_ = os.Setenv("DB_HOST", expected)
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db host should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-host=localhost",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "localhost"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db host should be %q but got %q", expected, viper.GetString(name))
}

func TestDBPort(t *testing.T) {
	name := "db.port"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPort(flgs)

	expected := 5432
	assert.Equal(t, expected, viper.GetInt(name), "Default db port should be %v but got %v", expected, viper.GetInt(name))

	expected = 5433
	_ = os.Setenv("DB_PORT", strconv.FormatInt(int64(expected), 10))
	assert.Equal(t, expected, viper.GetInt(name), "After setting env var db port should be %v but got %v", expected, viper.GetInt(name))

	args := []string{
		"--db-port=5442",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 5442
	assert.Equal(t, expected, viper.GetInt(name), "After setting flag db port should be %v but got %v", expected, viper.GetInt(name))
}

func TestDBPoolSize(t *testing.T) {
	name := "db.poolsize"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPoolSize(flgs)

	expected := 0
	assert.Equal(t, expected, viper.GetInt(name), "Default db pool size should be %v but got %v", expected, viper.GetInt(name))

	expected = 1
	_ = os.Setenv("DB_POOLSIZE", strconv.FormatInt(int64(expected), 10))
	assert.Equal(t, expected, viper.GetInt(name), "After setting env var db port should be %v but got %v", expected, viper.GetInt(name))

	args := []string{
		"--db-poolsize=2",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 2
	assert.Equal(t, expected, viper.GetInt(name), "After setting flag db poolsize should be %v but got %v", expected, viper.GetInt(name))
}

func TestNewOptions(t *testing.T) {
	vipr := viper.New()
	vipr.Set(DBHostViperKey, "host-test")
	vipr.Set(DBPortViperKey, 5432)
	vipr.Set(DBDatabaseViperKey, "db-test")
	vipr.Set(DBUserViperKey, "user-test")
	vipr.Set(DBPasswordViperKey, "password")
	vipr.Set(DBTLSSSLModeViperKey, "verify-full")
	vipr.Set(DBTLSCertViperKey, tlstestutils.OneLineRSACertPEMA)
	vipr.Set(DBTLSKeyViperKey, tlstestutils.OneLineRSAKeyPEMA)
	vipr.Set(DBTLSCAViperKey, tlstestutils.OneLineRSACertPEMB)

	cfg := NewConfig(vipr)
	opts, err := cfg.PGOptions()
	require.NoError(t, err)
	assert.Equal(t, "host-test:5432", opts.Addr)
	assert.Equal(t, "db-test", opts.Database)
	assert.Equal(t, "user-test", opts.User)
	assert.Equal(t, "password", opts.Password)
}
