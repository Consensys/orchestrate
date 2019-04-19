package sarama

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKafkaAddresses(t *testing.T) {
	name := "kafka.addresses"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	KafkaAddresses(flgs)

	expected := []string{
		"localhost:9092",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default")

	os.Setenv("KAFKA_ADDRESS", "localhost:9192")
	expected = []string{
		"localhost:9192",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From Environment Variable")
	os.Unsetenv("KAFKA_ADDRESS")

	args := []string{
		"--kafka-address=127.0.0.1:9091",
		"--kafka-address=127.0.0.2:9091",
	}
	flgs.Parse(args)

	expected = []string{
		"127.0.0.1:9091",
		"127.0.0.2:9091",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From flag")
}
