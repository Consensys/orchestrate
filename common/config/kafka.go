package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(kafkaAddressViperKey, kafkaAddressDefault)
	viper.BindEnv(kafkaAddressViperKey, kafkaAddressEnv)
}

var (
	kafkaAddressFlag     = "kafka-address"
	kafkaAddressViperKey = "kafka.addresses"
	kafkaAddressDefault  = []string{"localhost:9092"}
	kafkaAddressEnv      = "KAFKA_ADDRESS"
)

// KafkaAddresses register flag for Kafka server addresses
func KafkaAddresses(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Kafka server to connect to.
Environment variable: %q`, kafkaAddressEnv)
	f.StringSlice(kafkaAddressFlag, kafkaAddressDefault, desc)
	viper.BindPFlag(kafkaAddressViperKey, f.Lookup(kafkaAddressFlag))
}
