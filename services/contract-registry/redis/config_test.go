package redis

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitFlags(flgs)

	config := Config()

	// Check the default values
	assert.Equal(t, redisURIDefault, config.URI, "Wrong default for redis URI value")
	assert.Equal(t, redisMaxIdleConnDefault, config.MaxIdle, "Wrong default for redis maxIdleConn")
	assert.Equal(t, redisMaxActiveConnDefault, config.MaxActive, "Wrong default for MaxActiveConn")
	assert.Equal(t, redisMaxConnLifetimeDefault, config.MaxConnLifetime, "Wrong default for redis MaxConnLifetime")
	assert.Equal(t, redisIdleTimeoutDefault, config.IdleTimeout, "Wrong default for redis IdleTimeout")
	assert.Equal(t, redisWaitDefault, config.Wait, "Wrong default for redis Wait")

	// Custom values to pass either as flags or env
	customURI := "abcd:6788"
	customMaxIdleConn := "1679536"
	customMaxActiveConn := "68609800"
	customMaxConnLifetime := "1s"
	customIdleConnTimeout := "1s"
	customWait := "true"

	// Pass the custom values as environment variables
	_ = os.Setenv(redisURIEnv, customURI)
	_ = os.Setenv(redisMaxIdleConnEnv, customMaxIdleConn)
	_ = os.Setenv(redisMaxActiveConnEnv, customMaxActiveConn)
	_ = os.Setenv(redisMaxConnLifetimeEnv, customMaxConnLifetime)
	_ = os.Setenv(redisIdleTimeoutEnv, customIdleConnTimeout)
	_ = os.Setenv(redisWaitEnv, customWait)

	// Reset the config object with the new viper fields
	config = Config()

	// Check the default values
	assert.Equal(t, customURI, config.URI, "Wrong value for redis URI value")
	assert.Equal(t, 1679536, config.MaxIdle, "Wrong value for redis maxIdleConn")
	assert.Equal(t, 68609800, config.MaxActive, "Wrong value for MaxActiveConn")
	assert.Equal(t, time.Duration(1)*time.Second, config.MaxConnLifetime, "Wrong value for redis MaxConnLifetime")
	assert.Equal(t, time.Duration(1)*time.Second, config.IdleTimeout, "Wrong value for redis IdleTimeout")
	assert.Equal(t, true, config.Wait, "Wrong value for redis Wait")

}
