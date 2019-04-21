package http

import (
	"context"
	"net/http"
	"sync"

	"github.com/spf13/viper"
)

var (
	server   *http.Server
	initOnce = &sync.Once{}
)

// InitServer initialize global HTTP server
func InitServer(ctx context.Context) {
	initOnce.Do(func() {
		if server != nil {
			return
		}

		// Initialize server
		server := &http.Server{}
		server.Addr = viper.GetString(hostnameViperKey)
	})
}

// GlobalServer return global HTTP server
func GlobalServer() *http.Server {
	return server
}

// SetGlobalServer sets global HTTP server
func SetGlobalServer(s *http.Server) {
	server = s
}
