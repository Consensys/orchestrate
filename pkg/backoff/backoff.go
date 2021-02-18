package backoff

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type BackOff interface {
	backoff.BackOff
	NewBackOff() backoff.BackOff
}

type backOff struct {
	backoff.BackOff
	newBackOffFunc func() backoff.BackOff
}

func (b backOff) NewBackOff() backoff.BackOff {
	return b.newBackOffFunc()
}

func ConstantBackOffWithMaxRetries(d time.Duration, maxRetries uint64) BackOff {
	newBckOff := func() backoff.BackOff {
		return backoff.WithMaxRetries(backoff.NewConstantBackOff(d), maxRetries)
	}
	return &backOff{newBckOff(), newBckOff}
}

func IncrementalBackOff(initInterval, interval, elapsed time.Duration) BackOff {
	newBckOff := func() backoff.BackOff {
		bckOff := backoff.NewExponentialBackOff()
		bckOff.InitialInterval = initInterval
		bckOff.MaxInterval = interval
		bckOff.MaxElapsedTime = elapsed
		return bckOff
	}
	return &backOff{newBckOff(), newBckOff}
}

func IncrementalBackOffWithMaxRetries(initInterval, interval time.Duration, maxRetries uint64) BackOff {
	newBckOff := func() backoff.BackOff {
		bckOff := backoff.NewExponentialBackOff()
		bckOff.InitialInterval = initInterval
		bckOff.MaxInterval = interval
		return backoff.WithMaxRetries(bckOff, maxRetries)
	}
	return &backOff{newBckOff(), newBckOff}
}
