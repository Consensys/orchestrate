package httpclient

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	initOnce = &sync.Once{}
	client   *http.Client
)

// Init initialize global gRPC server
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client == nil {
			client = NewClient()
		}

		log.Infof("http-client: ready")
	})
}

// GlobalClient gets global HTTP client
func GlobalClient() *http.Client {
	return client
}

// SetGlobalClient sets global HTTP client
func SetGlobalClient(c *http.Client) {
	client = c
}
