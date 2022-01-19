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
	"github.com/oxtoacart/bpool"
	"github.com/stretchr/testify/assert"
	"github.com/traefik/traefik/v2/pkg/testhelpers"
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
	w := httptest.NewRecorder()
	req := testhelpers.MustNewRequest(http.MethodGet, "http://foo.bar/", nil)

	res := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Header:     http.Header{},
		Request:    req,
	}

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

func TestProxyForward(t *testing.T) {
	proxyURLs := []string{"http://example-chain.es"}
	servers := make([]*dynamic.Server, 0)
	for _, chainURL := range proxyURLs {
		servers = append(servers, &dynamic.Server{
			URL: chainURL,
		})
	}

	
	w := httptest.NewRecorder()
	req := testhelpers.MustNewRequest(http.MethodGet, "http://proxy.es", nil)
	res := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Header:     http.Header{},
		Request:    req,
	}

	proxy, err := New(
		&dynamic.ReverseProxy{
			PassHostHeader: utils.Bool(false),
			LoadBalancer: &dynamic.LoadBalancer{
				Servers: servers,
			},
		},
		&staticTransport{res},
		testBpool,
		func(r *http.Response) error {
			assert.Equal(t, "proxy.es", r.Request.URL.Host)
			assert.Equal(t, "", r.Request.URL.RawPath)
			return nil
		},
	)

	assert.NoError(t, err)
	proxy.ServeHTTP(w, req)
}
