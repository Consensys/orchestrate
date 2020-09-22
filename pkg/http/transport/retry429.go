package transport

import (
	"net/http"
	"strconv"
	"time"
)

// AuthHeaderTransport is an internal Transport for Orchestrate
type Retry429Transport struct {
	T http.RoundTripper
}

func NewRetry429Transport() Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &Retry429Transport{
			T: nxt,
		}
	}
}

func (t *Retry429Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		resp, err := t.T.RoundTrip(req)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter, _ := strconv.ParseInt(
				resp.Header.Get("Retry-After"),
				10, 64,
			)

			if retryAfter > 0 {
				select {
				case <-time.After(time.Duration(1000000000 * retryAfter)):
					continue
				case <-req.Context().Done():
					return nil, req.Context().Err()
				}
			}
		}

		return resp, nil
	}
}
