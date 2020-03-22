package jwt

import (
	"context"
	"crypto/rsa"
	"sync"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/certificate"
)

var (
	checker  *JWT
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if checker != nil {
			return
		}

		conf := NewConfig()

		//
		rawCert := viper.GetString(CertificateViperKey)
		if rawCert == "" {
			log.Infof("jwt: no certificate provided")
		} else {
			// Decode certificate provided in configuration
			cert, err := certificate.DecodeStringToCertificate(rawCert)
			if err != nil {
				log.WithError(err).Fatalf("jwt: invalid certificate")
			}

			// Cast certificate into an RSA public key
			pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
			if !ok {
				log.Fatalf("jwt: certificate is not an RSA public key")
			}

			conf.Key = func(token *jwt.Token) (interface{}, error) { return pubKey, nil }
		}

		checker = New(conf)
	})
}

// GlobalChecker returns global Authentication Manager
func GlobalChecker() *JWT {
	return checker
}

// SetGlobalAuth sets global Authentication Manager
func SetGlobalChecker(c *JWT) {
	checker = c
	log.Debug("authentication manager: set")
}
