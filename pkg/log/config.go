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
	viper.SetDefault(TimestampViperKey, timestampDefault)
	_ = viper.BindEnv(TimestampViperKey, timestampEnv)
}

const (
	levelFlag     = "log-level"
	LevelViperKey = "log.level"
	levelDefault  = "info"
	levelEnv      = "LOG_LEVEL"
)

const (
	formatFlag     = "log-format"
	FormatViperKey = "log.format"
	formatDefault  = "text"
	formatEnv      = "LOG_FORMAT"
)

const (
	timestampFlag     = "log-timestamp"
	TimestampViperKey = "log.timestamp"
	timestampDefault  = false
	timestampEnv      = "LOG_TIMESTAMP"
)

var ECSJsonFormatter = &logrus.JSONFormatter{
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "@timestamp",
		logrus.FieldKeyLevel: "log.level",
		logrus.FieldKeyMsg:   "message",
	},
}

func Flags(f *pflag.FlagSet) {
	level(f)
	format(f)
	timestamp(f)
}

// Level register flag for Level
func level(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log level (one of %q).
Environment variable: %q`, []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}, levelEnv)
	f.String(levelFlag, levelDefault, desc)
	_ = viper.BindPFlag(LevelViperKey, f.Lookup(levelFlag))
}

// Format register flag for Log Format
func format(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Log formatter (one of %q).
Environment variable: %q`, []string{"text", "json"}, formatEnv)
	f.String(formatFlag, formatDefault, desc)
	_ = viper.BindPFlag(FormatViperKey, f.Lookup(formatFlag))
}

func timestamp(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable logging with timestamp (only TEXT format).
Environment variable: %q`, timestampEnv)
	f.Bool(timestampFlag, timestampDefault, desc)
	_ = viper.BindPFlag(TimestampViperKey, f.Lookup(timestampFlag))
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Level:     vipr.GetString(LevelViperKey),
		Format:    vipr.GetString(FormatViperKey),
		Timestamp: vipr.GetBool(TimestampViperKey),
	}
}

type Config struct {
	Level     string
	Format    string
	Timestamp bool
}

func (cfg *Config) ToTraefik() *traefiktypes.TraefikLog {
	tCfg := &traefiktypes.TraefikLog{
		Level: cfg.Level,
	}
	switch cfg.Format {
	case "json":
		tCfg.Format = "json"
	default:
		tCfg.Format = "common"
	}

	return tCfg
}

// ConfigureLogger configures logger
// It sets Traefik global logger so it should be called only once per process
func ConfigureLogger(cfg *Config, logger *logrus.Logger) error {
	if cfg == nil {
		return nil
	}

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
		logrus.SetFormatter(ECSJsonFormatter)
	default:
		formatter := &logrus.TextFormatter{
			PadLevelText: true,
		}

		if cfg.Timestamp {
			formatter.FullTimestamp = true
			formatter.DisableTimestamp = false
		}

		logger.SetFormatter(formatter)
	}

	// TODO: implement internal mechanism for extracting logger for context
	// here we are modifying a global variable so ConfigureLogger should be called once
	traefiklog.SetLogger(logger)
	return nil
}
