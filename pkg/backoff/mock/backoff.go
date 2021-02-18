package mock

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type MockBackoff struct {
	hasRetried bool
}

func (b *MockBackoff) Reset() {}

func (b *MockBackoff) NextBackOff() time.Duration {
	b.hasRetried = true
	return backoff.Stop
}

func (b *MockBackoff) HasRetried() bool {
	return b.hasRetried
}

type MockIntervalBackoff struct {
	hasRetried bool
}

func (b *MockIntervalBackoff) Reset() {}

func (b *MockIntervalBackoff) NewBackOff() backoff.BackOff {
	return b
}

func (b *MockIntervalBackoff) NextBackOff() time.Duration {
	b.hasRetried = true
	return backoff.DefaultInitialInterval
}

func (b *MockIntervalBackoff) HasRetried() bool {
	return b.hasRetried
}
