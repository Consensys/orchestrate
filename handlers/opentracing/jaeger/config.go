package jaeger

import (
	"fmt"

	"github.com/uber/jaeger-client-go/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hostViperKey, hostDefault)
	viper.BindEnv(hostViperKey, hostEnv)
	viper.SetDefault(portViperKey, portDefault)
	viper.BindEnv(portViperKey, portEnv)
	viper.SetDefault(samplerParamViperKey, samplerParamDefault)
	viper.BindEnv(samplerParamViperKey, samplerParamEnv)
	viper.SetDefault(samplerTypeViperKey, samplerTypeDefault)
	viper.BindEnv(samplerTypeViperKey, samplerTypeEnv)
	viper.SetDefault(serviceNameViperKey, serviceNameDefault)
	viper.BindEnv(serviceNameViperKey, serviceNameEnv)
	viper.SetDefault(disabledViperKey, disabledDefault)
	viper.BindEnv(disabledViperKey, disabledEnv)
	viper.SetDefault(logSpansViperKey, logSpansDefault)
	viper.BindEnv(logSpansViperKey, logSpansEnv)
}

// NewConfig create new Jaeger configuration
func NewConfig() *config.Configuration {
	return &config.Configuration{
		ServiceName: viper.GetString("jaeger.service.name"),
		Disabled:    viper.GetBool("jaeger.disabled"),
		Sampler: &config.SamplerConfig{
			Type:  viper.GetString("jaeger.sampler.type"),
			Param: viper.GetFloat64("jaeger.sampler.param"),
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           viper.GetBool("jaeger.reporter.logspans"),
			LocalAgentHostPort: fmt.Sprintf("%v:%v", viper.GetString("jaeger.host"), viper.GetString("jaeger.port")),
		},
	}
}

var (
	hostFlag     = "jaeger-host"
	hostViperKey = "jaeger.host"
	hostDefault  = "jaeger"
	hostEnv      = "JAEGER_HOST"
)

// Host register a flag for Jaeger host
func Host(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger host.
Environment variable: %q`, hostEnv)
	f.String(hostFlag, hostDefault, desc)
	viper.BindPFlag(hostViperKey, f.Lookup(hostFlag))
}

var (
	portFlag     = "jaeger-port"
	portViperKey = "jaeger.port"
	portDefault  = 6831
	portEnv      = "JAEGER_PORT"
)

// Port register a flag for Jaeger port
func Port(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger port
Environment variable: %q`, portEnv)
	f.Int(portFlag, portDefault, desc)
	viper.BindPFlag(portViperKey, f.Lookup(portFlag))
}

// TODO : adding all binding Flag and to trigger Jaegger
var (
	serviceNameViperKey = "jaeger.service.name"
	serviceNameDefault  = "jaeger"
	serviceNameEnv      = "JAEGER_SERVICE_NAME"
)

// TODO : adding all binding Flag and to trigger Jaegger
var (
	disabledViperKey = "jaeger.disabled"
	disabledDefault  = true
	disabledEnv      = "JAEGER_DISABLED"
)

// TODO : adding all binding Flag and to trigger Jaegger
var (
	logSpansViperKey = "jaeger.reporter.logspans"
	logSpansDefault  = true
	logSpansEnv      = "JAEGER_REPORTER_LOGSPANS"
)

var (
	samplerParamFlag     = "jaeger-sampler-param"
	samplerParamViperKey = "jaeger.sampler.param"
	samplerParamDefault  = 1
	samplerParamEnv      = "JAEGER_SAMPLER_PARAM"
)

// SamplerParam register a flag for jaeger
func SamplerParam(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger sampler
Environment variable: %q`, samplerParamEnv)
	f.Int(samplerParamFlag, samplerParamDefault, desc)
	viper.BindPFlag(samplerParamViperKey, f.Lookup(samplerParamFlag))
}

var (
	samplerTypeFlag     = "jaeger-sampler-type"
	samplerTypeViperKey = "jaeger.sampler.type"
	samplerTypeDefault  = "const"
	samplerTypeEnv      = "JAEGER_SAMPLER_TYPE"
)

// SamplerType register a flag for jaeger
func SamplerType(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger sampler
Environment variable: %q`, samplerTypeEnv)
	f.String(samplerTypeFlag, samplerTypeDefault, desc)
	viper.BindPFlag(samplerTypeViperKey, f.Lookup(samplerTypeFlag))
}
