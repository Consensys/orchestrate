// +build unit

package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func NewRedisMock() *miniredis.Miniredis {
	mredis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return mredis
}

var testKey = "test-key"

func TestNonceManager(t *testing.T) {
	mredis := NewRedisMock()
	conf := &Config{
		Expiration: 1,
		Host:       mredis.Host(),
		Port:       mredis.Port(),
	}

	pool, _ := NewPool(conf)
	nm = NewNonceManager(pool, conf)

	n, ok, err := nm.GetLastAttributed(testKey)
	assert.NoError(t, err, "When manager is empty: GetLastAttributed should not error")
	assert.False(t, ok, "When manager is empty: GetLastAttributed should not find nonce")
	assert.Equal(t, uint64(0), n, "When manager is empty: GetLastAttributed should returned null nonce")

	err = nm.SetLastAttributed(testKey, 10)
	assert.NoError(t, err, "When manager is empty: SetLastAttributed should not error")

	n, ok, err = nm.GetLastAttributed(testKey)
	assert.NoError(t, err, "When last attributed has been set: GetLastAttributed should not error")
	assert.True(t, ok, "When last attributed has been set: GetLastAttributed should find nonce")
	assert.Equal(t, uint64(10), n, "When last attributed has been set: GetLastAttributed should returned non zero nonce")

	err = nm.IncrLastAttributed(testKey)
	assert.NoError(t, err, "When last attributed has been set: IncrLastAttributed should not error")

	n, _, _ = nm.GetLastAttributed(testKey)
	assert.Equal(t, uint64(11), n, "When last attributed has been incremented: GetLastAttributed should returned incremented nonce")
}

func TestNonceNonceSender(t *testing.T) {
	mredis := NewRedisMock()
	conf := &Config{
		Expiration: 1,
		Host:       mredis.Host(),
		Port:       mredis.Port(),
	}

	pool, _ := NewPool(conf)
	nm = NewNonceManager(pool, conf)

	n, ok, err := nm.GetLastSent(testKey)
	assert.NoError(t, err, "When manager is empty: GetLastSent should not error")
	assert.False(t, ok, "When manager is empty: GetLastSent should not find nonce")
	assert.Equal(t, uint64(0), n, "When manager is empty: GetLastSent should returned null nonce")

	err = nm.SetLastSent(testKey, 10)
	assert.NoError(t, err, "When manager is empty: SetLastSent should not error")

	n, ok, err = nm.GetLastSent(testKey)
	assert.NoError(t, err, "When last sent has been set: GetLastSent should not error")
	assert.True(t, ok, "When last sent has been set: GetLastSent should find nonce")
	assert.Equal(t, uint64(10), n, "When last sent has been set: GetLastSent should returned non zero nonce")

	err = nm.IncrLastSent(testKey)
	assert.NoError(t, err, "When last sent has been set: IncrLastAttributed should not error")

	n, _, _ = nm.GetLastSent(testKey)
	assert.Equal(t, uint64(11), n, "When last sent has been incremented: GetLastAttributed should returned incremented nonce")

	recovering, _ := nm.IsRecovering(testKey)
	assert.False(t, recovering, "When recovery status has not been updated: IsRecovering should be false")

	_ = nm.SetRecovering(testKey, true)

	recovering, _ = nm.IsRecovering(testKey)
	assert.True(t, recovering, "When recovery status has been set: IsRecovering should be true")

}
