package postgres

import (
	"net"
)

func Dialer(cfg *Config) *net.Dialer {
	return &net.Dialer{
		Timeout:   cfg.DialTimeout,
		KeepAlive: cfg.KeepAliveInterval,
	}
}
