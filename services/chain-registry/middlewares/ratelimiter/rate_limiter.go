package ratelimiter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/server"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *CooldownRateLimiter

	maxDelay     time.Duration
	defaultDelay time.Duration

	next http.Handler
}

func New(rlManager *Manager, maxDelay, defaultDelay, cooldown time.Duration, routerName string, next http.Handler) *RateLimiter {
	// Retrieve limiter from cache
	limiter, ok := rlManager.Get(routerName)
	if !ok {
		limiter = NewCooldownRateLimiter(rate.NewLimiter(defaultLimits[0], 100), cooldown)
		rlManager.Set(routerName, limiter)
	}

	return &RateLimiter{
		maxDelay:     maxDelay,
		defaultDelay: defaultDelay,
		limiter:      limiter,
		next:         next,
	}
}

type JSONRpcResponse struct {
	Version string        `json:"jsonrpc,omitempty"`
	Error   *JSONRpcError `json:"error,omitempty"`
}

type JSONRpcError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (rl *RateLimiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	resa := rl.limiter.Reserve()
	if !resa.OK() {
		http.Error(rw, "No bursty traffic allowed", http.StatusTooManyRequests)
		return
	}

	// Delay
	delay := resa.Delay()
	if delay > rl.maxDelay {
		resa.Cancel()
		rl.serveDelayError(rw, delay)
		return
	}
	time.Sleep(delay)

	// Wrap response writer to intercept 429 Errors
	rwinterceptor := server.NewResponseWriterInterceptor(rw, interceptor429Error)

	rl.next.ServeHTTP(rwinterceptor, req)

	if rwinterceptor.Interceptor() != nil {
		// We intercepted a 429 error
		// So we prepare for rate limit update
		limit := rate.Inf
		var delay time.Duration

		// We first try to decode response body
		decoder := json.NewDecoder(rwinterceptor.Interceptor().(io.Reader))
		resp := &JSONRpcResponse{}
		err := decoder.Decode(resp)
		if err == nil && resp.Error != nil {
			limit, delay = infura429ErrorLimit(resp.Error.Data)
		}

		// Set limit
		updated, oldLimit, newLimit := rl.limiter.SetLimit(limit, limit == rate.Inf)
		if updated {
			log.FromContext(req.Context()).
				WithField("limit.old", oldLimit).
				WithField("limit.new", newLimit).
				WithField("burst.new", rl.limiter.Burst()).
				Warning("Rate limit updated")
		}

		// Set delay
		if delay == 0 {
			retryAfter, _ := strconv.ParseInt(
				rwinterceptor.Interceptor().Header().Get("Retry-After"),
				10, 64,
			)

			if retryAfter != 0 {
				// Retry-After returns a number of second, we convert to a Duration (base unit Nanosecond)
				delay = time.Duration(1000000000 * retryAfter)
			} else {
				delay = rl.defaultDelay
			}
		}

		rl.serveDelayError(rw, delay)
	}
}

func (rl *RateLimiter) serveDelayError(rw http.ResponseWriter, delay time.Duration) {
	rw.Header().Set("Retry-After", fmt.Sprintf("%v", delay.Seconds()))
	rw.Header().Set("X-Retry-In", delay.String())
	rw.WriteHeader(http.StatusTooManyRequests)
	_, _ = rw.Write([]byte(http.StatusText(http.StatusTooManyRequests)))
}

func interceptor429Error(statusCode int, header http.Header) server.WriterInterceptor {
	if statusCode == http.StatusTooManyRequests {
		return server.NewBytesBufferInterceptor(header)
	}
	return nil
}

type Infura429Data struct {
	See  string `json:"see,omitempty"`
	Rate struct {
		CurrentRPS     float64 `json:"current_rps,omitempty"`
		AllowedRPS     float64 `json:"allowed_rps,omitempty"`
		BackoffSeconds float64 `json:"backoff_seconds,omitempty"`
	} `json:"rate,omitempty"`
}

func infura429ErrorLimit(data json.RawMessage) (rate.Limit, time.Duration) {
	infura429 := &Infura429Data{}
	err := json.Unmarshal(data, infura429)
	if err != nil || infura429.Rate.AllowedRPS == 0 {
		return rate.Inf, 0
	}
	return rate.Limit(0.9 * infura429.Rate.AllowedRPS), time.Duration(1000000000 * infura429.Rate.BackoffSeconds)
}
