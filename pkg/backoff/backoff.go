package backoff

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

func ConstantBackOffWithMaxRetries(d time.Duration, maxRetries uint64) backoff.BackOff {
	return backoff.WithMaxRetries(backoff.NewConstantBackOff(d), maxRetries)
}

func IncrementalBackOff(interval, elapsed time.Duration) backoff.BackOff {
	bckOff := backoff.NewExponentialBackOff()
	bckOff.MaxInterval = interval
	bckOff.MaxElapsedTime = elapsed
	return bckOff
}
