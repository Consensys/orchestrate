package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(jaegerHostViperKey, jaegerHostDefault)
	viper.BindEnv(jaegerHostViperKey, jaegerHostEnv)
	viper.SetDefault(jaegerPortViperKey, jaegerPortDefault)
	viper.BindEnv(jaegerPortViperKey, jaegerPortEnv)
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
	jaegerPortDefault  = 5775
	jaegerPortEnv      = "JAEGER_PORT"
)

// JaegerPort register a flag for Jaeger port
func JaegerPort(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger port
Environment variable: %q`, jaegerPortEnv)
	f.Int(jaegerPortFlag, jaegerPortDefault, desc)
	viper.BindPFlag(jaegerPortViperKey, f.Lookup(jaegerPortFlag))
}

var (
	jaegerSamplerFlag     = "jaeger-sampler"
	jaegerSamplerViperKey = "jaeger.sampler"
	jaegerSamplerDefault  = 0.5
	jaegerSamplerEnv      = "JAEGER_SAMPLER"
)

// JaegerSampler register a flag for jaegger sampler
func JaegerSampler(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger sampler
Environment variable: %q`, jaegerSamplerEnv)
	f.Float64(jaegerSamplerFlag, jaegerSamplerDefault, desc)
	viper.SetDefault(jaegerSamplerViperKey, jaegerSamplerDefault)
	viper.BindPFlag(jaegerSamplerViperKey, f.Lookup(jaegerSamplerFlag))
	viper.BindEnv(jaegerSamplerViperKey, jaegerSamplerEnv)
}
