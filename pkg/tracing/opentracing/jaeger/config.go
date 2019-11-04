package jaeger

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func init() {
	viper.SetDefault(hostViperKey, hostDefault)
	_ = viper.BindEnv(hostViperKey, hostEnv)
	viper.SetDefault(portViperKey, portDefault)
	_ = viper.BindEnv(portViperKey, portEnv)
	viper.SetDefault(serviceNameViperKey, serviceNameDefault)
	_ = viper.BindEnv(serviceNameViperKey, serviceNameEnv)
	viper.SetDefault(rpcMetricsViperKey, rpcMetricsDefault)
	_ = viper.BindEnv(rpcMetricsViperKey, rpcMetricsEnv)
	viper.SetDefault(enabledViperKey, enabledDefault)
	_ = viper.BindEnv(enabledViperKey, enabledEnv)
	viper.SetDefault(collectorURLViperKey, collectorURLDefault)
	_ = viper.BindEnv(collectorURLViperKey, collectorURLEnv)
	viper.SetDefault(userViperKey, userDefault)
	_ = viper.BindEnv(userViperKey, userEnv)
	viper.SetDefault(passwordViperKey, passwordDefault)
	_ = viper.BindEnv(passwordViperKey, passwordEnv)
	viper.SetDefault(samplerParamViperKey, samplerParamDefault)
	_ = viper.BindEnv(samplerParamViperKey, samplerParamEnv)
	viper.SetDefault(samplerTypeViperKey, samplerTypeDefault)
	_ = viper.BindEnv(samplerTypeViperKey, samplerTypeEnv)
	viper.SetDefault(logSpansViperKey, logSpansDefault)
	_ = viper.BindEnv(logSpansViperKey, logSpansEnv)
}

// InitFlags register Jaeger flags
func InitFlags(f *pflag.FlagSet) {
	Host(f)
	Port(f)
	ServiceName(f)
	RPCMetrics(f)
	Enabled(f)
	CollectorURL(f)
	User(f)
	Password(f)
	SamplerParam(f)
	SamplerType(f)
	LogSpans(f)
}

const (
	hostFlag     = "jaeger-agent-host"
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

const (
	portFlag     = "jaeger-agent-port"
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

const (
	serviceNameViperKey = "jaeger.service.name"
	serviceNameFlag     = "jaeger-service-name"
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

const (
	collectorURLViperKey = "jaeger.collector.url"
	collectorURLFlag     = "jaeger-collector-url"
	collectorURLDefault  = ""
	collectorURLEnv      = "JAEGER_COLLECTOR_URL"
)

// CollectorURL register a flag for Jaeger collector url
func CollectorURL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger collector url to send spans to
Environment variable: %q`, collectorURLEnv)
	f.String(collectorURLFlag, collectorURLDefault, desc)
	_ = viper.BindPFlag(collectorURLViperKey, f.Lookup(collectorURLFlag))
}

const (
	userViperKey = "jaeger.user"
	userFlag     = "jaeger-user"
	userDefault  = ""
	userEnv      = "JAEGER_USER"
)

// User register a flag for Jaeger User
func User(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger User
Environment variable: %q`, userEnv)
	f.String(userFlag, userDefault, desc)
	_ = viper.BindPFlag(userViperKey, f.Lookup(userFlag))
}

const (
	passwordViperKey = "jaeger.password"
	passwordFlag     = "jaeger-password"
	passwordDefault  = ""
	passwordEnv      = "JAEGER_PASSWORD"
)

// Password register a flag for Jaeger password
func Password(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger password
Environment variable: %q`, passwordEnv)
	f.String(passwordFlag, passwordDefault, desc)
	_ = viper.BindPFlag(passwordViperKey, f.Lookup(passwordFlag))
}

const (
	enabledViperKey = "jaeger.enabled"
	enabledFlag     = "jaeger-enabled"
	enabledDefault  = true
	enabledEnv      = "JAEGER_ENABLED"
)

// Enabled register a flag to enable Jaeger
func Enabled(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable Jaeger reporting
Environment variable: %q`, enabledEnv)
	f.Bool(enabledFlag, enabledDefault, desc)
	_ = viper.BindPFlag(enabledViperKey, f.Lookup(enabledFlag))
}

const (
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

const (
	logSpansViperKey = "jaeger.reporter.log.spans"
	logSpansFlag     = "jaeger-reporter-log-spans"
	logSpansDefault  = true
	logSpansEnv      = "JAEGER_REPORTER_LOG_SPANS"
)

// LogSpans register a flag for LogSpans Jaeger option
func LogSpans(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`LogSpans, when true, enables LoggingReporter that runs in parallel with the main reporter and logs all submitted spans.
Environment variable: %q`, logSpansEnv)
	f.Bool(logSpansFlag, logSpansDefault, desc)
	_ = viper.BindPFlag(logSpansViperKey, f.Lookup(logSpansFlag))
}

const (
	samplerParamFlag     = "jaeger-sampler-param"
	samplerParamViperKey = "jaeger.sampler.param"
	samplerParamDefault  = 1
	samplerParamEnv      = "JAEGER_SAMPLER_PARAM"
)

// SamplerParam register a flag for Jaeger
func SamplerParam(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Jaeger sampler
Environment variable: %q`, samplerParamEnv)
	f.Int(samplerParamFlag, samplerParamDefault, desc)
	_ = viper.BindPFlag(samplerParamViperKey, f.Lookup(samplerParamFlag))
}

const (
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
		Disabled:    !viper.GetBool(enabledViperKey),
		RPCMetrics:  viper.GetBool(rpcMetricsViperKey),
		Sampler: &config.SamplerConfig{
			Type:  viper.GetString(samplerTypeViperKey),
			Param: viper.GetFloat64(samplerParamViperKey),
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           viper.GetBool(logSpansViperKey),
			LocalAgentHostPort: fmt.Sprintf("%v:%v", viper.GetString(hostViperKey), viper.GetString(portViperKey)),
			CollectorEndpoint:  viper.GetString(collectorURLViperKey),
			User:               viper.GetString(userViperKey),
			Password:           viper.GetString(passwordViperKey),
		},
	}
}

// TracerFromConfig returns a wrapped tracer from config object
func TracerFromConfig(c *config.Configuration) (opentracing.Tracer, error) {
	metrics := prometheus.New()
	tracer, _, err := c.NewTracer(
		config.Logger(logger{entry: log.StandardLogger().WithFields(log.Fields{"system": "opentracing.jaeger"})}),
		config.Observer(rpcmetrics.NewObserver(metrics.Namespace(jaegermetrics.NSOptions{Name: c.ServiceName}), rpcmetrics.DefaultNameNormalizer)),
	)

	return tracer, err
}

// TracerFromViperConfig returns a wrapped Tracer object.
// Log fatal if it encounters an error
func TracerFromViperConfig() opentracing.Tracer {
	conf := NewConfig()
	tracer, err := TracerFromConfig(conf)
	if err != nil {
		log.Fatalf("Could not instantiate tracer object: %v", err)
	}

	return tracer
}
