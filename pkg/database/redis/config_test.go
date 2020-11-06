// +build unit

package redis

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRedisHost(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	URL(flgs)
	
	// #1: Default value check
	assert.Equal(t, hostDefault, viper.GetString(HostViperKey), "RedisURL # Default")

	// #2: ENV variable
	hostV := "127.0.0.1"
	_ = os.Setenv(hostEnv, hostV)
	assert.Equal(t, hostV, viper.GetString(HostViperKey), "RedisURL # ENV")
	_ = os.Unsetenv(hostEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s=%v", hostFlag, hostV),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, hostV, viper.GetString(HostViperKey), "RedisURL # CLI")
}

func TestRedisPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	URL(flgs)
	
	// #1: Default value check
	assert.Equal(t, portDefault, viper.GetString(PortViperKey), "RedisPort # Default")

	// #2: ENV variable
	portV := "6378"
	_ = os.Setenv(portEnv, portV)
	assert.Equal(t, portV, viper.GetString(PortViperKey), "RedisPort # ENV")
	_ = os.Unsetenv(portEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s=%v", portFlag, portV),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, portV, viper.GetString(PortViperKey), "RedisPort # CLI")
}

func TestRedisTLSEnable(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	TLSEnableFlag(flgs)
	
	// #1: Default value check
	assert.Equal(t, tlsEnableDefault, viper.GetBool(TLSEnableViperKey), "RedisTLSEnable # Default")

	// #2: ENV variable
	_ = os.Setenv(tlsEnableEnv, "1")
	assert.Equal(t, true, viper.GetBool(TLSEnableViperKey), "RedisTLSEnable # ENV")
	_ = os.Unsetenv(tlsEnableEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s", tlsEnableFlag),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, true, viper.GetBool(TLSEnableViperKey), "RedisTLSEnable # CLI")
}

func TestRedisTLSCA(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	TLSCAFlag(flgs)
	
	// #1: Default value check
	assert.Equal(t, tlsCADefault, viper.GetString(TLSCAViperKey), "RedisTLSCA # Default")

	// #2: ENV variable
	CAValue := "CAValue"
	_ = os.Setenv(tlsCAEnv, CAValue)
	assert.Equal(t, CAValue, viper.GetString(TLSCAViperKey), "RedisTLSCA # ENV")
	_ = os.Unsetenv(tlsCAEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s=%v", tlsCAFlag, CAValue),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, CAValue, viper.GetString(TLSCAViperKey), "RedisTLSCA# CLI")
}

func TestRedisTLSCert(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	TLSCertFlag(flgs)
	
	// #1: Default value check
	assert.Equal(t, tlsCertDefault, viper.GetString(TLSCertViperKey), "RedisTLSCA # Default")

	// #2: ENV variable
	CertValue := "CertValue"
	_ = os.Setenv(tlsCertEnv, CertValue)
	assert.Equal(t, CertValue, viper.GetString(TLSCertViperKey), "RedisTLSCA # ENV")
	_ = os.Unsetenv(tlsCertEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s=%v", tlsCertFlag, CertValue),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, CertValue, viper.GetString(TLSCertViperKey), "RedisTLSCA# CLI")
}

func TestRedisTLSKey(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	TLSKeyFlag(flgs)
	
	// #1: Default value check
	assert.Equal(t, tlsKeyDefault, viper.GetString(TLSKeyViperKey), "RedisTLSCA # Default")

	// #2: ENV variable
	CertKeyValue := "CertKeyValue"
	_ = os.Setenv(tlsKeyEnv, CertKeyValue)
	assert.Equal(t, CertKeyValue, viper.GetString(TLSKeyViperKey), "RedisTLSCA # ENV")
	_ = os.Unsetenv(tlsKeyEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s=%v", tlsKeyFlag, CertKeyValue),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, CertKeyValue, viper.GetString(TLSKeyViperKey), "RedisTLSCA# CLI")
}

func TestRedisTLSSkipVerify(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	SkipVerifyFlag(flgs)
	
	// #1: Default value check
	assert.Equal(t, tlsSkipVerifyDefault, viper.GetBool(TLSSkipVerifyViperKey), "RedisTLSCA # Default")

	// #2: ENV variable
	_ = os.Setenv(tlsSkipVerifyEnv, "1")
	assert.Equal(t, true, viper.GetBool(TLSSkipVerifyViperKey), "RedisTLSCA # ENV")
	_ = os.Unsetenv(tlsSkipVerifyEnv)

	// #2: CLI flag
	args := []string{
		fmt.Sprintf("--%s", tlsSkipVerifyFlag),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, true, viper.GetBool(TLSSkipVerifyViperKey), "RedisTLSCA# CLI")
}

func TestNewConfig_TLS(t *testing.T) {
	vipr := viper.New()
	cfg := NewConfig(vipr)
	assert.Nil(t, cfg.TLS)

	vipr.Set(TLSEnableViperKey, true)
	cfg = NewConfig(vipr)
	assert.NotNil(t, cfg.TLS)
	assert.Equal(t, cfg.TLS.ServerName, cfg.Host)
	assert.Empty(t, cfg.TLS.Certificates)
	assert.Empty(t, cfg.TLS.CAs)
	
	vipr.Set(TLSCAViperKey, "CACert")
	cfg = NewConfig(vipr)
	assert.NotEmpty(t, cfg.TLS.CAs)
	assert.Empty(t, cfg.TLS.Certificates)
	assert.False(t, cfg.TLS.InsecureSkipVerify)
	
	vipr.Set(TLSCertViperKey, "Cert")
	vipr.Set(TLSKeyViperKey, "CertKey")
	cfg = NewConfig(vipr)
	assert.NotEmpty(t, cfg.TLS.Certificates)
	assert.False(t, cfg.TLS.InsecureSkipVerify)
	
	vipr.Set(TLSSkipVerifyViperKey, true)
	cfg = NewConfig(vipr)
	assert.True(t, cfg.TLS.InsecureSkipVerify)
}
