package jaeger

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
)

func init() {
	viper.SetDefault(hostViperKey, hostDefault)
	_ = viper.BindEnv(hostViperKey, hostEnv)
	viper.SetDefault(portViperKey, portDefault)
	_ = viper.BindEnv(portViperKey, portEnv)
	viper.SetDefault(samplerParamViperKey, samplerParamDefault)
	_ = viper.BindEnv(samplerParamViperKey, samplerParamEnv)
	viper.SetDefault(samplerTypeViperKey, samplerTypeDefault)
	_ = viper.BindEnv(samplerTypeViperKey, samplerTypeEnv)
	viper.SetDefault(serviceNameViperKey, serviceNameDefault)
	_ = viper.BindEnv(serviceNameViperKey, serviceNameEnv)
	viper.SetDefault(endpointViperKey, endpointDefault)
	_ = viper.BindEnv(endpointViperKey, endpointEnv)
	viper.SetDefault(userViperKey, userDefault)
	_ = viper.BindEnv(userViperKey, userEnv)
	viper.SetDefault(passwordViperKey, passwordDefault)
	_ = viper.BindEnv(passwordViperKey, passwordEnv)
	viper.SetDefault(logSpansViperKey, logSpansDefault)
	_ = viper.BindEnv(logSpansViperKey, logSpansEnv)
	viper.SetDefault(disabledViperKey, disabledDefault)
	_ = viper.BindEnv(disabledViperKey, disabledEnv)
	viper.SetDefault(rpcMetricsViperKey, rpcMetricsDefault)
	_ = viper.BindEnv(rpcMetricsViperKey, rpcMetricsEnv)
	viper.SetDefault(logSpansViperKey, logSpansDefault)
	_ = viper.BindEnv(logSpansViperKey, logSpansEnv)
}

// InitFlags register Jaeger flags
func InitFlags(f *pflag.FlagSet) {
	Host(f)
	Port(f)
	ServiceName(f)
	RPCMetrics(f)
	Disabled(f)
	Endpoint(f)
	User(f)
	Password(f)
	SamplerParam(f)
	SamplerType(f)
}

var (
	hostFlag     = "jaeger-host"
	hostViperKey = "jaeger.agent.host"
	hostDefault  = "localhost"
	hostEnv      = "JAEGER_AGENT_HOST"
)

// Host register a flag for Jaeger host
func Host(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger host.
Environment variable: %q`, hostEnv)
	f.String(hostFlag, hostDefault, desc)
	_ = viper.BindPFlag(hostViperKey, f.Lookup(hostFlag))
}

var (
	portFlag     = "jaeger-port"
	portViperKey = "jaeger.agent.port"
	portDefault  = 6831
	portEnv      = "JAEGER_AGENT_PORT"
)

// Port register a flag for Jaeger port
func Port(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger port
Environment variable: %q`, portEnv)
	f.Int(portFlag, portDefault, desc)
	_ = viper.BindPFlag(portViperKey, f.Lookup(portFlag))
}

var (
	serviceNameViperKey = "jaeger.service.name"
	serviceNameFlag     = "jaeger-service"
	serviceNameDefault  = "jaeger"
	serviceNameEnv      = "JAEGER_SERVICE_NAME"
)

// ServiceName register a flag for Jaeger Service name
func ServiceName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger ServiceName to use on the tracer
Environment variable: %q`, serviceNameEnv)
	f.String(serviceNameFlag, serviceNameDefault, desc)
	_ = viper.BindPFlag(serviceNameViperKey, f.Lookup(serviceNameFlag))
}

var (
	endpointViperKey = "jaeger.endpoint"
	endpointFlag     = "jaeger-endpoint"
	endpointDefault  = ""
	endpointEnv      = "JAEGER_ENDPOINT"
)

// Endpoint register a flag for Jaeger Endpoint
func Endpoint(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger collector endpoint to send spans to
Environment variable: %q`, endpointEnv)
	f.String(endpointFlag, endpointDefault, desc)
	_ = viper.BindPFlag(endpointViperKey, f.Lookup(endpointFlag))
}

