package proxy

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestProxyCacheTTLRegistryDefault(t *testing.T) {
	Flags(pflag.NewFlagSet("test-proxy-default", pflag.ContinueOnError))

	assert.Equal(t, time.Duration(0), viper.GetDuration(CacheTTLViperKey))

	cfg := NewConfig()
	assert.Nil(t, cfg.ProxyCacheTTL)
}

func TestProxyCacheTTLRegistryENV(t *testing.T) {
	Flags(pflag.NewFlagSet("test-proxy-env", pflag.ContinueOnError))

	_ = os.Setenv(cacheTTLEnv, "1s")
	assert.Equal(t, 1*time.Second, viper.GetDuration(CacheTTLViperKey))
	cfg := NewConfig()
	assert.Equal(t, int64(1000), cfg.ProxyCacheTTL.Milliseconds())
	_ = os.Unsetenv(cacheTTLEnv)
}

func TestProxyCacheTTLRegistryFlag(t *testing.T) {
	flgs := pflag.NewFlagSet("test-proxy-flag", pflag.ContinueOnError)
	Flags(flgs)

	err := flgs.Parse([]string{fmt.Sprintf("--%s=%s", cacheTTLFlag, "1s")})
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, 1*time.Second, viper.GetDuration(CacheTTLViperKey))
	expected, _ := time.ParseDuration("1s")

	cfg := NewConfig()
	assert.Equal(t, expected.Milliseconds(), cfg.ProxyCacheTTL.Milliseconds())
}

func TestMaxIdleConnsPerHostDefault(t *testing.T) {
	Flags(pflag.NewFlagSet("test-max-idle-connections-per-host-default", pflag.ContinueOnError))

	assert.Equal(t, 50, viper.GetInt(MaxIdleConnsPerHostViperKey))

	cfg := NewConfig()
	assert.Equal(t, 50, cfg.ServersTransport.MaxIdleConnsPerHost)
}

func TestMaxIdleConnsPerHostENV(t *testing.T) {
	Flags(pflag.NewFlagSet("test-max-idle-connections-per-host-env", pflag.ContinueOnError))

	_ = os.Setenv(maxIdleConnsPerHostEnv, "2")
	assert.Equal(t, 2, viper.GetInt(MaxIdleConnsPerHostViperKey))

	cfg := NewConfig()
	assert.Equal(t, 2, cfg.ServersTransport.MaxIdleConnsPerHost)
	_ = os.Unsetenv(maxIdleConnsPerHostEnv)
}

func TestMaxIdleConnsPerHostFlag(t *testing.T) {
	flgs := pflag.NewFlagSet("test-max-idle-connections-per-host-flag", pflag.ContinueOnError)
	Flags(flgs)

	err := flgs.Parse([]string{fmt.Sprintf("--%s=%s", maxIdleConnsPerHostFlag, "3")})
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, 3, viper.GetInt(MaxIdleConnsPerHostViperKey))

	cfg := NewConfig()
	assert.Equal(t, 3, cfg.ServersTransport.MaxIdleConnsPerHost)
}
