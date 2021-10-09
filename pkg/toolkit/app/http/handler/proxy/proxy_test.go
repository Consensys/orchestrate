// +build unit

package proxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/traefik/traefik/v2/pkg/testhelpers"
	"github.com/oxtoacart/bpool"
)

type staticTransport struct {
	res *http.Response
}

func (t *staticTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.res, nil
}

func testCfg(passHost bool) *dynamic.ReverseProxy {
	return &dynamic.ReverseProxy{
		PassHostHeader: utils.Bool(passHost),
	}
}

var testBpool = bpool.NewBytePool(32, 1024)

func BenchmarkProxy(b *testing.B) {
	res := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	w := httptest.NewRecorder()
	req := testhelpers.MustNewRequest(http.MethodGet, "http://foo.bar/", nil)

	proxy, err := New(
		testCfg(false),
		&staticTransport{res},
		testBpool,
		nil,
	)

	if err != nil {
		b.Errorf("Could not build proxy: %v", err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		proxy.ServeHTTP(w, req)
	}
}
