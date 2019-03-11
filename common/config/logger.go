package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	logLevelFlag     = "log-level"
	logLevelViperKey = "log.level"
	logLevelDefault  = "debug"
	logLevelEnv      = "LOG_LEVEL"
)

// LogLevel register flag for LogLevel
func LogLevel(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}, logLevelEnv)
	f.String(logLevelFlag, logLevelDefault, desc)
	viper.BindPFlag(logLevelViperKey, f.Lookup(logLevelFlag))
	viper.BindEnv(logLevelViperKey, logLevelEnv)
}

var (
	logFormatFlag     = "log-format"
	logFormatViperKey = "log.format"
	logFormatDefault  = "text"
	logFormatEnv      = "LOG_FORMAT"
)

// LogFormat register flag for Log Format
func LogFormat(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, logFormatEnv)
	f.String(logFormatFlag, logFormatDefault, desc)
	viper.BindPFlag(logFormatViperKey, f.Lookup(logFormatFlag))
	viper.BindEnv(logFormatViperKey, logFormatEnv)
}

// ConfigureLogger configure logger
func ConfigureLogger() {
	switch viper.GetString(logFormatViperKey) {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	if logLevel, err := log.ParseLevel(viper.GetString(logLevelViperKey)); err != nil {
		log.Fatalf("Invalid log level: %v", err)
	} else {
		log.New()
		log.SetLevel(logLevel)
	}
}
