package kafka

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
)

func init() {
	viper.SetDefault(DisableExternalTxViperKey, disableExternalTxDefault)
	_ = viper.BindEnv(DisableExternalTxViperKey, disableExternalTxEnv)
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	disableExternalTx(f)
}

const (
	disableExternalTxFlag     = "disable-external-tx"
	DisableExternalTxViperKey = "disable.external.tx"
	disableExternalTxDefault  = false
	disableExternalTxEnv      = "DISABLE_EXTERNAL_TX"
)

// disableExternalTx register flag for Listener Start Default
func disableExternalTx(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Boolean: skip all tx that are not sent directly by Orchestrate
Environment variable: %q`, disableExternalTxEnv)
	f.Bool(disableExternalTxFlag, disableExternalTxDefault, desc)
	_ = viper.BindPFlag(DisableExternalTxViperKey, f.Lookup(disableExternalTxFlag))
}

type Config struct {
	TopicTxDecoder    string
	DisableExternalTx bool
}

func NewConfig() *Config {
	return &Config{
		TopicTxDecoder:    viper.GetString(broker.TxDecoderViperKey),
		DisableExternalTx: viper.GetBool(DisableExternalTxViperKey),
	}
}
