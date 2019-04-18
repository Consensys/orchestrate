package jaeger


import (
	"fmt"
	"github.com/uber/jaeger-client-go/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(jaegerHostViperKey, jaegerHostDefault)
	viper.BindEnv(jaegerHostViperKey, jaegerHostEnv)
	viper.SetDefault(jaegerPortViperKey, jaegerPortDefault)
	viper.BindEnv(jaegerPortViperKey, jaegerPortEnv)

	viper.SetDefault(jaegerSamplerParamViperKey, jaegerSamplerParamDefault)
	viper.BindEnv(jaegerSamplerParamViperKey, jaegerSamplerParamEnv)
	viper.SetDefault(jaegerSamplerTypeViperKey, jaegerSamplerTypeDefault)
	viper.BindEnv(jaegerSamplerTypeViperKey, jaegerSamplerTypeEnv)
	viper.SetDefault(jaegerServiceNameViperKey, jaegerServiceNameDefault)
	viper.BindEnv(jaegerServiceNameViperKey, jaegerServiceNameEnv)
	viper.SetDefault(jaegerDisabledViperKey, jaegerDisabledDefault)
	viper.BindEnv(jaegerDisabledViperKey, jaegerDisabledEnv)
	viper.SetDefault(jaegerLogSpansViperKey, jaegerLogSpansDefault)
	viper.BindEnv(jaegerLogSpansViperKey, jaegerLogSpansEnv)
}

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
	jaegerHostFlag     = "jaeger-host"
	jaegerHostViperKey = "jaeger.host"
	jaegerHostDefault  = "jaeger"
	jaegerHostEnv      = "JAEGER_HOST"
)

// JaegerHost register a flag for Jaeger host
func JaegerHost(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger host.
Environment variable: %q`, jaegerHostEnv)
	f.String(jaegerHostFlag, jaegerHostDefault, desc)
	viper.BindPFlag(jaegerHostViperKey, f.Lookup(jaegerHostFlag))
}

var (
	jaegerPortFlag     = "jaeger-port"
	jaegerPortViperKey = "jaeger.port"
	jaegerPortDefault  = 6831
	jaegerPortEnv      = "JAEGER_PORT"
)

// JaegerPort register a flag for Jaeger port
func JaegerPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger port
Environment variable: %q`, jaegerPortEnv)
	f.Int(jaegerPortFlag, jaegerPortDefault, desc)
	viper.BindPFlag(jaegerPortViperKey, f.Lookup(jaegerPortFlag))
}

// TODO : adding all binding Flag and to trigger Jaegger
var (
	jaegerServiceNameViperKey = "jaeger.service.name"
	jaegerServiceNameDefault  = "jaeger"
	jaegerServiceNameEnv      = "JAEGER_SERVICE_NAME"
)

// TODO : adding all binding Flag and to trigger Jaegger
var (
	jaegerDisabledViperKey = "jaeger.disabled"
	jaegerDisabledDefault  = true
	jaegerDisabledEnv      = "JAEGER_DISABLED"
)

// TODO : adding all binding Flag and to trigger Jaegger
var (
	jaegerLogSpansViperKey = "jaeger.reporter.logspans"
	jaegerLogSpansDefault  = true
	jaegerLogSpansEnv      = "JAEGER_REPORTER_LOGSPANS"
)

var (
	jaegerSamplerParamFlag     = "jaeger-sampler-param"
	jaegerSamplerParamViperKey = "jaeger.sampler.param"
	jaegerSamplerParamDefault  = 1
	jaegerSamplerParamEnv      = "JAEGER_SAMPLER_PARAM"
)

// JaegerSampler register a flag for jaeger
func JaegerSamplerParam(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger sampler
Environment variable: %q`, jaegerSamplerParamEnv)
	f.Int(jaegerSamplerParamFlag, jaegerSamplerParamDefault, desc)
	viper.BindPFlag(jaegerSamplerParamViperKey, f.Lookup(jaegerSamplerParamFlag))
}

var (
	jaegerSamplerTypeFlag     = "jaeger-sampler-type"
	jaegerSamplerTypeViperKey = "jaeger.sampler.type"
	jaegerSamplerTypeDefault  = "const"
	jaegerSamplerTypeEnv      = "JAEGER_SAMPLER_TYPE"
)

// JaegerSampler register a flag for jaeger
func JaegerSamplerType(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger sampler
Environment variable: %q`, jaegerSamplerTypeEnv)
	f.String(jaegerSamplerTypeFlag, jaegerSamplerTypeDefault, desc)
	viper.BindPFlag(jaegerSamplerTypeViperKey, f.Lookup(jaegerSamplerTypeFlag))
}