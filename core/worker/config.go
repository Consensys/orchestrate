package worker

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(workerSlotsViperKey, workerSlotsDefault)
	viper.BindEnv(workerSlotsViperKey, workerSlotsEnv)
	viper.SetDefault(workerPartitionsViperKey, workerPartitionsDefault)
	viper.BindEnv(workerPartitionsViperKey, workerPartitionsEnv)
	viper.SetDefault(workerTimeoutViperKey, workerTimeoutDefault)
	viper.BindEnv(workerTimeoutViperKey, workerTimeoutEnv)
}

// Config is worker configuration
type Config struct {
	Slots      int64
	Partitions int64
}

// Validate ensure configuration is valid
func (c *Config) Validate() error {
	if c.Slots <= 0 {
		return fmt.Errorf("At least one worker slot is required")
	}

	if c.Partitions <= 0 {
		return fmt.Errorf("At least one partition is required")
	}

	return nil
}

// NewConfig create new worker configuration
func NewConfig() Config {
	config := Config{}
	config.Slots = viper.GetInt64("worker.slots")
	config.Partitions = viper.GetInt64("worker.partitions")
	return config
}

// InitFlags register flags for worker
func InitFlags(f *pflag.FlagSet) {
	Slots(f)
	Partitions(f)
}

var (
	workerSlotsFlag     = "worker-slots"
	workerSlotsViperKey = "worker.slots"
	workerSlotsDefault  = uint(20)
	workerSlotsEnv      = "WORKER_SLOTS"
)

// Slots register flag for Kafka server addresses
func Slots(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum number of messages the worker can treat concurrently.
Environment variable: %q`, workerSlotsEnv)
	f.Uint(workerSlotsFlag, workerSlotsDefault, desc)
	viper.BindPFlag(workerSlotsViperKey, f.Lookup(workerSlotsFlag))
}

var (
	workerPartitionsFlag     = "worker-partitions"
	workerPartitionsViperKey = "worker.partitions"
	workerPartitionsDefault  = uint(50)
	workerPartitionsEnv      = "WORKER_PARTITIONS"
)

// Partitions register flag for Kafka server addresses
func Partitions(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Number of partitions spawned by worker to treat messages in parallel.
Environment variable: %q`, workerPartitionsEnv)
	f.Uint(workerPartitionsFlag, workerPartitionsDefault, desc)
	viper.BindPFlag(workerPartitionsViperKey, f.Lookup(workerPartitionsFlag))
}

var (
	workerTimeoutFlag     = "worker-timeout"
	workerTimeoutViperKey = "worker.timeout"
	workerTimeoutDefault  = 60 * time.Second
	workerTimeoutEnv      = "WORKER_TIMEOUT"
)

// Timeout register flag for Kafka server addresses
func Timeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Maximum time for a message to be handled a message
Environment variable: %q`, workerTimeoutEnv)
	f.Duration(workerTimeoutFlag, workerTimeoutDefault, desc)
	viper.BindPFlag(workerTimeoutViperKey, f.Lookup(workerTimeoutFlag))
}
