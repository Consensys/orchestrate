package postgres

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	tlstestutils "github.com/ConsenSys/orchestrate/pkg/toolkit/tls/testutils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	PGFlags(f)
}

func TestDBUser(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBUser(flgs)

	expected := dbUserDefault // nolint:goconst // reason
	assert.Equal(t, expected, viper.GetString(DBUserViperKey), "Default db user should be %q but got %q", expected, viper.GetString(DBUserViperKey))

	expected = "env-user"
	_ = os.Setenv(dbUserEnv, expected)
	assert.Equal(t, expected, viper.GetString(DBUserViperKey), "After setting env var db user should be %q but got %q", expected, viper.GetString(DBUserViperKey))

	args := []string{
		fmt.Sprintf("--%s=%s", dbUserFlag, "flag-user"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-user"
	assert.Equal(t, expected, viper.GetString(DBUserViperKey), "After setting flag db user should be %q but got %q", expected, viper.GetString(DBUserViperKey))
}

func TestDBPassword(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPassword(flgs)

	expected := dbPasswordDefault
	assert.Equal(t, expected, viper.GetString(DBPasswordViperKey), "Default db password should be %q but got %q", expected, viper.GetString(DBPasswordViperKey))

	expected = "env-password"
	_ = os.Setenv(dbPasswordEnv, expected)
	assert.Equal(t, expected, viper.GetString(DBPasswordViperKey), "After setting env var db password should be %q but got %q", expected, viper.GetString(DBPasswordViperKey))

	args := []string{
		fmt.Sprintf("--%s=%s", dbPasswordFlag, "flag-password"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-password"
	assert.Equal(t, expected, viper.GetString(DBPasswordViperKey), "After setting flag db password should be %q but got %q", expected, viper.GetString(DBPasswordViperKey))
}

func TestDBDatabase(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBDatabase(flgs)

	expected := dbDatabaseDefault
	assert.Equal(t, expected, viper.GetString(DBDatabaseViperKey), "Default db database should be %q but got %q", expected, viper.GetString(DBDatabaseViperKey))

	expected = "env-database"
	_ = os.Setenv(dbDatabaseEnv, expected)
	assert.Equal(t, expected, viper.GetString(DBDatabaseViperKey), "After setting env var db database should be %q but got %q", expected, viper.GetString(DBDatabaseViperKey))

	args := []string{
		fmt.Sprintf("--%s=%s", dbDatabaseFlag, "flag-database"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-database"
	assert.Equal(t, expected, viper.GetString(DBDatabaseViperKey), "After setting flag db database should be %q but got %q", expected, viper.GetString(DBDatabaseViperKey))
}

func TestDBHost(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBHost(flgs)

	expected := dbHostDefault
	_ = os.Setenv(dbHostEnv, expected)
	assert.Equal(t, expected, viper.GetString(DBHostViperKey), "After setting env var db host should be %q but got %q", expected, viper.GetString(DBHostViperKey))

	args := []string{
		fmt.Sprintf("--%s=%s", dbHostFlag, "localhost"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "localhost"
	assert.Equal(t, expected, viper.GetString(DBHostViperKey), "After setting flag db host should be %q but got %q", expected, viper.GetString(DBHostViperKey))
}

func TestDBPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPort(flgs)

	expected := dbPortDefault
	assert.Equal(t, expected, viper.GetInt(DBPortViperKey), "Default db port should be %v but got %v", expected, viper.GetInt(DBPortViperKey))

	expected = 5433
	_ = os.Setenv(dbPortEnv, strconv.FormatInt(int64(expected), 10))
	assert.Equal(t, expected, viper.GetInt(DBPortViperKey), "After setting env var db port should be %v but got %v", expected, viper.GetInt(DBPortViperKey))

	args := []string{
		fmt.Sprintf("--%s=%s", dbPortFlag, "5442"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 5442
	assert.Equal(t, expected, viper.GetInt(DBPortViperKey), "After setting flag db port should be %v but got %v", expected, viper.GetInt(DBPortViperKey))
}

func TestDBPoolSize(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPoolSize(flgs)

	expected := dbPoolSizeDefault
	assert.Equal(t, expected, viper.GetInt(DBPoolSizeViperKey), "Default db pool size should be %v but got %v", expected, viper.GetInt(DBPoolSizeViperKey))

	expected = 1
	_ = os.Setenv(dbPoolSizeEnv, strconv.FormatInt(int64(expected), 10))
	assert.Equal(t, expected, viper.GetInt(DBPoolSizeViperKey), "After setting env var db port should be %v but got %v", expected, viper.GetInt(DBPoolSizeViperKey))

	args := []string{
		fmt.Sprintf("--%s=%s", dbPoolSizeFlag, "2"),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 2
	assert.Equal(t, expected, viper.GetInt(DBPoolSizeViperKey), "After setting flag db poolsize should be %v but got %v", expected, viper.GetInt(DBPoolSizeViperKey))
}

func TestDBKeepAliveInterval(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBKeepAliveInterval(flgs)

	expected := dbKeepAliveDefault
	assert.Equal(t, expected, viper.GetDuration(DBKeepAliveKey), "Default db keep alive should be %v but got %v", expected, viper.GetDuration(DBKeepAliveKey))

	expected = time.Second
	_ = os.Setenv(dbKeepAliveEnv, expected.String())
	assert.Equal(t, expected, viper.GetDuration(DBKeepAliveKey), "After setting env var keep alive should be %v but got %v", expected, viper.GetDuration(DBKeepAliveKey))

	expected = 2 * time.Second
	args := []string{
		fmt.Sprintf("--%s=%s", dbKeepAliveFlag, expected.String()),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, expected, viper.GetDuration(DBKeepAliveKey), "After setting flag db keep aliave should be %v but got %v", expected, viper.GetDuration(DBKeepAliveKey))
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
