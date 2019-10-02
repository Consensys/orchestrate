package jaeger

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitFlags(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitFlags(flgs)

	assert.Equal(t, hostDefault, viper.GetString(hostViperKey), "Default")
	assert.Equal(t, portDefault, viper.GetInt(portViperKey), "Default")
	assert.Equal(t, serviceNameDefault, viper.GetString(serviceNameViperKey), "Default")
	assert.Equal(t, endpointDefault, viper.GetString(endpointViperKey), "Default")
	assert.Equal(t, userDefault, viper.GetString(userViperKey), "Default")
	assert.Equal(t, passwordDefault, viper.GetString(passwordViperKey), "Default")
	assert.Equal(t, disabledDefault, viper.GetBool(disabledViperKey), "Default")
	assert.Equal(t, rpcMetricsDefault, viper.GetBool(rpcMetricsViperKey), "Default")
	assert.Equal(t, logSpansDefault, viper.GetBool(logSpansViperKey), "Default")
	assert.Equal(t, samplerParamDefault, viper.GetInt(samplerParamViperKey), "Default")
	assert.Equal(t, samplerTypeDefault, viper.GetString(samplerTypeViperKey), "Default")
}

func TestHost(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Host(flgs)
	assert.Equal(t, hostDefault, viper.GetString(hostViperKey), "Default")

	_ = os.Setenv("JAEGER_AGENT_HOST", "env-jaeger")
	expected := "env-jaeger"
	assert.Equal(t, expected, viper.GetString(hostViperKey), "From Environment Variable")
	_ = os.Unsetenv("JAEGER_AGENT_HOST")

	args := []string{
		"--jaeger-host=flag-jaeger",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-jaeger"
	assert.Equal(t, expected, viper.GetString(hostViperKey), "From Flag")

	// As tests are run in the same context when in the same package,
	// hostFlah has to be reset to default value after update testing for exported_test.go to be successful
	e := flgs.Set(hostFlag, hostDefault)
	assert.NoError(t, e, "No error expected")
}

func TestPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Port(flgs)
	assert.Equal(t, portDefault, viper.GetInt(portViperKey), "Default")

	_ = os.Setenv("JAEGER_AGENT_PORT", "5778")
	expected := 5778
	assert.Equal(t, expected, viper.GetInt(portViperKey), "From Environment Variable")
	_ = os.Unsetenv("JAEGER_AGENT_PORT")

	args := []string{
		"--jaeger-port=5779",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 5779
	assert.Equal(t, expected, viper.GetInt(portViperKey), "From Flag")
}

func TestServiceName(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ServiceName(flgs)
	assert.Equal(t, serviceNameDefault, viper.GetString(serviceNameViperKey), "Default")

	expected := "Test-service"
	_ = os.Setenv(serviceNameEnv, expected)
	assert.Equal(t, expected, viper.GetString(serviceNameViperKey), "From Environment Variable")
	_ = os.Unsetenv(serviceNameEnv)

	expected = "Test-service-2"
	args := []string{
		fmt.Sprintf("--%v=%v", serviceNameFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, expected, viper.GetString(serviceNameViperKey), "From Flag")
}

func TestEndPoint(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Endpoint(flgs)
	assert.Equal(t, endpointDefault, viper.GetString(endpointViperKey), "Default")

	expected := "Test-endpoint"
	_ = os.Setenv(endpointEnv, expected)
	assert.Equal(t, expected, viper.GetString(endpointViperKey), "From Environment Variable")
	_ = os.Unsetenv(endpointEnv)

	expected = "Test-endpoint-2"
	args := []string{
		fmt.Sprintf("--%v=%v", endpointFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, expected, viper.GetString(endpointViperKey), "From Flag")
}

func TestUser(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	User(flgs)
	assert.Equal(t, userDefault, viper.GetString(userViperKey), "Default")

	expected := "Test-user"
	_ = os.Setenv(userEnv, expected)
	assert.Equal(t, expected, viper.GetString(userViperKey), "From Environment Variable")
	_ = os.Unsetenv(userEnv)

	expected = "Test-user-2"
	args := []string{
		fmt.Sprintf("--%v=%v", userFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, expected, viper.GetString(userViperKey), "From Flag")
}

func TestPassword(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Password(flgs)
	assert.Equal(t, passwordDefault, viper.GetString(passwordViperKey), "Default")

	expected := "Test-password"
	_ = os.Setenv(passwordEnv, expected)
	assert.Equal(t, expected, viper.GetString(passwordViperKey), "From Environment Variable")
	_ = os.Unsetenv(passwordEnv)

	expected = "Test-password-2"
	args := []string{
		fmt.Sprintf("--%v=%v", passwordFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, expected, viper.GetString(passwordViperKey), "From Flag")
}

func TestDisabled(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Disabled(flgs)
	assert.Equal(t, disabledDefault, viper.GetBool(disabledViperKey), "Default")

	expected := "true"
	_ = os.Setenv(disabledEnv, expected)
	assert.Equal(t, true, viper.GetBool(disabledViperKey), "From Environment Variable")
	_ = os.Unsetenv(disabledEnv)

	expected = "true"
	args := []string{
		fmt.Sprintf("--%v=%v", disabledFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, true, viper.GetBool(disabledViperKey), "From Flag")

	// As tests are run in the same context when in the same package,
	// disabledFlag has to be reset to default value after update testing for exported_test.go to be successful
	e := flgs.Set(disabledFlag, strconv.FormatBool(disabledDefault))
	assert.NoError(t, e, "No error expected")
}

func TestRpcMetrics(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RPCMetrics(flgs)
	assert.Equal(t, rpcMetricsDefault, viper.GetBool(rpcMetricsViperKey), "Default")

	expected := "true"
	_ = os.Setenv(rpcMetricsEnv, expected)
	assert.Equal(t, true, viper.GetBool(rpcMetricsViperKey), "From Environment Variable")
	_ = os.Unsetenv(rpcMetricsEnv)

	expected = "true"
	args := []string{
		fmt.Sprintf("--%v=%v", rpcMetricsFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, true, viper.GetBool(rpcMetricsViperKey), "From Flag")
}

func TestLogSpans(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogSpans(flgs)
	assert.Equal(t, logSpansDefault, viper.GetBool(logSpansViperKey), "Default")

	expected := "true"
	_ = os.Setenv(logSpansEnv, expected)
	assert.Equal(t, true, viper.GetBool(logSpansViperKey), "From Environment Variable")
	_ = os.Unsetenv(logSpansEnv)

	expected = "true"
	args := []string{
		fmt.Sprintf("--%v=%v", logSpansFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, true, viper.GetBool(logSpansViperKey), "From Flag")
}

func TestSamplerParam(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	SamplerParam(flgs)
	assert.Equal(t, samplerParamDefault, viper.GetInt(samplerParamViperKey), "Default")

	_ = os.Setenv("JAEGER_SAMPLER_PARAM", "0")
	expected := 0
	assert.Equal(t, expected, viper.GetInt(samplerParamViperKey), "From Environment Variable")
	_ = os.Unsetenv("JAEGER_HOST")

	args := []string{
		"--jaeger-sampler-param=0",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 0
	assert.Equal(t, expected, viper.GetInt(samplerParamViperKey), "From Flag")
}

func TestSamplerType(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	SamplerType(flgs)
	assert.Equal(t, samplerTypeDefault, viper.GetString(samplerTypeViperKey), "Default")

	expected := "probabilistic"
	_ = os.Setenv(samplerTypeEnv, expected)
	assert.Equal(t, expected, viper.GetString(samplerTypeViperKey), "From Environment Variable")
	_ = os.Unsetenv(samplerTypeEnv)

	expected = "rateLimiting"
	args := []string{
		fmt.Sprintf("--%v=%v", samplerTypeFlag, expected),
	}

	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, expected, viper.GetString(samplerTypeViperKey), "From Flag")
}
