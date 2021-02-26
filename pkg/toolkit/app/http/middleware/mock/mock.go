package mock

import (
	"net/http"
)

func NewMockMiddleware(mockHandler http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(rw, req)
			mockHandler.ServeHTTP(rw, req)
		})
	}
}
