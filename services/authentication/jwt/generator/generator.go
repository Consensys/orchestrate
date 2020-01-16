package generator

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	fositejwt "github.com/ory/fosite/token/jwt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"
)

type JWTGenerator struct {
	ClaimsNamespace string
	privateKey      *rsa.PrivateKey
}

func New(multiTenancyEnabled bool, namespace, pemPrivateKey string) (*JWTGenerator, error) {
	pkey, err := LoadRsaPrivateKeyFromVar(pemPrivateKey)
	if err != nil {
		return nil, err
	}
	return &JWTGenerator{
		ClaimsNamespace: namespace,
		privateKey:      pkey,
	}, nil
}

func (j *JWTGenerator) GenerateAccessTokenWithTenantID(tenantID string) (string, error) {
	customClaims := map[string]interface{}{
		j.ClaimsNamespace: &jwt.OrchestrateClaims{
			TenantID: tenantID,
		}}
	return j.GenerateAccessToken(customClaims)
}

func (j *JWTGenerator) GenerateAccessToken(customClaims map[string]interface{}) (tokenValue string, err error) {
	jwtGenerator := &oauth2.DefaultJWTStrategy{
		JWTStrategy: &fositejwt.RS256JWTStrategy{
			PrivateKey: j.privateKey,
		},
	}
	tokenRequest := &fosite.Request{
		GrantedAudience: []string{"https://auth0.com/api/v2/"},
		GrantedScope:    []string{"read:users", "update:users", "create:users"},
		Client: &fosite.DefaultClient{
			ID:     "App-test",
			Secret: []byte("mysecret"),
		},
		Session: &oauth2.JWTSession{
			JWTClaims: &fositejwt.JWTClaims{
				Issuer:    "Orchestrate",
				Subject:   "e2e-test",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				Extra:     customClaims,
			},
			JWTHeader: &fositejwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				fosite.AccessToken: time.Now().UTC().Add(time.Hour),
			},
		},
	}
	// The access Token contain already the signature
	accessToken, _, err := jwtGenerator.GenerateAccessToken(context.Background(), tokenRequest)
	return accessToken, err
}

func (j *JWTGenerator) GenerateIDToken(customClaims map[string]interface{}) (tokenValue string, err error) {

	jwtGenerator := &openid.DefaultStrategy{
		JWTStrategy: &fositejwt.RS256JWTStrategy{
			PrivateKey: j.privateKey,
		},
	}

	tokenRequest := fosite.NewAccessRequest(&openid.DefaultSession{
		Claims: &fositejwt.IDTokenClaims{
			Issuer:   "Orchestrate",
			Subject:  "e2e-test",
			Audience: []string{"https://auth0.com/api/v2/"},
			IssuedAt: time.Now().UTC(),
			Extra:    customClaims,
		},
		Headers: &fositejwt.Headers{},
	})

	idToken, err := jwtGenerator.GenerateIDToken(context.Background(), tokenRequest)

	return idToken, err
}

// Try to load the clear private key from file, if not found, log error and generate a new RSA Private Key
func LoadRSAPrivateKeyFromFile(rsaPrivateKeyLocation string) (*rsa.PrivateKey, error) {

	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		log.WithError(err).Errorf("No RSA private key found")
		return nil, err
	}

	privPem, _ := pem.Decode(priv)
	if privPem.Type != "PRIVATE KEY" {
		log.Errorf("RSA private key is of the wrong type: %s", privPem.Type)
		return nil, err
	}

	privPemBytes := privPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			log.Errorf("Unable to parse RSA private key, generating a temp one: %s", err)
			return nil, err
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Errorf("Unable to parse RSA private key, generating a temp one: %s", err)
		return nil, err
	}

	return privateKey, nil
}

func LoadRsaPrivateKeyFromVar(rawPrivateKey string) (*rsa.PrivateKey, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(rawPrivateKey)
	if err != nil {
		log.Errorf("Unable to Decode RSA private key")
		return nil, err
	}
	// Parse the private key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(decodedPrivateKey); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(decodedPrivateKey); err != nil { // note this returns type `interface{}`
			log.Errorf("Unable to parse RSA private key")
			return nil, err
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Errorf("Unable to parse RSA private key, generating a temp one")
		return nil, fmt.Errorf("invalid rsa private key")
	}
	return privateKey, nil
}
