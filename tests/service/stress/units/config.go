package units

import (
	"time"
)

type WorkloadConfig struct {
	nAccounts              int
	waitForEnvelopeTimeout time.Duration
}

func NewWorkloadConfig(nAccounts int, waitForEnvelopeTimeout time.Duration) *WorkloadConfig {
	return &WorkloadConfig{
		nAccounts:              nAccounts,
		waitForEnvelopeTimeout: waitForEnvelopeTimeout,
	}
}
