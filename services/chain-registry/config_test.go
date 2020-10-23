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

	assert.Equal(t, time.Duration(0), viper.GetDuration(CacheTTLViperKey))

	cfg := NewConfig(viper.New())
	assert.Nil(t, cfg.ProxyCacheTTL)
}

func TestProxyCacheTTLRegistryENV(t *testing.T) {
	flgs := pflag.NewFlagSet("test-chain-registry-1", pflag.ContinueOnError)
	cmdFlags(flgs)

	_ = os.Setenv(cacheTTLEnv, "1s")
	assert.Equal(t, 1*time.Second, viper.GetDuration(CacheTTLViperKey))
	cfg := NewConfig(viper.New())
	assert.Equal(t, int64(1000), cfg.ProxyCacheTTL.Milliseconds())
	_ = os.Unsetenv(cacheTTLEnv)
}

func TestProxyCacheTTLRegistryFlag(t *testing.T) {
	flgs := pflag.NewFlagSet("test-chain-registry-2", pflag.ContinueOnError)
	cmdFlags(flgs)
	args := []string{fmt.Sprintf("--%s=%s", cacheTTLFlag, "1s")}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, 1*time.Second, viper.GetDuration(CacheTTLViperKey))
	expected, _ := time.ParseDuration("1s")
	cfg := NewConfig(viper.New())
	assert.Equal(t, expected.Milliseconds(), cfg.ProxyCacheTTL.Milliseconds())
}

func TestMaxIdleConnsPerHostDefault(t *testing.T) {
	flgs := pflag.NewFlagSet("test-max-idle-connections-per-host-1", pflag.ContinueOnError)
	cmdFlags(flgs)

	assert.Equal(t, 50, viper.GetInt(MaxIdleConnsPerHostViperKey))

	cfg := NewConfig(viper.New())
	assert.Equal(t, 50, cfg.ServersTransport.MaxIdleConnsPerHost)
}

func TestMaxIdleConnsPerHostENV(t *testing.T) {
	flgs := pflag.NewFlagSet("test-max-idle-connections-per-host-1", pflag.ContinueOnError)
	cmdFlags(flgs)

	_ = os.Setenv(maxIdleConnsPerHostEnv, "2")
	assert.Equal(t, 2, viper.GetInt(MaxIdleConnsPerHostViperKey))
	cfg := NewConfig(viper.New())
	assert.Equal(t, 2, cfg.ServersTransport.MaxIdleConnsPerHost)
	_ = os.Unsetenv(maxIdleConnsPerHostEnv)
}

func TestMaxIdleConnsPerHostFlag(t *testing.T) {
	flgs := pflag.NewFlagSet("test-max-idle-connections-per-host-1", pflag.ContinueOnError)
	cmdFlags(flgs)
	args := []string{fmt.Sprintf("--%s=%s", maxIdleConnsPerHostFlag, "3")}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, 3, viper.GetInt(MaxIdleConnsPerHostViperKey))
	cfg := NewConfig(viper.New())
	assert.Equal(t, 3, cfg.ServersTransport.MaxIdleConnsPerHost)
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
}
