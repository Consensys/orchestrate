package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNonceSenderMemory(t *testing.T) {
	ns := NewNonceSender(time.Second)

	testKey := "nonce-sender-memory"
	n, ok, err := ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)

	err = ns.SetLastSent(testKey, 10)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	n, ok, err = ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, uint64(10), n)

	err = ns.IncrLastSent(testKey)
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	n, _, _ = ns.GetLastSent(testKey)
	assert.Equal(t, uint64(11), n)

	err = ns.DeleteLastSent(testKey)
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	n, ok, err = ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)

	err = ns.SetLastSent(testKey, 10)
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	n, ok, err = ns.GetLastSent(testKey)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), n)
}
