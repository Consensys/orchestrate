package ristretto

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/spf13/viper"
)

type Config struct {
	Cache    *ristretto.Config
	CacheTTL *time.Duration
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		Cache: &ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		},
	}

	cacheStr := vipr.GetDuration(CacheTTLViperKey)
	if cacheStr != time.Duration(0) {
		cfg.CacheTTL = &cacheStr
	}

	return cfg
}
