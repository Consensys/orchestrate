package http

import (
	"context"
	"net/http"

	"github.com/spf13/viper"
)

func init() {
	server = &http.Server{}
}

var server *http.Server

// InitServer initialize global HTTP server
func InitServer(ctx context.Context) {
	// Set server Address
	server.Addr = viper.GetString("http.hostname")
}

// GlobalServer return global HTTP server
func GlobalServer() *http.Server {
	return server
}

// SetGlobalServer sets global HTTP server
func SetGlobalServer(s *http.Server) {
	server = s
}
