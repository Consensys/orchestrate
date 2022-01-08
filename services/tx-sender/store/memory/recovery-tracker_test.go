// +build unit

package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryTrackerMemory(t *testing.T) {
	rt := NewNonceRecoveryTracker()

	testKey := "recovery-tracker-memory"
	n := rt.Recovering(testKey)
	assert.Equal(t, uint64(0), n)

	rt.Recover(testKey)
	rt.Recover(testKey)
	n = rt.Recovering(testKey)
	assert.Equal(t, uint64(2), n)

	rt.Recovered(testKey)
	n = rt.Recovering(testKey)
	assert.Equal(t, uint64(0), n)
}
