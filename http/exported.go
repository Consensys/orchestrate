package http

import (
	"context"
	"net/http"
	"sync"

	"github.com/spf13/viper"
)

var (
	mux      *http.ServeMux
	server   *http.Server
	initOnce = &sync.Once{}
)

func init() {
	mux = http.NewServeMux()
}

// Init initialize global HTTP server
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if server != nil {
			return
		}

		// Initialize server
		server = &http.Server{}
		server.Addr = viper.GetString(hostnameViperKey)
		server.Handler = mux
	})
}

// Enhance allows to register new handlers on Global Server ServeMux
func Enhance(enhancer ServeMuxEnhancer) {
	enhancer(mux)
}

// GlobalServer return global HTTP server
func GlobalServer() *http.Server {
	return server
}

// SetGlobalServer sets global HTTP server
func SetGlobalServer(s *http.Server) {
	server = s
}
