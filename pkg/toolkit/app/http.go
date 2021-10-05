package app

import (
	"net/http"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	http2 "github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/spf13/viper"
)

func NewHTTPClient(vipr *viper.Viper) *http.Client {
	cfg := http2.NewDefaultConfig()
	if vipr != nil {
		cfg.XAPIKey = vipr.GetString(key.APIKeyViperKey)
	}

	return http2.NewClient(cfg)
}
