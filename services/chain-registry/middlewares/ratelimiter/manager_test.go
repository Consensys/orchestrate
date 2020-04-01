// +build unit

package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestSetLimit(t *testing.T) {
	rl := NewCooldownRateLimiter(rate.NewLimiter(rate.Inf, 100), 100*time.Millisecond)

	updated, _, _ := rl.setLimit(rate.Limit(10), false)
	assert.True(t, updated, "#1 Limit should have been updated")
	assert.Equal(t, rate.Limit(10), rl.limiter.Limit(), "#1 Limit should be correct")
	assert.Equal(t, 5, rl.limiter.Burst(), "#1 Burst should have beeen set and correct")

	updated, _, _ = rl.setLimit(rate.Limit(100), false)
	assert.True(t, updated, "#2 Limit should not have been updated")
	assert.Equal(t, rate.Limit(100), rl.limiter.Limit(), "#2 Limit should be correct")
	assert.Equal(t, 10, rl.limiter.Burst(), "#2 Burst should have been set and correct")

	// Test default behavior
	updated, _, _ = rl.setLimit(rate.Inf, true)
	assert.True(t, updated, "#3 Limit should have been updated")
	assert.Equal(t, rate.Limit(1000), rl.limiter.Limit(), "#3 Limit should  be correct")
	assert.Equal(t, 50, rl.limiter.Burst(), "#3 Burst should have been set and correct")
}

func TestSetLimitWithCooldown(t *testing.T) {
	rl := NewCooldownRateLimiter(rate.NewLimiter(rate.Inf, 100), 100*time.Millisecond)

	updated, _, _ := rl.setLimitWithCooldown(rate.Limit(10), false)
	assert.True(t, updated, "#1 Limit should have been updated")
	assert.Equal(t, rate.Limit(10), rl.limiter.Limit(), "#1 Limit should have beeen set")
	assert.Equal(t, 5, rl.limiter.Burst(), "#1 Burst should have beeen set")

	updated, _, _ = rl.setLimitWithCooldown(rate.Limit(20), false)
	assert.False(t, updated, "#2 Limit should not have been updated")
	assert.Equal(t, rate.Limit(10), rl.limiter.Limit(), "#2 Limit should not have been set due to cooldown")

	// Exhaust cooldown
	time.Sleep(101 * time.Millisecond)

	updated, _, _ = rl.setLimitWithCooldown(rate.Limit(1000), false)
	assert.True(t, updated, "#1 Limit should have been updated")
	assert.Equal(t, rate.Limit(1000), rl.limiter.Limit(), "#3 Limit should not have been set")
	assert.Equal(t, 50, rl.limiter.Burst(), "#3 Burst should have been set")
}
