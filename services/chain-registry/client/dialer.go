package client

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

func DialWithDefaultOptions(url string) *HTTPClient {
	conf := NewConfig(url)
	return NewHTTPClient(http.NewClient(), conf)
}
