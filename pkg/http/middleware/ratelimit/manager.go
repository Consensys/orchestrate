package ratelimit

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"golang.org/x/time/rate"
)

var defaultLimits = []rate.Limit{
	rate.Inf,
	rate.Limit(1000),
	rate.Limit(100),
	rate.Limit(50),
	rate.Limit(10),
	rate.Limit(5),
	rate.Limit(1),
	rate.Limit(0),
}

// CooldownRateLimiter is a rate limiter that respects a cooldown between every updates
type CooldownRateLimiter struct {
	limiter *rate.Limiter

	cooldown time.Duration

	mux             *sync.Mutex
	firstUpdateOnce *sync.Once
	updateCooldown  *time.Ticker

	limitIdx int
	limits   []rate.Limit
}

func NewCooldownRateLimiter(limits []float64, cooldown time.Duration) *CooldownRateLimiter {
	rl := &CooldownRateLimiter{

		cooldown:        cooldown,
		mux:             &sync.Mutex{},
		firstUpdateOnce: &sync.Once{},
	}

	if len(limits) == 0 {
		rl.limits = defaultLimits
	} else {
		for _, l := range limits {
			rl.limits = append(rl.limits, rate.Limit(l))
		}
	}

	// Create limiter
	rl.limiter = rate.NewLimiter(rl.limits[0], 0)
	rl.setBurst()

	return rl
}

func (l *CooldownRateLimiter) Reserve() *rate.Reservation {
	return l.limiter.Reserve()
}

func (l *CooldownRateLimiter) Burst() int {
	return l.limiter.Burst()
}

func (l *CooldownRateLimiter) Limit() rate.Limit {
	return l.limiter.Limit()
}

func (l *CooldownRateLimiter) SetLimit(limit rate.Limit, useDefault bool) (updated bool, oldLimit, newLimit rate.Limit) {
	return l.setLimitWithCooldown(limit, useDefault)
}

func (l *CooldownRateLimiter) setLimitWithCooldown(limit rate.Limit, useDefault bool) (updated bool, oldLimit, newLimit rate.Limit) {
	// There is no proper to way to create a ticker with instant first tick
	// c.f https://github.com/golang/go/issues/17601
	// So we use sync.Once as a way around
	l.firstUpdateOnce.Do(func() {
		updated, oldLimit, newLimit = l.setLimit(limit, useDefault)
		l.updateCooldown = time.NewTicker(l.cooldown)
	})

	// SetLimit
	select {
	case <-l.updateCooldown.C:
		return l.setLimit(limit, useDefault)
	default:
		return
	}
}

func (l *CooldownRateLimiter) setLimit(limit rate.Limit, useDefault bool) (updated bool, oldLimit, newLimit rate.Limit) {
	l.mux.Lock()
	defer l.mux.Unlock()

	oldLimit = l.limiter.Limit()
	if useDefault {
		l.setLimitDefault()
	} else if limit != oldLimit {
		l.limiter.SetLimit(limit)
	}

	if oldLimit != l.limiter.Limit() {
		l.setBurst()
		return true, oldLimit, l.limiter.Limit()
	}

	return false, 0, 0
}

func (l *CooldownRateLimiter) setLimitDefault() {
	if l.limitIdx < len(defaultLimits)-1 {
		l.limitIdx++
		l.limiter.SetLimit(defaultLimits[l.limitIdx])
	}
}

func (l *CooldownRateLimiter) setBurst() {
	limit := l.limiter.Limit()

	// Set burst
	switch {
	case limit <= 10:
		l.limiter.SetBurst(5)
	case limit <= 100:
		l.limiter.SetBurst(10)
	case limit <= 1000:
		l.limiter.SetBurst(50)
	default:
		l.limiter.SetBurst(100)
	}
}

type Manager struct {
	cache *ristretto.Cache
}

func NewManager(cache *ristretto.Cache) *Manager {
	return &Manager{
		cache: cache,
	}
}

func (m *Manager) Get(key string) (*CooldownRateLimiter, bool) {
	if v, ok := m.cache.Get(key); ok {
		return v.(*CooldownRateLimiter), true
	}

	return nil, false
}

func (m *Manager) Set(key string, limiter *CooldownRateLimiter) bool {
	return m.cache.Set(key, limiter, 0)
}