var (
	userViperKey = "jaeger.user"
	userFlag     = "jaeger-user"
	userDefault  = ""
	userEnv      = "JAEGER_USER"
)

// User register a flag for Jaeger User
func User(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger collector User
Environment variable: %q`, userEnv)
	f.String(userFlag, userDefault, desc)
	_ = viper.BindPFlag(userViperKey, f.Lookup(userFlag))
}

var (
	passwordViperKey = "jaeger.password"
	passwordFlag     = "jaeger-password"
	passwordDefault  = ""
	passwordEnv      = "JAEGER_PASSWORD"
)

// Password register a flag for Jaeger Endpoint
func Password(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger collector password
Environment variable: %q`, passwordEnv)
	f.String(passwordFlag, passwordDefault, desc)
	_ = viper.BindPFlag(passwordViperKey, f.Lookup(passwordFlag))
}

var (
	disabledViperKey = "jaeger.disabled"
	disabledFlag     = "jaeger-disabled"
	disabledDefault  = false
	disabledEnv      = "JAEGER_DISABLED"
)

// Disabled register a flag to disable Jaeger
func Disabled(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Disable Jaeger reporting
Environment variable: %q`, disabledEnv)
	f.Bool(disabledFlag, disabledDefault, desc)
	_ = viper.BindPFlag(disabledViperKey, f.Lookup(disabledFlag))
}

var (
	rpcMetricsViperKey = "jaeger.rpc.metrics"
	rpcMetricsFlag     = "jaeger-rpc-metrics"
	rpcMetricsDefault  = false
	rpcMetricsEnv      = "JAEGER_RPC_METRICS"
)

// RPCMetrics register a flag to enable RPC Metrics
func RPCMetrics(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable Jaeger RPC metrics
Environment variable: %q`, rpcMetricsEnv)
	f.Bool(rpcMetricsFlag, rpcMetricsDefault, desc)
	_ = viper.BindPFlag(rpcMetricsViperKey, f.Lookup(rpcMetricsFlag))
}

var (
	logSpansViperKey = "jaeger.reporter.log.spans"
	logSpansFlag     = "jaeger-log-spans"
	logSpansDefault  = true
	logSpansEnv      = "JAEGER_REPORTER_LOG_SPANS"
)

// LogSpans register a flag for LogSpans Jaeger option
func LogSpans(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`When enabled Jager reporter to run in parallel
Environment variable: %q`, logSpansEnv)
	f.Bool(logSpansFlag, logSpansDefault, desc)
	_ = viper.BindPFlag(logSpansViperKey, f.Lookup(logSpansFlag))
}

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
	_ = viper.BindPFlag(samplerParamViperKey, f.Lookup(samplerParamFlag))
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
	_ = viper.BindPFlag(samplerTypeViperKey, f.Lookup(samplerTypeFlag))
}

// NewConfig create new Jaeger configuration
func NewConfig() *config.Configuration {
	return &config.Configuration{
		ServiceName: viper.GetString(serviceNameViperKey),
		Disabled:    viper.GetBool(disabledViperKey),
		RPCMetrics:  viper.GetBool(rpcMetricsViperKey),
		Sampler: &config.SamplerConfig{
			Type:  viper.GetString(samplerTypeViperKey),
			Param: viper.GetFloat64(samplerParamViperKey),
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           viper.GetBool(logSpansViperKey),
			LocalAgentHostPort: fmt.Sprintf("%v:%v", viper.GetString(hostViperKey), viper.GetString(portViperKey)),
			CollectorEndpoint:  viper.GetString(endpointViperKey),
			User:               viper.GetString(userViperKey),
			Password:           viper.GetString(passwordViperKey),
		},
	}
}
