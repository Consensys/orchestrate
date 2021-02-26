package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authjwt "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/jwt"
	authkey "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockHandler struct {
	served bool
}

func (h *MockHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.served = true
}

var (
	APIKey              = "test-key"
	bearer              = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly9hdXRoMC5jb20vYXBpL3YyLyJdLCJleHAiOjE1NzkxNjc0MTQsImh0dHA6Ly9vcmNoZXN0cmF0ZS5pbmZvIjp7InRlbmFudF9pZCI6ImI0OWVlMWJjLWYwZmEtNDMwZC04OWIyLWE0ZmQwZGM5ODkwNiIsInJvbGUiOiJ0ZXN0LXJvbGUifSwiaWF0IjoxNTc5MTYzODE0LCJpc3MiOiJPcmNoZXN0cmF0ZSIsImp0aSI6IjZlZmY3MzI0LTVkZTEtNDA2NS05NGNmLWU3ZWYzZTliYjg1MCIsIm5iZiI6MTU3OTE2MzgxNCwic2NwIjpbInJlYWQ6dXNlcnMiLCJ1cGRhdGU6dXNlcnMiLCJjcmVhdGU6dXNlcnMiXSwic3ViIjoiZTJlLXRlc3QifQ.fvlJcrCwbvj-W1VrfSzcn5F7LpsZ0xbOQTcCqVwwmyq8EOv5VwoV-geoX6tj4d0T2pew-6EK8DR-GrwXVjlo2LQQhYY_TRpnVHl1wDE1IvahExnh_0oPwpH3oKjsxbLPyM94bG-eIJGyInA3w-llCXR5WhOwccO4lKW4GaAXsj6TKGiowh_9HEw9jSN2y9OXGvUiE9_8n_5rp1Shp_vBMHJ-5usOozoaJdgl13Dln1YTqSl422CKb1UndBGRXayCfMpqnzLuURTYYspWOn3c6QTbjjMAm8ifZIl8rDrI8zl8j2FM1kHZt-5ZZe5zJv7rCGwPQviLnWQBqIVElJv6Tg"
	rawCert             = "MIIDYjCCAkoCCQC9pJWk7qdipjANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJGUjEOMAwGA1UEBwwFUGFyaXMxEjAQBgNVBAoMCUNvbnNlblN5czEQMA4GA1UECwwHUGVnYVN5czEuMCwGA1UEAwwlZTJlLXRlc3RzLm9yY2hlc3RyYXRlLmNvbnNlbnN5cy5wYXJpczAeFw0xOTEyMjcxNjI5MTdaFw0yMDEyMjYxNjI5MTdaMHMxCzAJBgNVBAYTAkZSMQ4wDAYDVQQHDAVQYXJpczESMBAGA1UECgwJQ29uc2VuU3lzMRAwDgYDVQQLDAdQZWdhU3lzMS4wLAYDVQQDDCVlMmUtdGVzdHMub3JjaGVzdHJhdGUuY29uc2Vuc3lzLnBhcmlzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo0NqWqI3TSi1uOBvCUquclWo4LcsYT21tNUXQ8YyqVYRSsiBv+ZKZBCjD8XklLPih40kFSe+r6DNca5/LH/okQIdc8nsQg+BLCkXeH2NFv+QYtPczAw4YhS6GVxJk3u9sfp8NavWBcQbD3MMDpehMOvhSl0zoP/ZlH6ErKHNtoQgUpPNVQGysNU21KpClmIDD/L1drsbq+rFiDrcVWaOLwGxr8SBd/0b4ngtcwH16RJaxcIXXT5AVia1CNdzmU5/AIg3OfgzvKn5AGrMZBsmGAiCyn4/P3PnuF81/WHukk5ETLnzOH+vC2elSmZ8y80HCGeqOiQ1rs66L936wX8cDwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCNcTs3n/Ps+yIZDH7utxTOaqpDTCB10MzPmb22UAal89couIT6R0fAu14p/LTkxdb2STDySsQY2/Lv6rPdFToHGUI9ZYOTYW1GOWkt1EAao9BzdsoJVwmTON6QnOBKy/9RxlhWP+XSWVsY0te6KYzS7rQyzQoJQeeBNMpUnjiQji9kKi5j9rbVMdjIb4HlmYrcE95ps+oFkyJoA1HLVytAeOjJPXGToNlv3k2UPJzOFUM0ujWWeBTyHMCmZ4RhlrfzDNffY5dlW82USjc5dBlzRyZalXSjhcVhK4asUodomVntrvCShp/8C9LpbQZ+ugFNE8J6neStWrhpRU9/sBJx"
	apikeyHeader        = "X-API-Key"
	authorizationHeader = "Authorization"
)

func TestAuth(t *testing.T) {
	testCases := []struct {
		desc                string
		path                string
		authorizationHeader string
		authorizationToken  string
		expectedStatusCode  int
		expectedServed      bool
	}{
		{
			"missing auth",
			"/b49ee1bc-f0fa-430d-89b2-a4fd0dc98906",
			"",
			"",
			http.StatusUnauthorized,
			false,
		},
		{
			"with API key",
			"/b49ee1bc-f0fa-430d-89b2-a4fd0dc98906",
			apikeyHeader,
			APIKey,
			http.StatusOK,
			true,
		},
		{
			"with JWT Token",
			"/b49ee1bc-f0fa-430d-89b2-a4fd0dc98906",
			authorizationHeader,
			"Bearer " + bearer,
			http.StatusOK,
			true,
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			// Create HTTP handler
			nextH := &MockHandler{}
			jwtChecker, err := authjwt.New(&authjwt.Config{
				ClaimsNamespace:      "http://orchestrate.info",
				SkipClaimsValidation: true,
				Certificate:          []byte(rawCert),
			})
			require.NoError(t, err)

			auth := New(
				jwtChecker,
				authkey.New(APIKey),
				true,
			)

			handler := mux.NewRouter()
			handler.PathPrefix("/").Handler(auth.Handler(nextH))

			req := httptest.NewRequest("GET", "http://example.com"+test.path, nil)
			if test.authorizationHeader == authorizationHeader {
				req.Header.Set(authorizationHeader, test.authorizationToken)
			}
			if test.authorizationHeader == apikeyHeader {
				req.Header.Set(apikeyHeader, test.authorizationToken)
			}

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)
			resp := w.Result()
			assert.Equal(t, test.expectedStatusCode, resp.StatusCode, "Status Code should be valid")
			assert.Equal(t, test.expectedServed, nextH.served, "Given handler should have been called or ignored")
		})
	}
}
