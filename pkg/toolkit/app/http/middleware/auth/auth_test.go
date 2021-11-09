package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt/mock"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/golang/mock/gomock"

	authjwt "github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt"
	authkey "github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockHandler struct {
	served bool
}

func (h *MockHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	h.served = true
}

const (
	APIKey              = "test-key"
	bearer              = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly9hdXRoMC5jb20vYXBpL3YyLyJdLCJleHAiOjE1NzkxNjc0MTQsImh0dHA6Ly9vcmNoZXN0cmF0ZS5pbmZvIjp7InRlbmFudF9pZCI6ImI0OWVlMWJjLWYwZmEtNDMwZC04OWIyLWE0ZmQwZGM5ODkwNiIsInJvbGUiOiJ0ZXN0LXJvbGUifSwiaWF0IjoxNTc5MTYzODE0LCJpc3MiOiJPcmNoZXN0cmF0ZSIsImp0aSI6IjZlZmY3MzI0LTVkZTEtNDA2NS05NGNmLWU3ZWYzZTliYjg1MCIsIm5iZiI6MTU3OTE2MzgxNCwic2NwIjpbInJlYWQ6dXNlcnMiLCJ1cGRhdGU6dXNlcnMiLCJjcmVhdGU6dXNlcnMiXSwic3ViIjoiZTJlLXRlc3QifQ.fvlJcrCwbvj-W1VrfSzcn5F7LpsZ0xbOQTcCqVwwmyq8EOv5VwoV-geoX6tj4d0T2pew-6EK8DR-GrwXVjlo2LQQhYY_TRpnVHl1wDE1IvahExnh_0oPwpH3oKjsxbLPyM94bG-eIJGyInA3w-llCXR5WhOwccO4lKW4GaAXsj6TKGiowh_9HEw9jSN2y9OXGvUiE9_8n_5rp1Shp_vBMHJ-5usOozoaJdgl13Dln1YTqSl422CKb1UndBGRXayCfMpqnzLuURTYYspWOn3c6QTbjjMAm8ifZIl8rDrI8zl8j2FM1kHZt-5ZZe5zJv7rCGwPQviLnWQBqIVElJv6Tg"
	apikeyHeader        = "X-API-Key"
	authorizationHeader = "Authorization"
)

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := mock.NewMockValidator(ctrl)

	mockValidator.EXPECT().ValidateToken(gomock.Any(), bearer).Return(&entities.UserClaims{TenantID: "tenant"}, nil)

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

			auth := New(authjwt.New(mockValidator), authkey.New(APIKey), true)

			handler := mux.NewRouter()
			handler.PathPrefix("/").Handler(auth.Handler(nextH))

			req := httptest.NewRequest("GET", "https://example.com"+test.path, nil)
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
