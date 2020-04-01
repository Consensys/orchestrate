// +build unit

package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testKey = "test-key"

func TestNonceNonceAttributor(t *testing.T) {
	nm = NewNonceManager()
	n, ok, err := nm.GetLastAttributed(testKey)
	assert.Nil(t, err, "When manager is empty: GetLastAttributed should not error")
	assert.False(t, ok, "When manager is empty: GetLastAttributed should not find nonce")
	assert.Equal(t, uint64(0), n, "When manager is empty: GetLastAttributed should returned null nonce")

	err = nm.SetLastAttributed(testKey, 10)
	assert.Nil(t, err, "When manager is empty: SetLastAttributed should not error")

	n, ok, err = nm.GetLastAttributed(testKey)
	assert.Nil(t, err, "When last attributed has been set: GetLastAttributed should not error")
	assert.True(t, ok, "When last attributed has been set: GetLastAttributed should find nonce")
	assert.Equal(t, uint64(10), n, "When last attributed has been set: GetLastAttributed should returned non zero nonce")

	err = nm.IncrLastAttributed(testKey)
	assert.Nil(t, err, "When last attributed has been set: IncrLastAttributed should not error")

	n, _, _ = nm.GetLastAttributed(testKey)
	assert.Equal(t, uint64(11), n, "When last attributed has been incremented: GetLastAttributed should returned incremented nonce")
}

func TestNonceManagerMultiGoRoutine(t *testing.T) {
	nm = NewNonceManager()
	wait := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		// As required by NonceManager each goroutine will be Incrementing nonces
		// on a dedicated key
		key := fmt.Sprintf("test-key-%v", i)
		_ = nm.SetLastAttributed(key, 0)
		wait.Add(1)
		go func(key string) {
			for j := 0; j < 100; j++ {
				_ = nm.IncrLastAttributed(key)
			}
			wait.Done()
		}(key)
	}

	wait.Wait()
	for i := 0; i < 10; i++ {
		n, _, _ := nm.GetLastAttributed(fmt.Sprintf("test-key-%v", i))
		assert.Equal(t, uint64(100), n, "Final nonce for %q should be correct", i)
	}
}

func TestNonceNonceSender(t *testing.T) {
	nm = NewNonceManager()
	testKey := "test-key"
	n, ok, err := nm.GetLastSent(testKey)
	assert.Nil(t, err, "When manager is empty: GetLastSent should not error")
	assert.False(t, ok, "When manager is empty: GetLastSent should not find nonce")
	assert.Equal(t, uint64(0), n, "When manager is empty: GetLastSent should returned null nonce")

	err = nm.SetLastSent(testKey, 10)
	assert.Nil(t, err, "When manager is empty: SetLastSent should not error")

	n, ok, err = nm.GetLastSent(testKey)
	assert.Nil(t, err, "When last sent has been set: GetLastSent should not error")
	assert.True(t, ok, "When last sent has been set: GetLastSent should find nonce")
	assert.Equal(t, uint64(10), n, "When last sent has been set: GetLastSent should returned non zero nonce")

	err = nm.IncrLastSent(testKey)
	assert.Nil(t, err, "When last sent has been set: IncrLastAttributed should not error")

	n, _, _ = nm.GetLastSent(testKey)
	assert.Equal(t, uint64(11), n, "When last sent has been incremented: GetLastAttributed should returned incremented nonce")

	recovering, _ := nm.IsRecovering(testKey)
	assert.False(t, recovering, "When recovery status has not been updated: IsRecovering should be false")

	_ = nm.SetRecovering(testKey, true)

	recovering, _ = nm.IsRecovering(testKey)
	assert.True(t, recovering, "When recovery status has been set: IsRecovering should be true")
}
