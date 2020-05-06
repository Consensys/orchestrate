package log

import (
	"fmt"
	"strings"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(LevelViperKey, levelDefault)
	_ = viper.BindEnv(LevelViperKey, levelEnv)
	viper.SetDefault(FormatViperKey, formatDefault)
	_ = viper.BindEnv(FormatViperKey, formatEnv)
}

const (
	levelFlag     = "log-level"
	LevelViperKey = "log.level"
	levelDefault  = "info"
	levelEnv      = "LOG_LEVEL"
)

// Level register flag for Level
func Level(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}, levelEnv)
	f.String(levelFlag, levelDefault, desc)
	_ = viper.BindPFlag(LevelViperKey, f.Lookup(levelFlag))
}

const (
	formatFlag     = "log-format"
	FormatViperKey = "log.format"
	formatDefault  = "text"
	formatEnv      = "LOG_FORMAT"
)

// Format register flag for Log Format
func Format(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, formatEnv)
	f.String(formatFlag, formatDefault, desc)
	_ = viper.BindPFlag(FormatViperKey, f.Lookup(formatFlag))
}

func NewConfig(vipr *viper.Viper) *traefiktypes.TraefikLog {
	return &traefiktypes.TraefikLog{
		Level:  vipr.GetString(LevelViperKey),
		Format: ToTraefikFormat(vipr.GetString(FormatViperKey)),
	}
}

func ToTraefikFormat(format string) string {
	switch format {
	case "json":
		return "json"
	default:
		return "common"
	}
}

// ConfigureLogger configures logger
// It sets Traefik global logger so it should be called only once per process
func ConfigureLogger(cfg *traefiktypes.TraefikLog, logger *logrus.Logger) error {
	if cfg != nil {
		if cfg.Level == "" {
			cfg.Level = "INFO"
		}

		// Set Level
		level, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
		if err != nil {
			return err
		}
		logger.SetLevel(level)

		// Set Formatter
		switch cfg.Format {
		case "json":
			logger.SetFormatter(&logrus.JSONFormatter{})
		default:
			logger.SetFormatter(&logrus.TextFormatter{})
		}

		// TODO: implement internal mechanism for extracting logger for context
		// here we are modifiying a global variable so ConfigureLogger should be called once
		traefiklog.SetLogger(logger)
	}

	return nil
}
