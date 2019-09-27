package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(logLevelViperKey, logLevelDefault)
	_ = viper.BindEnv(logLevelViperKey, logLevelEnv)
	viper.SetDefault(logFormatViperKey, logFormatDefault)
	_ = viper.BindEnv(logFormatViperKey, logFormatEnv)
}

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
	_ = viper.BindPFlag(logLevelViperKey, f.Lookup(logLevelFlag))
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
	_ = viper.BindPFlag(logFormatViperKey, f.Lookup(logFormatFlag))
}

// InitLogger Initialize logrus Logger
func InitLogger() {
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
