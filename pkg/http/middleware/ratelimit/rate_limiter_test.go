// +build unit

package ratelimit

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	"golang.org/x/time/rate"
)

func TestInfura429ErrorLimit(t *testing.T) {
	infura429 := &Infura429Data{
		See: "test-see",
	}
	infura429.Rate.CurrentRPS = 13.333
	infura429.Rate.AllowedRPS = 10.0
	infura429.Rate.BackoffSeconds = 30.0

	b, _ := json.Marshal(infura429)
	limit, delay := infura429ErrorLimit(json.RawMessage(b))
	assert.Equal(t, rate.Limit(9), limit, "Limit should be correct")
	assert.Equal(t, 30*time.Second, delay, "Delay should be correct")

	limit, delay = infura429ErrorLimit(json.RawMessage(`Rate Limit`))
	assert.Equal(t, rate.Inf, limit, "Limit should be correct")
	assert.Equal(t, time.Duration(0), delay, "Delay should be correct")

}

type Mock409Handler struct{}

func (h *Mock409Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusTooManyRequests)
}

func TestServeHTTPIntercept409(t *testing.T) {
	nextH := &Mock409Handler{}
	cfg := &dynamic.RateLimit{
		MaxDelay:     time.Second,
		DefaultDelay: 30 * time.Second,
	}

	limiter := NewCooldownRateLimiter([]float64{math.MaxFloat64}, 100*time.Millisecond)

	rl := New(limiter, cfg)
	h := rl.Handler(nextH)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("", "", nil)
	h.ServeHTTP(rec, req)

	// Response should be 409 with Retry-After header Set
	resp := rec.Result()
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "Resp status code should be correct")
	assert.Equal(t, "30", resp.Header.Get("Retry-After"), "Retry-After header should have been set")
	assert.Equal(t, "30s", resp.Header.Get("X-Retry-In"), "X-Retry-In header should have been set")

	// Limiter should have been updated
	assert.Equal(t, rate.Limit(1000), rl.limiter.Limit(), "Limit should have been updated")
}

type MockBurstHandler struct {
	counter *int32
}

func (h *MockBurstHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	atomic.AddInt32(h.counter, 1)
}

func TestServeHTTPBurst(t *testing.T) {
	var counter int32
	nextH := &MockBurstHandler{&counter}
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	builder := NewBuilder(
		NewManager(cache),
	)
	cfg := &dynamic.RateLimit{
		MaxDelay:     time.Second,
		DefaultDelay: 30 * time.Second,
		Cooldown:     100 * time.Millisecond,
		Limits:       []float64{5},
	}

	mid, _, _ := builder.Build(context.Background(), "test", cfg)
	rl := mid(nextH)

	rounds := 20
	records := make(chan *httptest.ResponseRecorder, rounds)
	wg := &sync.WaitGroup{}
	wg.Add(rounds)
	go func() {
		wg.Wait()
		close(records)
	}()

	for i := 0; i < rounds; i++ {
		go func() {
			req, _ := http.NewRequest("", "", nil)
			record := httptest.NewRecorder()
			rl.ServeHTTP(record, req)
			records <- record
			wg.Done()
		}()
	}

	for record := range records {
		resp := record.Result()
		if resp.StatusCode == http.StatusTooManyRequests {
			if resp.Header.Get("Retry-After") == "" {
				b, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err, "Body should be readable")
				assert.Equal(t, "No bursty traffic allowed", string(b), "Error message should be correct")
			}
		}
	}

	assert.Equal(t, counter, int32(5), "Handler should have been protected from burst")
}
