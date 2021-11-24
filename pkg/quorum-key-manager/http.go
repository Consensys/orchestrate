package quorumkeymanager

import (
	"crypto/tls"
	"net/http"

	http2 "github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/spf13/viper"
)

func NewHTTPClient(vipr *viper.Viper) (*http.Client, error) {
	cfg := http2.NewDefaultConfig()
	// Support user's JWT forwarding
	cfg.AuthHeaderForward = true

	cfg.InsecureSkipVerify = vipr.GetBool(tlsSkipVerifyViperKey)
	APIKey := vipr.GetString(AuthAPIKeyViperKey)
	if APIKey != "" {
		cfg.AuthHeaderForward = false
		cfg.Authorization = "Basic " + APIKey
	}

	certFile := vipr.GetString(AuthClientTLSCertViperKey)
	keyFile := vipr.GetString(AuthClientTLSKeyViperKey)
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		cfg.ClientCert = &cert
	}

	return http2.NewClient(cfg), nil
}
