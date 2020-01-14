package generator

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"time"

	"google.golang.org/grpc"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/token"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

type JWTGenerator struct {
	multitenancy    bool
	tenantNamespace string
	privateKey      *rsa.PrivateKey
}

func New(multiTenancyEnabled bool, tenantNamespace, pemPrivateKey string) *JWTGenerator {
	if multiTenancyEnabled {
		return &JWTGenerator{
			multitenancy:    multiTenancyEnabled,
			tenantNamespace: tenantNamespace,
			privateKey:      LoadRsaPrivateKeyFromVar(pemPrivateKey),
		}
	}
	return &JWTGenerator{
		multitenancy: multiTenancyEnabled,
	}
}

func (j *JWTGenerator) InjectAccessTokenIntoEnvelop(tenantID string, e *envelope.Envelope) error {
	if j.multitenancy {
		log.Debugf("inject AccessToken in the envelope")
		customClaims := map[string]interface{}{j.tenantNamespace + authentication.TenantIDKey: tenantID}
		accessToken, err := j.GenerateAccessToken(customClaims)
		if err != nil {
			log.Errorf("Unable to GenerateAccessToken: %s", err)
			return err
		}
		e.SetMetadataValue(token.OauthToken, accessToken)
	}
	return nil
}

func (j *JWTGenerator) InjectAccessTokenIntoGRPC(tenantID string) (grpc.CallOption, error) {
	var accessToken string
	var err error
	if j.multitenancy {
		log.Debugf("inject AccessToken in the gRPC")
		customClaims := map[string]interface{}{j.tenantNamespace + authentication.TenantIDKey: tenantID}
		accessToken, err = j.GenerateAccessToken(customClaims)
		if err != nil {
			log.Errorf("Unable to GenerateAccessToken: %s", err)
			return nil, err
		}
	} else {
		accessToken = ""
	}
	perRPCCredentials := token.NewJWTAccessFromEnvelope(accessToken)

	return grpc.PerRPCCredentials(perRPCCredentials), nil
}

func (j *JWTGenerator) GenerateAccessToken(customClaims map[string]interface{}) (tokenValue string, err error) {
	jwtGenerator := &oauth2.DefaultJWTStrategy{
		JWTStrategy: &jwt.RS256JWTStrategy{
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
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "Orchestrate",
				Subject:   "e2e-test",
				IssuedAt:  time.Now().UTC(),
				NotBefore: time.Now().UTC(),
				Extra:     customClaims,
			},
			JWTHeader: &jwt.Headers{
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
		JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: j.privateKey,
		},
	}

	tokenRequest := fosite.NewAccessRequest(&openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:   "Orchestrate",
			Subject:  "e2e-test",
			Audience: []string{"https://auth0.com/api/v2/"},
			IssuedAt: time.Now().UTC(),
			Extra:    customClaims,
		},
		Headers: &jwt.Headers{},
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

func LoadRsaPrivateKeyFromVar(rawPrivateKey string) *rsa.PrivateKey {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(rawPrivateKey)
	if err != nil {
		log.Errorf("Unable to Decode RSA private key")
		return nil
	}
	// Parse the private key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(decodedPrivateKey); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(decodedPrivateKey); err != nil { // note this returns type `interface{}`
			log.Errorf("Unable to parse RSA private key")
			return nil
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Errorf("Unable to parse RSA private key, generating a temp one")
		return nil
	}
	return privateKey
}
