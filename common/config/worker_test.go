package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestWorkerInTopic(t *testing.T) {
	name := "worker.in"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerInTopic(flgs, "TOPIC_IN", "test-in-topic")
	expected := "test-in-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("TOPIC_IN", "test-env-in-topic")
	expected = "test-env-in-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-in=test-flag-in-topic",
	}
	flgs.Parse(args)
	expected = "test-flag-in-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestWorkerOutTopic(t *testing.T) {
	name := "worker.out"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerOutTopic(flgs, "TOPIC_OUT", "test-out-topic")
	expected := "test-out-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("TOPIC_OUT", "test-env-out-topic")
	expected = "test-env-out-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-out=test-flag-out-topic",
	}
	flgs.Parse(args)
	expected = "test-flag-out-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestWorkerConsumerGroup(t *testing.T) {
	name := "worker.group"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerConsumerGroup(flgs, "CONSUMER_GROUP", "test-consumer-group")
	expected := "test-consumer-group"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("CONSUMER_GROUP", "test-env-consumer-group")
	expected = "test-env-consumer-group"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-group=test-flag-consumer-group",
	}
	flgs.Parse(args)
	expected = "test-flag-consumer-group"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestWorkerSlots(t *testing.T) {
	name := "worker.slots"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerSlots(flgs)
	expected := 20
	if viper.GetInt(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("WORKER_SLOTS", "125")
	expected = 125
	if viper.GetInt(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-slots=150",
	}
	flgs.Parse(args)
	expected = 150
	if viper.GetInt(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
