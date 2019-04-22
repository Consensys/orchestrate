package postgres

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const DEFAULT = "default"

func TestDBUser(t *testing.T) {
	name := "db.user"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBUser(flgs, DEFAULT)

	expected := DEFAULT
	assert.Equal(t, expected, viper.GetString(name), "Default db user should be %q but got %q", expected, viper.GetString(name))

	os.Setenv("DB_USER", "env-user")
	expected = "env-user"
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db user should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-user=flag-user",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = "flag-user"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db user should be %q but got %q", expected, viper.GetString(name))
}

func TestDBPassword(t *testing.T) {
	name := "db.password"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPassword(flgs, DEFAULT)

	expected := DEFAULT
	assert.Equal(t, expected, viper.GetString(name), "Default db password should be %q but got %q", expected, viper.GetString(name))

	os.Setenv("DB_PASSWORD", "env-password")
	expected = "env-password"
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db password should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-password=flag-password",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = "flag-password"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db password should be %q but got %q", expected, viper.GetString(name))
}

func TestDBDatabase(t *testing.T) {
	name := "db.database"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBDatabase(flgs, DEFAULT)

	expected := DEFAULT
	assert.Equal(t, expected, viper.GetString(name), "Default db database should be %q but got %q", expected, viper.GetString(name))

	os.Setenv("DB_DATABASE", "env-database")
	expected = "env-database"
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db database should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-database=flag-database",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = "flag-database"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db database should be %q but got %q", expected, viper.GetString(name))
}

func TestDBHost(t *testing.T) {
	name := "db.host"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBHost(flgs, "127.0.0.1")

	expected := "127.0.0.1"
	assert.Equal(t, expected, viper.GetString(name), "Default db host should be %q but got %q", expected, viper.GetString(name))

	os.Setenv("DB_HOST", "127.1.1.1")
	expected = "127.1.1.1"
	assert.Equal(t, expected, viper.GetString(name), "After setting env var db host should be %q but got %q", expected, viper.GetString(name))

	args := []string{
		"--db-host=localhost",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = "localhost"
	assert.Equal(t, expected, viper.GetString(name), "After setting flag db host should be %q but got %q", expected, viper.GetString(name))
}

func TestDBPort(t *testing.T) {
	name := "db.port"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	DBPort(flgs, 5432)

	expected := 5432
	assert.Equal(t, expected, viper.GetInt(name), "Default db port should be %v but got %v", expected, viper.GetInt(name))

	os.Setenv("DB_PORT", "5433")
	expected = 5433
	assert.Equal(t, expected, viper.GetInt(name), "After setting env var db port should be %v but got %v", expected, viper.GetInt(name))

	args := []string{
		"--db-port=5442",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = 5442
	assert.Equal(t, expected, viper.GetInt(name), "After setting flag db port should be %v but got %v", expected, viper.GetInt(name))
}
