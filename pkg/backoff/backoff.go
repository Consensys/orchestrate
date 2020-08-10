package backoff

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

func ConstantBackOffWithMaxRetries(d time.Duration, maxRetries uint64) backoff.BackOff {
	return backoff.WithMaxRetries(backoff.NewConstantBackOff(d), maxRetries)
}
