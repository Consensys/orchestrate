// +build unit

package metrics

import (
	"net/http"
	"testing"
)

func TestMetricsHandler(t *testing.T) {
	mux := http.NewServeMux()

	Enhancer(func() error { return nil }, func() error { return nil })(mux)
}
