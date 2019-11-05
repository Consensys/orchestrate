package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(LogLevelViperKey, logLevelDefault)
	_ = viper.BindEnv(LogLevelViperKey, logLevelEnv)
	viper.SetDefault(LogFormatViperKey, logFormatDefault)
	_ = viper.BindEnv(LogFormatViperKey, logFormatEnv)
}

const (
	logLevelFlag     = "log-level"
	LogLevelViperKey = "log.level"
	logLevelDefault  = "info"
	logLevelEnv      = "LOG_LEVEL"
)

// LogLevel register flag for LogLevel
func LogLevel(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}, logLevelEnv)
	f.String(logLevelFlag, logLevelDefault, desc)
	_ = viper.BindPFlag(LogLevelViperKey, f.Lookup(logLevelFlag))
}

const (
	logFormatFlag     = "log-format"
	LogFormatViperKey = "log.format"
	logFormatDefault  = "text"
	logFormatEnv      = "LOG_FORMAT"
)

// LogFormat register flag for Log Format
func LogFormat(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, logFormatEnv)
	f.String(logFormatFlag, logFormatDefault, desc)
	_ = viper.BindPFlag(LogFormatViperKey, f.Lookup(logFormatFlag))
}

// InitLogger Initialize logrus Logger
func InitLogger() {
	switch viper.GetString(LogFormatViperKey) {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	if logLevel, err := log.ParseLevel(viper.GetString(LogLevelViperKey)); err != nil {
		log.Fatalf("Invalid log level: %v", err)
	} else {
		log.New()
		log.SetLevel(logLevel)
	}
}
