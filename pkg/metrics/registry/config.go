package registry

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	metricsModulesFlag       = "metrics-modules"
	MetricsModulesViperKey   = "metrics.modules"
	metricsModulesEnv        = "METRICS_MODULES"
	MetricsModuleEnableAll   = "ENABLED"
	MetricsModuleDisabledAll = "DISABLED"
)

var (
	MetricsModulesDefault = []string{MetricsModuleEnableAll}
	GoMetricsModule       = "go"
	ProcessMetricsModule  = "process"
	HealthzMetricsModule  = "healthz"
)

func init() {
	viper.SetDefault(MetricsModulesViperKey, MetricsModulesDefault)
	_ = viper.BindEnv(MetricsModulesViperKey, metricsModulesEnv)
}

// Flags register flags for tx sentry
func Flags(f *pflag.FlagSet, modules ...string) {
	metricsModuleDesc := fmt.Sprintf(`List of metrics modules exposed. Available metric modules are %q, to enable all use %s or to disable all %s. 
Environment variable: %q`,
		append(modules, GoMetricsModule, ProcessMetricsModule, HealthzMetricsModule), MetricsModuleEnableAll, MetricsModuleDisabledAll, metricsModulesEnv)
	f.StringSlice(metricsModulesFlag, MetricsModulesDefault, metricsModuleDesc)
	_ = viper.BindPFlag(MetricsModulesViperKey, f.Lookup(metricsModulesFlag))
}

type Config struct {
	modules []string
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		modules: vipr.GetStringSlice(MetricsModulesViperKey),
	}
}

func (c *Config) IsActive(module string) bool {
	for _, m := range c.modules {
		if m == MetricsModuleDisabledAll {
			return false
		}

		if m == module || m == MetricsModuleEnableAll {
			return true
		}
	}

	return false
}

func (c *Config) Modules() []string {
	return c.modules
}
