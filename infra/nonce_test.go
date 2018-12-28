package infra

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStripeMutexConcurrent(t *testing.T) {
	mux := NewStripeMutex(10)
	counts := make([]int, 50)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mux.Lock(string(i % len(counts)))
			defer mux.Unlock(string(i % len(counts)))
			// We add some randomness in time execution
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Millisecond)
			counts[i%len(counts)]++
		}(i)
	}
	wg.Wait()

	for _, c := range counts {
		if c != rounds/len(counts) {
			t.Errorf("StripeMutex: Expected %v but got %v", rounds/len(counts), c)
		}
	}
}
