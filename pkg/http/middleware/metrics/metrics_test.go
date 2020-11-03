// +build unit
// +build !race

package metrics

import (
	"context"
	"testing"
	"crypto/tls"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/metrics/mock"
	mockmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

func TestMetrics(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	mockHTTP := mock.NewMockHTTPMetrics(ctrlr)

	b := NewBuilder(mockHTTP)
	ctx := httputil.WithEntryPoint(httputil.WithService(context.Background(), "service-test"), "entrypoint-test")

	metrics, _, err := b.Build(ctx, "", nil)
	require.NoError(t, err, "Build")

	mockHandler := mockhandler.NewMockHandler(ctrlr)

	reqsCounter := mockmetrics.NewMockCounter(ctrlr)
	tlsReqsCounter := mockmetrics.NewMockCounter(ctrlr)
	reqLatencyHistogram := mockmetrics.NewMockHistogram(ctrlr)
	openConnsGauge := mockmetrics.NewMockGauge(ctrlr)

	// First call with a non TLS request
	mockHTTP.EXPECT().RequestsCounter().Return(reqsCounter)
	mockHTTP.EXPECT().RequestsLatencyHistogram().Return(reqLatencyHistogram)
	mockHTTP.EXPECT().OpenConnsGauge().Return(openConnsGauge)

	reqsCounter.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", multitenancy.DefaultTenant, "method", "GET", "protocol", "http", "code", "200").
		Return(reqsCounter)
	reqsCounter.EXPECT().Add(float64(1))

	reqLatencyHistogram.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", multitenancy.DefaultTenant, "method", "GET", "protocol", "http", "code", "200").
		Return(reqLatencyHistogram)
	reqLatencyHistogram.EXPECT().Observe(gomock.Any())

	openConnsGauge.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", multitenancy.DefaultTenant, "method", "GET", "protocol", "http").
		Return(openConnsGauge)
	openConnsGauge1 := openConnsGauge.EXPECT().Add(float64(1))
	openConnsGauge.EXPECT().Add(float64(-1)).After(openConnsGauge1)

	mockHandler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	rw := httptest.NewRecorder()
	metrics(mockHandler).ServeHTTP(rw, req)
	rw.WriteHeader(http.StatusOK)

	// Second call with a TLS request
	mockHTTP.EXPECT().RequestsCounter().Return(reqsCounter)
	mockHTTP.EXPECT().TLSRequestsCounter().Return(tlsReqsCounter)
	mockHTTP.EXPECT().RequestsLatencyHistogram().Return(reqLatencyHistogram)
	mockHTTP.EXPECT().OpenConnsGauge().Return(openConnsGauge)

	reqsCounter.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", "tenant-test", "method", "GET", "protocol", "http", "code", "200").
		Return(reqsCounter)
	reqsCounter.EXPECT().Add(float64(1))

	tlsReqsCounter.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", "tenant-test", "tls_version", "1.2", "tls_cipher", "TLS_RSA_WITH_AES_128_GCM_SHA256").
		Return(tlsReqsCounter)
	tlsReqsCounter.EXPECT().Add(float64(1))

	reqLatencyHistogram.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", "tenant-test", "method", "GET", "protocol", "http", "code", "200").
		Return(reqLatencyHistogram)
	reqLatencyHistogram.EXPECT().Observe(gomock.Any())

	openConnsGauge.EXPECT().
		With("entrypoint", "entrypoint-test", "service", "service-test", "tenant_id", "tenant-test", "method", "GET", "protocol", "http").
		Return(openConnsGauge)
	openConnsGauge1 = openConnsGauge.EXPECT().Add(float64(1))
	openConnsGauge.EXPECT().Add(float64(-1)).After(openConnsGauge1)

	mockHandler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

	req, _ = http.NewRequest(http.MethodGet, "http://test.com", nil)
	req = req.WithContext(multitenancy.WithTenantID(req.Context(), "tenant-test"))
	req.TLS = &tls.ConnectionState{
		Version:     tls.VersionTLS12,
		CipherSuite: tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}
	rw = httptest.NewRecorder()
	metrics(mockHandler).ServeHTTP(rw, req)
}
