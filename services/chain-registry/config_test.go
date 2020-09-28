package chainregistry

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitRegistry(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cmdFlags(f)

	var expected []string
	assert.Equal(t, expected, viper.GetStringSlice(InitViperKey), "Default")

	_ = os.Setenv(initEnv, "test1 test2")
	assert.Equal(t, []string{"test1", "test2"}, viper.GetStringSlice(InitViperKey), "From Environment Variable")
	_ = os.Unsetenv(initEnv)
}

func TestProxyCacheTTLRegistryDefault(t *testing.T) {
	flgs := pflag.NewFlagSet("test-chain-registry-1", pflag.ContinueOnError)
	cmdFlags(flgs)

	assert.Equal(t, 0, viper.GetInt(CacheTTLViperKey))

	cfg := NewConfig(viper.New())
	assert.Nil(t, cfg.ProxyCacheTTL)
}

func TestProxyCacheTTLRegistryENV(t *testing.T) {
	flgs := pflag.NewFlagSet("test-chain-registry-1", pflag.ContinueOnError)
	cmdFlags(flgs)

	expected := 1000
	_ = os.Setenv(cacheTTLEnv, "1000")
	assert.Equal(t, expected, viper.GetInt(CacheTTLViperKey))
	cfg := NewConfig(viper.New())
	assert.Equal(t, int64(expected), cfg.ProxyCacheTTL.Milliseconds())
	_ = os.Unsetenv(cacheTTLEnv)
}

func TestProxyCacheTTLRegistryFlag(t *testing.T) {
	flgs := pflag.NewFlagSet("test-chain-registry-2", pflag.ContinueOnError)
	cmdFlags(flgs)
	args := []string{fmt.Sprintf("--%s=%d", cacheTTLFlag, 1000)}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, 1000, viper.GetInt(CacheTTLViperKey))
	expected, _ := time.ParseDuration("1s")
	cfg := NewConfig(viper.New())
	assert.Equal(t, expected.Milliseconds(), cfg.ProxyCacheTTL.Milliseconds())
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
}
