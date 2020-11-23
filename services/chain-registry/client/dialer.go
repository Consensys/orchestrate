package client

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
)

func DialWithDefaultOptions(url string) *HTTPClient {
	conf := NewConfig(url)
	return NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}
